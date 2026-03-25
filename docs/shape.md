# Shape

Shape type definitions used by ShapeMask for drawing geometric masks.

## ShapeType

```go
type ShapeType string
```

### Directional Shapes (8 directions)

| Type | Distance Metric | Description |
|------|----------------|-------------|
| `ShapeTypeLinear` | - | Directional (8 directions) |
| `ShapeTypeCircle` | Euclidean | Circular shape |
| `ShapeTypeDiamond` | Manhattan | Diamond/rhombus shape |
| `ShapeTypeRectangle` | Chebyshev | Rectangular shape |

### Special Shapes (using dedicated kernels)

| Type | Description |
|------|-------------|
| `ShapeTypeTriangle` | Triangle shape |
| `ShapeTypeStar5` | 5-pointed star |
| `ShapeTypeHeart` | Heart shape |
| `ShapeTypeCross` | Cross shape |

## Distance Metrics

| Metric | Formula | Shape |
|--------|---------|-------|
| Euclidean | `sqrt(dx² + dy²)` | Circle |
| Manhattan | `|dx| + |dy|` | Diamond |
| Chebyshev | `max(|dx|, |dy|)` | Rectangle |

## ShapeOptions

```go
type ShapeOptions struct {
    Direction Direction
    Color     color.Color
}
```

### Option Functions

```go
func WithShapeDirection(direction Direction) func(*ShapeOptions)
func WithShapeColor(color color.Color) func(*ShapeOptions)
```

## ShapeTypeFromString

```go
func ShapeTypeFromString(s string) ShapeType
```

Converts a string to ShapeType. Returns `ShapeTypeLinear` if no match found.

## FeatherMode

`FeatherMode` is defined in `ti` package but used by Mask, not Shape.

```go
type FeatherMode int

const (
    FeatherModeLinear     FeatherMode = 0
    FeatherModeConic      FeatherMode = 1
    FeatherModeSmoothstep FeatherMode = 2
    FeatherModeSigmoid    FeatherMode = 3
)
```

See [Mask](mask.md) for feather usage.

## Usage

```go
import "github.com/go-mixed/go-canvas/ti"

// Draw a circular mask
mask.DrawShape(ti.ShapeTypeCircle, 1.0)

// Draw a directional wipe with custom direction
mask.DrawShape(ti.ShapeTypeLinear, 1.0, ti.WithShapeDirection(ti.DirectionRight))

// Draw a star with color
mask.DrawShape(ti.ShapeTypeStar5, 1.0, ti.WithShapeColor(color.RGBA{255, 255, 255, 255}))
```

## Related

- [Mask](mask.md) - ShapeMask uses these types with FeatherMode
