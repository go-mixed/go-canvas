package render

import (
	"iter"
	"slices"
	"sync"
)

// Stage 舞台
type Stage struct {
	children []ISprite

	screen Screen
	mutex  sync.Mutex
}

type Screen struct {
	Sprite
}

func (s *Stage) Add(sprite ISprite) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.children = append(s.children, sprite)
}

func (s *Stage) Remove(sprite ISprite) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.children = slices.DeleteFunc(s.children, func(child ISprite) bool {
		return child == sprite
	})
}

func (s *Stage) Children() iter.Seq[ISprite] {
	return func(yield func(ISprite) bool) {
		for _, child := range s.children {
			if !yield(child) {
				return
			}
		}
	}
}

func (s *Stage) Render() {
	s.screen.FillColor(0xFFFFFFFF)
}
