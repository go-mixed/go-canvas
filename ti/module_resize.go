package ti

// resizeParams 缩放参数
type resizeParams struct {
	scaleX  float32
	scaleY  float32
	offsetX float32
	offsetY float32
}

// computeResizeScaleAndOffset 计算缩放比例和偏移量
func computeResizeScaleAndOffset(srcWidth, srcHeight, dstWidth, dstHeight float32, fillMode FillMode) resizeParams {
	scaleX := dstWidth / srcWidth
	scaleY := dstHeight / srcHeight

	var offsetX, offsetY float32

	switch fillMode {
	case FillModeStretch:
		// 直接拉伸，scaleX/scaleY 不变
		offsetX = 0
		offsetY = 0
	case FillModeFit:
		// 等比适应，可能有黑边
		scale := min(scaleX, scaleY)
		scaleX = scale
		scaleY = scale
		offsetX = (dstWidth - srcWidth*scaleX) * 0.5
		offsetY = (dstHeight - srcHeight*scaleY) * 0.5
	case FillModeFill:
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

// Resize 缩放纹理
func (m *AotModule) Resize(input *TiImage, output *TiImage, opts ResizeOptions) {
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
	case ScaleModeNearest:
		kernelName = "resize_nearest"
	case ScaleModeLinear:
		kernelName = "resize_bilinear"
	case ScaleModeCubic:
		kernelName = "resize_bicubic"
	case ScaleModeLanczos:
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
		Run()
}
