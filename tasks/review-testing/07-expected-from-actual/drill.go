// Package drill — C10/07 expected-from-actual.
package drill

// Discount applies a 10% discount to price.
func Discount(price int) int {
	return price - price*10/100
}
