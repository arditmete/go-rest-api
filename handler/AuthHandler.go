package handler

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

type AuthHandler struct {
}

var secretKey = []byte("golang")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (h *AuthHandler) GenerateJWT(username string) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (r *AuthHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		fmt.Println("AuthMiddleware token -> " + tokenString)
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println("AuthMiddleware")
		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil {
			// Add more prints for debugging
			fmt.Println("Token parsing error:", err)
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					// Token is malformed
					fmt.Println("Token is malformed")
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					// Token has expired or is not yet valid
					fmt.Println("Token has expired or is not yet valid")
				} else {
					// Other validation errors
					fmt.Println("Token validation error:", err)
				}
			} else {
				// Other parsing errors
				fmt.Println("Token parsing error:", err)
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err != nil {
			fmt.Println("AuthMiddleware error != null")
		}
		if !token.Valid {
			fmt.Println("AuthMiddleware token != Valid")
		}
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If the token is valid, call the next handler
		next.ServeHTTP(w, r)
	}
}
