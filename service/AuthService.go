package service

import (
	"awesomeProject/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"time"
)

type AuthService struct {
	userService *UserService
}

func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{userService: userService}
}

var secretKey = []byte("golang")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (r *AuthService) AuthMiddleware(rr *http.Request) error {
	tokenString := rr.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return errors.New("Token is malformed!")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return errors.New("Token has expired or is not yet valid!")
			} else {
				return errors.New("Token validation error!")
			}
		} else {
			return errors.New("Token parsing error!")
		}
		return err
	}
	if err != nil {
		return errors.New(err.Error())
	}
	if err != nil && !token.Valid {
		return errors.New("Token is not valid!")
	}
	return nil
}

func (h *AuthService) GenerateJWT(username string) (string, error) {
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

func (h *AuthService) LoginUser(w http.ResponseWriter, r *http.Request) (string, error) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body with error: %s", err), http.StatusInternalServerError)
		return "", err
	}
	var loginCredentials model.LoginRequest
	if err := json.Unmarshal(body, &loginCredentials); err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshalling JSON with error: %s", err), http.StatusInternalServerError)
		return "", err
	}
	ctx := context.WithValue(r.Context(), "loginCredentials", loginCredentials)
	r = r.WithContext(ctx)
	// Retrieve user based on the provided username
	user, err := h.userService.GetUser(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving User: %s", err), http.StatusInternalServerError)
		return "", err
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(loginCredentials.Password))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid credentials: %s", err), http.StatusInternalServerError)
		return "", err
	}
	token, err := h.GenerateJWT(loginCredentials.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating JWT: %s", err), http.StatusInternalServerError)
		return "", err
	}
	response := map[string]string{"token": token}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Eror encoding response with error: %s", err), http.StatusInternalServerError)
	}
	return token, nil
}
