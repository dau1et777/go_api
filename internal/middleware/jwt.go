package middleware

import (
	"net/http"

	"go-api/cmd/server/auth"
)

func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := auth.ExtractToken(r)
		if tokenString == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		claims, err := auth.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Token is valid, call the next handler
		r.Header.Set("user_id", string(rune(claims.UserID)))
		r.Header.Set("email", claims.Email)
		next(w, r)
	}
}
