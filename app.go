package main

import (
	"awesomeProject/handler"
	_ "awesomeProject/handler"
	"awesomeProject/repository"
	"awesomeProject/service"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
)

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/posts")
	if err != nil {
		_ = fmt.Errorf("Connection to database is refused")
		panic(err.Error())
	}
	defer db.Close()

	postRepository := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepository)
	postHandler := handler.NewPostHandler(postService)

	// Initialize handlers with the service and repository dependencies
	getPostsHandler := postHandler.GetPostsHandler
	createPostHandler := postHandler.CreatePostHandler
	getPostHandler := postHandler.GetPostHandler
	updatePostHandler := postHandler.UpdatePostHandler
	deletePostHandler := postHandler.DeletePostHandler

	router := mux.NewRouter()
	router.HandleFunc("/posts", getPostsHandler).Methods("GET")
	router.HandleFunc("/posts", createPostHandler).Methods("POST")
	router.HandleFunc("/posts/{id}", getPostHandler).Methods("GET")
	router.HandleFunc("/posts/{id}", updatePostHandler).Methods("PUT")
	router.HandleFunc("/posts/{id}", deletePostHandler).Methods("DELETE")
	fmt.Println("Listening to port 8000!")
	http.ListenAndServe(":8000", router)
}
