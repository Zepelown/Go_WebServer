package util

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Zepelown/Go_WebServer/config"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/dto"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProvideTokenCookie(w http.ResponseWriter, config config.EnvConfig, userID primitive.ObjectID, username string) error {
	expirationTime := time.Now().Add(5 * time.Minute)
	// 클레임 생성
	claims := &dto.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   userID.Hex(),
		},
	}

	// 클레임을 사용하여 토큰 생성 (알고리즘: HS256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 비밀 키로 토큰 서명
	tokenString, err := token.SignedString([]byte(config.JwtSecretKey))
	if err != nil {
		return fmt.Errorf("토큰 서명 실패: %w", err)
	}

	// 생성된 토큰을 클라이언트 쿠키에 설정
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		Path:     "/",
	})
	log.Println("쿠키 설저 완료")
	return nil
}
