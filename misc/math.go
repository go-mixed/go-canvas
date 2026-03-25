package misc

import (
	"math"

	"golang.org/x/exp/constraints"
)

const Epsilon = 1e-6

// NumberEqual 判断两个数字是否相等，考虑浮点数精度问题
func NumberEqual[A, B constraints.Integer | constraints.Float](a A, b B, epsilon float32) bool {
	af := float64(a)
	bf := float64(b)
	return af == bf || math.Abs(af-bf) <= float64(epsilon)
}

// Clamp 返回值在 [0, 1] 范围内
func Clamp[V ~float32 | ~float64](v V) V {
	return max(min(v, 1.0), 0.0)
}

// Abs 返回 v 的绝对值
func Abs[V ~float32 | ~float64 | ~int | ~int8 | ~int16 | ~int32 | ~int64](x V) V {
	if x < 0 {
		return -x
	}
	return x
}
