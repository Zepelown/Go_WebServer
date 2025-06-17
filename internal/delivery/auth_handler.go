package handler // 핸들러도 별도 패키지로 관리

import (
	"encoding/json"
	"net/http"

	"github.com/Zepelown/Go_WebServer/pkg/appcontext"
)

func AuthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// 여기! appcontext의 헬퍼 함수를 사용하여 컨텍스트에서 값을 가져옵니다.
	claims, ok := appcontext.GetUserClaims(r.Context())
	if !ok {
		// 미들웨어를 통과했다면 이 에러는 거의 발생하지 않지만, 안전장치로 둡니다.
		http.Error(w, "Could not retrieve user info from context", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"username": claims.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
