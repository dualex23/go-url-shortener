package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const UserIDKey contextKey = "userID"

// Claims структура для хранения информации в JWT
type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

// GenerateToken создает новый JWT для пользователя
func GenerateToken(userID string, JWTtoken []byte) (string, error) {
	expirationTime := time.Now().Add(72 * time.Hour)
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTtoken)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken проверяет JWT на валидность
func ValidateToken(tokenString string, JWTtoken []byte) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWTtoken, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
