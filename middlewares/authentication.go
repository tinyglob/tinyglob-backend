package middlewares

import (
	"net/http"
	"tinyglob-backend/types"
)

// AuthMiddleware is a middleware that checks the presence and validity of JWT tokens in the Authorization header.
func AuthMiddleware(next http.Handler) http.Handler {
	// Define the middleware function
	middlewareFunc := func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Request does not contain an access token", http.StatusUnauthorized)
			return
		}

		err := types.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	}

	// Return the middleware function
	return http.HandlerFunc(middlewareFunc)
}
