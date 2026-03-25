package effect

import (
	"github.com/go-mixed/go-canvas/misc"
)

// EasingFunction 缓动函数类型，接受 [0,1] 范围的输入，返回缓动后的值
// 返回值可以超出 [0,1] 范围以产生弹性效果
type EasingFunction func(float32) float32

// linear 线性缓动函数
func linear(x float32) float32 {
	return x
}

// ease 等同于 cubic-bezier(0.25, 0.1, 0.25, 1.0)
func ease(x float32) float32 {
	return cubicBezier(0.25, 0.1, 0.25, 1.0)(x)
}

// ease-in 等同于 cubic-bezier(0.42, 0, 1.0, 1.0)
func easeIn(x float32) float32 {
	return cubicBezier(0.42, 0, 1.0, 1.0)(x)
}

// ease-out 等同于 cubic-bezier(0, 0, 0.58, 1.0)
func easeOut(x float32) float32 {
	return cubicBezier(0, 0, 0.58, 1.0)(x)
}

// ease-in-out 等同于 cubic-bezier(0.42, 0, 0.58, 1.0)
func easeInOut(x float32) float32 {
	return cubicBezier(0.42, 0, 0.58, 1.0)(x)
}

// cubicBezier 创建三次贝塞尔曲线缓动函数
// 参数 p1x, p1y, p2x, p2y 是控制点坐标
// 参考坐标可以超出 [0,1] 范围产生弹性效果
//
// Example:
//
//	bounce := cubicBezier(0.68, -0.55, 0.265, 1.55)
//	value := bounce(0.5)
func cubicBezier(p1x, p1y, p2x, p2y float32) EasingFunction {
	return newCubicBezier(p1x, p1y, p2x, p2y)
}

// newCubicBezier 内部实现：三次贝塞尔曲线
//
// 使用二分法求解贝塞尔曲线，与浏览器实现保持一致
// 贝塞尔曲线公式: B(t) = (1-t)³P₀ + 3(1-t)²tP₁ + 3(1-t)t²P₂ + t³P₃
// 其中 P₀ = (0,0), P₃ = (1,1) 为固定端点
func newCubicBezier(p1x, p1y, p2x, p2y float32) EasingFunction {
	// 预计算常数以减少重复计算
	// 由于 p1x, p2x 用于 x 方程，p1y, p2y 用于 y 方程

	return func(x float32) float32 {
		// 边界情况优化
		if x <= 0 {
			return 0
		}
		if x >= 1 {
			return 1
		}

		// 求解对应的 t 值，然后计算 y 坐标
		t := solveTForX(x, p1x, p2x)
		return bezierY(t, p1y, p2y)
	}
}

// bezierX 计算贝塞尔曲线的 x 坐标
// B(t) = 3(1-t)²t*p1x + 3(1-t)t²*p2x + t³
func bezierX(t, p1x, p2x float32) float32 {
	t2 := t * t
	t3 := t2 * t
	mt := 1 - t
	mt2 := mt * mt
	return 3*mt2*t*p1x + 3*mt*t2*p2x + t3
}

// bezierY 计算贝塞尔曲线的 y 坐标
// B(t) = 3(1-t)²t*p1y + 3(1-t)t²*p2y + t³
func bezierY(t, p1y, p2y float32) float32 {
	t2 := t * t
	t3 := t2 * t
	mt := 1 - t
	mt2 := mt * mt
	return 3*mt2*t*p1y + 3*mt*t2*p2y + t3
}

// solveTForX 使用二分法求解 bezierX(t) = x 时的 t 值
//
// Args:
//
//	x: 目标 x 坐标 [0, 1]
//	p1x, p2x: 贝塞尔曲线控制点 x 坐标
//
// Returns:
//
//	对应的参数 t [0, 1]
func solveTForX(x, p1x, p2x float32) float32 {
	// 边界情况
	if x <= 0 {
		return 0
	}
	if x >= 1 {
		return 1
	}

	// 二分法求解
	var t0 float32 = 0.0
	var t1 float32 = 1.0
	var epsilon float32 = 1e-6
	maxIterations := 20

	for i := 0; i < maxIterations; i++ {
		tMid := (t0 + t1) / 2
		xMid := bezierX(tMid, p1x, p2x)

		if misc.Abs(xMid-x) < epsilon {
			return tMid
		}

		if xMid < x {
			t0 = tMid
		} else {
			t1 = tMid
		}
	}

	return (t0 + t1) / 2
}

// 预定义常用缓动函数名称映射
var easingFunctions = map[string]EasingFunction{
	"linear":      linear,
	"ease":        ease,
	"ease-in":     easeIn,
	"ease-out":    easeOut,
	"ease-in-out": easeInOut,
}

// GetEasingFunction 根据名称获取缓动函数
func GetEasingFunction(name string) EasingFunction {
	fn, ok := easingFunctions[name]
	if !ok {
		fn = DefaultEasingFunction
	}
	return fn
}

var DefaultEasingFunction = ease
