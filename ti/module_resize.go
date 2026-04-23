package ti

import "github.com/go-mixed/go-canvas/ctypes"

// resizeRegions 内部计算结果
type resizeRegions struct {
	srcX, srcY, srcW, srcH float32
	dstX, dstY, dstW, dstH float32
}

// computeResizeRegions 根据 FillMode 自动计算 src/dst 区块；
// 若显式传入 srcRegion/dstRegion 则直接使用，完全覆盖自动计算。
func computeResizeRegions(
	srcWidth, srcHeight, dstWidth, dstHeight float32,
	fillMode ctypes.FillMode,
	srcRegion, dstRegion ctypes.Rectangle[int],
) resizeRegions {
	r := resizeRegions{}

	if !srcRegion.Empty() {
		r.srcX, r.srcY, r.srcW, r.srcH = float32(srcRegion.X()), float32(srcRegion.Y()), float32(srcRegion.Width()), float32(srcRegion.Height())
	} else {
		r.srcX, r.srcY, r.srcW, r.srcH = 0, 0, srcWidth, srcHeight
	}

	if !dstRegion.Empty() {
		r.dstX, r.dstY, r.dstW, r.dstH = float32(dstRegion.X()), float32(dstRegion.Y()), float32(dstRegion.Width()), float32(dstRegion.Height())
	} else {
		// 根据 FillMode 计算 dst 区块（src 区块已确定）
		sw, sh := r.srcW, r.srcH
		switch fillMode {
		case ctypes.FillModeStretch:
			r.dstX, r.dstY, r.dstW, r.dstH = 0, 0, dstWidth, dstHeight
		case ctypes.FillModeFit:
			scale := min(dstWidth/sw, dstHeight/sh)
			fitW := sw * scale
			fitH := sh * scale
			r.dstX = (dstWidth - fitW) * 0.5
			r.dstY = (dstHeight - fitH) * 0.5
			r.dstW = fitW
			r.dstH = fitH
		case ctypes.FillModeFill:
			scale := max(dstWidth/sw, dstHeight/sh)
			cropW := dstWidth / scale
			cropH := dstHeight / scale
			// src 区块中心裁剪
			r.srcX = r.srcX + (sw-cropW)*0.5
			r.srcY = r.srcY + (sh-cropH)*0.5
			r.srcW = cropW
			r.srcH = cropH
			r.dstX, r.dstY, r.dstW, r.dstH = 0, 0, dstWidth, dstHeight
		default:
			r.dstX, r.dstY, r.dstW, r.dstH = 0, 0, dstWidth, dstHeight
		}
	}

	return r
}

// AsyncResize 缩放纹理。
// srcRegion/dstRegion 非 nil 时直接使用指定区块，忽略 opts.FillMode 的自动计算。
func (m *AotModule) AsyncResize(input *ctypes.TiImage, output *ctypes.TiImage, opts ctypes.ResizeOptions, srcRegion, dstRegion ctypes.Rectangle[int]) {
	srcShape := input.Shape()
	dstShape := output.Shape()
	srcWidth := float32(srcShape[0])
	srcHeight := float32(srcShape[1])
	dstWidth := float32(dstShape[0])
	dstHeight := float32(dstShape[1])

	r := computeResizeRegions(srcWidth, srcHeight, dstWidth, dstHeight, opts.FillMode, srcRegion, dstRegion)

	var kernelName string
	switch opts.ScaleMode {
	case ctypes.ScaleModeNearest:
		kernelName = "resize_nearest"
	case ctypes.ScaleModeLinear:
		kernelName = "resize_bilinear"
	case ctypes.ScaleModeCubic:
		kernelName = "resize_bicubic"
	case ctypes.ScaleModeLanczos:
		kernelName = "resize_lanczos"
	}

	kernel := m.getCache(kernelName)
	kernel.Launch().
		ArgNdArray(input).
		ArgNdArray(output).
		ArgFloat32(r.srcX).ArgFloat32(r.srcY).ArgFloat32(r.srcW).ArgFloat32(r.srcH).
		ArgFloat32(r.dstX).ArgFloat32(r.dstY).ArgFloat32(r.dstW).ArgFloat32(r.dstH).
		RunAsync()
}

func (m *AotModule) Resize(input *ctypes.TiImage, output *ctypes.TiImage, opts ctypes.ResizeOptions, srcRegion, dstRegion ctypes.Rectangle[int]) {
	m.AsyncResize(input, output, opts, srcRegion, dstRegion)
	m.runtime.Wait()
}
