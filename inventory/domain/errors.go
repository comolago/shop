package domain

import (
	"fmt"
)

type ErrHandler struct {
	Number      uintptr`json:"number,string"`
	Component   string `json:"component,string"`
	Function    string `json:"function,string"`
	Description string `json:"description,string"`
}

var errors = [...]string{
	1: "", // EXTERNAL_LIBRARY - used to handle external library errors
	2: "Item Id is empty - please provide a value",
	3: "Item Name is empty - please provide a value",
}

func (e ErrHandler) Error() string {
	var message string
	if 0 <= int(e.Number) && int(e.Number) < len(errors) {
		message = fmt.Sprintf("%s | %s", e.Component, e.Function)
		if int(e.Number) == 1 {
			message = fmt.Sprintf("%s | %s | %s | ", e.Component, e.Function, e.Description)
		} else {
			message = fmt.Sprintf("%s | %s | %s | %s", e.Component, e.Function, errors[e.Number], e.Description)
		}

	} else {
		message = fmt.Sprintf("%s | %s | Undefined error: error number is %d | %s", e.Component, e.Function, e.Number, e.Description)
	}
	return message
}
