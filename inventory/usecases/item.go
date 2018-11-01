// Application Business Rules
package usecases

import (
   "context"
   "encoding/json"
   "net/http"
   "github.com/comolago/shop/inventory/domain"
   "github.com/go-kit/kit/endpoint"
   "github.com/gorilla/mux"
   "strconv"
)

// Response type with a Item type message
type ItemResponse struct {
   Msg   domain.Item `json:"msg,omitempty"`
   Err *domain.ErrHandler `json:",omitempty"`
}

// Create GetItemById endpoint
func MakeGetItemEndpoint(svc domain.InventoryHandler) endpoint.Endpoint {
   return func(_ context.Context, request interface{}) (interface{}, error) {
      var req domain.Item
      var resp domain.Item
      var err *domain.ErrHandler
      req = request.(domain.Item)
      if req.Id >=0  {
         resp, err = svc.GetItemById(req.Id)
         if err != nil {
            return nil, err
         }
      } else {
         return nil, &domain.ErrHandler{3, "func ", "MakeGetItemEndpoint", ""}
      }
      return ItemResponse{resp, nil}, nil
   }
}

// Decode GetItemById requests
func DecodeGetItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
   vars := mux.Vars(r)
   requestType, ok := vars["type"]
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

// Create AddItem endpoint
func MakeAddItemEndpoint(svc domain.InventoryHandler) endpoint.Endpoint {
   return func(_ context.Context, request interface{}) (interface{}, error) {
      req := request.(domain.Item)
      v, err := svc.AddItem(req)
      if err != nil {
         //return StringResponse{v, err}, nil
         return nil, err
      }
      return StringResponse{v, nil}, nil
   }
}

// Decode AddItem requests
func DecodeAddItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
   var request domain.Item
   if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
      return nil, err
   }
   return request, nil
}

// Create DelItemById endpoint
func MakeDelItemEndpoint(svc domain.InventoryHandler) endpoint.Endpoint {
   return func(_ context.Context, request interface{}) (interface{}, error) {
      var req domain.Item
      var resp string
      var err *domain.ErrHandler
      req = request.(domain.Item)
      if req.Id >=0  {
         resp, err = svc.DelItemById(req.Id)
         if err != nil {
            return nil, err
         }
      } else {
         return nil, err
      }
      return StringResponse{resp, nil}, nil
   }
}

// Decode DelItemById requests
func DecodeDelItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
   vars := mux.Vars(r)
   /*requestType, ok := vars["type"]
   if !ok {
      return nil, &domain.ErrHandler{1, "func ", "DecodeDelItemRequest", requestType}
   }
   if requestType == "id" {*/
      tmpid, ok := vars["id"]
      if !ok {
         return nil, &domain.ErrHandler{2, "func ", "DecodeDelItemRequest", ""}
      }
      id, _ := strconv.Atoi(tmpid)
      return domain.Item{
         Id: id,
      }, nil
   /*} else {
      return nil, &domain.ErrHandler{9, "func ", "DecodeDelItemRequest", ""}
   }*/
}

