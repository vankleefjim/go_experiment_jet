package thelper

import (
	"time"

	"github.com/google/go-cmp/cmp"
)

func FixMonotonicTimePtr() cmp.Option {
	return cmp.Comparer(func(x, y *time.Time) bool {
		if (x == nil) || (y == nil) {
			return x == y
		}

		return TimeEqual(*x, *y)
	})
}

func TimeEqual(x, y time.Time) bool {
	return x.Truncate(time.Millisecond).Equal(y.Truncate(time.Millisecond))
}
