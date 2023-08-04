package validators

import (
	"auth/pkg/model"
	"fmt"
	"strings"
	"unicode"
)

type PasswordValidator struct {
	MinLength  int
	MinNumeric int
	MinSpecial int
	MinUpper   int
	MinLower   int
}

func NewPasswordValidator(configuration model.Password) *PasswordValidator {
	return &PasswordValidator{
		MinLength:  configuration.MinLength,
		MinNumeric: configuration.MinNumeric,
		MinSpecial: configuration.MinSpecial,
		MinUpper:   configuration.MinUpperCase,
		MinLower:   configuration.MinLowerCase,
	}
}

func (v PasswordValidator) Validate(value any) error {
	password, ok := value.(string)
	if !ok {
		return fmt.Errorf("input value must be a string")
	}
	var length, numerics, upper, lower, special int
	for _, c := range password {
		length++
		switch {
		case unicode.IsNumber(c):
			numerics++
		case unicode.IsUpper(c):
			upper++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special++
		case unicode.IsLetter(c) || c == ' ':
			lower++
		}
	}

	sb := strings.Builder{}
	if length < v.MinLength {
		sb.WriteString(fmt.Sprintf("must be at least %v characters\n", v.MinLength))
	}
	if numerics < v.MinNumeric {
		sb.WriteString(fmt.Sprintf("must containt at least %v numeric characters\n", v.MinNumeric))
	}
	if upper < v.MinUpper {
		sb.WriteString(fmt.Sprintf("must containt at least %v uppercase characters\n", v.MinUpper))
	}
	if lower < v.MinLower {
		sb.WriteString(fmt.Sprintf("must containt at least %v lowercase characters\n", v.MinLower))
	}
	if special < v.MinSpecial {
		sb.WriteString(fmt.Sprintf("must containt at least %v sepcial characters\n", v.MinSpecial))
	}

	if sb.Len() > 0 {
		return fmt.Errorf("password validation error: %s", sb.String())
	}
	return nil

}
