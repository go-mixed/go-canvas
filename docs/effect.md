# Effect

The `Effect` module provides the base interface and utilities for visual effects in go-canvas.

## Symbol Overview

```
IEffect (interface)
└── Apply(sprite render.ISpriteOperator, progress float32)

Effect
├── Fields: options, direction
├── newEffect(inOut EffectInOut, options effectOptions) *Effect
└── Methods:
    └── getEaseProgress(progress float32) float32
```

## IEffect Interface

All effects implement the `IEffect` interface:

```go
type IEffect interface {
    Apply(sprite render.ISpriteOperator, progress float32)
}
```

**Parameters:**
- `sprite` - The sprite to apply the effect to (must implement `render.ISpriteOperator`)
- `progress` - Effect progress from 0.0 (start) to 1.0 (complete)

## Effect Struct

Base effect struct used by all effect implementations.

### Fields

```go
type Effect struct {
    options   effectOptions
    direction EffectInOut
}
```

### EffectInOut

```go
type EffectInOut uint32

const (
    EffectIn  EffectInOut = 0
    EffectOut EffectInOut = 1
)
```

- `EffectIn` - Effect plays forward (0 to 1)
- `EffectOut` - Effect plays backward (1 to 0)

## getEaseProgress

```go
func (e *Effect) getEaseProgress(progress float32) float32
```

Applies easing to the progress value based on the configured easing function and direction.

## Related

- [Transition Effects](transition.md) - Fade, Slide, Zoom, Wipe, Rotate
- [Ken Burns Effect](kenburns.md) - Pan and zoom effect
- [Effect Options](effect_options.md) - Easing and direction options
