package render

import (
	"slideshow/ti"
)

type Mask struct {
	rect      ti.Rectangle[float32]
	texture   *ti.TiMask
	distField *ti.TiGrid // 复用的距离场
	renderer  *Renderer
}

var _ IMask = (*Mask)(nil)

func NewMask(renderer *Renderer, width, height uint32) (*Mask, error) {
	texture, err := ti.NewTiMask(renderer.Runtime(), width, height)
	if err != nil {
		return nil, err
	}

	distField, err := ti.NewTiGrid(renderer.Runtime(), width, height)
	if err != nil {
		texture.Release()
		return nil, err
	}

	return &Mask{
		rect:      ti.Rect[float32](0, 0, float32(width), float32(height)),
		texture:   texture,
		distField: distField,
		renderer:  renderer,
	}, nil
}

// FillWithTexture 将纹理填充到 Mask，并应用羽化
func (s *Mask) FillWithTexture(texture *ti.TiImage) {
	//  将图像转换为遮罩（提取 alpha 通道）
	s.renderer.Module().ImageToMask(texture, s.texture)
}

func (s *Mask) ApplyFeather(featherRadius uint32, featherMode ti.FeatherMode) {
	// 计算距离场（使用复用的 distField）
	s.renderer.Module().ComputeDistanceField(s.texture, s.distField)

	// 应用羽化
	s.renderer.Module().ComputeFeather(s.distField, s.texture, float32(featherRadius), featherMode)
}

func (s *Mask) Release() {
	if s.texture != nil {
		s.texture.Release()
	}
	if s.distField != nil {
		s.distField.Release()
	}
}

func (s *Mask) Texture() *ti.TiMask {
	return s.texture
}
