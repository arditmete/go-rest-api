package service

import (
	"awesomeProject/model"
	"awesomeProject/repository"
	"net/http"
)

type UserService struct {
	Repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{Repository: repository}
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
