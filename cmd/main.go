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
	"github.com/Zepelown/Go_WebServer/pkg/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
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
	postCollection := client.Database("webJungleDB").Collection("posts")
	commentCollection := client.Database("webJungleDB").Collection("comments")

	userRepo := repository.NewMongoUserRepository(userCollection)
	postRepo := repository.NewMongoPostRepository(postCollection)
	commentRepo := repository.NewMongoCommentRepository(commentCollection)

	userUsecase := usecase.NewUserUsecase(userRepo)
	postUsecase := usecase.NewPostUsecase(postRepo, userRepo)
	commentUsecase := usecase.NewCommentUsecase(postRepo, userRepo, commentRepo)

	userHandler := handler.NewUserHandler(userUsecase)
	postHandler := handler.NewPostHandler(postUsecase, userUsecase, commentUsecase)

	mux := http.NewServeMux()

	mux.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		userHandler.Login(w, r, config)
	})
	mux.HandleFunc("/users/register", func(w http.ResponseWriter, r *http.Request) {
		userHandler.Register(w, r)
	})
	mux.Handle("/users/auth", middleware.JwtAuthMiddleware(http.HandlerFunc(userHandler.FindUserByToken), config))
	mux.Handle("/main", middleware.JwtAuthMiddleware(http.HandlerFunc(postHandler.LoadAllPosts), config))
	mux.HandleFunc("POST /post", middleware.JwtAuthMiddleware(http.HandlerFunc(postHandler.WritePost), config).ServeHTTP)
	mux.HandleFunc("GET /post/{id}", middleware.JwtAuthMiddleware(http.HandlerFunc(postHandler.LoadOnePost), config).ServeHTTP)
	mux.HandleFunc("POST /post/{postId}/comment", middleware.JwtAuthMiddleware(http.HandlerFunc(postHandler.WriteComment), config).ServeHTTP)

	fmt.Println("웹 서버가 8080 포트에서 실행됩니다.")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		Debug:            true,
	})
	handler := c.Handler(mux)
	list_err := http.ListenAndServe(config.ServerPortUrl, handler)
	if list_err != nil {
		panic(list_err)
	}
}
