package collections

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Slices_Map(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		given []int
		want  []int
	}{
		"some_entries": {
			given: []int{1, 2, 3},
			want:  []int{1, 4, 9},
		},
		"empty": {
			given: []int{},
			want:  []int{},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// given
			// when
			got := Map(tc.given, func(i int) int { return i * i })
			// then
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("Diff(-want+got): %s\n", diff)
			}
		})
	}
}
