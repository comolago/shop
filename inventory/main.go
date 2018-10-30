package main
	//"github.com/prometheus/client_golang/prometheus/promhttp"
//	httptransport "github.com/go-kit/kit/transport/http"
 //  "github.com/gorilla/mux"

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"

	"github.com/comolago/shop/inventory/domain"
	"github.com/comolago/shop/inventory/infrastructure"
	"github.com/comolago/shop/inventory/usecases"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
   "fmt"
"os/signal"
"syscall"
"time"
ratelimitkit "github.com/go-kit/kit/ratelimit"

"golang.org/x/time/rate"

"golang.org/x/net/context"
)

func main() {


ctx := context.Background()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
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

	var svc domain.InventoryHandler
	svc = domain.Inventory{nil,new(infrastructure.PostgresqlDb)}
        svc.Open()
        svc.GetItemById(1)
	svc = infrastructure.LoggingMiddleware(logger)(svc)
	svc = infrastructure.Metrics(requestCount, requestLatency)(svc)

        limit := rate.NewLimiter(rate.Every(35*time.Millisecond), 100)

        e := usecases.MakeGetItemEndpoint(svc)
	e = ratelimitkit.NewErroringLimiter(limit)(e)
	endpoint := usecases.Endpoints{
		GetItemEndpoint: e,
	}


	/*addItemHandler := httptransport.NewServer(
		usecases.MakeAddItemEndpoint(svc),
		usecases.DecodeAddItemRequest,
		usecases.EncodeResponse,
	)*/
	/*getItemHandler := httptransport.NewServer(
		usecases.MakeGetItemEndpoint(svc),
		usecases.DecodeGetItemRequest,
		usecases.EncodeResponse,
	)*/

//	defer db.Close()
	//http.Handle("/items/get", getItemHandler)
	//http.Handle("/items/add", addItemHandler)
	//http.Handle("/metrics", promhttp.Handler())
	//http.ListenAndServe(":8080", nil)

errChan := make(chan error)

//r := mux.NewRouter()
//r.Methods("GET").Path("/items/get/{type}/{id}").Handler(getItemHandler)
//r.Methods("POST").Path("/items/add").Handler(addItemHandler)

r := usecases.MakeHttpHandler(ctx, endpoint, logger)


go func() {
		fmt.Println("Starting server at port 8080")
		handler := r
		errChan <- http.ListenAndServe(":8080", handler)
	}()


	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
fmt.Println(<- errChan)

}


