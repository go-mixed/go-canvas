package ti

import "image/color"

// FillColor 填充纹理
func (m *AotModule) FillColor(texture *TiImage, c color.Color) {
	kernel := m.getCache("fill_color")

	kernel.Launch().ArgNdArray(texture).ArgVectorFloat32(Color2TiColor(c)...).Run()
}

// Blur 对纹理进行模糊处理
func (m *AotModule) Blur(input *TiImage, output *TiImage, mode BlurMode, radius int32) {
	var kernelName string
	switch mode {
	case BlurModeBox:
		kernelName = "blur_box"
	case BlurModeGaussian:
		kernelName = "blur_gaussian"
	case BlurModeMosaic:
		kernelName = "blur_mosaic"
	}

	kernel := m.getCache(kernelName)
	kernel.Launch().
		ArgNdArray(input).
		ArgNdArray(output).
		ArgInt32(radius).
		Run()
}
