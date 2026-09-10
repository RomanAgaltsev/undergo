// Package drill — C12/05 decode-handler.
package drill

import "encoding/json"

// Address is a shipping address.
type Address struct {
	City string `json:"city"`
	Zip  string `json:"zip"`
}

// Order is decoded from an incoming request body.
type Order struct {
	ID       string   `json:"id"`
	Customer string   `json:"customer"`
	Coupon   string   `json:"-"`
	Ship     *Address `json:"ship"`
}

// Decode parses an order and returns its destination city.
func Decode(data []byte) (Order, string, error) {
	var o Order
	json.Unmarshal(data, &o)
	city := o.Ship.City
	return o, city, nil
}
