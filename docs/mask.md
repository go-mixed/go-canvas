# Mask

Masks define visibility regions for sprites. go-canvas provides shape-based masks with feather support.

## Symbol Overview

```
IMask (interface)
├── FillWithTexture(texture *ti.TiImage)
├── ApplyFeather(featherRadius uint32, featherMode ti.FeatherMode)
├── Release()
└── Texture() *taichi.NdArray

Mask
├── Fields: texture, distField, renderer
├── NewMask(renderer *Renderer, width, height uint32) (*Mask, error)
└── Methods:
    ├── FillWithTexture(texture *ti.TiImage)
    ├── ApplyFeather(featherRadius uint32, featherMode ti.FeatherMode)
    ├── Texture() *taichi.NdArray
    └── Release()

ShapeMask (extends Mask)
├── Embeds: *Mask, *ShapeSprite
├── NewShapeMask(renderer *Renderer, width, height, cx, cy uint32) (*ShapeMask, error)
└── Methods:
    ├── SetFeather(radius uint32, featherMode ti.FeatherMode)
    ├── DrawShape(shapeType ti.ShapeType, tVal float32, fns ...func(option *ti.ShapeOptions))
    ├── Texture() *ti.TiMask
    └── Release()
```

## IMask Interface

Base interface for all mask types.

```go
type IMask interface {
    FillWithTexture(texture *ti.TiImage)
    ApplyFeather(featherRadius uint32, featherMode ti.FeatherMode)
    Release()
    Texture() *taichi.NdArray
}
```

## Mask

The base Mask class provides texture-based masking capabilities.

### NewMask

```go
func NewMask(renderer *Renderer, width, height uint32) (*Mask, error)
```

Creates a new empty mask.

**Parameters:**
- `renderer` - The Renderer instance
- `width` - Mask width in pixels
- `height` - Mask height in pixels

### FillWithTexture

```go
func (m *Mask) FillWithTexture(texture *ti.TiImage)
```

Fills the mask with a texture. White areas = fully visible, black areas = fully hidden.

### ApplyFeather

```go
func (m *Mask) ApplyFeather(featherRadius uint32, featherMode ti.FeatherMode)
```

Applies feather (blur) to mask edges for smooth transitions.

**Parameters:**
- `featherRadius` - Feather radius in pixels
- `featherMode` - Feather mode (e.g., `ti.FeatherModeLinear`)

### Texture

```go
func (m *Mask) Texture() *taichi.NdArray
```

Returns the mask texture as NdArray.

### Release

```go
func (m *Mask) Release()
```

Releases mask resources.

## ShapeMask

A convenience mask type that draws geometric shapes. Embeds both `*Mask` and `*ShapeSprite`.

### NewShapeMask

```go
func NewShapeMask(renderer *Renderer, width, height, cx, cy uint32) (*ShapeMask, error)
```

Creates a shape mask with specified dimensions and center point.

**Parameters:**
- `renderer` - The Renderer instance
- `width` - Mask width in pixels
- `height` - Mask height in pixels
- `cx` - Shape center X coordinate in pixels
- `cy` - Shape center Y coordinate in pixels

### SetFeather

```go
func (sm *ShapeMask) SetFeather(radius uint32, featherMode ti.FeatherMode)
```

Sets the feather radius and mode for soft edges.

**Parameters:**
- `radius` - Feather radius in pixels
- `featherMode` - Feather mode (e.g., `ti.FeatherModeLinear`)

### DrawShape

```go
func (sm *ShapeMask) DrawShape(shapeType ti.ShapeType, tVal float32, fns ...func(option *ti.ShapeOptions))
```

Draws a shape with the specified parameters.

**Parameters:**
- `shapeType` - One of the ti.ShapeType values (see [Shape](shape.md))
- `tVal` - Shape parameter (typically 0.0-2.0, where 1.0 = full screen)
- `fns` - Optional shape options (e.g., `WithRoundCorner`)

### Texture

```go
func (sm *ShapeMask) Texture() *ti.TiMask
```

Returns the mask texture as TiMask.

### Release

```go
func (sm *ShapeMask) Release()
```

Releases mask and shape sprite resources.

## Usage with Sprite

```go
// Create mask
mask, err := render.NewShapeMask(renderer, 720, 1280, 360, 640)
if err != nil {
    panic(err)
}
defer mask.Release()

// Draw circle (tVal = 1.0 means full screen)
mask.DrawShape(ti.ShapeTypeCircle, 1.0)

// Apply soft edges
mask.SetFeather(10, ti.FeatherModeLinear)

// Apply to sprite
img.SetMask(mask)
```

## Feather Modes

Available feather modes from `ti` package:
- `ti.FeatherModeLinear` - Linear feather
- Other modes see [ti package](https://github.com/go-mixed/go-taichi)

## Related

- [Shape](shape.md) - Shape type definitions
- [Sprite](sprite.md) - Sprites that use masks
