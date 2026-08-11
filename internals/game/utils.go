package game

import (
	"cmp"
	"slices"
)

func sort(Events []event) {
	slices.SortFunc(Events, func(A event, B event) int {
		return cmp.Compare(A.PriorityValue, B.PriorityValue)
	})
}
