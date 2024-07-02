package auth

import (
	"time"

	"github.com/dualex23/go-url-shortener/internal/app/config"
	"github.com/golang-jwt/jwt/v4"
) // Секретный ключ для подписи токена

// Claims структура для хранения информации в JWT
type Claims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

// GenerateToken создает новый JWT для пользователя
func GenerateToken(userID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Токен истекает через 24 часа
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.AppParseFlags().JWTkey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken проверяет JWT на валидность
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.AppParseFlags().JWTkey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
