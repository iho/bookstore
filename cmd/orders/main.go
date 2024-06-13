package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/iho/bookstore/internal/orders"
	"github.com/iho/bookstore/protos/gen/orders/v1/ordersv1connect"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	ordersService := orders.NewOrdersService(client)
	mux := http.NewServeMux()
	mux.Handle(ordersv1connect.NewOrdersServiceHandler(ordersService))

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

	fmt.Println("Starting server on :9999")

	if err := http.ListenAndServe(":9999", h2c.NewHandler(mux, &http2.Server{})); err != nil {
		log.Fatal(err)
	}
}
