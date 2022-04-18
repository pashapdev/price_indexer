package entities

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type Ticker string

const (
	BTCUSDTicker Ticker = "BTC_USD"
)

type TickerPrice struct {
	Ticker Ticker
	Time   time.Time
	Price  string // decimal value. example: "0", "10", "12.2", "13.2345122"
}

func (t *TickerPrice) DecimalPrice() (decimal.Decimal, error) {
	price, err := decimal.NewFromString(t.Price)
	if err != nil {
		return decimal.Zero, fmt.Errorf("price %s is not valid: %w", t.Price, err)
	}
	return price, nil
}

type TickerPriceWithErr struct {
	TickerPrice
	Err error
}

type TickerWithDecimalPrice struct {
	Time  time.Time
	Price decimal.Decimal
}
