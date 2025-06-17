package middleware // 미들웨어는 별도 패키지로 관리하는 것이 좋습니다.

import (
	"net/http"

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
		// 1. 요청에서 'token' 쿠키를 가져옵니다.
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// 쿠키가 없는 경우, 401 Unauthorized 응답
				http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
				return
			}
			// 다른 종류의 에러인 경우 (예: 잘못된 요청)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// 2. 쿠키 값(토큰 문자열)을 가져옵니다.
		tokenString := c.Value
		claims := &dto.Claims{} // 기존에 사용하시던 Claims 구조체

		// 3. 토큰을 파싱하고 유효성을 검증합니다.
		// jwt.ParseWithClaims는 서명, 만료 시간(exp) 등을 모두 검증해줍니다[8].
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
