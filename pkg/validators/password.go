package validators

import (
	"auth/pkg/config"
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

func NewPasswordValidator(configuration config.Password) *PasswordValidator {
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

	errs := make([]string, 0)

	if length < v.MinLength {
		errs = append(errs, fmt.Sprintf("must be at least %v characters", v.MinLength))
	}
	if numerics < v.MinNumeric {

		errs = append(errs, fmt.Sprintf("must containt at least %v numeric characters", v.MinNumeric))
	}
	if upper < v.MinUpper {
		errs = append(errs, fmt.Sprintf("must containt at least %v uppercase characters", v.MinUpper))
	}
	if lower < v.MinLower {
		errs = append(errs, fmt.Sprintf("must containt at least %v lowercase characters", v.MinLower))
	}
	if special < v.MinSpecial {
		errs = append(errs, fmt.Sprintf("must containt at least %v sepcial characters", v.MinSpecial))
	}

	if len(errs) > 0 {
		return fmt.Errorf("password validation error: %s", strings.Join(errs, ", "))
	}
	return nil

}
