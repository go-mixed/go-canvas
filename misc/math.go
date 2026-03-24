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
