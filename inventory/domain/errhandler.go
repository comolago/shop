// Business Rules
package domain

import (
   "fmt"
)

// Error Handling Struc
// You can specify:
// unique error Number,
// Component (the receiver),
// the Function/Method,
// Description
type ErrHandler struct {
   Number      uintptr`json:"number,string"`
   Component   string `json:"component,string"`
   Function    string `json:"function,string"`
   Description string `json:"description,string"`
}

// Error() method implementation so to satisfy Golang standard error interface
// When error Number is equal to 1 only Description is rendered into return message
// elsewhere error message contains also the relate error message defined in errors array of strings
// See errors.go to view the list of errors and number 
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
