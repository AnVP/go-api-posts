package middleware

import (
	"context"
	"go-api-posts/configs"
	"go-api-posts/pkg/jwt"
	"net/http"
	"strings"
)

type key string

const (
	ContextEmailKey = key("email")
)

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		isValid, data := jwt.NewJWT(config.Auth.Token).Parse(token)
		if !isValid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
