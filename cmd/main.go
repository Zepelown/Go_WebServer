package main

// JSON 데이터 처리
import (
	"context"
	"fmt"
	"log"
	"net/http"

	systemConfig "github.com/Zepelown/Go_WebServer/config"
	handler "github.com/Zepelown/Go_WebServer/internal/delivery"
	"github.com/Zepelown/Go_WebServer/internal/repository"
	"github.com/Zepelown/Go_WebServer/internal/usecase"
	"github.com/joho/godotenv"
	"github.com/sesaquecruz/go-env-loader/pkg/env"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	var config systemConfig.EnvConfig
	if err := env.LoadEnv(&config); err != nil {
		log.Fatalf("error loading environment variables: %v", err)
	}
	client := systemConfig.InitMongoDb(config.DbUrl)
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	userCollection := client.Database("webJungleDB").Collection("users")
	userRepo := repository.NewMongoUserRepository(userCollection)

	userUsecase := usecase.NewUserUsecase(userRepo)

	userHandler := handler.NewUserHandler(userUsecase)

	postCollection := client.Database("webJungleDB").Collection("posts")
	postRepo := repository.NewMongoPostRepository(postCollection)
	postUsecase := usecase.NewPostUsecase(postRepo, userRepo)
	postHandler := handler.NewPostHandler(postUsecase, userUsecase)

	mux := http.NewServeMux()

	mux.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		userHandler.Login(w, r)
	})
	mux.HandleFunc("/users/register", func(w http.ResponseWriter, r *http.Request) {
		userHandler.Register(w, r)
	})
	mux.HandleFunc("/users/find", func(w http.ResponseWriter, r *http.Request) {
		userHandler.FindUserById(w, r)
	})
	mux.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		postHandler.LoadAllPosts(w, r)
	})
	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		postHandler.WritePost(w, r)
	})
	mux.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		postHandler.LoadOnePost(w, r)
	})

	fmt.Println("웹 서버가 8080 포트에서 실행됩니다.")
	server := systemConfig.CorsMiddleware(mux)
	list_err := http.ListenAndServe(config.ServerPortUrl, server)
	if list_err != nil {
		panic(list_err)
	}
}
