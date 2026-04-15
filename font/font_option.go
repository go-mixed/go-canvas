package font

import (
	"github.com/go-mixed/go-canvas/internel/misc"
)

type FontOptions struct {
	logger misc.Logger
	dpi    float64
}

func FontOpt() *FontOptions {
	return &FontOptions{
		logger: nil,
		dpi:    96.0,
	}
}

func (o *FontOptions) SetLogger(logger misc.Logger) *FontOptions {
	o.logger = logger
	return o
}

func (o *FontOptions) SetDpi(dpi float64) *FontOptions {
	o.dpi = dpi
	return o
}
