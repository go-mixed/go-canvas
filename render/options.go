package render

import (
	"github.com/go-mixed/go-canvas/internel/misc"
)

type stageOptions struct {
	enabledRAWImage bool
	logger          misc.Logger
}

type stageOptFunc func(o *stageOptions)

// WithRawImage 是否在渲染时，同时输出到BGRA的图片流中。这样可以通过 GetBgraImage 来获取 RAW 图片数据。
func WithRawImage(v bool) stageOptFunc {
	return func(o *stageOptions) {
		o.enabledRAWImage = v
	}
}

func WithLogger(logger misc.Logger) stageOptFunc {
	return func(o *stageOptions) {
		o.logger = logger
	}
}
