package usecases

import (
	"context"
	"encoding/json"
	"net/http"
	"github.com/comolago/shop/inventory/domain"

)

type Response struct {
	Msg   string `json:"msg,omitempty"`
	//Err string `json:"err,omitempty"` // errors don't define JSON marshaling
	Err *domain.ErrHandler `json:",omitempty"`
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
