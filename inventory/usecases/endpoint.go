// Application Business Rules
package usecases

import (
   "context"
   "encoding/json"
   "net/http"
   "github.com/go-kit/kit/endpoint"
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
   DelItemEndpoint endpoint.Endpoint
   AuthEndpoint endpoint.Endpoint
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

func EncodeAuthResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
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


