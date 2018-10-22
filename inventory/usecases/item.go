package usecases

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/comolago/shop/inventory/domain"
	"github.com/go-kit/kit/endpoint"
)

// For each method, we define request and response structs
type AddItemRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

// Endpoints are a primary abstraction in go-kit. An endpoint represents a single RPC (method in our service interface)
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

func MakeCountEndpoint(svc domain.InventoryInt) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v := svc.Count(req.S)
		return countResponse{v}, nil
	}
}

func DecodeAddItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
