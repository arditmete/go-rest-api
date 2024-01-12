package repository

import (
	"awesomeProject/model"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetUsers() ([]model.User, error) {
	result, err := r.DB.Query("SELECT * from users")
	var users []model.User
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var user model.User
		err := result.Scan(&user.ID, &user.Username)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}
	return users, err
}

func (rr *UserRepository) GetUser(w http.ResponseWriter, r *http.Request) (model.User, error) {
	w.Header().Set("Content-Type", "application/json")
	loginCredentials, ok := r.Context().Value("loginCredentials").(model.LoginRequest)
	if !ok {
		return model.User{}, fmt.Errorf("User information not found in context!")
	}
	result, err := rr.DB.Query("SELECT * FROM users WHERE username = ?", loginCredentials.Username)
	if err != nil {
		return model.User{}, fmt.Errorf("Error querying database: %v", err)
	}
	defer result.Close()

	var user model.User
	if result.Next() == false {
		return model.User{}, errors.New("User not found")
	}
	err = result.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("Error scanning database result: %v", err)
	}
	return user, nil
}

func (r *UserRepository) CreateUser(w http.ResponseWriter, rr *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	// Prepare SQL statement
	stmt, err := r.DB.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing SQL statement: %v", err)
	}

	// Read request body
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %v", err)
	}

	// Unmarshal JSON
	keyVal := make(map[string]string)
	if err := json.Unmarshal(body, &keyVal); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	// Extract username and password
	username := keyVal["username"]
	password := keyVal["password"]

	// Encrypt password
	passwordEncrypted, err := bcrypt.GenerateFromPassword([]byte(password), 2)
	if err != nil {
		return fmt.Errorf("error encrypting password: %v", err)
	}

	// Check if the username already exists
	var existingUsername string
	err = r.DB.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&existingUsername)
	if err == nil {
		// User already exists, return an error
		return fmt.Errorf("user already exists")
	} else if err != sql.ErrNoRows {
		// An unexpected error occurred
		return fmt.Errorf("error checking existing username: %v", err)
	}
	// Execute the INSERT statement
	_, err = stmt.Exec(username, passwordEncrypted)
	if err != nil {
		return fmt.Errorf("error executing SQL statement: %v", err)
	}

	return nil
}

func (r *UserRepository) DeleteUser(w http.ResponseWriter, rr *http.Request) error {
	params := mux.Vars(rr)
	_, err := r.GetUserByID(w, rr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving user: %s", err), http.StatusInternalServerError)
		return err
	}
	stmt, err := r.DB.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	return err
}

func (r *UserRepository) UpdateUser(w http.ResponseWriter, rr *http.Request) error {
	params := mux.Vars(rr)
	_, err := r.GetUserByID(w, rr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving user: %s", err), http.StatusInternalServerError)
		return err
	}
	stmt, err := r.DB.Prepare("UPDATE users SET username = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newUsername := keyVal["username"]
	_, err = stmt.Exec(newUsername, params["id"])
	if err != nil {
		panic(err.Error())
	}
	return err
}

func (h *UserRepository) GetUserByID(w http.ResponseWriter, r *http.Request) (model.User, error) {
	var user model.User
	params := mux.Vars(r)
	result, err := h.DB.Query("SELECT * FROM users WHERE id = ?", params["id"])
	if err != nil {
		return user, err
	}
	defer result.Close()
	if result.Next() == false {
		return user, errors.New("User not found")
	}
	err = result.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}
