package commons

import (
	"github.com/go-playground/validator/v10"
	"strings"
	"unicode"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
	Validate.RegisterValidation("password", PasswordValidation)
}

// PasswordValidation ...
func PasswordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	hasDigit := false
	hasLower := false
	hasUpper := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:'\"<>,.?/~`", char):
			hasSpecial = true
		}
	}

	return hasDigit && hasLower && hasUpper && hasSpecial
}
