package validators

import (
	autherror "auth/pkg/errors"
	"auth/pkg/models"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserValidator_Validate(t *testing.T) {

	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{"Valid input", models.User{"yann", "password"}, false},
		{"No username", models.User{"", "password"}, true},
		{"No password", models.User{"yann", ""}, true},
		{"No username and password", models.User{"", ""}, true},
		{"Nil input", models.User{"", ""}, true},
		{"Bad input", "username", true},
		{"Bad password", models.User{"yann", "a"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := UserValidator{
				StructValidator:   validator.New(),
				PasswordValidator: PasswordValidator{MinLength: 4},
			}
			err := v.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				require.ErrorAs(t, err, &autherror.ValidationErr{})
			}
		})
	}
}
