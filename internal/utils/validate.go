package utils

import (
	"regexp"
	"slices"

	"github.com/BetOnz-Company/betonz-go/internal/product"

	"github.com/go-playground/validator/v10"
)

func ValidateUsername(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile("^[a-zA-Z0-9_]+$")
	return regex.MatchString(fl.Field().String())
}

func ValidateBankAccountNumber(fl validator.FieldLevel) bool {
	parent := fl.Parent()
	bankName := parent.FieldByName("BankName").String()

	if bankName == "KBZPAY" || bankName == "OK_DOLLAR" || bankName == "WAVE_PAY" {
		regex := regexp.MustCompile("^\\d\\d \\d\\d\\d \\d\\d\\d \\d\\d\\d$")
		return regex.MatchString(fl.Field().String())
	}
	if bankName == "KBZ" {
		regex := regexp.MustCompile("^\\d\\d\\d\\d \\d\\d\\d\\d \\d\\d\\d\\d \\d\\d\\d\\d \\d$")
		return regex.MatchString(fl.Field().String())
	}
	if bankName == "CB" {
		regex := regexp.MustCompile("^\\d\\d\\d\\d \\d\\d\\d\\d \\d\\d\\d\\d \\d\\d\\d\\d$")
		return regex.MatchString(fl.Field().String())
	}
	if bankName == "AGD" || bankName == "AYA" || bankName == "YOMA" {
		regex := regexp.MustCompile("^\\d\\d\\d\\d \\d\\d\\d\\d \\d\\d\\d\\d \\d\\d\\d\\d")
		return regex.MatchString(fl.Field().String())
	}

	return true
}

func ValidateProduct(fl validator.FieldLevel) bool {
	i := product.Product(fl.Field().Int())
	return i == product.MainWallet || slices.Contains(product.AllProducts, i)
}
