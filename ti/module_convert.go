package ti

import "image/color"

func (m *AotModule) AsyncTiImageToBgra(input *TiImage, output *BgraImage) {
	kernel := m.getCache("ti_image_to_bgra")
	kernel.Launch().ArgNdArray(input).ArgNdArray(output).RunAsync()
}

func (m *AotModule) TiImageToBgra(input *TiImage, output *BgraImage) {
	m.AsyncTiImageToBgra(input, output)
	m.runtime.Wait()
}

// AsyncFillColor 填充纹理
func (m *AotModule) AsyncFillColor(texture *TiImage, c color.Color) {
	kernel := m.getCache("fill_color")

	kernel.Launch().ArgNdArray(texture).ArgVectorFloat32(Color2TiColor(c)...).RunAsync()
}

func (m *AotModule) FillColor(texture *TiImage, c color.Color) {
	m.AsyncFillColor(texture, c)
	m.runtime.Wait()
}

// AsyncBlur 对纹理进行模糊处理
func (m *AotModule) AsyncBlur(input *TiImage, output *TiImage, mode BlurMode, radius int32) {
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
		RunAsync()
}

func (m *AotModule) Blur(input *TiImage, output *TiImage, mode BlurMode, radius int32) {
	m.AsyncBlur(input, output, mode, radius)
	m.runtime.Wait()
}
