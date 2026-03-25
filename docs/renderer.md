# Renderer

The `Renderer` is the entry point for creating rendering contexts in go-canvas. It manages the Taichi runtime and shader module.

## Symbol Overview

```
Renderer
├── Fields: runtime, module
├── NewRenderer(runtime) (*Renderer, error)
└── Methods:
    ├── Release()
    ├── Runtime() taichi.Runtime
    └── Module() *ti.AotModule
```

## Constructor

### NewRenderer

```go
func NewRenderer(runtime taichi.Runtime) (*Renderer, error)
```

Creates a new Renderer instance with the specified Taichi runtime.

**Parameters:**
- `runtime` - Taichi runtime instance (e.g., created with `taichi.NewRuntime(taichi.ArchCuda, ...)`)

**Returns:**
- `*Renderer` - Renderer instance
- `error` - Error if creation fails

**Example:**

```go
runtime, err := taichi.NewRuntime(taichi.ArchCuda, taichi.WithCacheTcm(true))
if err != nil {
    panic(err)
}
defer runtime.Release()

renderer, err := render.NewRenderer(runtime)
if err != nil {
    panic(err)
}
defer renderer.Release()
```

## Methods

### Release

```go
func (r *Renderer) Release()
```

Releases all resources held by the renderer. Must be called when done.

### Runtime

```go
func (r *Renderer) Runtime() *taichi.Runtime
```

Returns the underlying Taichi runtime instance.

### Module

```go
func (r *Renderer) Module() *ti.AotModule
```

Returns the loaded shader module.

## Related

- [Stage](stage.md) - The main canvas that uses Renderer
- [Sprite](sprite.md) - Visual elements rendered by the Stage
