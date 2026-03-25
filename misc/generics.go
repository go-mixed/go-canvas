package misc

import (
	"time"

	"golang.org/x/exp/constraints"
)

// Scalar is a type parameter constraint for functions accepting basic types.
//
// It represents the supported basic types this package can cast to.
type Scalar interface {
	string | bool | constraints.Float | constraints.Integer | time.Time | time.Duration
}
