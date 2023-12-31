package jwt

import (
	"auth/pkg/config"
	"auth/pkg/models"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type TokenGenerator interface {
	Generate(user models.User) (string, error)
}

type generator struct {
	signingMethod jwt.SigningMethod
	signedString  string
	issuer        string
	audience      string
	expDuration   time.Duration
}

func NewTokenGenerator(tokenConfig config.Token) TokenGenerator {
	return &generator{
		signingMethod: jwt.GetSigningMethod(tokenConfig.SigningMethod),
		signedString:  tokenConfig.SignedKey,
		expDuration:   time.Minute * time.Duration(tokenConfig.ExpDuration),
		issuer:        tokenConfig.Issuer,
		audience:      tokenConfig.Audience,
	}
}

// Generate generates a token from the models.User
func (g *generator) Generate(user models.User) (string, error) {
	token := jwt.New(g.signingMethod)

	claims := token.Claims.(jwt.MapClaims)

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
