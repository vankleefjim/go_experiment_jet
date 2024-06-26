package thelper

import (
	"time"

	"github.com/google/go-cmp/cmp"
)

func FixMonotonicTimePtr() cmp.Option {
	return cmp.Comparer(func(x, y *time.Time) bool {
		// TODO make sure this code is called somehow
		panic("AHH")
		if (x == nil) != (y == nil) {
			return false
		}
		if x == nil {
			return true
		}

		return x.Equal(*y)
	})
}
