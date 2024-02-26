package verify

import (
	"fmt"
	"net/http"
	"strings"
)

func JwtMiddleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		token, err := extractToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		claims, err := verifyToken(token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Print(claims)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("token is not provided")
	}
	token, isBearer := strings.CutPrefix(authHeader, "Bearer ")
	if !isBearer {
		return "", fmt.Errorf("invalid auth token")
	}
	fmt.Println(token)
	return "", nil
}

func verifyToken(token string) (string, error) {
	return "", nil
}
