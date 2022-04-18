package heaps

import "github.com/shopspring/decimal"

// An DecimalMinHeap is a min-heap of decimal.
type DecimalMinHeap []decimal.Decimal

func (h DecimalMinHeap) Len() int           { return len(h) }
func (h DecimalMinHeap) Less(i, j int) bool { return h[i].LessThan(h[j]) }
func (h DecimalMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *DecimalMinHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(decimal.Decimal))
}

func (h *DecimalMinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
