# go-canvas Documentation

Detailed documentation for go-canvas modules and components.

## Core Modules

### Renderer & Stage

| Document | Description |
|----------|-------------|
| [renderer.md](renderer.md) | Entry point for rendering contexts |
| [stage.md](stage.md) | Main canvas that holds and renders sprites |

### Sprite System

| Document | Description |
|----------|-------------|
| [sprite.md](sprite.md) | Base visual element with transform properties |
| [image.md](image.md) | Sprite that loads images from files |

### Mask System

| Document | Description |
|----------|-------------|
| [mask.md](mask.md) | Visibility regions for sprites |
| [shape.md](shape.md) | Shape types for ShapeMask |

### Effects

| Document | Description |
|----------|-------------|
| [effect.md](effect.md) | Base effect interface and utilities |
| [effect_options.md](effect_options.md) | Easing functions and effect configuration |
| [transition.md](transition.md) | Transition effects (fade, slide, zoom, wipe, rotate) |
| [kenburns.md](kenburns.md) | Ken Burns pan and zoom effect |

## Visual Guides

See the main [README.md](../README.md) for visual explanations of concepts.

| File | Description |
|------|-------------|
| `bounding_box_visualization.png` | Sprite rotation and bounding box diagram |
| `normalized_coords_visualization.png` | Normalized coordinate system |
| `bounding_box.py` | Python script to generate bounding box visualization |
| `normalized_coords.py` | Python script to generate coordinate visualization |

## Module Hierarchy

```
Renderer
    └── Stage
            └── Sprite
                    ├── ImageSprite
                    └── Mask (ShapeMask)
```

## Quick Links

- [Main README](../README.md)
- [Render Package](../render/)
- [Effect Package](../effect/)
- [Taichi Bindings](../ti/)
