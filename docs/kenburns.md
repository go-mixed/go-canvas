# Ken Burns Effect

The Ken Burns effect is a pan and zoom animation that adds dynamic movement to static images.

## Symbol Overview

```
KenBurnsEffect
├── Embeds: *Effect
├── NewKenBurnsEffect(inOut EffectInOut, opts ...optionFn) IEffect
└── Methods:
    ├── Apply(sprite render.ISpriteOperator, progress float32)
    └── calculatePan(width, height, progress float32) (x, y float32)
```

## Constructor

### NewKenBurnsEffect

```go
func NewKenBurnsEffect(inOut EffectInOut, opts ...optionFn) IEffect
```

Creates a new Ken Burns effect with the specified options.

**Default options:**
- ZoomStart: 1.0
- ZoomEnd: 1.2
- PanIntensity: 0.1

## Methods

### Apply

```go
func (e *KenBurnsEffect) Apply(sprite render.ISpriteOperator, progress float32)
```

Applies the Ken Burns effect to a sprite at the given progress (0.0 to 1.0).

The effect combines:
- **Zoom** - Scale changes from ZoomStart to ZoomEnd
- **Pan** - Position shifts based on calculated direction

### calculatePan

```go
func (e *KenBurnsEffect) calculatePan(width, height, progress float32) (float32, float32)
```

Calculates the pan offset for the current progress.

**Returns:**
- `x` - X offset
- `y` - Y offset

## Available Directions

Uses `ti.Direction` enum:

| Name | Direction | Description |
|------|-----------|-------------|
| `DirectionTop` | Up | Pan upward |
| `DirectionBottom` | Down | Pan downward |
| `DirectionLeft` | Left | Pan leftward |
| `DirectionRight` | Right | Pan rightward |
| `DirectionTopLeft` | Up-Left | Pan toward upper-left |
| `DirectionTopRight` | Up-Right | Pan toward upper-right |
| `DirectionBottomLeft` | Down-Left | Pan toward lower-left |
| `DirectionBottomRight` | Down-Right | Pan toward lower-right |
| `DirectionCenter` | Center | No pan (zoom only) |

## Usage

```go
import (
    "github.com/go-mixed/go-canvas/effect"
    "github.com/go-mixed/go-canvas/ti"
)

// Create Ken Burns effect
kb := effect.NewKenBurnsEffect(effect.EffectIn,
    effect.WithZoomRange(1.0, 1.3),
    effect.WithPanIntensity(0.15),
    effect.WithDirection(ti.DirectionTopLeft),
)

// Apply during animation loop
for i := 0; i <= 100; i++ {
    progress := float32(i) / 100.0
    kb.Apply(sprite, progress)
}
```

## Visual Description

```
Start (progress=0)              End (progress=1)
┌─────────────┐               ┌─────────────┐
│             │               │   ┌───┐     │
│   ┌─────┐   │     ───►      │   │   │     │
│   │     │   │               │   └───┘     │
│   └─────┘   │               │             │
│             │               └─────────────┘
└─────────────┘               (zoomed in, panned)
```

## Related

- [Effect](effect.md) - Base effect interface
- [Effect Options](effect_options.md) - Easing and option functions
- [Transition Effects](transition.md) - Other transition types
