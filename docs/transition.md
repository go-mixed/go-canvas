# Transition Effects

Transition effects animate sprites during state changes like slides, fades, zooms, and wipes.

## Symbol Overview

```
EffectFn
└── Type: func(inOut EffectInOut, opts ...optionFn) IEffect

Factory Functions:
├── NewKenBurnsEffect(inOut EffectInOut, opts ...optionFn) IEffect
├── NewFadeEffect(inOut EffectInOut, opts ...optionFn) IEffect
├── NewRotateEffect(inOut EffectInOut, opts ...optionFn) IEffect
├── NewSlideEffect(inOut EffectInOut, opts ...optionFn) IEffect
├── NewZoomEffect(inOut EffectInOut, opts ...optionFn) IEffect
└── NewWipeEffect(inOut EffectInOut, opts ...optionFn) IEffect

Effect Types:
├── KenBurnsEffect (extends Effect)
├── FadeEffect (extends Effect)
├── RotateEffect (extends Effect)
├── SlideEffect (extends Effect)
├── ZoomEffect (extends Effect)
└── WipeEffect (extends Effect)

GetTransitionEffect
└── func GetTransitionEffect(name string) (EffectFn, error)
```

## Effect Factory Functions

All factory functions follow the same signature:

```go
func NewXXXEffect(inOut EffectInOut, opts ...optionFn) IEffect
```

### NewKenBurnsEffect

```go
func NewKenBurnsEffect(inOut EffectInOut, opts ...optionFn) IEffect
```

Creates a Ken Burns (pan + zoom) effect.

**Default options:**
- ZoomStart: 1.0, ZoomEnd: 1.2
- PanIntensity: 0.1

### NewFadeEffect

```go
func NewFadeEffect(inOut EffectInOut, opts ...optionFn) IEffect
```

Creates a fade effect (opacity transition).

### NewRotateEffect

```go
func NewRotateEffect(inOut EffectInOut, opts ...optionFn) IEffect
```

Creates a rotation effect.

**Default options:**
- AngleStart: 0, AngleEnd: 360
- ScaleStart: 0.5, ScaleEnd: 1
- Direction: DirectionRight

### NewSlideEffect

```go
func NewSlideEffect(inOut EffectInOut, opts ...optionFn) IEffect
```

Creates a slide effect (position transition).

**Default options:**
- Direction: DirectionRight

### NewZoomEffect

```go
func NewZoomEffect(inOut EffectInOut, opts ...optionFn) IEffect
```

Creates a zoom effect.

**Default options:**
- ZoomStart: 0.5, ZoomEnd: 1.0

### NewWipeEffect

```go
func NewWipeEffect(inOut EffectInOut, opts ...optionFn) IEffect
```

Creates a wipe effect (directional reveal using ShapeSprite as mask).

**Default options:**
- ShapeType: ShapeTypeRectangle
- Direction: DirectionRight

## GetTransitionEffect

```go
func GetTransitionEffect(name string) (EffectFn, error)
```

Creates a transition effect factory by name.

**Supported names:**

| Name | Effect Type | Description |
|------|-------------|-------------|
| `pan_left` | KenBurns | Pan to left |
| `pan_right` | KenBurns | Pan to right |
| `pan_top` | KenBurns | Pan to top |
| `pan_bottom` | KenBurns | Pan to bottom |
| `pan_top_left` | KenBurns | Pan to top-left |
| `pan_top_right` | KenBurns | Pan to top-right |
| `pan_bottom_left` | KenBurns | Pan to bottom-left |
| `pan_bottom_right` | KenBurns | Pan to bottom-right |
| `pan_center` | KenBurns | Stay at center |
| `rotate` | RotateEffect | Rotation effect |
| `slide` | SlideEffect | Slide effect |
| `zoom` | ZoomEffect | Zoom effect |
| `wipe` | WipeEffect | Rectangle wipe |
| `fade` | FadeEffect | Fade effect |
| `heart` | WipeEffect | Heart shape wipe |
| `star5` | WipeEffect | 5-point star wipe |
| `cross` | WipeEffect | Cross shape wipe |
| `linear` | WipeEffect | Linear wipe |
| `circle` | WipeEffect | Circle wipe |
| `diamond` | WipeEffect | Diamond wipe |
| `rectangle` | WipeEffect | Rectangle wipe |
| `triangle` | WipeEffect | Triangle wipe |

## Usage

```go
// Using factory directly
fade := effect.NewFadeEffect(effect.EffectIn, effect.WithEasing("ease-out"))
fade.Apply(sprite, progress)

// Using GetTransitionEffect
effectFn, _ := effect.GetTransitionEffect("zoom")
zoom := effectFn(effect.EffectOut, effect.WithZoomRange(1.0, 1.5))
zoom.Apply(sprite, progress)
```

## Applying Effects

```go
// Create effect
effect := effect.NewKenBurnsEffect(effect.EffectIn,
    effect.WithZoomRange(1.0, 1.3),
    effect.WithPanIntensity(0.15),
    effect.WithDirection(ti.DirectionTopLeft),
)

// Apply to sprite (progress: 0.0 to 1.0)
for i := 0; i <= 100; i++ {
    progress := float32(i) / 100.0
    effect.Apply(sprite, progress)
}
```

## Effect Direction

```go
// EffectIn: plays 0 -> 1 (fade in, zoom in, etc.)
effect := NewFadeEffect(EffectIn)
effect.Apply(sprite, 0.0)   // alpha = 0
effect.Apply(sprite, 1.0)   // alpha = 1

// EffectOut: plays 1 -> 0 (fade out, zoom out, etc.)
effect := NewFadeEffect(EffectOut)
effect.Apply(sprite, 0.0)   // alpha = 1
effect.Apply(sprite, 1.0)   // alpha = 0
```

## Related

- [Effect](effect.md) - Base effect interface
- [Ken Burns](kenburns.md) - Pan and zoom effect
- [Effect Options](effect_options.md) - Easing and option functions
