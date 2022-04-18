package endpoints

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/pashapdev/price_indexer/examples/internal/subscriber"
	"github.com/pashapdev/price_indexer/pkg/aggregator"
	"github.com/pashapdev/price_indexer/pkg/entities"

	"github.com/shopspring/decimal"
)

func MakeSubscriberListHandler(pool *subscriber.Pool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handler404(w)
			return
		}
		respondWithJSON(w, http.StatusOK, pool.List())
	}
}

func MakeSubscriberHandler(pool *subscriber.Pool, aggregator *aggregator.Aggregator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createSubscriber(pool, aggregator, w, r)
			return
		case http.MethodDelete:
			subscriberStop(pool, w, r)
			return
		default:
			handler404(w)
			return
		}
	}
}

type Subscriber struct {
	Index int             `json:"index"`
	Price decimal.Decimal `json:"price"`
}

func createSubscriber(pool *subscriber.Pool, aggregator *aggregator.Aggregator, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handler404(w)
		return
	}

	var req Subscriber
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("failed to decode body:", err)
		respondWithError(w, http.StatusBadRequest, "failed to decode body")
		return
	}

	newSubscriber := subscriber.NewSubscribe(req.Index, entities.BTCUSDTicker, req.Price)
	if err := pool.Add(newSubscriber); err != nil {
		log.Println("index should ne unique:", err)
		respondWithError(w, http.StatusBadRequest, "failed to decode body")
		return
	}
	aggregator.Add(context.TODO(), newSubscriber)

	respondWithJSON(w, http.StatusCreated, struct{}{})
}

func subscriberStop(pool *subscriber.Pool, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		handler404(w)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("id %s should be int: %s\n", r.URL.Query().Get("id"), err.Error())
		respondWithError(w, http.StatusBadRequest, "id should be int")
		return
	}

	pool.Delete(id)
	respondWithJSON(w, http.StatusOK, struct{}{})
}

func MakeDefaultHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler404(w)
	}
}
