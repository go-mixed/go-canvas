package misc

import (
	"math"
)

const Epsilon = 1e-6

// NumberEqual 判断两个数字是否相等，考虑浮点数精度问题
func NumberEqual[A, B Integer | Float](a A, b B, epsilon float32) bool {
	af := float64(a)
	bf := float64(b)
	return af == bf || math.Abs(af-bf) <= float64(epsilon)
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

// FloorToInt 返回 float64 的向下取整整数值。
// FloorToInt returns floor(v) as int.
func FloorToInt(v float64) int {
	return int(math.Floor(v))
}
