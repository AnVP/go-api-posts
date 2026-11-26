package post

import (
	"errors"
	"go-api-posts/configs"
	"go-api-posts/internal/user"
	"go-api-posts/pkg/di"
	"go-api-posts/pkg/event"
	"go-api-posts/pkg/middleware"
	"go-api-posts/pkg/req"
	"go-api-posts/pkg/res"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type PostHandlerDeps struct {
	PostRepository *PostRepository
	UserRepository di.IUserRepository
	EventBus       *event.EventBus
	Config         *configs.Config
}
type PostHandler struct {
	PostRepository *PostRepository
	UserRepository di.IUserRepository
	EventBus       *event.EventBus
}

func NewPostHandler(router *http.ServeMux, deps PostHandlerDeps) {
	handler := &PostHandler{
		PostRepository: deps.PostRepository,
		UserRepository: deps.UserRepository,
		EventBus:       deps.EventBus,
	}
	router.HandleFunc("GET /post/{id}", handler.GetPostById())
	router.Handle("POST /post", middleware.IsAuthed(handler.Create(), deps.Config))
	router.Handle("DELETE /post/{id}", middleware.IsAuthed(handler.DeletePost(), deps.Config))
	router.Handle("PATCH /post/{id}", middleware.IsAuthed(handler.UpdatePost(), deps.Config))
	router.HandleFunc("GET /posts", handler.GetPosts())
}

func (handler *PostHandler) userFromContext(r *http.Request) (*user.User, error) {
	email, ok := r.Context().Value(middleware.ContextEmailKey).(string)
	if !ok {
		return nil, errors.New("email not found in context")
	}
	return handler.UserRepository.FindByEmail(email)
}

func (handler *PostHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[PostCreateRequest](&w, r)
		if err != nil {
			return
		}
		currentUser, err := handler.userFromContext(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		post := NewPost(body.Image, body.Description, currentUser.ID)
		createdPost, err := handler.PostRepository.Create(post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, createdPost, http.StatusOK)
	}
}

func (handler *PostHandler) GetPostById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		post, err := handler.PostRepository.GetById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		go handler.EventBus.Send(event.Event{
			Type: event.PostVisitedEvent,
			Data: post.ID,
		})
		res.Json(w, post, http.StatusOK)
	}
}

func (handler *PostHandler) DeletePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		post, err := handler.PostRepository.GetById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		currentUser, err := handler.userFromContext(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if post.UserID != currentUser.ID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		err = handler.PostRepository.Delete(post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, post, http.StatusOK)
	}
}

func (handler *PostHandler) UpdatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[PostUpdateRequest](&w, r)
		if err != nil {
			return
		}
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		post, err := handler.PostRepository.GetById(idString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		currentUser, err := handler.userFromContext(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if post.UserID != currentUser.ID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		updatedPost, err := handler.PostRepository.Update(&Post{
			Model:       gorm.Model{ID: uint(id)},
			Image:       body.Image,
			Description: body.Description,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, updatedPost, http.StatusOK)
	}
}

func (handler *PostHandler) GetPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		posts := handler.PostRepository.GetPosts(limit, offset)
		count := handler.PostRepository.Count()
		res.Json(w, GetPostsResponse{
			Posts: posts,
			Count: count,
		}, http.StatusOK)
	}
}
