package service

import (
	"awesomeProject/model"
	"awesomeProject/repository"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

type UserService struct {
	Repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{Repository: repository}
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(username string) (string, error) {
	// Set expiration time for the token (e.g., 1 hour)
	expirationTime := time.Now().Add(1 * time.Hour)

	// Create JWT claims
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the token with claims and sign it with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString("ardit")
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) GetUsers() ([]model.User, error) {
	return s.Repository.GetUsers()
}

func (s *UserService) GetUser(w http.ResponseWriter, r *http.Request) (model.User, error) {
	return s.Repository.GetUser(w, r)
}

func (s *UserService) GetUserById(w http.ResponseWriter, r *http.Request) (model.User, error) {
	return s.Repository.GetUserByID(w, r)
}

func (s *UserService) CreateUser(w http.ResponseWriter, r *http.Request) error {
	return s.Repository.CreateUser(w, r)
}

func (s *UserService) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	return s.Repository.UpdateUser(w, r)
}

func (s *UserService) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	return s.Repository.DeleteUser(w, r)
}
