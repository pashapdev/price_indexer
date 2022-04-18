package heaps_test

import (
	"container/heap"
	"testing"

	"github.com/pashapdev/price_indexer/internal/heaps"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestMaxDecimalHeap(t *testing.T) {
	h := &heaps.DecimalMaxHeap{
		decimal.NewFromInt(2),
		decimal.NewFromInt(1),
		decimal.NewFromInt(5)}
	heap.Init(h)

	heap.Push(h, decimal.NewFromInt(3))

	expectedValues := []decimal.Decimal{
		decimal.NewFromInt(5),
		decimal.NewFromInt(3),
		decimal.NewFromInt(2),
		decimal.NewFromInt(1)}

	require.Len(t, expectedValues, h.Len())

	for _, expectedValue := range expectedValues {
		actualValue := heap.Pop(h).(decimal.Decimal)
		require.True(t, actualValue.Equal(expectedValue))
	}
}
