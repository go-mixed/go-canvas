package ctypes

import (
	"image/color"
	"math"

	"github.com/go-mixed/go-canvas/internel/misc"
)

type Point[T misc.Integer | misc.Float] struct {
	X, Y T
}

// Pt is shorthand for [Point]{X, Y}.
func Pt[T misc.Integer | misc.Float](X, Y T) Point[T] {
	return Point[T]{X, Y}
}

// Add returns the vector p+q.
func (p Point[T]) Add(q Point[T]) Point[T] {
	return Point[T]{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Point[T]) Sub(q Point[T]) Point[T] {
	return Point[T]{p.X - q.X, p.Y - q.Y}
}

// Mul returns the vector p*k.
func (p Point[T]) Mul(k T) Point[T] {
	return Point[T]{p.X * k, p.Y * k}
}

// Div returns the vector p/k.
func (p Point[T]) Div(k T) Point[T] {
	return Point[T]{p.X / k, p.Y / k}
}

// In reports whether p is in r.
func (p Point[T]) In(r Rectangle[T]) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

// Mod returns the point q in r such that p.X-q.X is a multiple of r's width
// and p.Y-q.Y is a multiple of r's height.
func (p Point[T]) Mod(r Rectangle[T]) Point[T] {
	w, h := r.Dx(), r.Dy()
	p = p.Sub(r.Min)
	p.X = T(math.Mod(float64(p.X), float64(w)))
	if p.X < 0 {
		p.X += w
	}
	p.Y = T(math.Mod(float64(p.Y), float64(h)))
	if p.Y < 0 {
		p.Y += h
	}
	return p.Add(r.Min)
}

// Eq reports whether p and q are equal.
func (p Point[T]) Eq(q Point[T]) bool {
	return p == q
}

type Rectangle[T misc.Integer | misc.Float] struct {
	Min, Max Point[T]
}

// RectXY is shorthand for [Rectangle]{Pt(x0, y0), [Pt](x1, y1)}. The returned
// rectangle has minimum and maximum coordinates swapped if necessary so that
// it is well-formed.
func RectXY[T misc.Integer | misc.Float](x0, y0, x1, y1 T) Rectangle[T] {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle[T]{Point[T]{x0, y0}, Point[T]{x1, y1}}
}

// RectWH is shorthand for [Rectangle]{Pt(x, y), Pt(x+w, y+h)}.
func RectWH[T misc.Integer | misc.Float](x, y, w, h T) Rectangle[T] {
	return RectXY(x, y, x+w, y+h)
}

// ToRect convert a Rectangle[S] to Rectangle[D]
func ToRect[D, S misc.Integer | misc.Float](srcRect Rectangle[S]) Rectangle[D] {
	return Rectangle[D]{
		Point[D]{D(srcRect.Min.X), D(srcRect.Min.Y)},
		Point[D]{D(srcRect.Max.X), D(srcRect.Max.Y)},
	}
}

func (r Rectangle[T]) X() T {
	return r.Min.X
}

func (r Rectangle[T]) Y() T {
	return r.Min.Y
}

// MoveTo returns a NEW Rectangle with moving the entire rectangle to x, y. (Min=x, y, Max=Width+x, Height+y)
func (r Rectangle[T]) MoveTo(x, y T) Rectangle[T] {
	deltaX, deltaY := x-r.X(), y-r.Y()
	return r.Add(Pt(deltaX, deltaY))
}

// Width the width of the Rectangle
func (r Rectangle[T]) Width() T {
	return r.Dx()
}

// Height the height of the Rectangle
func (r Rectangle[T]) Height() T {
	return r.Dy()
}

// Dx returns the width(x distance) of the Rectangle
func (r Rectangle[T]) Dx() T {
	return r.Max.X - r.Min.X
}

// Dy returns the height(y distance) of the Rectangle
func (r Rectangle[T]) Dy() T {
	return r.Max.Y - r.Min.Y
}

func (r Rectangle[T]) Size() Point[T] {
	return Point[T]{r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y}
}

// Resize the rectangle
func (r Rectangle[T]) Resize(width T, height T) Rectangle[T] {
	r.Max.X = r.Min.X + width
	r.Max.Y = r.Min.Y + height
	return r
}

// Add return A NEW Rectangle with the MIN/MAX point added to the rectangle.
func (r Rectangle[T]) Add(delta Point[T]) Rectangle[T] {
	return Rectangle[T]{
		Point[T]{r.Min.X + delta.X, r.Min.Y + delta.Y},
		Point[T]{r.Max.X + delta.X, r.Max.Y + delta.Y},
	}
}

// Sub return A NEW Rectangle with the MIN/MAX point subtracted from the rectangle.
func (r Rectangle[T]) Sub(delta Point[T]) Rectangle[T] {
	return Rectangle[T]{
		Point[T]{r.Min.X - delta.X, r.Min.Y - delta.Y},
		Point[T]{r.Max.X - delta.X, r.Max.Y - delta.Y},
	}
}

// Inset
func (r Rectangle[T]) Inset(n T) Rectangle[T] {
	if r.Dx() < 2*n {
		r.Min.X = (r.Min.X + r.Max.X) / 2
		r.Max.X = r.Min.X
	} else {
		r.Min.X += n
		r.Max.X -= n
	}
	if r.Dy() < 2*n {
		r.Min.Y = (r.Min.Y + r.Max.Y) / 2
		r.Max.Y = r.Min.Y
	} else {
		r.Min.Y += n
		r.Max.Y -= n
	}
	return r
}

func (r Rectangle[T]) Intersect(s Rectangle[T]) Rectangle[T] {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	// Letting r0 and s0 be the values of r and s at the time that the method
	// is called, this next line is equivalent to:
	//
	// if max(r0.Min.X, s0.Min.X) >= min(r0.Max.X, s0.Max.X) || likewiseForY { etc }
	if r.Empty() {
		return Rectangle[T]{}
	}
	return r
}

// Union returns the smallest rectangle that contains both r and s.
func (r Rectangle[T]) Union(s Rectangle[T]) Rectangle[T] {
	if r.Empty() {
		return s
	}
	if s.Empty() {
		return r
	}
	if r.Min.X > s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y > s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X < s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y < s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	return r
}

// Empty reports whether the rectangle contains no points.
func (r Rectangle[T]) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// Eq reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.
func (r Rectangle[T]) Eq(s Rectangle[T]) bool {
	return r == s || r.Empty() && s.Empty()
}

// Overlaps reports whether r and s have a non-empty intersection.
func (r Rectangle[T]) Overlaps(s Rectangle[T]) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

// In reports whether every point in r is in s.
func (r Rectangle[T]) In(s Rectangle[T]) bool {
	if r.Empty() {
		return true
	}
	// Note that r.Max is an exclusive bound for r, so that r.In(s)
	// does not require that r.Max.In(s).
	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y
}

// Canon returns the canonical version of r. The returned rectangle has minimum
// and maximum coordinates swapped if necessary so that it is well-formed.
func (r Rectangle[T]) Canon() Rectangle[T] {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

// At implements the [Image] interface.
func (r Rectangle[T]) At(x, y T) color.Color {
	if (Point[T]{x, y}).In(r) {
		return color.Opaque
	}
	return color.Transparent
}

// RGBA64At implements the [RGBA64Image] interface.
func (r Rectangle[T]) RGBA64At(x, y T) color.RGBA64 {
	if (Point[T]{x, y}).In(r) {
		return color.RGBA64{0xffff, 0xffff, 0xffff, 0xffff}
	}
	return color.RGBA64{}
}

// Bounds implements the [Image] interface.
func (r Rectangle[T]) Bounds() Rectangle[T] {
	return r
}

// ColorModel implements the [Image] interface.
func (r Rectangle[T]) ColorModel() color.Model {
	return color.Alpha16Model
}

func (r Rectangle[T]) Center() Point[T] {
	return Point[T]{r.Min.X + r.Dx()/2,
		r.Min.Y + r.Dy()/2}
}
