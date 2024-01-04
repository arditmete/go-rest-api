package service

import (
	"awesomeProject/model"
	"awesomeProject/repository"
	"net/http"
)

type PostService struct {
	Repository *repository.PostRepository
}

func NewPostService(repository *repository.PostRepository) *PostService {
	return &PostService{Repository: repository}
}

func (s *PostService) GetPosts() ([]model.Post, error) {
	return s.Repository.GetPosts()
}

func (s *PostService) GetPost(w http.ResponseWriter, r *http.Request) (model.Post, error) {
	return s.Repository.GetPost(w, r)
}

func (s *PostService) CreatePost(w http.ResponseWriter, r *http.Request) error {
	return s.Repository.CreatePost(w, r)
}

func (s *PostService) UpdatePost(w http.ResponseWriter, r *http.Request) error {
	return s.Repository.UpdatePost(w, r)
}

func (s *PostService) DeletePost(w http.ResponseWriter, r *http.Request) error {
	return s.Repository.DeletePost(w, r)
}
