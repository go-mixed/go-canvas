package font

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	xfont "golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

const fontIndexCacheVersion = 7

type fontIndexCache struct {
	Version int                         `json:"version"`
	Entries map[string][]fontIndexEntry `json:"entries"`
}

type fontIndexEntry struct {
	Family    string              `json:"family"`
	SubFamily string              `json:"sub_family"`
	Bold      xfont.Weight        `json:"bold"`
	Italic    bool                `json:"italic"`
	Path      string              `json:"path"`
	FaceIndex int                 `json:"face_index,omitempty"`
	Size      int64               `json:"size"`
	ModTime   int64               `json:"mod_time_unix_nano"`
	Coverage  []cacheUnicodeRange `json:"coverage,omitempty"`
}

type cacheUnicodeRange struct {
	Start uint32 `json:"start"`
	End   uint32 `json:"end"`
}

// loadFonts 从目录枚举字体并使用 JSON 索引缓存加速元信息加载。
// loadFonts enumerates fonts and uses a JSON index cache to speed up metadata loading.
func (fs *FontLibrary) loadFonts(userFontPaths ...string) map[string][]*FontInfo {
	cachePath := fs.defaultFontIndexCachePath()
	paths := ListFont(append(userFontPaths, GetSystemFontDirectories()...))
	cache := fs.loadFontIndexCacheFile(cachePath)
	fontInfos := make(map[string][]*FontInfo)
	nextEntries := make(map[string][]fontIndexEntry, len(paths))
	changed := false

	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			changed = true
			continue
		}

		fs.logf("[font-library]load font: %s", path)

		if cachedEntries, ok := cache.Entries[path]; ok &&
			len(cachedEntries) > 0 &&
			cachedEntries[0].Size == stat.Size() &&
			cachedEntries[0].ModTime == stat.ModTime().UnixNano() &&
			cachedEntries[0].Family != "" {
			for _, entry := range cachedEntries {
				info := fs.fontInfoFromEntry(entry)
				key := info.Family
				if key == "" {
					key = info.Family
				}
				fontInfos[key] = append(fontInfos[key], info)
			}
			nextEntries[path] = cachedEntries
			continue
		}

		infos, err := fs.ReadTTCFontInfos(path)
		if err != nil {
			changed = true
			continue
		}
		if len(infos) == 0 {
			changed = true
			continue
		}
		entries := make([]fontIndexEntry, 0, len(infos))
		for _, info := range infos {
			key := info.Family
			if key == "" {
				key = info.Family
			}
			fontInfos[key] = append(fontInfos[key], info)
			entries = append(entries, fontIndexEntry{
				Family:    info.Family,
				SubFamily: info.SubFamily,
				Bold:      info.Bold,
				Italic:    info.Italic,
				Path:      info.FontPath,
				FaceIndex: info.FaceIndex,
				Size:      stat.Size(),
				ModTime:   stat.ModTime().UnixNano(),
				Coverage:  toCacheUnicodeRanges(info.coverageRanges),
			})
		}
		nextEntries[path] = entries
		changed = true
	}

	if !changed && len(cache.Entries) == len(nextEntries) {
		return fontInfos
	}
	fs.saveFontIndexCacheFile(cachePath, fontIndexCache{
		Version: fontIndexCacheVersion,
		Entries: nextEntries,
	})
	return fontInfos
}

// defaultFontIndexCachePath 返回字体索引缓存文件路径。
// defaultFontIndexCachePath returns the cache file path for font index.
func (fs *FontLibrary) defaultFontIndexCachePath() string {
	base, err := os.UserCacheDir()
	if err != nil || base == "" {
		base = os.TempDir()
	}
	return filepath.Join(base, "go-canvas", "font_index.json")
}

