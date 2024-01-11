package model

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password []byte `json:"password"`
}
