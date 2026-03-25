# go-canvas

A GPU-accelerated multimedia rendering engine built on [go-taichi](https://github.com/go-mixed/go-taichi).

## Features

- **GPU-Accelerated Rendering** - Leverages Taichi's high-performance GPU compute capabilities
- **Scene Graph Architecture** - Hierarchical structure with Stage, Scene, and Sprite objects
- **Mask & Shape System** - Powerful masking with feather support and various shape primitives
- **Rich Transform Support** - Scale, rotation, translation, and alpha blending
- **Multiple Backend Support** - CUDA, DirectX 12, and more Taichi backends
- **Effects System** - Ken Burns, fade, slide, zoom, wipe, and more transition effects

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                          Stage                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   Sprite    │  │   Sprite    │  │   Sprite    │  ...    │
│  │  (Image)    │  │   (Text)    │  │ (Spectrum)  │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│         │                │                │                 │
│         ▼                ▼                ▼                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │    Mask     │  │    Mask     │  │    Mask     │  ...    │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```


## Quick start

- Go 1.25+


### 1. Download Runtime

Download `runtime.7z` from [Taichi GitHub Releases](https://github.com/go-mixed/go-taichi/releases):

### 2. Extract to Your Project

Create a `lib/` directory and extract runtime files:

```
your_project/
└── lib/
    └── windows/          # (or linux/, darwin/)
        ├── taichi_c_api.dll
        ├── runtime_x64.bc
        ├── runtime_cuda.bc
        ├── runtime_dx12.bc
        └── slim_libdevice.10.bc
```

### 3. Set TI_LIB_DIR Environment Variable

```powershell
# Windows PowerShell
$env:TI_LIB_DIR = "C:\path\to\your\project\lib\windows"
go run main.go
```

```bash
# Linux
export TI_LIB_DIR=/path/to/your/project/lib/linux
go run main.go
```

```bash
# macOS
export TI_LIB_DIR=/path/to/your/project/lib/darwin
go run main.go
```

## Core Modules

### Rendering Pipeline

| Module | Description |
|--------|-------------|
| [Renderer](docs/renderer.md) | Entry point for rendering contexts |
| [Stage](docs/stage.md) | Main canvas holding and rendering sprites |

### Sprite System

| Module | Description |
|--------|-------------|
| [Sprite](docs/sprite.md) | Base visual element with transforms |
| [ImageSprite](docs/image.md) | Sprite that loads images from files |

### Mask System

| Module | Description |
|--------|-------------|
| [Mask](docs/mask.md) | Visibility regions for sprites |
| [Shape](docs/shape.md) | Shape types for ShapeMask |

### Effects

| Module | Description |
|--------|-------------|
| [Effect](docs/effect.md) | Base effect interface |
| [Effect Options](docs/effect_options.md) | Easing functions and configuration |
| [Transition](docs/transition.md) | Fade, slide, zoom, wipe, rotate effects |
| [Ken Burns](docs/kenburns.md) | Pan and zoom animation |

## Backends

go-canvas supports multiple Taichi backends:

- `taichi.ArchCuda` - NVIDIA CUDA
- `taichi.ArchDx12` - DirectX 12
- `taichi.ArchVulkan` - Vulkan (future)
- `taichi.ArchMetal` - Metal (future)

## Documentation

Detailed documentation is available in [/docs](docs/):

- [docs/SUMMARY.md](docs/SUMMARY.md) - Documentation index
- `docs/renderer.md` - Renderer module
- `docs/stage.md` - Stage module
- `docs/sprite.md` - Sprite base class
- `docs/image.md` - ImageSprite
- `docs/mask.md` - Mask system
- `docs/shape.md` - Shape types
- `docs/effect.md` - Effect base
- `docs/effect_options.md` - Easing and options
- `docs/transition.md` - Transition effects
- `docs/kenburns.md` - Ken Burns effect

### Visual Guides

- `docs/bounding_box_visualization.png` - Bounding box system
- `docs/normalized_coords_visualization.png` - Normalized coordinate system
- `docs/bounding_box.py` - Bounding box visualization code
- `docs/normalized_coords.py` - Coordinate visualization code

## Project Structure

```
go-canvas/
├── main.go           # Entry point
├── go.mod            # Go module
├── render/           # Core rendering package
│   ├── render.go     # Renderer
│   ├── stage.go      # Stage (canvas)
│   ├── sprite.go     # Sprite base
│   ├── image.go      # Image sprite
│   ├── mask.go       # Mask system
│   ├── shape.go      # Shape definitions
│   └── const.go      # Constants
├── effect/           # Effects package
│   ├── effect.go     # Base effect
│   ├── effect_options.go
│   ├── transition.go  # Transition effects
│   ├── kenburns.go   # Ken Burns effect
│   ├── easing.go     # Easing functions
│   └── const.go
├── ti/               # Taichi bindings
├── examples/         # Example assets
└── docs/             # Detailed documentation
```


## License

MIT
