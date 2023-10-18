package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/kellemNegasi/bank-system/util"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	currency, ok := fl.Field().Interface().(string)
	if ok {
		return util.IsValidCurrency(currency)
	}
	return false
}
