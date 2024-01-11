package handler

import (
	"awesomeProject/model"
	"awesomeProject/service"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
)

type LoginHandler struct {
	Service *service.UserService
	Auth    *AuthHandler
}

func NewLoginHandler(userService *service.UserService, authHandler *AuthHandler) *LoginHandler {
	return &LoginHandler{userService, authHandler}
}

func (h *LoginHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body with error: %s", err), http.StatusInternalServerError)
		return
	}
	var loginCredentials model.LoginRequest
	if err := json.Unmarshal(body, &loginCredentials); err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshalling JSON with error: %s", err), http.StatusInternalServerError)
		return
	}
	ctx := context.WithValue(r.Context(), "loginCredentials", loginCredentials)
	r = r.WithContext(ctx)
	// Retrieve user based on the provided username
	user, err := h.Service.GetUser(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving User: %s", err), http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(loginCredentials.Password))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid credentials: %s", err), http.StatusInternalServerError)
		return
	}
	token, err := h.Auth.GenerateJWT(loginCredentials.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating JWT: %s", err), http.StatusInternalServerError)
		return
	}
	response := map[string]string{"token": token}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Eror encoding response with error: %s", err), http.StatusInternalServerError)
	}
}
