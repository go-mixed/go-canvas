package render

import (
	"github.com/go-mixed/go-canvas/ti"
	"github.com/go-mixed/go-taichi/taichi"
	"github.com/pkg/errors"
)

// Renderer 基于现有AOT模块的简化精灵渲染器
type Renderer struct {
	runtime *taichi.Runtime
	module  *ti.AotModule
}

// NewRenderer 创建渲染器
func NewRenderer(arch taichi.Arch) (*Renderer, error) {
	runtime, err := taichi.NewRuntime(arch, taichi.WithCacheTcm(true))
	if err != nil {
		return nil, errors.Wrapf(err, "Create Taichi Runtime failed.")
	}

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

	if sr.runtime != nil {
		sr.runtime.Release()
		sr.runtime = nil
	}
}

func (sr *Renderer) Runtime() *taichi.Runtime {
	return sr.runtime
}

func (sr *Renderer) Module() *ti.AotModule {
	return sr.module
}
