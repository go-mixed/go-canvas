package render

import (
	"github.com/go-mixed/go-canvas/ti"
	"github.com/go-mixed/go-taichi/taichi"
)

type Mask struct {
	texture   *ti.TiMask
	distField *ti.TiGrid // 距离场
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
		texture:   texture,
		distField: distField,
		renderer:  renderer,
	}, nil
}

// FillWithTexture 将纹理填充到 Mask
func (m *Mask) FillWithTexture(texture *ti.TiImage) {
	//  将图像转换为遮罩（提取 alpha 通道）
	m.renderer.Module().ImageToMask(texture, m.texture)
}

func (m *Mask) ApplyFeather(featherRadius uint32, featherMode ti.FeatherMode) {
	// 计算距离场（使用复用的 distField）
	m.renderer.Module().ComputeDistanceField(m.texture, m.distField)

	// 应用羽化
	m.renderer.Module().ComputeFeather(m.distField, m.texture, float32(featherRadius), featherMode)
}

func (m *Mask) Release() {
	if m.texture != nil {
		m.texture.Release()
	}
	if m.distField != nil {
		m.distField.Release()
	}
}

func (m *Mask) Texture() *taichi.NdArray {
	return m.texture
}

type ShapeMask struct {
	*Mask
	*ShapeSprite

	featherRadius uint32
	featherMode   ti.FeatherMode
}

var _ IMask = (*ShapeMask)(nil)
var _ IElement = (*ShapeMask)(nil)

func NewShapeMask(renderer *Renderer, width, height uint32, cx, cy uint32) (*ShapeMask, error) {
	mask, err := NewMask(renderer, width, height)
	if err != nil {
		return nil, err
	}

	shapeSprite, err := NewShapeSprite(renderer, width, height, cx, cy)
	if err != nil {
		mask.Release()
		return nil, err
	}

	return &ShapeMask{
		Mask:        mask,
		ShapeSprite: shapeSprite.(*ShapeSprite),

		featherRadius: 0,
		featherMode:   ti.FeatherModeLinear,
	}, nil
}

func (m *ShapeMask) SetFeather(radius uint32, featherMode ti.FeatherMode) {
	m.featherRadius = radius
	m.featherMode = featherMode
}

func (m *ShapeMask) DrawShape(shapeType ti.ShapeType, tVal float32, fns ...func(option *ti.ShapeOptions)) {
	m.ShapeSprite.DrawShape(shapeType, tVal, fns...)
	m.Mask.FillWithTexture(m.ShapeSprite.Texture())
	if m.featherRadius > 0 {
		m.Mask.ApplyFeather(m.featherRadius, m.featherMode)
	}
}

func (m *ShapeMask) Texture() *ti.TiMask {
	return m.Mask.Texture()
}

func (m *ShapeMask) Release() {
	m.Mask.Release()
	m.ShapeSprite.Release()
}
