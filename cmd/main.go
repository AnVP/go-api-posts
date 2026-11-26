package main

import (
	"fmt"
	"go-api-posts/configs"
	"go-api-posts/internal/auth"
	"go-api-posts/internal/post"
	"go-api-posts/internal/stat"
	"go-api-posts/internal/user"
	"go-api-posts/pkg/db"
	"go-api-posts/pkg/event"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	fmt.Println(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	postRepository := post.NewPostRepository(database)
	userRepository := user.NewUserRepository(database)
	statRepository := stat.NewStatRepository(database)

	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepository,
	})

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	post.NewPostHandler(router, post.PostHandlerDeps{
		PostRepository: postRepository,
		UserRepository: userRepository,
		EventBus:       eventBus,
		Config:         conf,
	})
	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statRepository,
		Config:         conf,
	})

	go statService.AddClick()

	httpServer := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server...")
	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Printf("Server err %v\n", err)
		return
	}
}
