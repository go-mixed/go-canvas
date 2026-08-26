package render

import (
	"sync"

	"github.com/go-mixed/go-canvas/internel/misc"
)

// 生成动画方法的函数，会在调用时记录sprite的attribute
type AnimateFn func(sprite IElement) TickingFn

// 执行动画的函数，progress=[0, 1]
type TickingFn func(progress float32)

type animationItem struct {
	animateFn AnimateFn
	tickFn    TickingFn

	//from            *ctypes.Attribute
	//target          *ti.TargetAttribute
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
func (a *spriteAnimator) enqueue(animateFn AnimateFn, startAtFrameIndex, durationFrames int) {
	if animateFn == nil || durationFrames <= 0 {
		return
	} else if startAtFrameIndex < 0 {
		startAtFrameIndex = 0
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()

	item := &animationItem{
		animateFn:       animateFn,
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
		if front != nil && front.Value != nil && front.Value.started && front.Value.tickFn != nil {
			front.Value.tickFn(0.0)
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
			item.tickFn = item.animateFn(a.sprite)
			item.started = true
			if item.tickFn == nil {
				a.queue.PopFront()
				continue
			}
		}

		elapsed := frameIndex - item.startFrameIndex
		if elapsed >= item.durationFrames {
			item.tickFn(1.0)
			a.queue.PopFront()
			continue
		}

		progress := misc.Clamp(float32(elapsed) / float32(item.durationFrames))
		item.tickFn(progress)
		return true
	}
}
