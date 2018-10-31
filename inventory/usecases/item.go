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
)

type GetItemResponse struct {
   Msg   domain.Item `json:",omitempty"`
   Err *domain.ErrHandler `json:",omitempty"`
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

