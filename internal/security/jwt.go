package security

import (
	"time"

	customerrors "github.com/DKhorkov/hmtm-sso/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(value any, secretKey string, ttl time.Duration, algorithm string) (string, error) {
	token := jwt.New(jwt.GetSigningMethod(algorithm))
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", &customerrors.JWTClaimsError{}
	}

	claims["value"] = value
	claims["exp"] = time.Now().Add(ttl).Unix()
	return token.SignedString([]byte(secretKey))
}

func ParseJWT(tokenString, secretKey string) (any, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, &customerrors.InvalidJWTError{}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, &customerrors.JWTClaimsError{}
	}

	return claims["value"], nil
}
