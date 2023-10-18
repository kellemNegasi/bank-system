package util

const (
	USD = "USD"
	GBP = "GBP"
	EUR = "EUR"
)

func IsValidCurrency(currency string) bool {
	switch currency {
	case USD, EUR, GBP:
		return true
	}
	return false
}
