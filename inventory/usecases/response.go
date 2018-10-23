package usecases

import (
	"context"
	"encoding/json"
	"net/http"

)

type Response struct {
	Msg   string `json:"msg"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
