package indexer

import (
	"context"
	"time"

	"github.com/pashapdev/price_indexer/pkg/entities"

	"github.com/shopspring/decimal"
)

type (
	source interface {
		AggregatedChannel() <-chan entities.TickerWithDecimalPrice
	}

	calculator interface {
		Add(price decimal.Decimal)
		Calculate() (decimal.Decimal, error)
		Clear()
	}

	saver interface {
		Save(currentMinute time.Time, data string)
	}
)

type indexer struct {
	source     source
	calculator calculator
	saver      saver
	timeWindow *timeWindow
}

func New(source source, cacalculator calculator, saver saver) *indexer {
	return &indexer{
		source:     source,
		calculator: cacalculator,
		saver:      saver,
		timeWindow: &timeWindow{},
	}
}

func (idx *indexer) Run(ctx context.Context) {
	diff := 60 - time.Now().UTC().Second()
	ticker := time.NewTicker(time.Duration(diff) * time.Second)
	idx.timeWindow.Next(time.Now().Add(-time.Duration(diff) * time.Second))
	idx.calculator.Clear()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-idx.source.AggregatedChannel():
			if idx.timeWindow.Include(msg.Time) {
				idx.calculator.Add(msg.Price)
			}
		case currentTime := <-ticker.C:
			idx.save()
			idx.calculator.Clear()
			idx.timeWindow.Next(currentTime)
			ticker.Reset(time.Minute)
		}
	}
}

func (idx *indexer) save() {
	data := ""
	price, err := idx.calculator.Calculate()
	if err != nil {
		data = err.Error()
	} else {
		data = price.String()
	}

	idx.saver.Save(idx.timeWindow.min, data)
}
