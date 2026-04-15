package ctypes

import "image/color"

// BlurMode 模糊模式
type BlurMode int32

const (
	BlurModeBox      BlurMode = 0 // 普通模糊
	BlurModeGaussian BlurMode = 1 // 高斯模糊
	BlurModeMosaic   BlurMode = 2 // 马赛克
)

type ImageWriter interface {
	Set(x, y int, c color.Color)
}
