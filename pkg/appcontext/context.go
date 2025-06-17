package appcontext

import (
	"context"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/dto"
)

// 1. 컨텍스트 키 타입을 정의합니다.
// string 같은 기본 타입을 직접 사용하면 다른 패키지에서 정의한 키와 충돌할 수 있으므로,
// 사용자 정의 타입을 만드는 것이 안전합니다[5].
type contextKey string

// 2. 컨텍스트 키를 대문자로 시작하여 다른 패키지에서 접근할 수 있도록 Export 합니다.
const UserContextKey = contextKey("userClaims")

// 3. 컨텍스트에 사용자 정보를 설정하는 헬퍼 함수
// context.WithValue는 새로운 context를 반환한다는 점을 기억해야 합니다[3].
func SetUserClaims(ctx context.Context, claims *dto.Claims) context.Context {
	return context.WithValue(ctx, UserContextKey, claims)
}

// 4. 컨텍스트에서 사용자 정보를 가져오는 헬퍼 함수
func GetUserClaims(ctx context.Context) (*dto.Claims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*dto.Claims)
	return claims, ok
}
