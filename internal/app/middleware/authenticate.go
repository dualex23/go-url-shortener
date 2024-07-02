package middleware

import (
	"context"
	"net/http"

	"github.com/dualex23/go-url-shortener/internal/app/auth"
	"github.com/dualex23/go-url-shortener/internal/app/config"
	"github.com/golang-jwt/jwt/v4"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// Если кука не установлена
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			// Для других ошибок
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		tokenString := cookie.Value
		claims := &auth.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.AppParseFlags().JWTkey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Добавить ID пользователя в контекст запроса
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
