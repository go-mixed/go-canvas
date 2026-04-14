package misc

import (
	"math"
)

const Epsilon = 1e-6

// NumberEqual 判断两个数字是否相等，考虑浮点数精度问题
func NumberEqual[A, B Number](a A, b B, epsilon float64) bool {
	af := float64(a)
	bf := float64(b)
	return af == bf || math.Abs(af-bf) <= epsilon
}

// Clamp 返回值在 [0, 1] 范围内
func Clamp[V Float](v V) V {
	return max(min(v, 1.0), 0.0)
}

// Abs 返回 v 的绝对值
func Abs[V Float | Signed](x V) V {
	if x < 0 {
		return -x
	}
	return x
}

// Deg2Rad 将角度转换为弧度
func Deg2Rad[V Float](deg V) V {
	return deg * math.Pi / 180.0
}

// Rad2Deg 将弧度转换为角度
func Rad2Deg[V Float](rad V) V {
	return rad * 180.0 / math.Pi
}

// Floor 返回 向下取整数值。
// Floor returns floor(v) as T.
func Floor[T Number, P Number](v P) T {
	return T(math.Floor(float64(v)))
}

// Ceil 返回 向上取整数值
// Ceil returns ceil(v) as T.
func Ceil[T Number, P Number](v P) T {
	return T(math.Ceil(float64(v)))
}

// Lerp performs linear interpolation from a to b by t.
// For integer types, result is rounded to nearest integer.
// For float types, result keeps fractional precision.
func Lerp[T Number](a, b T, t float32) T {
	v := float64(a) + (float64(b)-float64(a))*float64(t)

	var zero T
	switch any(zero).(type) {
	case float32, float64:
		return T(v)
	default:
		return T(math.Round(v))
	}
}

const NaNInt = math.MaxInt

// IsNaNInt 是否未定义int
func IsNaNInt(v int) bool {
	return v == NaNInt
}
