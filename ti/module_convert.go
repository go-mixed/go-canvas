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
