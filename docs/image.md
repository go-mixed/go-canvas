# ImageSprite

`ImageSprite` is a specialized Sprite that loads and displays images from files.

## Symbol Overview

```
ImageSprite (extends Sprite)
├── NewImageSprite(renderer, path) (ISprite, error)
└── Inherits all Sprite methods
```

## Constructor

### NewImageSprite

```go
func NewImageSprite(renderer *Renderer, path string) (ISprite, error)
```

Creates a new ImageSprite by loading an image from the specified file path.

**Parameters:**
- `renderer` - The Renderer instance
- `path` - Path to the image file (supports common formats like JPEG, PNG, etc.)

**Returns:**
- `ISprite` - Image sprite instance
- `error` - Error if loading fails

**Example:**

```go
img, err := render.NewImageSprite(renderer, "photo.jpg")
if err != nil {
    panic(err)
}
defer img.Release()

// Scale to fit stage
img.ResizeTo(720, 1280)
```

## Usage with Stage

```go
stage, err := render.NewStage(renderer, 720, 1280)
if err != nil {
    panic(err)
}
defer stage.Release()

img, err := render.NewImageSprite(renderer, "photo.jpg")
if err != nil {
    panic(err)
}
defer img.Release()

img.ResizeTo(720, 1280)
stage.Add(img)
```

## Applying Mask

```go
mask, err := render.NewShapeMask(renderer, 720, 1280, 360, 640)
if err != nil {
    panic(err)
}
defer mask.Release()

mask.DrawShape(ti.ShapeTypeCircle, 0.5)
img.SetMask(mask)
```

## Related

- [Sprite](sprite.md) - Base sprite class
- [Mask](mask.md) - Mask system
- [Stage](stage.md) - Stage that renders sprites
