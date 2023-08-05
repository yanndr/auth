package jwt

import (
	"auth/pkg/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_generator_Generate(t *testing.T) {

	g := &generator{
		signingMethod: jwt.SigningMethodHS256,
		signedString:  "signedstring",
		issuer:        "test",
		audience:      "audience",
		expDuration:   5,
	}
	user := models.User{Username: "test"}

	token, err := g.Generate(user)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}
