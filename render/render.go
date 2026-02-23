package render

import (
	"fmt"

	"github.com/go-mixed/go-taichi/taichi"
)

// SimpleSpriteRenderer 基于现有AOT模块的简化精灵渲染器
type SimpleSpriteRenderer struct {
	runtime *taichi.Runtime
	module  *taichi.AotModule
}

// NewSimpleSpriteRenderer 创建简化精灵渲染器
// 使用 examples/aot_module 中已有的基础 kernels
func NewSimpleSpriteRenderer(runtime *taichi.Runtime) (*SimpleSpriteRenderer, error) {
	// 根据运行时后端选择对应的 AOT 模块
	var modulePath string
	switch runtime.Arch() {
	case taichi.ArchCuda:
		modulePath = "./cuda.tcm"
	case taichi.ArchVulkan:
		modulePath = "./vulkan.tcm"
	default:
		// 默认使用 vulkan
		modulePath = "./vulkan.tcm"
	}

	// 加载对应的 AOT 模块
	module, err := taichi.LoadAotModule(runtime, modulePath)
	if err != nil {
		return nil, fmt.Errorf("加载 AOT 模块失败: %w\n模块路径: %s", err, modulePath)
	}

	return &SimpleSpriteRenderer{
		runtime: runtime,
		module:  module,
	}, nil
}

// Release 释放渲染器资源
func (sr *SimpleSpriteRenderer) Release() {
	if sr.module != nil {
		sr.module.Release()
		sr.module = nil
	}
}
