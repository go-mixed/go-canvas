package font

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/sfnt"
)

const fontIndexCacheVersion = 2

type fontIndexCache struct {
	Version int                       `json:"version"`
	Entries map[string]fontIndexEntry `json:"entries"`
}

type fontIndexEntry struct {
	Family    string     `json:"family"`
	SubFamily string     `json:"sub_family"`
	Bold      FontWeight `json:"bold"`
	Italic    bool       `json:"italic"`
	Path      string     `json:"path"`
	Size      int64      `json:"size"`
	ModTime   int64      `json:"mod_time_unix_nano"`
}

// LoadFonts 读取给的目录+系统目录的所有字体（带磁盘索引缓存）
// LoadFonts reads all fonts from user/system directories with disk index cache.
// userFontPaths 可以为字体路径、或者包含字体的目录
// userFontPaths can be either font files or directories containing fonts.
func LoadFonts(userFontPaths ...string) map[string][]*FontInfo {
	fs := &FontLibrary{
		indexCachePath: defaultFontIndexCachePath(),
	}
	return fs.loadFonts(userFontPaths...)
}

// loadFonts 从目录枚举字体并使用 JSON 索引缓存加速元信息加载。
// loadFonts enumerates fonts and uses a JSON index cache to speed up metadata loading.
func (fs *FontLibrary) loadFonts(userFontPaths ...string) map[string][]*FontInfo {
	paths := ListFont(append(userFontPaths, GetSystemFontDirectories()...))
	cache := loadFontIndexCacheFile(fs.indexCachePath)
	fontInfos := make(map[string][]*FontInfo)
	nextEntries := make(map[string]fontIndexEntry, len(paths))
	changed := false

	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			changed = true
			continue
		}

		if cached, ok := cache.Entries[path]; ok &&
			cached.Size == stat.Size() &&
			cached.ModTime == stat.ModTime().UnixNano() &&
			cached.Family != "" {
			info := fontInfoFromEntry(cached)
			fontInfos[info.Family] = append(fontInfos[info.Family], info)
			nextEntries[path] = cached
			continue
		}

		info, err := ReadFontInfo(path)
		if err != nil {
			changed = true
			continue
		}
		fontInfos[info.Family] = append(fontInfos[info.Family], info)
		nextEntries[path] = fontIndexEntry{
			Family:    info.Family,
			SubFamily: info.SubFamily,
			Bold:      info.Bold,
			Italic:    info.Italic,
			Path:      info.FontPath,
			Size:      stat.Size(),
			ModTime:   stat.ModTime().UnixNano(),
		}
		changed = true
	}

	if !changed && len(cache.Entries) == len(nextEntries) {
		return fontInfos
	}
	saveFontIndexCacheFile(fs.indexCachePath, fontIndexCache{
		Version: fontIndexCacheVersion,
		Entries: nextEntries,
	})
	return fontInfos
}

// defaultFontIndexCachePath 返回字体索引缓存文件路径。
// defaultFontIndexCachePath returns the cache file path for font index.
func defaultFontIndexCachePath() string {
	base, err := os.UserCacheDir()
	if err != nil || base == "" {
		base = os.TempDir()
	}
	return filepath.Join(base, "go-canvas", "font_index_v2.json")
}

// loadFontIndexCacheFile 读取 JSON 索引缓存，失败时返回空缓存。
// loadFontIndexCacheFile reads cache from JSON and falls back to an empty cache on failure.
func loadFontIndexCacheFile(path string) fontIndexCache {
	empty := fontIndexCache{
		Version: fontIndexCacheVersion,
		Entries: make(map[string]fontIndexEntry),
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return empty
	}
	var cache fontIndexCache
	if err = json.Unmarshal(data, &cache); err != nil {
		return empty
	}
	if cache.Version != fontIndexCacheVersion || cache.Entries == nil {
		return empty
	}
	return cache
}

// saveFontIndexCacheFile 将字体索引写入 JSON 文件（失败忽略）。
// saveFontIndexCacheFile writes font index cache JSON to disk (best effort).
func saveFontIndexCacheFile(path string, cache fontIndexCache) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}
	data, err := json.Marshal(cache)
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o644)
}

func fontInfoFromEntry(entry fontIndexEntry) *FontInfo {
	return &FontInfo{
		Family:    entry.Family,
		SubFamily: entry.SubFamily,
		Bold:      entry.Bold,
		Italic:    entry.Italic,
		FontPath:  entry.Path,
	}
}

