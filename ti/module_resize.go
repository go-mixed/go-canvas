package ti

import "github.com/go-mixed/go-canvas/ctypes"

// resizeParams 缩放参数
type resizeParams struct {
	scaleX  float32
	scaleY  float32
	offsetX float32
	offsetY float32
}

// computeResizeScaleAndOffset 计算缩放比例和偏移量
func computeResizeScaleAndOffset(srcWidth, srcHeight, dstWidth, dstHeight float32, fillMode ctypes.FillMode) resizeParams {
	scaleX := dstWidth / srcWidth
	scaleY := dstHeight / srcHeight

	var offsetX, offsetY float32

	switch fillMode {
	case ctypes.FillModeStretch:
		// 直接拉伸，scaleX/scaleY 不变
		offsetX = 0
		offsetY = 0
	case ctypes.FillModeFit:
		// 等比适应，可能有黑边
		scale := min(scaleX, scaleY)
		scaleX = scale
		scaleY = scale
		offsetX = (dstWidth - srcWidth*scaleX) * 0.5
		offsetY = (dstHeight - srcHeight*scaleY) * 0.5
	case ctypes.FillModeFill:
		// 等比填充，可能裁剪
		scale := max(scaleX, scaleY)
		scaleX = scale
		scaleY = scale
		offsetX = (dstWidth - srcWidth*scaleX) * 0.5
		offsetY = (dstHeight - srcHeight*scaleY) * 0.5
	}

	return resizeParams{
		scaleX:  scaleX,
		scaleY:  scaleY,
		offsetX: offsetX,
		offsetY: offsetY,
	}
}

// AsyncResize 缩放纹理
func (m *AotModule) AsyncResize(input *ctypes.TiImage, output *ctypes.TiImage, opts ctypes.ResizeOptions) {
	// 获取源和目标的尺寸
	srcShape := input.Shape()
	dstShape := output.Shape()
	srcWidth := float32(srcShape[0])
	srcHeight := float32(srcShape[1])
	dstWidth := float32(dstShape[0])
	dstHeight := float32(dstShape[1])

	// 计算缩放参数
	params := computeResizeScaleAndOffset(srcWidth, srcHeight, dstWidth, dstHeight, opts.FillMode)

	// 根据 ScaleMode 调用对应的 kernel
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
		ArgFloat32(srcWidth).ArgFloat32(srcHeight).
		ArgFloat32(dstWidth).ArgFloat32(dstHeight).
		ArgFloat32(params.scaleX).ArgFloat32(params.scaleY).
		ArgFloat32(params.offsetX).ArgFloat32(params.offsetY).
		RunAsync()
}

func (m *AotModule) Resize(input *ctypes.TiImage, output *ctypes.TiImage, opts ctypes.ResizeOptions) {
	m.AsyncResize(input, output, opts)
	m.runtime.Wait()
}
