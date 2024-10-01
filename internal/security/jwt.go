package security

import (
	"github.com/DKhorkov/hmtm-sso/pkg/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/pkg/errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(user *entities.User, secretKey string, ttl time.Duration, algorithm string) (string, error) {
	token := jwt.New(jwt.GetSigningMethod(algorithm))
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", &customerrors.JWTClaimsError{}
	}

	claims["userID"] = user.ID
	claims["exp"] = time.Now().Add(ttl).Unix()
	return token.SignedString([]byte(secretKey))
}

func ParseJWT(tokenString, secretKey string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return 0, &customerrors.InvalidJWTError{}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, &customerrors.JWTClaimsError{}
	}

	return int(claims["userID"].(float64)), nil
}
