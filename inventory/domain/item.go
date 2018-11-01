// Business Rules
package domain

// an Item element
type Item struct {
   //Id   int `json:"id,string"`
   Id   int `json:"id"`
   Name string `json:"name"`
   Quantity int `json:"quantity"`
}
