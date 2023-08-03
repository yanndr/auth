package validators

import (
	"auth/pkg/model"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type ValidationErr struct {
	err error
}

func (v ValidationErr) Error() string {
	return fmt.Sprintf("validation error: %s", v.err)
}

func NewValidationErr(err error) ValidationErr {
	return ValidationErr{err: err}
}

type Validator interface {
	Validate(value any) error
}

type UserValidator struct {
	Validator         *validator.Validate
	PasswordValidator PasswordValidator
}

func (v UserValidator) Validate(value any) error {
	user, ok := value.(model.User)
	if !ok {
		return NewValidationErr(fmt.Errorf("the input value is not a model.User"))
	}
	err := v.Validator.Struct(user)
	if err != nil {
		return NewValidationErr(err)
	}
	err = v.PasswordValidator.Validate(user.Password)
	if err != nil {
		return NewValidationErr(err)
	}

	return nil
}
