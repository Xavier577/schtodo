package token

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenVerifyOptions struct {
	Secret      string
	SignedToken string
}

var (
	ErrExpiredToken = errors.New("TOKEN_EXPIRED")
	ErrParseToken   = errors.New("TOKEN_PARSE_ERR")
)

func Verify(tokenVerifyOptions *TokenVerifyOptions) (bool, JWTClaim, error) {

	var isValid bool
	var verificationError error
	var jwtClaim JWTClaim

	if tokenVerifyOptions.Secret == "" {
		tokenVerifyOptions.Secret = os.Getenv("TOKEN_SECRET")
	}

	token, err := jwt.ParseWithClaims(
		tokenVerifyOptions.SignedToken,
		&JWTClaim{},
		func(t *jwt.Token) (any, error) {
			return []byte(tokenVerifyOptions.Secret), nil
		},
	)

	if err != nil {
		isValid = false
		verificationError = err
	}

	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		isValid = true
		jwtClaim = *claims
	} else if !ok {
		verificationError = ErrParseToken
	} else if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		verificationError = ErrExpiredToken
	}

	return isValid, jwtClaim, verificationError
}
