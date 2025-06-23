package middleware

import (
	"context"
	"curd-web-go/config"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from cookie
		cookie, err := r.Cookie("Token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr := cookie.Value

		// Parse token
		token, err := jwt.ParseWithClaims(tokenStr, &config.ClaimsJWT{}, func(token *jwt.Token) (interface{}, error) {
			return config.Jwt_Secret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		// Get email from token (Username di claims)
		claims := token.Claims.(*config.ClaimsJWT)
		ctx := context.Background()

		// Check in Redis
		savedToken, err := config.RedisClient.Get(ctx, claims.Username).Result()
		if err != nil || savedToken != tokenStr {
			http.Error(w, "Token expired or not found", http.StatusUnauthorized)
			return
		}

		// Next handler
		next.ServeHTTP(w, r)
	})
}
