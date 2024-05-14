package util

import (
	"cmp"
	"slices"
)

func SortBy[S ~[]E, E any, P cmp.Ordered](s S, predicate func(e E) P) {
	slices.SortFunc(s, func(a E, b E) int {
		pa, pb := predicate(a), predicate(b)
		if pa < pb {
			return -1
		} else if pa > pb {
			return 1
		} else {
			return 0
		}
	})
}
