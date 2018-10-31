// Application Business Rules
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
   "github.com/comolago/shop/inventory/domain"
)

// Response type with a string message
type StringResponse struct {
   Msg   string `json:"msg,omitempty"`
   Err *domain.ErrHandler `json:",omitempty"`
}

// Struct with all the exposed endpoints
type Endpoints struct {
   GetItemEndpoint endpoint.Endpoint
   AddItemEndpoint endpoint.Endpoint
}

// Create a Mux router with all the endpoints
func MakeHttpHandler(_ context.Context, endpoint Endpoints, logger log.Logger) http.Handler {
   r := mux.NewRouter()
   options := []httptransport.ServerOption{
      httptransport.ServerErrorLogger(logger),
      httptransport.ServerErrorEncoder(encodeError),
   }
   r.Methods("GET").Path("/items/get/{type}/{id}").Handler(httptransport.NewServer(
      endpoint.GetItemEndpoint,
      DecodeGetItemRequest,
      EncodeItemResponse,
      options...,
   ))
   r.Methods("POST").Path("/items/add").Handler(httptransport.NewServer(
      endpoint.AddItemEndpoint,
      DecodeAddItemRequest,
      EncodeStringResponse,
      options...,
   ))
   r.Methods("GET").Path("/metrics").Handler(promhttp.Handler())
   return r
}

// Encode the error message
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

// Encode a response of StringResponse type
func EncodeStringResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
   e := response.(StringResponse).Err
   if e != nil {
      encodeError(ctx, e, w)
      return nil
   }
   w.Header().Set("Content-Type", "application/json; charset=utf-8")
   return json.NewEncoder(w).Encode(response)
}

// Encode a response of ItemResponse type
func EncodeItemResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
   e := response.(ItemResponse).Err
   if e != nil {
      encodeError(ctx, e, w)
      return nil
   }
   w.Header().Set("Content-Type", "application/json; charset=utf-8")
   return json.NewEncoder(w).Encode(response)
}
