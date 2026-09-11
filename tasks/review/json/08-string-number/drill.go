// Package drill — C12/08 string-number.
package drill

import "encoding/json"

// Money is decoded from a payment API that sends the amount as a JSON string,
// e.g. {"amount":"100"}.
type Money struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

// ParseMoney decodes a Money value.
func ParseMoney(data []byte) (Money, error) {
	var m Money
	if err := json.Unmarshal(data, &m); err != nil {
		return Money{}, err
	}
	return m, nil
}
