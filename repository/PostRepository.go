package repository

import (
	"awesomeProject/model"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type PostRepository struct {
	DB *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (r *PostRepository) GetPosts() ([]model.Post, error) {
	result, err := r.DB.Query("SELECT id, title from posts")
	var posts []model.Post
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var post model.Post
		err := result.Scan(&post.ID, &post.Title)
		if err != nil {
			panic(err.Error())
		}
		posts = append(posts, post)
	}
	return posts, err
}

func (r *PostRepository) GetPost(w http.ResponseWriter, rr *http.Request) (model.Post, error) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(rr)
	result, err := r.DB.Query("SELECT id, title FROM posts WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var post model.Post
	if result.Next() == false {
		return model.Post{}, errors.New("Post not found")
	}
	err = result.Scan(&post.ID, &post.Title)
	if err != nil {
		panic(err.Error())
	}
	return post, err
}

func (r *PostRepository) CreatePost(w http.ResponseWriter, rr *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := r.DB.Prepare("INSERT INTO posts(title) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]
	err = r.DB.QueryRow("SELECT title FROM posts WHERE title = ?", title).Scan(&title)
	if err == nil {
		http.Error(w, "Title already exists", http.StatusConflict)
		return err
	} else if err != sql.ErrNoRows {
		panic(err.Error())
	}
	_, err = stmt.Exec(title)
	if err != nil {
		panic(err.Error())
	}
	return err
}

func (r *PostRepository) DeletePost(w http.ResponseWriter, rr *http.Request) error {
	params := mux.Vars(rr)
	_, err := r.getPostByID(params["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving post: %s", err), http.StatusInternalServerError)
		return err
	}
	stmt, err := r.DB.Prepare("DELETE FROM posts WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	return err
}

func (r *PostRepository) UpdatePost(w http.ResponseWriter, rr *http.Request) error {
	params := mux.Vars(rr)
	_, err := r.getPostByID(params["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving post: %s", err), http.StatusInternalServerError)
		return err
	}
	stmt, err := r.DB.Prepare("UPDATE posts SET title = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newTitle := keyVal["title"]
	_, err = stmt.Exec(newTitle, params["id"])
	if err != nil {
		panic(err.Error())
	}
	return err
}

func (r *PostRepository) getPostByID(postID string) (model.Post, error) {
	var post model.Post
	result, err := r.DB.Query("SELECT id, title FROM posts WHERE id = ?", postID)
	if err != nil {
		return post, err
	}
	defer result.Close()
	if result.Next() == false {
		return post, errors.New("Post not found")
	}
	err = result.Scan(&post.ID, &post.Title)
	if err != nil {
		return post, err
	}
	return post, nil
}
