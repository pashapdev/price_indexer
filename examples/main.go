package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pashapdev/price_indexer/examples/internal/endpoints"
	"github.com/pashapdev/price_indexer/examples/internal/subscriber"
	"github.com/pashapdev/price_indexer/internal/config"
	"github.com/pashapdev/price_indexer/internal/saver"
	"github.com/pashapdev/price_indexer/pkg/aggregator"
	"github.com/pashapdev/price_indexer/pkg/entities"
	"github.com/pashapdev/price_indexer/pkg/indexer"
	"github.com/pashapdev/price_indexer/pkg/median"
)

func main() {
	log.Println("reading config...")
	cfg := config.New()
	if err := cfg.Validate(); err != nil {
		log.Println("failed to validate config:", err)
		os.Exit(1)
	}

	a := aggregator.New(entities.BTCUSDTicker)
	defer a.Close()

	pool := subscriber.NewPool()
	s := saver.New()
	idx := indexer.New(a, median.New(), s)

	go func(ctx context.Context) {
		idx.Run(ctx)
	}(context.Background())

	http.HandleFunc("/api/v1/subscribers", endpoints.MakeSubscriberListHandler(pool))
	http.HandleFunc("/api/v1/subscriber", endpoints.MakeSubscriberHandler(pool, a))
	http.HandleFunc("*", endpoints.MakeDefaultHandler())
	server := &http.Server{
		Addr: cfg.Address,
	}

	go func() {
		log.Println("starting server...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("failed to start server:", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	shutdownTimeout := 5 * time.Second //nolint:revive
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	log.Println("try to shutdown server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Println("failed to shutdown server:", err)
		os.Exit(1)
	}
}
