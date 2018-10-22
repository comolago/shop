package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"

	"github.com/comolago/shop/inventory/domain"
	"github.com/comolago/shop/inventory/infrastructure"
	"github.com/comolago/shop/inventory/usecases"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Transports expose the service to the network. In this first example we utilize JSON over HTTP.
func main() {
	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	// Logging domain.

	//declare metrics
	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "shop",
		Subsystem: "inventory",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "shop",
		Subsystem: "inventory",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var svc domain.InventoryInt
	svc = domain.Inventory{}
	svc = infrastructure.LoggingMiddleware(logger)(svc)
	svc = infrastructure.Metrics(requestCount, requestLatency)(svc)

	//rlbucket := ratelimit.NewBucket(1*time.Second, 5)

	addItemHandler := httptransport.NewServer(
		usecases.MakeAddItemEndpoint(svc),
		usecases.DecodeAddItemRequest,
		usecases.EncodeResponse,
	)

	countHandler := httptransport.NewServer(
		usecases.MakeCountEndpoint(svc),
		usecases.DecodeCountRequest,
		usecases.EncodeResponse,
	)

	http.Handle("/items/add", addItemHandler)
	http.Handle("/count", countHandler)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
	//log.Fatal(http.ListenAndServe(":8080", nil))
}
