# Sprite

`Sprite` is the base visual element in go-canvas. All renderable objects (images, text, spectrum) implement the `ISprite` interface.

## Symbol Overview

```
ISpriteOperator (interface)
├── X(), Y() float32
├── Width(), Height() float32
├── Scale(), Rotation(), Alpha() float32
├── CenterX(), CenterY() float32
├── Texture() *ti.TiImage
├── SetX(x float32) ISprite
├── SetY(y float32) ISprite
├── MoveTo(x, y float32) ISprite
├── SetScale(scale float32) ISprite
├── SetScaleTo(width, height float32) ISprite
├── SetRotation(rotation float32) ISprite
└── SetAlpha(alpha float32) ISprite

ISprite (interface, extends ISpriteOperator)
├── ISpriteOperator
├── SetMask(mask IMask) ISprite
├── Mask() IMask
├── SetTexture(texture *ti.TiImage) ISprite
├── FillTexture(color color.Color)
├── BoundingBox(parentWidth, parentHeight float32) ti.Rectangle[float32]
├── ResizeTo(width, height uint32) ISprite
├── Renderer() *Renderer
└── Release()

Sprite (struct, implements ISprite)
├── NewSprite(renderer *Renderer, texture *ti.TiImage) ISprite
├── NewBlockSprite(renderer *Renderer, width, height uint32) (ISprite, error)
└── Methods: (see ISprite interface)
```

## Constructors

### NewSprite

```go
func NewSprite(renderer *Renderer, texture *ti.TiImage) ISprite
```

Creates a new sprite with the specified texture.

### NewBlockSprite

```go
func NewBlockSprite(renderer *Renderer, width, height uint32) (ISprite, error)
```

Creates a sprite with a solid color block (transparent by default).

## Transform Methods

### Position

```go
func (s *Sprite) SetX(x float32) ISprite
func (s *Sprite) X() float32
func (s *Sprite) SetY(y float32) ISprite
func (s *Sprite) Y() float32
func (s *Sprite) MoveTo(x, y float32) ISprite
```

### Scale

```go
func (s *Sprite) SetScale(scale float32) ISprite
func (s *Sprite) Scale() float32
func (s *Sprite) SetScaleTo(width, height float32) ISprite
```

### Rotation

```go
func (s *Sprite) SetRotation(rotation float32) ISprite
func (s *Sprite) Rotation() float32
```

Rotation is specified in **radians** (not degrees).

### Alpha

```go
func (s *Sprite) SetAlpha(alpha float32) ISprite
func (s *Sprite) Alpha() float32
```

Alpha is in range [0, 1], where 0 is fully transparent and 1 is fully opaque.

## Geometry Methods

### Dimensions

```go
func (s *Sprite) Width() float32
func (s *Sprite) Height() float32
func (s *Sprite) ResizeTo(width, height uint32) ISprite
```

### Center

```go
func (s *Sprite) CenterX() float32
func (s *Sprite) CenterY() float32
```

### Bounding Box

```go
func (s *Sprite) BoundingBox(parentWidth, parentHeight float32) ti.Rectangle[float32]
```

Returns the axis-aligned bounding box after applying transforms (scale, rotation).

**Parameters:**
- `parentWidth` - Parent display area width
- `parentHeight` - Parent display area height

## Texture Methods

```go
func (s *Sprite) SetTexture(texture *ti.TiImage) ISprite
func (s *Sprite) FillTexture(color color.Color)
func (s *Sprite) Texture() *ti.TiImage
```

## Mask Methods

```go
func (s *Sprite) SetMask(mask IMask) ISprite
func (s *Sprite) Mask() IMask
```

See [Mask](mask.md) for more details.

## Renderer

```go
func (s *Sprite) Renderer() *Renderer
```

Returns the sprite's renderer instance.

## Release

```go
func (s *Sprite) Release()
```

Releases sprite resources (must be called to avoid GPU memory leaks).

## Related

- [ImageSprite](image.md) - Image-specific sprite
- [Mask](mask.md) - Mask system for sprites
- [Stage](stage.md) - Stage that renders sprites
