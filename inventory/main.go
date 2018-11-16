package main

import (
   "fmt"
   "os"
   "os/signal"
   "syscall"
   "time"
   lg "log"
   "net/http"
   "golang.org/x/time/rate"
   "golang.org/x/net/context"
   "github.com/go-kit/kit/log"
   kitprometheus "github.com/go-kit/kit/metrics/prometheus"
   ratelimitkit "github.com/go-kit/kit/ratelimit"
   stdprometheus "github.com/prometheus/client_golang/prometheus"
   "github.com/comolago/shop/inventory/domain"
   "github.com/comolago/shop/inventory/infrastructure"
   "github.com/comolago/shop/inventory/usecases"

)

func main() {
   ctx := context.Background()
   errChan := make(chan error)

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
   err:=svc.Open()
   if err != nil {
      lg.Fatal(err)
   }
   svc = infrastructure.LoggingMiddleware(logger)(svc)
   svc = infrastructure.Metrics(requestCount, requestLatency)(svc)

   getItemEndpointRateLimit := rate.NewLimiter(rate.Every(35*time.Millisecond), 100)
   getItemEndpoint := usecases.MakeGetItemEndpoint(svc)
   getItemEndpoint = ratelimitkit.NewErroringLimiter(getItemEndpointRateLimit)(getItemEndpoint)

   addItemEndpointRateLimit := rate.NewLimiter(rate.Every(35*time.Millisecond), 100)
   addItemEndpoint := usecases.MakeAddItemEndpoint(svc)
   addItemEndpoint = ratelimitkit.NewErroringLimiter(addItemEndpointRateLimit)(addItemEndpoint)

   delItemEndpointRateLimit := rate.NewLimiter(rate.Every(35*time.Millisecond), 100)
   delItemEndpoint := usecases.MakeDelItemEndpoint(svc)
   delItemEndpoint = ratelimitkit.NewErroringLimiter(delItemEndpointRateLimit)(delItemEndpoint)


   key := []byte("supersecret")

   var auth domain.AuthHandler
   auth = infrastructure.AuthService{key, svc.GetDBHandler()}
   auth = infrastructure.LoggingAuthMiddleware(logger)(auth)
   //auth = instrumentingAuthMiddleware{requestAuthCount, requestAuthLatency, auth}
   
   authEndpointRateLimit := rate.NewLimiter(rate.Every(35*time.Millisecond), 100)
   authEndpoint := usecases.MakeAuthEndpoint(auth)
   authEndpoint = ratelimitkit.NewErroringLimiter(authEndpointRateLimit)(authEndpoint)

   endpoint := usecases.Endpoints{
      GetItemEndpoint: getItemEndpoint,
      AddItemEndpoint: addItemEndpoint,
      DelItemEndpoint: delItemEndpoint,
      AuthEndpoint: authEndpoint,
   }

   httpHandler := usecases.MakeMux(ctx, endpoint, auth, logger)

   go func() {
      fmt.Println("Starting server at port 8080")
      handler := httpHandler
      errChan <- http.ListenAndServe(":8080", handler)
   }()

   go func() {
      c := make(chan os.Signal, 1)
      signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
      errChan <- fmt.Errorf("%s", <-c)
   }()
   fmt.Println(<- errChan)
}