func (f *FontInfo) GetTrueTypeFont() (*truetype.Font, error) {
	if f.TruetypeFont != nil {
		return f.TruetypeFont, nil
	}

	data, err := os.ReadFile(f.FontPath)
	if err != nil {
		return nil, err
	}

	tf, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}
	f.TruetypeFont = tf
	return tf, nil
}

func ReadFontInfo(path string) (*FontInfo, error) {
	info := &FontInfo{
		FontPath: path,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	sf, err := parseSFNTFont(data)
	if err != nil {
		return nil, err
	}

	// 从 name table 中选择最可靠的 Family 名称，并过滤乱码候选。
	// Pick the most reliable family name from name table and filter mojibake-like candidates.
	info.Family = pickFontFamilyName(sf, path)

	// 判断粗体和斜体（优先 PreferredSubfamily）
	// Determine weight/italic by subfamily (prefer PreferredSubfamily).
	info.SubFamily = pickFontSubFamilyName(sf)
	if info.SubFamily != "" {
		info.Italic = isItalic(info.SubFamily)
	}

	// 从 fontStyles 匹配粗细数值
	info.Bold = matchWeight(info.SubFamily)
	return info, nil
}

// parseSFNTFont 解析 TTF/OTF/TTC 的第一个字体面。
// parseSFNTFont parses the first face from TTF/OTF/TTC data.
func parseSFNTFont(data []byte) (*sfnt.Font, error) {
	if coll, err := sfnt.ParseCollection(data); err == nil {
		return coll.Font(0)
	}
	return sfnt.Parse(data)
}

// pickFontFamilyName 从多个 NameID 中挑选可用的 family 名称。
// pickFontFamilyName chooses a usable family name from multiple NameID candidates.
func pickFontFamilyName(f *sfnt.Font, path string) string {
	buf := &sfnt.Buffer{}
	candidates := []string{
		sfntNameOrEmpty(f, buf, sfnt.NameIDTypographicFamily),
		sfntNameOrEmpty(f, buf, sfnt.NameIDFamily),
		sfntNameOrEmpty(f, buf, sfnt.NameIDFull),
		sfntNameOrEmpty(f, buf, sfnt.NameIDPostScript),
	}
	for _, name := range candidates {
		name = strings.TrimSpace(name)
		if isLikelyBadFontName(name) {
			continue
		}
		return name
	}
	// 最后兜底：文件名（无扩展名）
	// Last resort: file stem.
	stem := strings.TrimSpace(stripExtension(filepath.Base(path)))
	if stem == "" {
		return "fallback"
	}
	return stem
}

// pickFontSubFamilyName 选择可用的 subfamily 名称。
// pickFontSubFamilyName chooses a usable subfamily name.
func pickFontSubFamilyName(f *sfnt.Font) string {
	buf := &sfnt.Buffer{}
	candidates := []string{
		sfntNameOrEmpty(f, buf, sfnt.NameIDTypographicSubfamily),
		sfntNameOrEmpty(f, buf, sfnt.NameIDSubfamily),
	}
	for _, name := range candidates {
		name = strings.TrimSpace(name)
		if isLikelyBadFontName(name) {
			continue
		}
		return name
	}
	return ""
}

func sfntNameOrEmpty(f *sfnt.Font, b *sfnt.Buffer, id sfnt.NameID) string {
	v, err := f.Name(b, id)
	if err != nil {
		return ""
	}
	return v
}

// isLikelyBadFontName 判断名称是否疑似乱码/无效。
// isLikelyBadFontName returns true when the name looks like mojibake/invalid.
func isLikelyBadFontName(name string) bool {
	if name == "" {
		return true
	}
	question := 0
	for _, r := range name {
		if r == '?' || r == '\uFFFD' {
			question++
		}
	}
	// 全是问号，或问号占比过高，视为无效名称。
	// Treat all-? or very high ? ratio as invalid.
	if question == len([]rune(name)) {
		return true
	}
	return float64(question)/float64(len([]rune(name))) > 0.6
}

func isItalic(subFamily string) bool {
	lower := strings.ToLower(subFamily)
	for style := range italicStyles {
		if strings.Contains(lower, style) {
			return true
		}
	}
	return false
}

func matchWeight(subFamily string) FontWeight {
	if subFamily == "" {
		return FontWeightRegular // 默认 Regular
	}
	lower := strings.ToLower(subFamily)
	var maxWeight FontWeight = 0
	for style, weight := range fontStyles {
		if strings.Contains(lower, style) && weight > maxWeight {
			maxWeight = weight
		}
	}
	if maxWeight == 0 {
		maxWeight = FontWeightRegular
	}
	return maxWeight
}
