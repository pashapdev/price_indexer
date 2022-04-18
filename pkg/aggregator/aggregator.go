package aggregator

import (
	"context"
	"log"
	"sync"

	"github.com/pashapdev/price_indexer/pkg/entities"
)

type priceStreamSubscriber interface {
	SubscribePriceStream() chan entities.TickerPriceWithErr
}

type Aggregator struct {
	ticker entities.Ticker
	out    chan entities.TickerWithDecimalPrice
	stop   chan struct{}
	wg     sync.WaitGroup
}

func New(ticker entities.Ticker) *Aggregator {
	out := make(chan entities.TickerWithDecimalPrice)
	stop := make(chan struct{})
	return &Aggregator{
		ticker: ticker,
		out:    out,
		stop:   stop,
	}
}

func (a *Aggregator) AggregatedChannel() <-chan entities.TickerWithDecimalPrice {
	return a.out
}

func (a *Aggregator) Add(ctx context.Context, subscribers ...priceStreamSubscriber) {
	for _, subscriber := range subscribers {
		a.wg.Add(1)
		go func(in <-chan entities.TickerPriceWithErr) {
			defer a.wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case <-a.stop:
					return

				case tickerPrice, open := <-in:
					if !open {
						return
					}
					if tickerPrice.Err == nil {
						price, err := tickerPrice.DecimalPrice()
						if err != nil {
							log.Println(err)
						} else {
							tickerWithDecimalPrice := entities.TickerWithDecimalPrice{Time: tickerPrice.Time, Price: price}
							a.out <- tickerWithDecimalPrice
						}
					} else {
						log.Println("error from subscriber: %w", tickerPrice.Err)
						return
					}
				}
			}
		}(subscriber.SubscribePriceStream())
	}
}

func (a *Aggregator) Close() {
	close(a.stop)
	a.wg.Wait()
	close(a.out)
}
