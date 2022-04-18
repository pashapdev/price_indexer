package subscriber

import (
	"fmt"
	"sync"

	"github.com/shopspring/decimal"
)

type (
	Subscriber struct {
		Index int             `json:"index"`
		Price decimal.Decimal `json:"price"`
	}
)

type Pool struct {
	subscribers map[int]*subscriber
	mx          sync.Mutex
}

func NewPool() *Pool {
	return &Pool{
		subscribers: make(map[int]*subscriber),
	}
}

func (p *Pool) List() []Subscriber {
	p.mx.Lock()
	defer p.mx.Unlock()
	res := make([]Subscriber, len(p.subscribers))
	i := 0
	for _, subscriber := range p.subscribers {
		res[i].Index = subscriber.index
		res[i].Price = subscriber.price
		i++
	}

	return res
}

func (p *Pool) Add(subscriber *subscriber) error {
	p.mx.Lock()
	defer p.mx.Unlock()
	if _, exist := p.subscribers[subscriber.GetIndex()]; exist {
		return fmt.Errorf("subscriber with id %d already exists", subscriber.GetIndex())
	}
	p.subscribers[subscriber.GetIndex()] = subscriber
	return nil
}

func (p *Pool) Delete(subscriberID int) {
	p.mx.Lock()
	defer p.mx.Unlock()
	p.subscribers[subscriberID].Close()
	delete(p.subscribers, subscriberID)
}
