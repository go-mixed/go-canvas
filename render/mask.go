package render

import (
	"github.com/go-mixed/go-canvas/ti"
	"github.com/go-mixed/go-taichi/taichi"
)

type Mask struct {
	rect      ti.Rectangle[float32]
	texture   *ti.TiMask
	distField *ti.TiGrid // 距离场
	parent    IMaskParent
	// 真实的实例
	instance IMask
}

var _ IMask = (*Mask)(nil)

func BuildMask[T IMask](parent IMaskParent, texture *ti.TiMask, instanceCreator func(*Mask) (T, error)) (T, error) {
	shape := texture.Shape()
	width, height := shape[0], shape[1]
	distField, err := ti.NewTiGrid(parent.Renderer().Runtime(), width, height)
	var nilT T
	if err != nil {
		return nilT, err
	}

	m := &Mask{
		texture:   texture,
		distField: distField,
		parent:    parent,
	}

	instance, err := instanceCreator(m)
	if err != nil {
		return nilT, err
	}
	parent.AddMask(instance)
	m.instance = instance
	return instance, nil
}

// FillWithTexture 将纹理填充到 Mask
func (m *Mask) FillWithTexture(texture *ti.TiImage) {
	//  将图像转换为遮罩（提取 alpha 通道）
	m.parent.Renderer().Module().ImageToMask(texture, m.texture)
}

func (m *Mask) ApplyFeather(featherRadius uint32, featherMode ti.FeatherMode) {
	// 计算距离场（使用复用的 distField）
	m.parent.Renderer().Module().ComputeDistanceField(m.texture, m.distField)

	// 应用羽化
	m.parent.Renderer().Module().ComputeFeather(m.distField, m.texture, float32(featherRadius), featherMode)
}

func (m *Mask) Release() {
	if m.texture != nil {
		m.texture.Release()
		m.texture = nil
	}
	if m.distField != nil {
		m.distField.Release()
		m.distField = nil
	}
}

func (m *Mask) Texture() *taichi.NdArray {
	return m.texture
}

func (m *Mask) RemoveFromParent() {
	if m.instance != nil {
		m.parent.RemoveMask(m.instance)
	}
}

type ShapeMask struct {
	*Mask
	*ShapeSprite

	featherRadius uint32
	featherMode   ti.FeatherMode
}

func (m *ShapeMask) RemoveMask(mask IMask) {
	//TODO implement me
	panic("implement me")
}

var _ IMask = (*ShapeMask)(nil)
var _ IElement = (*ShapeMask)(nil)

func NewShapeMask(parent IMaskParent, width, height uint32, cx, cy uint32) (*ShapeMask, error) {

	texture, err := ti.NewTiMask(parent.Renderer().Runtime(), width, height)
	if err != nil {
		return nil, err
	}

	return BuildMask(parent, texture, func(mask *Mask) (*ShapeMask, error) {
		shapeSprite, err := NewShapeSprite(SelfRelease(parent.Renderer()), width, height, cx, cy)
		if err != nil {
			mask.Release()
			return nil, err
		}

		return &ShapeMask{
			Mask:        mask,
			ShapeSprite: shapeSprite,

			featherRadius: 0,
			featherMode:   ti.FeatherModeLinear,
		}, nil
	})

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

func (m *ShapeMask) release() {
	if m.Mask != nil {
		m.Mask.Release()
		m.Mask = nil
	}
	if m.ShapeSprite != nil {
		m.ShapeSprite.Release()
		m.ShapeSprite = nil
	}
}
