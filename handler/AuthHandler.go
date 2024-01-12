package handler

import (
	"awesomeProject/service"
	"fmt"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (r *AuthHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, rr *http.Request) {
		err := r.authService.AuthMiddleware(rr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error in authentication with error: %s", err), http.StatusUnauthorized)
			return
		}
		next(w, rr)
	}
}

func (r *AuthHandler) LoginUserHandler(w http.ResponseWriter, rr *http.Request) {
	_, err := r.authService.LoginUser(w, rr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in login with error: %s", err), http.StatusInternalServerError)
		return
	}
}
