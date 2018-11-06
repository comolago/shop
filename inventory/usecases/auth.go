package usecases

import (
   "context"
   "net/http"
   "encoding/json"
   "github.com/go-kit/kit/endpoint"
   "github.com/comolago/shop/inventory/infrastructure"
)

func MakeAuthEndpoint(svc infrastructure.AuthHandler) endpoint.Endpoint {
   return func(ctx context.Context, request interface{}) (interface{}, error) {
      req := request.(infrastructure.AuthRequest)
      token, err := svc.Auth(req.ClientID, req.ClientSecret)
      if err != nil {
         return nil, err
      }
      return infrastructure.AuthResponse{token, ""}, nil
   }
}

func DecodeAuthRequest(_ context.Context, r *http.Request) (interface{}, error) {
   var request infrastructure.AuthRequest
   if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
      return nil, err
   }
   return request, nil
}

func AuthErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
        code := http.StatusUnauthorized
        msg := err.Error()

        w.WriteHeader(code)
        json.NewEncoder(w).Encode(infrastructure.AuthResponse{Token: "", Err: msg})
}
