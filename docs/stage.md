# Stage

The `Stage` is the main canvas that holds all visual elements (Sprites). It orchestrates rendering of the entire scene.

## Symbol Overview

```
Stage
├── Fields: children, renderer, screen, mutex
├── NewStage(renderer, width, height) (*Stage, error)
└── Methods:
    ├── Add(sprite)
    ├── Remove(sprite)
    ├── Children() []Sprite
    ├── Render()
    ├── Texture() *ti.TiImage
    └── Release()
```

## Constructor

### NewStage

```go
func NewStage(renderer *Renderer, width, height int) (*Stage, error)
```

Creates a new Stage (canvas) with the specified dimensions.

**Parameters:**
- `renderer` - The Renderer instance
- `width` - Canvas width in pixels
- `height` - Canvas height in pixels

**Returns:**
- `*Stage` - Stage instance
- `error` - Error if creation fails

## Methods

### Add

```go
func (s *Stage) Add(sprite Sprite) error
```

Adds a sprite to the stage. The sprite will be rendered in the order it was added (later additions render on top).

### Remove

```go
func (s *Stage) Remove(sprite Sprite) error
```

Removes a sprite from the stage.

### Children

```go
func (s *Stage) Children() []Sprite
```

Returns all sprites currently on the stage.

### Render

```go
func (s *Stage) Render()
```

Renders the entire stage. Must be called before accessing the texture.

### Texture

```go
func (s *Stage) Texture() *ti.TiImage
```

Returns the rendered texture after calling Render(). Use with `ti.SaveTiImageToFile()` to save to disk.

### Release

```go
func (s *Stage) Release()
```

Releases all resources held by the stage.

## Usage Pattern

```go
stage, err := render.NewStage(renderer, 720, 1280)
if err != nil {
    panic(err)
}
defer stage.Release()

stage.Add(sprite1)
stage.Add(sprite2)

stage.Render()
texture := stage.Texture()

err = ti.SaveTiImageToFile(texture, "output.png")
```

## Related

- [Renderer](renderer.md) - Creates the Stage
- [Sprite](sprite.md) - Visual elements added to Stage
