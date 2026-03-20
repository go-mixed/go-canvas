package render

import (
	"slideshow/ti"

	"github.com/go-mixed/go-taichi/taichi"
	"github.com/pkg/errors"
)

// Renderer 基于现有AOT模块的简化精灵渲染器
type Renderer struct {
	runtime *taichi.Runtime
	module  *ti.AotModule
}

// NewRenderer 创建简化精灵渲染器
// 使用 examples/aot_module 中已有的基础 kernels
func NewRenderer(runtime *taichi.Runtime) (*Renderer, error) {
	// 加载对应的 AOT 模块
	module, err := ti.LoadAotModule(runtime)
	if err != nil {
		return nil, errors.Wrapf(err, "Load TCM Module failed.")
	}

	return &Renderer{
		runtime: runtime,
		module:  module,
	}, nil
}

// Release 释放渲染器资源
func (sr *Renderer) Release() {
	if sr.module != nil {
		sr.module.Release()
		sr.module = nil
	}
}
