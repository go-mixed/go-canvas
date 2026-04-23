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

// AsyncCopy 将 src 的指定区块直接复制到 dst 的指定区块，无插值。
// 区块为零值时表示使用整张图。
func (m *AotModule) AsyncCopy(src, dst *ctypes.TiImage, srcRegion, dstRegion ctypes.Rectangle[int]) {
	srcShape := src.Shape()
	dstShape := dst.Shape()

	sr := srcRegion
	if sr.Empty() {
		sr = ctypes.RectWH(0, 0, int(srcShape[0]), int(srcShape[1]))
	}
	dr := dstRegion
	if dr.Empty() {
		dr = ctypes.RectWH(0, 0, int(dstShape[0]), int(dstShape[1]))
	}

	kernel := m.getCache("copy_region")
	kernel.Launch().
		ArgNdArray(src).
		ArgNdArray(dst).
		ArgInt32(int32(sr.X())).ArgInt32(int32(sr.Y())).ArgInt32(int32(sr.Width())).ArgInt32(int32(sr.Height())).
		ArgInt32(int32(dr.X())).ArgInt32(int32(dr.Y())).ArgInt32(int32(dr.Width())).ArgInt32(int32(dr.Height())).
		RunAsync()
}

func (m *AotModule) Copy(src, dst *ctypes.TiImage, srcRegion, dstRegion ctypes.Rectangle[int]) {
	m.AsyncCopy(src, dst, srcRegion, dstRegion)
	m.runtime.Wait()
}

// AsyncRenderBorder 在已有内容的 dst 上叠加 border 并做圆角裁剪。
// dst 中的 content/padding 区域像素保持不动，border 区域填对应颜色，圆角外强制透明。
func (m *AotModule) AsyncRenderBorder(dst *ctypes.TiImage, border ctypes.Border) {
	kernel := m.getCache("render_border")

	kernel.Launch().
		ArgNdArray(dst).
		ArgVectorFloat32(float32(border.TopWidth), float32(border.RightWidth), float32(border.BottomWidth), float32(border.LeftWidth)).
		ArgVectorFloat32(ctypes.Color2TiColor(ctypes.OrTransparentColor(border.TopColor))...).
		ArgVectorFloat32(ctypes.Color2TiColor(ctypes.OrTransparentColor(border.RightColor))...).
		ArgVectorFloat32(ctypes.Color2TiColor(ctypes.OrTransparentColor(border.BottomColor))...).
		ArgVectorFloat32(ctypes.Color2TiColor(ctypes.OrTransparentColor(border.LeftColor))...).
		ArgVectorFloat32(float32(border.TopLeftRadius), float32(border.TopRightRadius), float32(border.BottomRightRadius), float32(border.BottomLeftRadius)).
		RunAsync()
}

func (m *AotModule) RenderBorder(dst *ctypes.TiImage, border ctypes.Border) {
	m.AsyncRenderBorder(dst, border)
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
