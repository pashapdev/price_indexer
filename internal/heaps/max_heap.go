package heaps

import "github.com/shopspring/decimal"

// An DecimalMaxHeap is a max-heap of decimal.
type DecimalMaxHeap []decimal.Decimal

func (h DecimalMaxHeap) Len() int           { return len(h) }
func (h DecimalMaxHeap) Less(i, j int) bool { return h[j].LessThan(h[i]) }
func (h DecimalMaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *DecimalMaxHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(decimal.Decimal))
}

func (h *DecimalMaxHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
