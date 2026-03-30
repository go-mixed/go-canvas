package render

import (
	"image/color"

	"github.com/go-mixed/go-canvas/misc"
	"github.com/go-mixed/go-canvas/ti"
)

type Container struct {
	*Sprite

	children *misc.List[ISprite]

	childOffsetX int
	childOffsetY int
}

var _ IContainer = (*Container)(nil)
var _ IParent = (*Container)(nil)
var _ IMaskParent = (*Container)(nil)

// NewContainer 创建容器，只能添加子精灵
// 容器中的texture为空白，当Render时，会将子精灵、容器渲染在上面
func NewContainer(parent IParent, attribute *ti.Attribute) (*Container, error) {
	texture, err := ti.NewTiImage(parent.Renderer().Runtime(), uint32(attribute.Width()), uint32(attribute.Height()))
	if err != nil {
		return nil, err
	}

	return BuildSprite(parent, attribute, texture, func(s *Sprite) (*Container, error) {
		return &Container{
			Sprite:   s,
			children: misc.NewList[ISprite](),
		}, nil
	})
}

func (c *Container) AddChild(sprite ISprite) {
	c.LockForUpdate(func() {

		if c.children.Index(func(item ISprite) bool {
			return item == sprite
		}) > 0 {
			return
		}

		c.children.PushBack(sprite)
	}, func() bool { return true })
}

func (c *Container) RemoveChild(sprite ISprite) {
	c.LockForUpdate(func() {

		c.children.RemoveAll(func(child ISprite) bool {
			return child == sprite
		})

		// 递归释放
		sprite.Release()
	}, func() bool { return true })
}

func (c *Container) RemoveFromParent() {
	if c.parent != nil {
		c.parent.RemoveChild(c)
	}
}

func (c *Container) Children() *misc.List[ISprite] {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.children
}

func (c *Container) IsDirty() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, child := range c.children.Range() {
		if child.IsDirty() {
			return true
		}
	}

	return c.Sprite.IsDirty()
}

func (c *Container) ClientRect() ti.Rectangle[float32] {
	rect := c.Sprite.ClientRect()
	for _, child := range c.children.Range() {
		bbox := child.ClientRect()
		rect = rect.Union(bbox)
	}
	return rect
}

func (c *Container) Render() {
	defer func() {
		c.SetDirty(false)
	}()

	if !c.IsDirty() {
		return
	}

	// 置空为透明
	c.parent.Renderer().Module().FillColor(c.texture, color.Transparent)

	c.mutex.Lock()
	children := c.children
	c.mutex.Unlock()

	w, h := c.attribute.Width(), c.attribute.Height()
	relativeContainerRect := ti.RectWH(0, 0, float32(w), float32(h))

	for _, child := range children.Range() {
		// 渲染子级
		child.Render()
		childTexture := child.Texture()

		// 得到子项（旋转、缩放、平移）之后真实坐标
		bbox := child.ClientRect()
		// 添加滚动条
		bbox = bbox.Add(ti.Pt[float32](float32(c.childOffsetX), float32(c.childOffsetY)))
		// 计算和当前容器的交集
		bbox = bbox.Intersect(relativeContainerRect)

		if bbox.Dx() == 0 || bbox.Dy() == 0 {
			continue
		}

		childAttribute := child.Attribute()

		// 暂时只支持第一个mask
		masks := child.Masks()
		var mask IMask
		if masks.Len() > 0 {
			mask = masks.At(0)
		}

		options := ti.RenderLayerOptions{
			X:        float32(childAttribute.X() + c.childOffsetX),
			Y:        float32(childAttribute.Y() + c.childOffsetY),
			Width:    float32(childAttribute.Width()),
			Height:   float32(childAttribute.Height()),
			Cx:       float32(childAttribute.Cx()),
			Cy:       float32(childAttribute.Cy()),
			ScaleX:   childAttribute.ScaleX(),
			ScaleY:   childAttribute.ScaleY(),
			Rotation: childAttribute.Rotation(),
			Alpha:    childAttribute.Alpha(),
			MinX:     int32(bbox.Min.X),
			MaxX:     int32(bbox.Max.X),
			MinY:     int32(bbox.Min.Y),
			MaxY:     int32(bbox.Max.Y),
		}

		if mask != nil {
			c.parent.Renderer().Module().RenderLayerWithMask(
				childTexture,
				mask.Texture(),
				c.texture,
				options,
			)
		} else {
			c.parent.Renderer().Module().RenderLayerNoMask(
				childTexture,
				c.texture,
				options,
			)
		}
	}
}

func (c *Container) Release() {
	if c.Sprite != nil {
		c.Sprite.Release()
	}

	for _, child := range c.children.Range() {
		child.Release()
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Sprite = nil
	c.children.Clear()
}

func (c *Container) ScrollTop(y int) {
	c.LockForUpdate(func() {
		if y < 0 {
			y = 0
		}

		cr := c.ClientRect()
		ch := int(cr.Height())
		if y > ch {
			y = ch
		}

		c.childOffsetY = -y
	}, func() bool {
		return y != c.childOffsetY
	})
}

func (c *Container) ScrollLeft(x int) {
	c.LockForUpdate(func() {
		if x < 0 {
			x = 0
		}
		cr := c.ClientRect()
		ch := int(cr.Height())
		if x > ch {
			x = ch
		}

		c.childOffsetX = -x
	}, func() bool {
		return x != c.childOffsetX
	})
}
