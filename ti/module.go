package ti

//go:generate uv run ./ti/kernel/generate_tcm.py

import (
	_ "embed"

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
		"fill_color",
		"resize_nearest", "resize_bilinear", "resize_bicubic", "resize_lanczos",
		"blur_box", "blur_gaussian", "blur_mosaic",

		"render_layer_no_mask",
		"render_layer_with_mask",
		"compute_normalized_coords",
		"compute_circle",
		"compute_diamond",
		"compute_rect",
		"compute_directional",
		"compute_triangle",
		"compute_star",
		"compute_heart",
		"compute_cross",
		"image_to_mask",
		"compute_distance_field",
		"feather_linear",
		"feather_conic",
		"feather_smoothstep",
		"feather_sigmoid",
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
