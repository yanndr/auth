package validators

import (
	autherror "auth/pkg/errors"
	"auth/pkg/models"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(value any) error
}

type UserValidator struct {
	Validator         *validator.Validate
	PasswordValidator Validator
}

func (v UserValidator) Validate(value any) error {
	user, ok := value.(models.User)
	if !ok {
		return autherror.NewValidationErr(fmt.Errorf("the input value is not a model.User"))
	}
	err := v.Validator.Struct(user)
	if err != nil {
		return autherror.NewValidationErr(err)
	}
	err = v.PasswordValidator.Validate(user.Password)
	if err != nil {
		return autherror.NewValidationErr(err)
	}

	return nil
}