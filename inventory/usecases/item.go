package usecases

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/comolago/shop/inventory/domain"
	"github.com/go-kit/kit/endpoint"
)

type AddItemRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func MakeAddItemEndpoint(svc domain.InventoryInt) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(AddItemRequest)
		v, err := svc.AddItem(req.Id, req.Name)
		if err != nil {
			return Response{v, err.Error()}, nil
		}
		return Response{v, ""}, nil
	}
}

func DecodeAddItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
