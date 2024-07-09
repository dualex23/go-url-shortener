package middleware

import (
	"context"
	"net/http"

	"github.com/dualex23/go-url-shortener/internal/app/auth"
	"github.com/dualex23/go-url-shortener/internal/app/logger"
	"github.com/golang-jwt/jwt/v4"
)

func Authenticate(jwtKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil {
				logger.GetLogger().Error("Error retrieving cookie:", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			logger.GetLogger().Infoln("func:", "Authenticate", "cookie:", cookie.Value)

			tokenString := cookie.Value
			claims := &auth.Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil {
				logger.GetLogger().Error("Error parsing token:", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				logger.GetLogger().Error("Invalid token")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), auth.UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
