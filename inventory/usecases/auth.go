package usecases

import (
   "context"
   "net/http"
   "encoding/json"
   "github.com/go-kit/kit/endpoint"
   "github.com/comolago/shop/inventory/domain"
)

func MakeAuthEndpoint(svc domain.AuthHandler) endpoint.Endpoint {
   return func(ctx context.Context, request interface{}) (interface{}, error) {
      req := request.(domain.AuthRequest)
      token, err := svc.Auth(req.ClientID, req.ClientSecret)
      if err != nil {
         return nil, err
      }
      return domain.AuthResponse{token, ""}, nil
   }
}

func DecodeAuthRequest(_ context.Context, r *http.Request) (interface{}, error) {
   var request domain.AuthRequest
   if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
      return nil, err
   }
   return request, nil
}

func AuthErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
        code := http.StatusUnauthorized
        msg := err.Error()

        w.WriteHeader(code)
        json.NewEncoder(w).Encode(domain.AuthResponse{Token: "", Err: msg})
}
