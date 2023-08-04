package jwt

import (
	"auth/pkg/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Generator interface {
	GenerateJWT(user *model.User) (string, error)
}

type generator struct {
	signingMethod jwt.SigningMethod
	signedString  string
	issuer        string
	expDuration   time.Duration
}

func NewGenerator(SigningMethod jwt.SigningMethod, signedString, issuer string, expDuration time.Duration) Generator {
	return generator{
		signingMethod: SigningMethod,
		signedString:  signedString,
		expDuration:   expDuration,
		issuer:        issuer,
	}
}

func (g generator) GenerateJWT(user *model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["sub"] = user.Username
	claims["iss"] = g.issuer
	claims["exp"] = time.Now().Add(g.expDuration).Unix()
	claims["iat"] = time.Now().Unix()
	tokenString, err := token.SignedString([]byte(g.signedString))

	if err != nil {
		_ = fmt.Errorf("something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}
