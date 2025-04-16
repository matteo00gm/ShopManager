package utils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = os.Getenv("JWT_SECRET")

func VerifyToken(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		//checking if type of the value stored in Method is of type *jwt.SigningMethodHMAC
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, errors.New("could not parse token")
	}

	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		return nil, errors.New("invalid token")
	}
	return parsedToken, nil
}

func GetParamsFromToken(token *jwt.Token, params ...string) ([]any, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}
	var parsedParams []any
	for _, param := range params {
		parsedParams = append(parsedParams, claims[param])
	}
	return parsedParams, nil
}
