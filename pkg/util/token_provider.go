package util

import (
	"net/http"
	"time"

	"github.com/Zepelown/Go_WebServer/config"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/dto"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProvideToken(w http.ResponseWriter, config config.EnvConfig, userID primitive.ObjectID, username string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24)

	claims := &dto.Claims{
		Username: username, // DB에서 가져온 username
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   userID.Hex(), // DB에서 가져온 _id
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JwtSecretKey))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return "", err
	}
	return tokenString, nil
}
