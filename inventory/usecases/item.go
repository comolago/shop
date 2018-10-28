package usecases

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/comolago/shop/inventory/domain"
	"github.com/go-kit/kit/endpoint"
)

func MakeAddItemEndpoint(svc domain.InventoryHandler) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(domain.Item)
		v, err := svc.AddItem(req.Id, req.Name)
		if err != nil {
			return Response{v, err}, nil
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
