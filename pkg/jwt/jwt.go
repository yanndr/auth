package jwt

import (
	"auth/pkg/config"
	"auth/pkg/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Generator interface {
	GenerateJWT(user *models.User) (string, error)
}

type generator struct {
	signingMethod jwt.SigningMethod
	signedString  string
	issuer        string
	audience      string
	expDuration   time.Duration
}

func NewGenerator(tokenConfig config.Token) Generator {
	return &generator{
		signingMethod: jwt.GetSigningMethod(tokenConfig.SigningMethod),
		signedString:  tokenConfig.SignedKey,
		expDuration:   time.Minute * time.Duration(tokenConfig.ExpDuration),
		issuer:        tokenConfig.Issuer,
		audience:      tokenConfig.Audience,
	}
}

func (g *generator) GenerateJWT(user *models.User) (string, error) {
	token := jwt.New(g.signingMethod)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["sub"] = user.Username
	claims["iss"] = g.issuer
	claims["aud"] = g.audience
	claims["exp"] = time.Now().Add(g.expDuration).Unix()
	claims["iat"] = time.Now().Unix()
	tokenString, err := token.SignedString([]byte(g.signedString))

	if err != nil {
		return "", fmt.Errorf("error creating the token: %w", err)
	}

	return tokenString, nil
}
