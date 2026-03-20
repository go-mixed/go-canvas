package ti

//go:generate uv run ./ti/kernel/generate_tcm.py

import (
	_ "embed"
	"image/color"

	"github.com/go-mixed/go-taichi/taichi"
	"github.com/pkg/errors"
)

//go:embed "tcm/cpu.tcm"
var cpu []byte

//go:embed "tcm/cuda.tcm"
var cuda []byte

//go:embed "tcm/vulkan.tcm"
var vulkan []byte

type AotModule struct {
	module *taichi.AotModule
	cache  map[string]*taichi.Kernel
}

// LoadAotModule 读取 AOT 模块
func LoadAotModule(runtime *taichi.Runtime) (*AotModule, error) {
	var data []byte
	switch runtime.Arch() {
	case taichi.ArchX64:
		data = cpu
	case taichi.ArchCuda:
		data = cuda
	case taichi.ArchVulkan:
		fallthrough
	default:
		// 默认使用 vulkan
		data = vulkan
	}

	// 加载对应的 AOT 模块
	m, err := taichi.LoadAotModule(runtime, data)
	if err != nil {
		return nil, errors.Wrapf(err, "Load TCM Module failed.")
	}

	// 获取对应的 kernels
	var modules = []string{
		"fill_texture",
		"cv_image_to_ti",
		"render_layer_no_mask",
		"render_layer_with_mask",
	}

	cache := make(map[string]*taichi.Kernel)
	for _, name := range modules {
		cache[name], err = m.GetKernel(name)
		if err != nil {
			return nil, errors.Wrapf(err, "[taichi]aot kernel %s not found", name)
		}
	}

	return &AotModule{
		module: m,
		cache:  cache,
	}, nil
}

func (m *AotModule) Release() {
	m.module.Release()
}

func (m *AotModule) getCache(name string) *taichi.Kernel {
	return m.cache[name]
}

// FillTexture 填充纹理
func (m *AotModule) FillTexture(texture *TiImage, c color.Color) {
	kernel := m.getCache("fill_texture")
	r, g, b, a := Color2TiColor(c)

	kernel.Launch().ArgNdArray(texture).ArgFloat32(r).ArgFloat32(g).ArgFloat32(b).ArgFloat32(a).Run()
}

// CvToTiImage 将 CvImage (h, w, [b, g, r]) 转换为 TiImage (w, h, [r, g, b, a])
func (m *AotModule) CvToTiImage(inputImage *CvImage, outputImage *TiImage) {
	kernel := m.getCache("cv_image_to_ti")
	kernel.Launch().ArgNdArray(inputImage).ArgNdArray(outputImage).Run()
}

// RenderLayerOptions 渲染层选项
type RenderLayerOptions struct {
	X, Y       float32 // 相对屏幕的偏移
	Cx, Cy     float32 // 纹理中心坐标
	Scale      float32 // 缩放倍数
	Rotation   float32 // 旋转弧度
	Alpha      float32 // 透明度 0.0-1.0
	Width      float32 // 纹理宽度
	Height     float32 // 纹理高度
	MinX, MaxX int32   // 包围盒 x 范围
	MinY, MaxY int32   // 包围盒 y 范围
}

// RenderLayerNoMask 渲染层（无遮罩）
func (m *AotModule) RenderLayerNoMask(texture *TiImage, screen *TiImage, opts RenderLayerOptions) {
	kernel := m.getCache("render_layer_no_mask")
	kernel.Launch().
		ArgNdArray(texture).
		ArgFloat32(opts.X).ArgFloat32(opts.Y).
		ArgFloat32(opts.Cx).ArgFloat32(opts.Cy).
		ArgFloat32(opts.Scale).ArgFloat32(opts.Rotation).ArgFloat32(opts.Alpha).
		ArgFloat32(opts.Width).ArgFloat32(opts.Height).
		ArgInt32(opts.MinX).ArgInt32(opts.MaxX).ArgInt32(opts.MinY).ArgInt32(opts.MaxY).
		ArgNdArray(screen).
		Run()
}

// RenderLayerWithMask 渲染层（带遮罩）
func (m *AotModule) RenderLayerWithMask(texture, mask, screen *TiImage, opts RenderLayerOptions) {
	kernel := m.getCache("render_layer_with_mask")
	kernel.Launch().
		ArgNdArray(texture).
		ArgFloat32(opts.X).ArgFloat32(opts.Y).
		ArgFloat32(opts.Cx).ArgFloat32(opts.Cy).
		ArgFloat32(opts.Scale).ArgFloat32(opts.Rotation).ArgFloat32(opts.Alpha).
		ArgFloat32(opts.Width).ArgFloat32(opts.Height).
		ArgInt32(opts.MinX).ArgInt32(opts.MaxX).ArgInt32(opts.MinY).ArgInt32(opts.MaxY).
		ArgNdArray(mask).
		ArgNdArray(screen).
		Run()
}
