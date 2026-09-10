// Package drill — C12/01 unexported-field.
package drill

import "encoding/json"

// Account is serialized to JSON for the accounts API.
type Account struct {
	name    string `json:"name"`
	Balance int    `json:"balance"`
}

// NewAccount builds an account.
func NewAccount(name string, balance int) Account {
	return Account{name: name, Balance: balance}
}

// Encode renders the account as JSON.
func Encode(a Account) ([]byte, error) {
	return json.Marshal(a)
}