// loadFontIndexCacheFile 读取 JSON 索引缓存，失败时返回空缓存。
// loadFontIndexCacheFile reads cache from JSON and falls back to an empty cache on failure.
func (fs *FontLibrary) loadFontIndexCacheFile(path string) fontIndexCache {
	empty := fontIndexCache{
		Version: fontIndexCacheVersion,
		Entries: make(map[string][]fontIndexEntry),
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
func (fs *FontLibrary) saveFontIndexCacheFile(path string, cache fontIndexCache) {
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

func (fs *FontLibrary) fontInfoFromEntry(entry fontIndexEntry) *FontInfo {
	return &FontInfo{
		Family:         entry.Family,
		SubFamily:      entry.SubFamily,
		Bold:           entry.Bold,
		Italic:         entry.Italic,
		FontPath:       entry.Path,
		FaceIndex:      entry.FaceIndex,
		coverageRanges: fromCacheUnicodeRanges(entry.Coverage),
	}
}

// parseCoverageRangesFromCMAP 仅用于“扫描阶段”，从字体 cmap 解析字符覆盖区间。
// parseCoverageRangesFromCMAP is scan-time only; it extracts coverage ranges from cmap.
func (fs *FontLibrary) parseCoverageRangesFromCMAP(data []byte) ([]unicodeRange, bool) {
	return fs.parseCoverageRangesFromCMAPForFace(data, 0)
}

// parseCoverageRangesFromCMAPForFace 仅用于扫描阶段，按 face 索引解析覆盖区间。
// parseCoverageRangesFromCMAPForFace is scan-time only; it extracts ranges for the given face index.
func (fs *FontLibrary) parseCoverageRangesFromCMAPForFace(data []byte, faceIndex int) ([]unicodeRange, bool) {
	if faceIndex < 0 {
		return nil, false
	}
	// TTC (TrueType Collection) does not start with a direct sfnt table directory.
	// We need to resolve to a concrete face offset first.
	if len(data) >= 4 && string(data[:4]) == "ttcf" {
		return fs.parseCoverageRangesFromTTC(data, faceIndex)
	}
	if faceIndex > 0 {
		return nil, false
	}
	return fs.parseCoverageRangesFromSFNT(data)
}

// parseCoverageRangesFromSFNT 从单个 SFNT 字体读取 cmap 覆盖区间。
// parseCoverageRangesFromSFNT extracts cmap coverage ranges from one SFNT font.
func (fs *FontLibrary) parseCoverageRangesFromSFNT(data []byte) ([]unicodeRange, bool) {
	return fs.parseCoverageRangesFromSFNTAt(data, 0)
}

// parseCoverageRangesFromSFNTAt 从指定 sfnt 起点读取 cmap 覆盖区间。
// parseCoverageRangesFromSFNTAt extracts cmap coverage ranges from sfnt at offset.
func (fs *FontLibrary) parseCoverageRangesFromSFNTAt(data []byte, sfntOffset int) ([]unicodeRange, bool) {
	cmapOff, cmapLen, ok := fs.findSFNTTableAt(data, sfntOffset, "cmap")
	if !ok || cmapLen < 4 {
		return nil, false
	}
	if int(cmapOff+cmapLen) > len(data) {
		return nil, false
	}
	cmap := data[cmapOff : cmapOff+cmapLen]
	if len(cmap) < 4 {
		return nil, false
	}
	numTables := int(binary.BigEndian.Uint16(cmap[2:4]))
	if len(cmap) < 4+numTables*8 {
		return nil, false
	}

	// Prefer format 12 (full Unicode), then format 4.
	bestFmt12 := -1
	bestFmt4 := -1
	for i := 0; i < numTables; i++ {
		rec := cmap[4+i*8 : 12+i*8]
		subOff := int(binary.BigEndian.Uint32(rec[4:8]))
		if subOff+2 > len(cmap) || subOff < 0 {
			continue
		}
		format := binary.BigEndian.Uint16(cmap[subOff : subOff+2])
		if format == 12 && bestFmt12 < 0 {
			bestFmt12 = subOff
		} else if format == 4 && bestFmt4 < 0 {
			bestFmt4 = subOff
		}
	}
	if bestFmt12 >= 0 {
		if out, ok := fs.parseFormat12Ranges(cmap[bestFmt12:]); ok {
			return out, true
		}
	}
	if bestFmt4 >= 0 {
		if out, ok := fs.parseFormat4Ranges(cmap[bestFmt4:]); ok {
			return out, true
		}
	}
	return nil, false
}

// parseCoverageRangesFromTTC 从 TTC 中选取第一个可解析子字体并读取覆盖区间。
// parseCoverageRangesFromTTC picks the first parsable face from TTC and extracts coverage ranges.
func (fs *FontLibrary) parseCoverageRangesFromTTC(data []byte, faceIndex int) ([]unicodeRange, bool) {
	// TTC header: tag(4) + version(4) + numFonts(4) + offsets[numFonts](4 each)
	if len(data) < 12 || string(data[:4]) != "ttcf" {
		return nil, false
	}
	numFonts := int(binary.BigEndian.Uint32(data[8:12]))
	if numFonts <= 0 {
		return nil, false
	}
	offsetsBase := 12
	if len(data) < offsetsBase+numFonts*4 {
		return nil, false
	}
	if faceIndex < numFonts {
		off := int(binary.BigEndian.Uint32(data[offsetsBase+faceIndex*4 : offsetsBase+faceIndex*4+4]))
		if off > 0 && off < len(data) {
			if out, ok := fs.parseCoverageRangesFromSFNTAt(data, off); ok {
				return out, true
			}
		}
	}
	// Fallback: if target face is malformed, try subsequent faces.
	for i := 0; i < numFonts; i++ {
		if i == faceIndex {
			continue
		}
		off := int(binary.BigEndian.Uint32(data[offsetsBase+i*4 : offsetsBase+i*4+4]))
		if off <= 0 || off >= len(data) {
			continue
		}
		if out, ok := fs.parseCoverageRangesFromSFNTAt(data, off); ok {
			return out, true
		}
	}
	return nil, false
}

func (fs *FontLibrary) findSFNTTable(data []byte, tag string) (offset, length int, ok bool) {
	return fs.findSFNTTableAt(data, 0, tag)
}

// findSFNTTableAt 在指定 sfnt 起点查找表，返回相对整个文件的绝对偏移。
// findSFNTTableAt finds a table from sfntOffset and returns absolute file offset.
func (fs *FontLibrary) findSFNTTableAt(data []byte, sfntOffset int, tag string) (offset, length int, ok bool) {
	if len(tag) != 4 || sfntOffset < 0 || len(data) < sfntOffset+12 {
		return 0, 0, false
	}
	sfntData := data[sfntOffset:]
	numTables := int(binary.BigEndian.Uint16(sfntData[4:6]))
	if len(sfntData) < 12+numTables*16 {
		return 0, 0, false
	}
	for i := 0; i < numTables; i++ {
		rec := sfntData[12+i*16 : 28+i*16]
		if string(rec[0:4]) != tag {
			continue
		}
		off := int(binary.BigEndian.Uint32(rec[8:12]))
		l := int(binary.BigEndian.Uint32(rec[12:16]))
		if off < 0 || l < 0 {
			return 0, 0, false
		}
		// TTC 常用“相对文件起点”的偏移；也兼容“相对 sfnt 起点”的实现。
		// TTC often stores file-absolute offsets; keep sfnt-relative fallback for compatibility.
		if off+l <= len(data) {
			return off, l, true
		}
		absOff := sfntOffset + off
		if absOff+l <= len(data) {
			return absOff, l, true
		}
		return 0, 0, false
	}
	return 0, 0, false
}

func (fs *FontLibrary) parseFormat12Ranges(sub []byte) ([]unicodeRange, bool) {
	if len(sub) < 16 {
		return nil, false
	}
	nGroups := int(binary.BigEndian.Uint32(sub[12:16]))
	if len(sub) < 16+nGroups*12 {
		return nil, false
	}
	out := make([]unicodeRange, 0, nGroups)
	for i := 0; i < nGroups; i++ {
		g := sub[16+i*12 : 28+i*12]
		start := rune(binary.BigEndian.Uint32(g[0:4]))
		end := rune(binary.BigEndian.Uint32(g[4:8]))
		startGlyph := binary.BigEndian.Uint32(g[8:12])
		if end < start {
			continue
		}
		// startGlyph==0 means first codepoint maps to .notdef; shift start if possible.
		if startGlyph == 0 {
			if start == end {
				continue
			}
			start++
		}
		out = append(out, unicodeRange{start: start, end: end})
	}
	return out, true
}

func (fs *FontLibrary) parseFormat4Ranges(sub []byte) ([]unicodeRange, bool) {
	if len(sub) < 16 {
		return nil, false
	}
	segCount := int(binary.BigEndian.Uint16(sub[6:8]) / 2)
	if segCount <= 0 {
		return nil, false
	}
	endCodesOff := 14
	startCodesOff := endCodesOff + segCount*2 + 2
	idDeltaOff := startCodesOff + segCount*2
	idRangeOff := idDeltaOff + segCount*2
	if idRangeOff+segCount*2 > len(sub) {
		return nil, false
	}

	out := make([]unicodeRange, 0, segCount)
	for i := 0; i < segCount; i++ {
		end := int(binary.BigEndian.Uint16(sub[endCodesOff+i*2 : endCodesOff+i*2+2]))
		start := int(binary.BigEndian.Uint16(sub[startCodesOff+i*2 : startCodesOff+i*2+2]))
		if start > end || end == 0xFFFF {
			continue
		}
		// Build precise ranges inside this segment (avoid false positives).
		in := false
		rStart := 0
		prev := 0
		for cp := start; cp <= end; cp++ {
			gid, ok := fs.glyphIndexFromFormat4(sub, segCount, i, cp)
			supported := ok && gid != 0
			if supported {
				if !in {
					rStart = cp
					prev = cp
					in = true
				} else if cp == prev+1 {
					prev = cp
				} else {
					out = append(out, unicodeRange{start: rune(rStart), end: rune(prev)})
					rStart = cp
					prev = cp
				}
			} else if in {
				out = append(out, unicodeRange{start: rune(rStart), end: rune(prev)})
				in = false
			}
		}
		if in {
			out = append(out, unicodeRange{start: rune(rStart), end: rune(prev)})
		}
	}
	return out, true
}

func (fs *FontLibrary) glyphIndexFromFormat4(sub []byte, segCount, segIdx, cp int) (uint16, bool) {
	endCodesOff := 14
	startCodesOff := endCodesOff + segCount*2 + 2
	idDeltaOff := startCodesOff + segCount*2
	idRangeOff := idDeltaOff + segCount*2

	start := int(binary.BigEndian.Uint16(sub[startCodesOff+segIdx*2 : startCodesOff+segIdx*2+2]))
	end := int(binary.BigEndian.Uint16(sub[endCodesOff+segIdx*2 : endCodesOff+segIdx*2+2]))
	if cp < start || cp > end {
		return 0, false
	}
	idDelta := int16(binary.BigEndian.Uint16(sub[idDeltaOff+segIdx*2 : idDeltaOff+segIdx*2+2]))
	idRangeOffset := int(binary.BigEndian.Uint16(sub[idRangeOff+segIdx*2 : idRangeOff+segIdx*2+2]))

	if idRangeOffset == 0 {
		return uint16(cp + int(idDelta)), true
	}
	// Address of glyph index entry:
	// ptr = &idRangeOffset[i] + idRangeOffset[i] + 2*(cp-startCode[i])
	offsetWordPos := idRangeOff + segIdx*2
	glyphPos := offsetWordPos + idRangeOffset + 2*(cp-start)
	if glyphPos < 0 || glyphPos+2 > len(sub) {
		return 0, false
	}
	gid := binary.BigEndian.Uint16(sub[glyphPos : glyphPos+2])
	if gid == 0 {
		return 0, true
	}
	return uint16(int(gid) + int(idDelta)), true
}

func (fs *FontLibrary) mergeRanges(in []unicodeRange) []unicodeRange {
	if len(in) == 0 {
		return nil
	}
	sort.Slice(in, func(i, j int) bool {
		if in[i].start == in[j].start {
			return in[i].end < in[j].end
		}
		return in[i].start < in[j].start
	})
	out := make([]unicodeRange, 0, len(in))
	cur := in[0]
	for i := 1; i < len(in); i++ {
		r := in[i]
		if r.start <= cur.end+1 {
			if r.end > cur.end {
				cur.end = r.end
			}
			continue
		}
		out = append(out, cur)
		cur = r
	}
	out = append(out, cur)
	return out
}

// GetOpenTypeFont 按 FontPath+FaceIndex 读取并缓存 OpenType 字体面。
// GetOpenTypeFont loads and caches the OpenType face by FontPath+FaceIndex.
func (f *FontInfo) GetOpenTypeFont() (*opentype.Font, error) {
	if f.OpenTypeFont != nil {
		return f.OpenTypeFont, nil
	}

	data, err := os.ReadFile(f.FontPath)
	if err != nil {
		return nil, err
	}

	var of *opentype.Font
	if coll, cerr := opentype.ParseCollection(data); cerr == nil {
		of, err = coll.Font(f.FaceIndex)
		if err != nil {
			return nil, err
		}
	} else {
		of, err = opentype.Parse(data)
		if err != nil {
			return nil, err
		}
	}
	f.OpenTypeFont = of
	return of, nil
}

// ReadTTCFontInfos 读取文件内所有可用 face 的字体信息（TTC 可返回多个）。
// ReadTTCFontInfos reads all usable faces from a file (TTC may return multiple).
func (fs *FontLibrary) ReadTTCFontInfos(path string) ([]*FontInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	sfs, err := parseSFNTFonts(data)
	if err != nil {
		return nil, err
	}
	infos := make([]*FontInfo, 0, len(sfs))
	for _, sf := range sfs {
		info := &FontInfo{FontPath: path}
		info.FaceIndex = sf.faceIndex
		// 从 name table 中选择最可靠的 Family 名称，并过滤乱码候选。
		// Pick the most reliable family name from name table and filter mojibake-like candidates.
		info.Family = normalizeFamilyName(pickFontFamilyName(sf.font, path))

		// 判断粗体和斜体（优先 PreferredSubfamily）
		// Determine weight/italic by subfamily (prefer PreferredSubfamily).
		info.SubFamily = pickFontSubFamilyName(sf.font)
		if info.SubFamily != "" {
			info.Italic = isItalic(info.SubFamily)
		}
		// 从 fontStyles 匹配粗细数值
		info.Bold = matchWeight(info.SubFamily)

		// 扫描阶段读取 coverage 范围并写入缓存；运行阶段只查范围，不再读取字体。
		// Read coverage ranges during scanning and persist to cache.
		// Runtime only queries these ranges and never parses coverage again.
		if ranges, ok := fs.parseCoverageRangesFromCMAPForFace(data, sf.faceIndex); ok {
			info.coverageRanges = fs.mergeRanges(ranges)
		} else {
			info.coverageRanges = nil
		}
		infos = append(infos, info)
		fs.logf("[font-library]read %s from %s", info.Family, path)
	}
	if len(infos) == 0 {
		return nil, os.ErrInvalid
	}
	return infos, nil
}

func toCacheUnicodeRanges(in []unicodeRange) []cacheUnicodeRange {
	if len(in) == 0 {
		return nil
	}
	out := make([]cacheUnicodeRange, 0, len(in))
	for _, r := range in {
		if r.end < r.start {
			continue
		}
		out = append(out, cacheUnicodeRange{
			Start: uint32(r.start),
			End:   uint32(r.end),
		})
	}
	return out
}

func fromCacheUnicodeRanges(in []cacheUnicodeRange) []unicodeRange {
	if len(in) == 0 {
		return nil
	}
	out := make([]unicodeRange, 0, len(in))
	for _, r := range in {
		if r.End < r.Start {
			continue
		}
		out = append(out, unicodeRange{
			start: rune(r.Start),
			end:   rune(r.End),
		})
	}
	return out
}

// parseSFNTFonts 解析字体文件中的所有可用面；非 TTC 时返回单个面。
// parseSFNTFonts parses all available faces from a font file; non-TTC returns one face.
type sfntFaceRef struct {
	font      *sfnt.Font
	faceIndex int
}

func parseSFNTFonts(data []byte) ([]sfntFaceRef, error) {
	if coll, err := sfnt.ParseCollection(data); err == nil {
		n := coll.NumFonts()
		out := make([]sfntFaceRef, 0, n)
		for i := 0; i < n; i++ {
			f, ferr := coll.Font(i)
			if ferr != nil {
				continue
			}
			out = append(out, sfntFaceRef{font: f, faceIndex: i})
		}
		if len(out) > 0 {
			return out, nil
		}
	}
	f, err := sfnt.Parse(data)
	if err != nil {
		return nil, err
	}
	return []sfntFaceRef{{font: f, faceIndex: 0}}, nil
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

func matchWeight(subFamily string) xfont.Weight {
	if subFamily == "" {
		return xfont.WeightNormal // 默认 Regular
	}
	lower := strings.ToLower(subFamily)
	var maxWeight xfont.Weight = xfont.WeightNormal
	for style, weight := range xfontStyles {
		if strings.Contains(lower, style) && weight > maxWeight {
			maxWeight = weight
		}
	}
	if maxWeight == xfont.WeightNormal {
		maxWeight = xfont.WeightNormal
	}
	return maxWeight
}
