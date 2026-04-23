package render

import (
	"image/color"

	"github.com/go-mixed/go-canvas/ctypes"
	"github.com/go-mixed/go-canvas/internel/misc"
	"github.com/go-mixed/go-canvas/ti"
	"github.com/go-mixed/go-taichi/taichi"
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
func NewContainer(parent IParent, attribute *ctypes.Attribute) (*Container, error) {
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
		}) >= 0 {
			return
		}

		c.children.PushBack(sprite)
	}, func() ctypes.DirtyMode { return ctypes.DirtyModeChildren | ctypes.DirtyModeComposite })
}

func (c *Container) RemoveChild(sprite ISprite) {
	c.LockForUpdate(func() {

		c.children.RemoveAll(func(child ISprite) bool {
			return child == sprite
		})

		// 递归释放
		sprite.Release()
	}, func() ctypes.DirtyMode { return ctypes.DirtyModeChildren | ctypes.DirtyModeComposite })
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

func (c *Container) ClientRect() ctypes.Rectangle[int] {
	rect := c.Sprite.ClientRect()
	for _, child := range c.children.Range() {
		bbox := child.ClientRect()
		rect = rect.Union(bbox)
	}
	return rect
}

// HasAnimationAt returns true when container or any child has animation
// to be evaluated at the given absolute frame.
func (c *Container) HasAnimationAt(frameIndex int) bool {
	if c.Sprite.HasAnimationAt(frameIndex) {
		return true
	}
	for _, child := range c.children.Range() {
		if child.HasAnimationAt(frameIndex) {
			return true
		}
	}
	return false
}

func (c *Container) Render(frameIndex int) error {
	defer func() {
		c.SetDirty(ctypes.DirtyModeNone)
	}()

	c.TickAnimation(frameIndex)

	c.mutex.Lock()
	children := c.children
	c.mutex.Unlock()

	for _, child := range children.Range() {
		child.TickAnimation(frameIndex)
	}

	if !c.IsDirty() {
		return nil
	}

	// 置空为透明
	c.Renderer().Module().AsyncFillColor(c.texture, color.Transparent)

	w, h := c.attribute.Width(), c.attribute.Height()
	relativeContainerRect := ctypes.RectWH(0, 0, w, h)

	for _, child := range children.Range() {
		// 渲染子级
		if err := child.Render(frameIndex); err != nil {
			return err
		}
		childTexture := child.Texture()

		// 得到子项（旋转、缩放、平移）之后真实坐标
		bbox := child.ClientRect()
		// 添加滚动条
		bbox = bbox.Add(ctypes.Pt(c.childOffsetX, c.childOffsetY))
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
			Width:    float32(childAttribute.ClientWidth()),
			Height:   float32(childAttribute.ClientHeight()),
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
			c.parent.Renderer().Module().AsyncRenderLayerWithMask(
				childTexture,
				mask.Texture(),
				c.texture,
				options,
			)
		} else {
			c.parent.Renderer().Module().AsyncRenderLayerNoMask(
				childTexture,
				c.texture,
				options,
			)
		}
	}
	return nil
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

func (c *Container) AddGarbageTexture(texture *taichi.NdArray) {
	c.Sprite.AddGarbageTexture(texture)
}

func (c *Container) ReleaseGarbageTextures() {
	c.Sprite.ReleaseGarbageTextures()

	for _, child := range c.children.Range() {
		child.ReleaseGarbageTextures()
	}
}

func (c *Container) ScrollTop(y int) {
	c.LockForUpdate(func() {
		if y < 0 {
			y = 0
		}

		cr := c.ClientRect()
		ch := cr.Height()
		if y > ch {
			y = ch
		}

		c.childOffsetY = -y
	}, func() ctypes.DirtyMode {
		if y != c.childOffsetY {
			return ctypes.DirtyModeComposite
		}
		return ctypes.DirtyModeNone
	})
}

func (c *Container) ScrollLeft(x int) {
	c.LockForUpdate(func() {
		if x < 0 {
			x = 0
		}
		cr := c.ClientRect()
		ch := cr.Height()
		if x > ch {
			x = ch
		}

		c.childOffsetX = -x
	}, func() ctypes.DirtyMode {
		if x != c.childOffsetX {
			return ctypes.DirtyModeComposite
		}
		return ctypes.DirtyModeNone
	})
}
