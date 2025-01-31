package utils

func GetValidCurrencies() []string {
	return []string{"USD", "EUR", "IRR"}
}

// Optional: Case-insensitive check helper
func IsValidCurrency(input string) bool {
	validCurrencies := GetValidCurrencies()

	for _, currency := range validCurrencies {
		if input == currency {
			return true
		}
	}
	return false
}
