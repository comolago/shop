package usecases

import (
   "context"
   "encoding/json"
   "net/http"
   "github.com/go-kit/kit/endpoint"
   "github.com/gorilla/mux"
   "github.com/go-kit/kit/log"
   httptransport "github.com/go-kit/kit/transport/http"
   "github.com/prometheus/client_golang/prometheus/promhttp"
)

type Endpoints struct {
   GetItemEndpoint endpoint.Endpoint
   AddItemEndpoint endpoint.Endpoint
}

func MakeHttpHandler(_ context.Context, endpoint Endpoints, logger log.Logger) http.Handler {
   r := mux.NewRouter()
   options := []httptransport.ServerOption{
      httptransport.ServerErrorLogger(logger),
      httptransport.ServerErrorEncoder(encodeError),
   }
   r.Methods("GET").Path("/items/get/{type}/{id}").Handler(httptransport.NewServer(
      endpoint.GetItemEndpoint,
      DecodeGetItemRequest,
      EncodeResponse,
      options...,
   ))
   r.Methods("POST").Path("/items/add").Handler(httptransport.NewServer(
      endpoint.AddItemEndpoint,
      DecodeAddItemRequest,
      EncodeResponse,
      options...,
   ))
   r.Methods("GET").Path("/metrics").Handler(promhttp.Handler())
   return r
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
   if err == nil {
      panic("encodeError with nil error")
   }
   w.Header().Set("Content-Type", "application/json; charset=utf-8")
   w.WriteHeader(http.StatusInternalServerError)
   json.NewEncoder(w).Encode(map[string]interface{}{
      "error": err.Error(),
   })
}
