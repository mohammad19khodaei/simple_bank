package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/mohammad19khodaei/simple_bank/utils"
)

var CurrencyValidator validator.Func = func(fl validator.FieldLevel) bool {
	inputCurrency, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	return utils.IsValidCurrency(inputCurrency)
}
