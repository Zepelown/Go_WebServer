package main

// JSON 데이터 처리
import (
	"context"
	"fmt"
	"net/http"

	config "github.com/Zepelown/Go_WebServer/config"
	handler "github.com/Zepelown/Go_WebServer/internal/delivery"
	"github.com/Zepelown/Go_WebServer/internal/repository"
	"github.com/Zepelown/Go_WebServer/internal/usecase"
)

func main() {
	// 1. MongoDB 연결 설정
	client := config.InitMongoDb()
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// 2. 사용할 데이터베이스와 컬렉션 가져오기
	userCollection := client.Database("webJungleDB").Collection("users")
	userRepo := repository.NewMongoUserRepository(userCollection)

	userUsecase := usecase.NewUserUsecase(userRepo)

	userHandler := handler.NewUserHandler(userUsecase)

	postCollection := client.Database("webJungleDB").Collection("posts")
	postRepo := repository.NewMongoPostRepository(postCollection)
	postUsecase := usecase.NewPostUsecase(postRepo)
	postHandler := handler.NewPostHandler(postUsecase)

	mux := http.NewServeMux()

	mux.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		userHandler.Login(w, r)
	})
	mux.HandleFunc("/users/register", func(w http.ResponseWriter, r *http.Request) {
		userHandler.Register(w, r)
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
	server := config.CorsMiddleware(mux)
	list_err := http.ListenAndServe(":8080", server)
	if list_err != nil {
		panic(list_err)
	}
}
