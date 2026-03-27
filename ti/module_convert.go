package ti

import "image/color"

// FillTexture 填充纹理
func (m *AotModule) FillTexture(texture *TiImage, c color.Color) {
	kernel := m.getCache("fill_texture")

	kernel.Launch().ArgNdArray(texture).ArgVectorFloat32(Color2TiColor(c)...).Run()
}

// CvToTiImage 将 CvImage (h, w, [b, g, r]) 转换为 TiImage (w, h, [r, g, b, a])
func (m *AotModule) CvToTiImage(inputImage *CvImage, outputImage *TiImage) {
	kernel := m.getCache("cv_image_to_ti")
	kernel.Launch().ArgNdArray(inputImage).ArgNdArray(outputImage).Run()
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
