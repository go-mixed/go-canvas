package render

import (
	"image/color"
	"slices"

	"github.com/go-mixed/go-canvas/ti"
)

type Container struct {
	*Sprite

	children []ISprite
}

var _ IContainer = (*Container)(nil)

// NewContainer 创建容器，只能添加子精灵
// 容器中的texture只将子精灵、容器渲染在上面
func NewContainer(renderer *Renderer, width, height uint32) (IContainer, error) {
	texture, err := ti.NewTiImage(renderer.Runtime(), width, height)
	if err != nil {
		return nil, err
	}

	return &Container{
		Sprite: NewSprite(renderer, texture),
	}, nil
}

func (c *Container) Add(sprite ISprite) {
	c.lockForUpdate(func() {
		c.children = append(c.children, sprite)
	}, func() bool { return true })
}

func (c *Container) Remove(sprite ISprite) {
	c.lockForUpdate(func() {
		c.children = slices.DeleteFunc(c.children, func(child ISprite) bool {
			return child == sprite
		})
	}, func() bool { return true })
}

func (c *Container) Children() []ISprite {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.children
}

func (c *Container) IsDirty() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, child := range c.children {
		if child.IsDirty() {
			return true
		}
	}

	// 如果只是修改容器的x,y等，无需重新渲染当前容器的texture，
	// 因为当前Render()是将子精灵渲染到c.texture上，而修改x,y等参数之后，会在父节点中处理
	return false
}

func (c *Container) ClientRect() ti.Rectangle[float32] {
	rect := c.Sprite.ClientRect()
	for _, child := range c.children {
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
	c.renderer.Module().FillTexture(c.texture, color.Transparent)

	c.mutex.Lock()
	children := c.children
	c.mutex.Unlock()

	w, h := c.Width(), c.Height()

	for _, child := range children {
		// 渲染子级
		child.Render()
		childTexture := child.Texture()

		bbox := child.ClippedRect(w, h)

		mask := child.Mask()

		options := ti.RenderLayerOptions{
			X:        child.X(),
			Y:        child.Y(),
			Width:    child.Width(),
			Height:   child.Height(),
			Cx:       child.Cx(),
			Cy:       child.Cy(),
			ScaleX:   child.ScaleX(),
			ScaleY:   child.ScaleY(),
			Rotation: child.Rotation(),
			Alpha:    child.Alpha(),
			MinX:     int32(bbox.Min.X),
			MaxX:     int32(bbox.Max.X),
			MinY:     int32(bbox.Min.Y),
			MaxY:     int32(bbox.Max.Y),
		}

		if mask != nil {
			c.renderer.Module().RenderLayerWithMask(
				childTexture,
				mask.Texture(),
				c.texture,
				options,
			)
		} else {
			c.renderer.Module().RenderLayerNoMask(
				childTexture,
				c.texture,
				options,
			)
		}
	}
}
