package render

import (
	"github.com/go-mixed/go-canvas/ctypes"
	"github.com/go-mixed/go-canvas/ti"
	"github.com/go-mixed/go-taichi/taichi"
)

type Mask struct {
	rect      ctypes.Rectangle[float32]
	texture   *ctypes.TiMask
	distField *ctypes.TiGrid // 距离场
	parent    IMaskParent
	// 真实的实例
	instance IMask
	isDirty  bool

	featherRadius uint32
	featherMode   ctypes.FeatherMode
}

var _ IMask = (*Mask)(nil)

func BuildMask[T IMask](parent IMaskParent, texture *ctypes.TiMask, instanceCreator func(*Mask) (T, error)) (T, error) {
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

		featherRadius: 0,
		featherMode:   ctypes.FeatherModeLinear,
	}

	instance, err := instanceCreator(m)
	if err != nil {
		return nilT, err
	}
	parent.AddMask(instance)
	m.instance = instance
	return instance, nil
}

func (m *Mask) SetInstance(instance IMask) {
	m.instance = instance
}

func (m *Mask) IsDirty() bool {
	return m.isDirty
}

func (m *Mask) SetDirty(val bool) {
	m.isDirty = val
}

// FillWithTexture 将纹理填充到 Mask
func (m *Mask) FillWithTexture(texture *ctypes.TiImage) {
	//  将图像转换为遮罩（提取 alpha 通道）
	m.parent.Renderer().Module().AsyncImageToMask(texture, m.texture)
	m.isDirty = true
}

func (m *Mask) SetFeather(featherRadius uint32, featherMode ctypes.FeatherMode) {
	m.featherRadius = featherRadius
	m.featherMode = featherMode
	m.isDirty = true
}

func (m *Mask) Render(frameIndex int) {
	defer func() {
		m.SetDirty(false)
	}()

	if !m.IsDirty() {
		return
	}

	if m.featherRadius > 0 {
		// 计算距离场（使用复用的 distField）
		m.parent.Renderer().Module().AsyncComputeDistanceField(m.texture, m.distField)
		// 应用羽化
		m.parent.Renderer().Module().AsyncComputeFeather(m.distField, m.texture, float32(m.featherRadius), m.featherMode)
	}
}

func (m *Mask) RemoveFromParent() {
	if m.instance != nil {
		m.parent.RemoveMask(m.instance)
	}
}

func (m *Mask) Texture() *taichi.NdArray {
	return m.texture
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

type ShapeMask struct {
	mask *Mask
	*ShapeSprite
}

var _ IMask = (*ShapeMask)(nil)
var _ IElement = (*ShapeMask)(nil)
var _ IShape = (*ShapeMask)(nil)

func NewShapeMask(parent IMaskParent, attribute *ctypes.Attribute) (*ShapeMask, error) {

	texture, err := ti.NewTiMask(parent.Renderer().Runtime(), uint32(attribute.Width()), uint32(attribute.Height()))
	if err != nil {
		return nil, err
	}

	return BuildMask(parent, texture, func(mask *Mask) (*ShapeMask, error) {
		shapeSprite, err := NewShapeSprite(SelfRelease(parent.Renderer()), attribute)
		if err != nil {
			mask.Release()
			return nil, err
		}

		m := &ShapeMask{
			mask:        mask,
			ShapeSprite: shapeSprite,
		}
		// 重新设置示例
		shapeSprite.SetInstance(m)
		return m, nil
	})
}

// IsDirty 覆盖父类 IsDirty 方法
func (m *ShapeMask) IsDirty() bool {
	return m.mask.IsDirty() || m.ShapeSprite.IsDirty()
}

// Render 覆盖父类 Render 方法
func (m *ShapeMask) Render(frameIndex int) {
	if !m.IsDirty() {
		return
	}

	m.ShapeSprite.Render(frameIndex)

	// FillWithTexture 会触发mask的dirty，下面的Render才会真正执行
	m.mask.FillWithTexture(m.ShapeSprite.Texture())
	m.mask.Render(frameIndex)
}

func (m *ShapeMask) SetFeather(radius uint32, featherMode ctypes.FeatherMode) {
	m.mask.SetFeather(radius, featherMode)
}

func (m *ShapeMask) FillWithTexture(texture *ctypes.TiImage) {
	m.mask.FillWithTexture(texture)
}

// Texture 覆盖父类 Texture 方法
func (m *ShapeMask) Texture() *ctypes.TiMask {
	return m.mask.Texture()
}

// Release 覆盖父类 Release 方法
func (m *ShapeMask) Release() {
	if m.mask != nil {
		m.mask.Release()
		m.mask = nil
	}
	if m.ShapeSprite != nil {
		m.ShapeSprite.Release()
		m.ShapeSprite = nil
	}
}
