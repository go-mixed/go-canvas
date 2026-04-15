package render

import (
	"sync"

	"github.com/go-mixed/go-canvas/ctypes"
	"github.com/go-mixed/go-canvas/internel/misc"
	"github.com/go-mixed/go-canvas/ti"
)

type animationItem struct {
	targetFn ti.TargetAttributeFn

	from            *ctypes.Attribute
	target          *ti.TargetAttribute
	startFrameIndex int
	durationFrames  int

	started bool
}

// spriteAnimator 维护单个精灵的串行动画队列。
// 每个动画段在开始帧时捕获起始属性，并按绝对帧号插值到目标属性。
type spriteAnimator struct {
	sprite IElement

	mutex   sync.Mutex
	queue   *misc.List[*animationItem]
	stopped bool
}

func newSpriteAnimator(sprite IElement) *spriteAnimator {
	return &spriteAnimator{
		sprite: sprite,
		queue:  misc.NewList[*animationItem](),
	}
}

func (a *spriteAnimator) setSprite(sprite IElement) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.sprite = sprite
}

// enqueue 追加动画段，startAtFrameIndex 与 durationFrames 均为帧单位。
func (a *spriteAnimator) enqueue(targetFn ti.TargetAttributeFn, startAtFrameIndex, durationFrames int) {
	if targetFn == nil || durationFrames <= 0 {
		return
	} else if startAtFrameIndex < 0 {
		startAtFrameIndex = 0
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()

	item := &animationItem{
		targetFn:        targetFn,
		startFrameIndex: startAtFrameIndex,
		durationFrames:  durationFrames,
	}
	a.queue.PushBack(item)
	a.stopped = false
}

// clear 清空所有待执行动画段。
func (a *spriteAnimator) clear() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.queue.Clear()
}

// hasPending 返回是否仍有待执行动画段。
func (a *spriteAnimator) hasPending() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return !a.stopped && a.queue.Len() > 0
}

// hasAnimationAt 返回给定绝对帧号下是否需要执行动画更新。
// 仅检查队头动画段（队列为串行语义）。
func (a *spriteAnimator) hasAnimationAt(frameIndex int) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.stopped || a.queue.Len() == 0 {
		return false
	}
	if frameIndex < 0 {
		frameIndex = 0
	}

	for it := a.queue.Front(); it != nil; it = it.Next() {
		item := it.Value
		if item == nil {
			continue
		}
		return frameIndex >= item.startFrameIndex
	}
	return false
}

// stop 停止动画推进；reset=true 时回滚到当前段起始属性。
func (a *spriteAnimator) stop(reset bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if reset {
		front := a.queue.Front()
		if front != nil && front.Value != nil && front.Value.started && front.Value.from != nil && front.Value.target != nil {
			applyModifiedFieldsLerp(a.sprite, front.Value.from, front.Value.target, 0.0)
		}
	}

	a.stopped = true
}

// tick 使用绝对帧号推进动画，返回是否仍有动画待执行。
func (a *spriteAnimator) tick(frameIndex int) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.stopped || a.queue.Len() == 0 {
		return false
	}
	if frameIndex < 0 {
		frameIndex = 0
	}

	for {
		front := a.queue.Front()
		if front == nil {
			return false
		}
		item := front.Value
		if item == nil {
			a.queue.PopFront()
			continue
		}

		if frameIndex < item.startFrameIndex {
			return true
		}
		if !item.started {
			item.from, item.target = item.targetFn(*a.sprite.Attribute())
			item.started = true
			if item.target == nil || item.target.Attribute == nil {
				a.queue.PopFront()
				continue
			}
		}

		elapsed := frameIndex - item.startFrameIndex
		if elapsed >= item.durationFrames {
			applyModifiedFieldsLerp(a.sprite, item.from, item.target, 1.0)
			a.queue.PopFront()
			continue
		}

		progress := misc.Clamp(float32(elapsed) / float32(item.durationFrames))
		eased := item.target.Easing(progress)
		applyModifiedFieldsLerp(a.sprite, item.from, item.target, eased)
		return true
	}
}

// applyModifiedFieldsLerp 按 modifiedFields 将 from->to 插值并应用到目标精灵。
func applyModifiedFieldsLerp(dst IAttribute, from *ctypes.Attribute, to *ti.TargetAttribute, t float32) {
	if from == nil || to == nil || to.Attribute == nil {
		return
	}

	hasWidth := to.IsModified(ti.ModifiedFieldWidth)
	hasHeight := to.IsModified(ti.ModifiedFieldHeight)
	if hasWidth || hasHeight {
		w := from.Width()
		h := from.Height()
		if hasWidth {
			w = misc.Lerp(from.Width(), to.Width(), t)
			if w <= 0 {
				w = 1
			}
		}
		if hasHeight {
			h = misc.Lerp(from.Height(), to.Height(), t)
			if h <= 0 {
				h = 1
			}
		}
		_ = dst.Resize(w, h)
	}

	if to.IsModified(ti.ModifiedFieldX) || to.IsModified(ti.ModifiedFieldY) {
		x := from.X()
		y := from.Y()
		if to.IsModified(ti.ModifiedFieldX) {
			x = misc.Lerp(from.X(), to.X(), t)
		}
		if to.IsModified(ti.ModifiedFieldY) {
			y = misc.Lerp(from.Y(), to.Y(), t)
		}
		dst.MoveTo(x, y)
	}

	if to.IsModified(ti.ModifiedFieldCx) {
		dst.SetCx(misc.Lerp(from.Cx(), to.Cx(), t))
	}
	if to.IsModified(ti.ModifiedFieldCy) {
		dst.SetCy(misc.Lerp(from.Cy(), to.Cy(), t))
	}

	if to.IsModified(ti.ModifiedFieldScaleX) || to.IsModified(ti.ModifiedFieldScaleY) {
		scaleX := from.ScaleX()
		scaleY := from.ScaleY()
		if to.IsModified(ti.ModifiedFieldScaleX) {
			scaleX = misc.Lerp(from.ScaleX(), to.ScaleX(), t)
		}
		if to.IsModified(ti.ModifiedFieldScaleY) {
			scaleY = misc.Lerp(from.ScaleY(), to.ScaleY(), t)
		}
		dst.SetScale(scaleX, scaleY)
	}

	if to.IsModified(ti.ModifiedFieldRotation) {
		dst.SetRotation(misc.Lerp(from.Rotation(), to.Rotation(), t))
	}
	if to.IsModified(ti.ModifiedFieldAlpha) {
		dst.SetAlpha(misc.Lerp(from.Alpha(), to.Alpha(), t))
	}

	if to.IsModified(ti.ModifiedFieldShape) && to.ShapeOpts != nil {
		if shapeSprite, ok := dst.(IShape); ok {
			opts := to.ShapeOpts
			tVal := misc.Lerp(opts.StartT, opts.EndT, t)
			shapeSprite.DrawShape(opts.ShapeType, tVal, opts.ShapeOptions)
		}
	}
}
