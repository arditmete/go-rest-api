package handler

import (
	"awesomeProject/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type PostHandler struct {
	Service *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{postService}
}

func (h *PostHandler) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	posts, err := h.Service.GetPosts()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting posts: %s", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	post, err := h.Service.GetPost(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting post with Id: %s with error: %s", params["id"], err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	err := h.Service.CreatePost(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating post with error: %s", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Post is created successfully!")
}

func (h *PostHandler) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	err := h.Service.UpdatePost(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating post witht Id: %s with error: %s", params["id"], err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Post is updated successfully!")
}

func (h *PostHandler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	err := h.Service.DeletePost(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting post witht Id: %s with error: %s", params["id"], err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Post is deleted successfully!")
}
