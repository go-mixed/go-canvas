# Effect Options

Configuration options for effects.

## Symbol Overview

```
EffectInOut
├── Type: uint32
└── Values: EffectIn (0), EffectOut (1)

EasingFunction
├── Type: func(float32) float32
└── Predefined: linear, ease, ease-in, ease-out, ease-in-out

EffectFn
└── Type: func(inOut EffectInOut, opts ...optionFn) IEffect

effectOptions
├── Fields: panOptions, rotateOptions, wipeOptions, zoomOptions, easingFn, direction
└── Internal use

optionFn
└── Type: func(options *effectOptions)
```

## EffectInOut

```go
type EffectInOut uint32

const (
    EffectIn  EffectInOut = 0
    EffectOut EffectInOut = 1
)
```

- `EffectIn` - Effect plays from 0 to 1 (fade in, zoom in, etc.)
- `EffectOut` - Effect plays from 1 to 0 (fade out, zoom out, etc.)

## EasingFunction

```go
type EasingFunction func(float32) float32
```

Easing functions transform the linear progress value. Input and output are in range [0, 1], but output can exceed this range for elastic effects.

### Predefined Easing Functions

| Name | Description |
|------|-------------|
| `linear` | Linear, no easing |
| `ease` | Equivalent to cubic-bezier(0.25, 0.1, 0.25, 1.0) |
| `ease-in` | Equivalent to cubic-bezier(0.42, 0, 1.0, 1.0) |
| `ease-out` | Equivalent to cubic-bezier(0, 0, 0.58, 1.0) |
| `ease-in-out` | Equivalent to cubic-bezier(0.42, 0, 0.58, 1.0) |

### GetEasingFunction

```go
func GetEasingFunction(name string) EasingFunction
```

Gets easing function by name. Returns `DefaultEasingFunction` if name not found.

### Custom Easing

Easing functions have the signature `func(float32) float32`:

```go
// Example - custom quadratic ease
func myEase(t float32) float32 {
    return t * t
}
```

### Cubic Bezier

```go
func cubicBezier(p1x, p1y, p2x, p2y float32) EasingFunction
```

Creates a custom cubic bezier easing function. Control points can exceed [0,1] range for elastic effects.

```go
// Bounce effect example
bounce := cubicBezier(0.68, -0.55, 0.265, 1.55)
value := bounce(0.5)
```

## EffectFn

```go
type EffectFn func(inOut EffectInOut, opts ...optionFn) IEffect
```

Factory function type for creating effects.

## optionFn

```go
type optionFn func(options *effectOptions)
```

Functional options for configuring effects.

### Available Options

| Option | Description |
|--------|-------------|
| `WithDirection(direction)` | Set pan/wipe direction |
| `WithDirectionStr(name)` | Set direction by name string |
| `WithZoomRange(min, max)` | Set zoom start/end values |
| `WithPanIntensity(float)` | Set pan intensity |
| `WithEasing(name)` | Set easing function by name |
| `WithRotateAngle(start, end)` | Set rotation angle range |
| `WithRotateScale(start, end)` | Set rotation scale range |
| `WithShapeType(shapeType)` | Set wipe shape type |
| `WithShapeTypeStr(name)` | Set wipe shape by name string |

## Usage

```go
// Create Ken Burns effect with custom options
effect := NewKenBurnsEffect(EffectIn,
    WithZoomRange(1.0, 1.3),
    WithPanIntensity(0.15),
    WithEasing("ease-in-out"),
    WithDirection(ti.DirectionTopLeft),
)

// Apply effect to sprite
effect.Apply(sprite, progress)
```

## Related

- [Effect](effect.md) - Base effect module
- [Ken Burns Effect](kenburns.md) - Uses effect options
- [Transition Effects](transition.md) - All transition effects
