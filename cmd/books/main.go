package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/iho/bookstore/internal/books"
	"github.com/iho/bookstore/protos/gen/books/v1/booksv1connect"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	redis "github.com/redis/go-redis/v9"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	booksService := books.NewBooksService(rdb)

	mux := http.NewServeMux()
	mux.Handle(booksv1connect.NewBooksServiceHandler(booksService))
	fmt.Println("Starting server on :9090")

	reg := prometheus.NewRegistry()

	// Add Go module build info.
	reg.MustRegister(collectors.NewBuildInfoCollector())
	reg.MustRegister(collectors.NewGoCollector(
		collectors.WithGoCollectorRuntimeMetrics(collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")}),
	))

	// Expose the registered metrics via HTTP.
	mux.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))

	if err := http.ListenAndServe(":9090", h2c.NewHandler(mux, &http2.Server{})); err != nil {
		log.Fatal(err)
	}
}
