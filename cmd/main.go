package main

// JSON 데이터 처리
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	handler "github.com/Zepelown/Go_WebServer/internal/delivery"
	"github.com/Zepelown/Go_WebServer/internal/repository"
	"github.com/Zepelown/Go_WebServer/internal/usecase"
)

func main() {
	// 1. MongoDB 연결 설정
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// 2. 사용할 데이터베이스와 컬렉션 가져오기
	userCollection := client.Database("webJungleDB").Collection("users")
	userRepo := repository.NewMongoUserRepository(userCollection)

	userUsecase := usecase.NewUserUsecase(userRepo)

	userHandler := handler.NewUserHandler(userUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		userHandler.Login(w, r)
	})
	mux.HandleFunc("/users/register", func(w http.ResponseWriter, r *http.Request) {
		userHandler.Register(w, r)
	})

	fmt.Println("웹 서버가 8080 포트에서 실행됩니다.")

	list_err := http.ListenAndServe(":8080", mux)
	if list_err != nil {
		panic(list_err)
	}
}
