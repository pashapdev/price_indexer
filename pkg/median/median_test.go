package median_test

import (
	"testing"

	"github.com/pashapdev/price_indexer/pkg/median"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestFairPricer(t *testing.T) {
	medianPricer := median.New()

	testCases := []struct {
		name          string
		addValue      decimal.Decimal
		expectedValue decimal.Decimal
	}{
		{
			name:          "add 1",
			addValue:      decimal.NewFromInt(1),
			expectedValue: decimal.NewFromInt(1),
		},
		{
			name:          "add 2",
			addValue:      decimal.NewFromInt(2),
			expectedValue: decimal.NewFromFloat32(1.5),
		},
		{
			name:          "add 6",
			addValue:      decimal.NewFromInt(6),
			expectedValue: decimal.NewFromInt(2),
		},
		{
			name:          "add 3",
			addValue:      decimal.NewFromInt(3),
			expectedValue: decimal.NewFromFloat32(2.5),
		},
		{
			name:          "add 100000",
			addValue:      decimal.NewFromInt(100000),
			expectedValue: decimal.NewFromInt(3),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			medianPricer.Add(testCase.addValue)
			actualValue, err := medianPricer.Calculate()
			require.NoError(t, err)
			require.True(t, testCase.expectedValue.Equal(actualValue))
		})
	}
	medianPricer.Clear()
	medianPricer.Add(decimal.NewFromInt(1))
	actualValue, err := medianPricer.Calculate()
	require.NoError(t, err)
	require.True(t, decimal.NewFromInt(1).Equal(actualValue))
}
