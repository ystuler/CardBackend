package util

import (
	"back/internal/models"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

const signingKey = "mySecret"

func GenerateJWT(userModel *models.User) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   strconv.Itoa(userModel.ID),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signingKey))
}

func ParseToken(accessToken string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
