package usecases

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/comolago/shop/inventory/domain"
	"github.com/go-kit/kit/endpoint"
   "github.com/gorilla/mux"
   "strconv"
   "fmt"
"github.com/go-kit/kit/log"
httptransport "github.com/go-kit/kit/transport/http"

"github.com/prometheus/client_golang/prometheus/promhttp"
)

type GetItemResponse struct {
   Msg   domain.Item `json:",omitempty"`
   Err *domain.ErrHandler `json:",omitempty"`
}

// endpoints wrapper
type Endpoints struct {
	GetItemEndpoint endpoint.Endpoint
	AddItemEndpoint endpoint.Endpoint
}


// Make Http Handler
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




func MakeGetItemEndpoint(svc domain.InventoryHandler) endpoint.Endpoint {
   return func(_ context.Context, request interface{}) (interface{}, error) {
      var req domain.Item
      var resp domain.Item
      var err *domain.ErrHandler
      req = request.(domain.Item)
      if req.Id >=0  {
         resp, err = svc.GetItemById(req.Id)
         fmt.Println(resp.Id)
         if err != nil {
            return nil, err
         }
      } else {
         return nil, &domain.ErrHandler{3, "func ", "MakeGetItemEndpoint", ""}
      }
      return GetItemResponse{resp, nil}, nil
   }
}

func DecodeGetItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
   fmt.Println("qui")
   vars := mux.Vars(r)
   requestType, ok := vars["type"]
   fmt.Println(requestType)
   if !ok {
      return nil, &domain.ErrHandler{1, "func ", "DecodeGetItemRequest", requestType}
   }
   if requestType == "id" {
      tmpid, ok := vars["id"]
      if !ok {
         return nil, &domain.ErrHandler{2, "func ", "DecodeGetItemRequest", ""}
      }
      id, _ := strconv.Atoi(tmpid)
      return domain.Item{
         Id: id,
      }, nil
   } else {
      return nil, &domain.ErrHandler{9, "func ", "DecodeGetItemRequest", ""}
   }
}

func MakeAddItemEndpoint(svc domain.InventoryHandler) endpoint.Endpoint {
   return func(_ context.Context, request interface{}) (interface{}, error) {
   fmt.Println("qua")
      req := request.(domain.Item)
      v, err := svc.AddItem(req.Id, req.Name)
      if err != nil {
         //return Response{v, err}, nil
         return nil, err
      }
      return Response{v, nil}, nil
   }
}


func DecodeAddItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
   var request domain.Item
   if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
      return nil, err
   }
   return request, nil
}

