package ti

import (
	"image/color"

	"github.com/go-mixed/go-canvas/ctypes"
)

func (m *AotModule) AsyncTiImageToBgra(input *ctypes.TiImage, output *ctypes.BgraImage) {
	kernel := m.getCache("ti_image_to_bgra")
	kernel.Launch().ArgNdArray(input).ArgNdArray(output).RunAsync()
}

func (m *AotModule) TiImageToBgra(input *ctypes.TiImage, output *ctypes.BgraImage) {
	m.AsyncTiImageToBgra(input, output)
	m.runtime.Wait()
}

// AsyncFillColor 填充纹理
func (m *AotModule) AsyncFillColor(texture *ctypes.TiImage, c color.Color) {
	kernel := m.getCache("fill_color")

	kernel.Launch().ArgNdArray(texture).ArgVectorFloat32(ctypes.Color2TiColor(c)...).RunAsync()
}

func (m *AotModule) FillColor(texture *ctypes.TiImage, c color.Color) {
	m.AsyncFillColor(texture, c)
	m.runtime.Wait()
}

// AsyncBlur 对纹理进行模糊处理
func (m *AotModule) AsyncBlur(input *ctypes.TiImage, output *ctypes.TiImage, mode ctypes.BlurMode, radius int32) {
	var kernelName string
	switch mode {
	case ctypes.BlurModeBox:
		kernelName = "blur_box"
	case ctypes.BlurModeGaussian:
		kernelName = "blur_gaussian"
	case ctypes.BlurModeMosaic:
		kernelName = "blur_mosaic"
	}

	kernel := m.getCache(kernelName)
	kernel.Launch().
		ArgNdArray(input).
		ArgNdArray(output).
		ArgInt32(radius).
		RunAsync()
}

func (m *AotModule) Blur(input *ctypes.TiImage, output *ctypes.TiImage, mode ctypes.BlurMode, radius int32) {
	m.AsyncBlur(input, output, mode, radius)
	m.runtime.Wait()
}
