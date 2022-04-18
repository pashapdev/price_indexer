package subscriber

import (
	"sync"
	"time"

	"github.com/pashapdev/price_indexer/pkg/entities"

	"github.com/shopspring/decimal"
)

type subscriber struct {
	index  int
	price  decimal.Decimal
	ch     chan entities.TickerPriceWithErr
	ticker entities.Ticker
	stop   chan struct{}
	wg     sync.WaitGroup
}

func NewSubscribe(index int, ticker entities.Ticker, price decimal.Decimal) *subscriber {
	ch := make(chan entities.TickerPriceWithErr)
	stop := make(chan struct{})
	return &subscriber{
		ch:     ch,
		index:  index,
		ticker: ticker,
		price:  price,
		stop:   stop,
	}
}

func (s *subscriber) SubscribePriceStream() chan entities.TickerPriceWithErr {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.stop:
				return
			case <-time.After(1 * time.Minute):
				s.ch <- entities.TickerPriceWithErr{
					TickerPrice: entities.TickerPrice{
						Price:  s.price.String(),
						Ticker: s.ticker,
						Time:   time.Now(),
					},
				}
			}
		}
	}()

	return s.ch
}

func (s *subscriber) GetIndex() int {
	return s.index
}

func (s *subscriber) Close() {
	close(s.stop)
	s.wg.Wait()
	close(s.ch)
}
