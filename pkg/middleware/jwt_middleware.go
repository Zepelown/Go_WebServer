package middleware // 미들웨어는 별도 패키지로 관리하는 것이 좋습니다.

import (
	"net/http"
	"strings"

	"github.com/Zepelown/Go_WebServer/config"
	"github.com/Zepelown/Go_WebServer/pkg/appcontext"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/dto"
	"github.com/golang-jwt/jwt/v5"
)

// 새로운 컨텍스트 키 타입 정의 (충돌 방지)
type contextKey string

const userContextKey = contextKey("userClaims")

func JwtAuthMiddleware(next http.Handler, config config.EnvConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. 요청 헤더에서 'Authorization' 값을 가져옵니다.
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		// 2. "Bearer <token>" 형식에서 토큰 부분만 추출합니다.
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}
		tokenString := headerParts[1]

		// 3. 토큰 파싱 및 검증 로직은 이전과 동일합니다.
		claims := &dto.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JwtSecretKey), nil
		})

		if err != nil || !token.Valid {
			// 토큰이 유효하지 않은 경우 (서명 불일치, 만료 등)
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}
		// 여기! appcontext의 헬퍼 함수를 사용하여 컨텍스트에 값을 저장합니다.
		ctx := appcontext.SetUserClaims(r.Context(), claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
