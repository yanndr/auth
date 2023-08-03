package validators

import (
	"fmt"
	"strings"
	"unicode"
)

type PasswordValidator struct {
	MinLength           int
	MinSpecialCharacter int
	MinNumeric          int
	MinSpecial          int
	MinUpper            int
	MinLower            int
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
		sb.WriteString(fmt.Sprintf("must be at least %v characters", v.MinLength))
	}
	if numerics < v.MinNumeric {
		sb.WriteString(fmt.Sprintf("must containt at least %v numeric characters", v.MinNumeric))
	}
	if upper < v.MinUpper {
		sb.WriteString(fmt.Sprintf("must containt at least %v uppercase characters", v.MinUpper))
	}

	if sb.Len() > 0 {
		return fmt.Errorf("password validation error: %s", sb.String())
	}
	return nil

}
