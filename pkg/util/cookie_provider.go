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
	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		Path:     "/", // 경로를 명시적으로 '/'로 설정하는 것이 좋습니다.
	}

	if config.AppEnv == "production" {
		cookie.SameSite = http.SameSiteNoneMode
		cookie.Secure = true
		// cookie.Domain = "your-actual-domain.com" // 배포 시 실제 도메인 설정
	} else {
		// ★★★★★ 수정 포인트 2: 로컬 개발 환경일 때, 쿠키가 적용될 도메인을 명시적으로 설정합니다.
		cookie.Domain = "my-dev-domain.com"
		// SameSite=Lax와 Secure=false도 명시적으로 설정해주는 것이 안전합니다.
		cookie.SameSite = http.SameSiteLaxMode
		cookie.Secure = false
	}
	http.SetCookie(w, cookie)
	log.Println("쿠키 설정 완료")
	return nil
}
