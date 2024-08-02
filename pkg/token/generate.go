package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTPayload map[string]any

type TokenGenOptions struct {
	Payload    JWTPayload
	Secret     string
	ExpiryDate time.Time
}

type JWTClaim struct {
	Payload JWTPayload `json:"payload"`
	jwt.RegisteredClaims
}

func GenerateHS256Token(tokenGenOptions TokenGenOptions) (string, error) {

	if tokenGenOptions.Secret == "" {
		tokenGenOptions.Secret = os.Getenv("TOKEN_SECRET")
	}

	claims := &JWTClaim{
		Payload: tokenGenOptions.Payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: tokenGenOptions.ExpiryDate},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tokenGenOptions.Secret))

}
