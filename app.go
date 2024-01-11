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

	type MiddlewareFunc func(http.Handler) http.Handler
	// ConvertToMiddlewareFunc converts UserRepository's AuthMiddleware to MiddlewareFunc

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/posts")
	if err != nil {
		_ = fmt.Errorf("Connection to database is refused")
		panic(err.Error())
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		_ = fmt.Errorf("Failed to ping database")
		panic(err.Error())
	}
	//create tables if not exist
	err = repository.InitializeDatabase(db)
	if err != nil {
		panic(err.Error())
	}

	postRepository := repository.NewPostRepository(db)
	userRepository := repository.NewUserRepository(db)
	postService := service.NewPostService(postRepository)
	userService := service.NewUserService(userRepository)
	postHandler := handler.NewPostHandler(postService)
	authHandler := &handler.AuthHandler{}
	loginHandler := handler.NewLoginHandler(userService, authHandler)
	userHandler := handler.NewUserHandler(userService)

	// Initialize handlers with the service and repository dependencies
	getPostsHandler := postHandler.GetPostsHandler
	createPostHandler := postHandler.CreatePostHandler
	getPostHandler := postHandler.GetPostHandler
	updatePostHandler := postHandler.UpdatePostHandler
	deletePostHandler := postHandler.DeletePostHandler
	loginUserHandler := loginHandler.LoginUserHandler
	getUsersHandler := userHandler.GetUsersHandler
	createUserHandler := userHandler.CreateUserHandler
	getUserHandler := userHandler.GetUserHandler
	updateUserHandler := userHandler.UpdateUserHandler
	deleteUserHandler := userHandler.DeleteUserHandler

	router := mux.NewRouter()

	router.HandleFunc("/login", loginUserHandler).Methods("POST")
	//posts
	postsRouter := router.PathPrefix("/posts").Subrouter()
	postsRouter.Use(ConvertToMiddlewareFunc(authHandler))
	postsRouter.HandleFunc("", getPostsHandler).Methods("GET")
	postsRouter.HandleFunc("", createPostHandler).Methods("POST")
	postsRouter.HandleFunc("/{id}", getPostHandler).Methods("GET")
	postsRouter.HandleFunc("/{id}", updatePostHandler).Methods("PUT")
	postsRouter.HandleFunc("/{id}", deletePostHandler).Methods("DELETE")
	//users
	usersRouter := router.PathPrefix("/users").Subrouter()
	postsRouter.Use(ConvertToMiddlewareFunc(authHandler))
	usersRouter.HandleFunc("/", getUsersHandler).Methods("GET")
	usersRouter.HandleFunc("/", createUserHandler).Methods("POST")
	usersRouter.HandleFunc("//{id}", getUserHandler).Methods("GET")
	usersRouter.HandleFunc("//{id}", updateUserHandler).Methods("PUT")
	usersRouter.HandleFunc("//{id}", deleteUserHandler).Methods("DELETE")
	fmt.Println("Listening to port 8000!")
	http.ListenAndServe(":8000", router)
}

func ConvertToMiddlewareFunc(authHandler *handler.AuthHandler) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return authHandler.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}
