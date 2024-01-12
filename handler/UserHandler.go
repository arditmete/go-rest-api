package handler

import (
	"awesomeProject/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := h.Service.GetUsers()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting users: %s", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	user, err := h.Service.GetUserById(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting user with Id: %s with error: %s", params["id"], err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := h.Service.CreateUser(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user with error: %s", err), http.StatusInternalServerError)
		return
	} else {
		json.NewEncoder(w).Encode("User is created successfully!")
	}
}

func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	err := h.Service.UpdateUser(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating user witht Id: %s with error: %s", params["id"], err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("User is updated successfully!")
}

func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	err := h.Service.DeleteUser(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting user with Id: %s with error: %s", params["id"], err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("User is deleted successfully!")
}
