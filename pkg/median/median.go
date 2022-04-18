// Step 1: Add next item to one of the heaps
//
// if next item is smaller than maxHeap root add it to maxHeap,
// else add it to minHeap
//
// Step 2: Balance the heaps (after this step heaps will be either balanced or
// one of them will contain 1 more item)
//
// if number of elements in one of the heaps is greater than the other by
// more than 1, remove the root element from the one containing more elements and
// add to the other one

package median

import (
	"container/heap"
	"errors"

	"github.com/pashapdev/price_indexer/internal/heaps"

	"github.com/shopspring/decimal"
)

type medianPricer struct {
	minHeap *heaps.DecimalMinHeap
	maxHeap *heaps.DecimalMaxHeap
}

func New() *medianPricer {
	medianPricer := &medianPricer{}
	medianPricer.Clear()
	return medianPricer
}

func (p *medianPricer) Add(price decimal.Decimal) {
	if p.maxHeap.Len() == 0 {
		heap.Push(p.maxHeap, price)
		return
	}

	if price.LessThan((*p.maxHeap)[0]) {
		heap.Push(p.maxHeap, price)
	} else {
		heap.Push(p.minHeap, price)
	}

	for p.minHeap.Len()-p.maxHeap.Len() > 1 {
		root := heap.Pop(p.minHeap)
		heap.Push(p.maxHeap, root)
	}

	for p.maxHeap.Len()-p.minHeap.Len() > 1 {
		root := heap.Pop(p.maxHeap)
		heap.Push(p.minHeap, root)
	}
}

func (p *medianPricer) Calculate() (decimal.Decimal, error) {
	if p.minHeap.Len() == 0 && p.maxHeap.Len() == 0 {
		return decimal.Zero, errors.New("no data to make decision")
	}

	if p.minHeap.Len() == p.maxHeap.Len() {
		return (*p.minHeap)[0].Add((*p.maxHeap)[0]).Mul(decimal.NewFromFloat32(0.5)), nil
	}

	if p.minHeap.Len() > p.maxHeap.Len() {
		return (*p.minHeap)[0], nil
	} else {
		return (*p.maxHeap)[0], nil
	}
}

func (p *medianPricer) Clear() {
	p.minHeap = &heaps.DecimalMinHeap{}
	p.maxHeap = &heaps.DecimalMaxHeap{}
	heap.Init(p.minHeap)
	heap.Init(p.maxHeap)
}
