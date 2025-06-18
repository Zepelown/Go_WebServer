package handler

import (
	"encoding/json"
	"net/http" // 웹 서버 기능

	"github.com/Zepelown/Go_WebServer/config"
	"github.com/Zepelown/Go_WebServer/internal/usecase"
	"github.com/Zepelown/Go_WebServer/pkg/appcontext"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/request"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/response"
	"github.com/Zepelown/Go_WebServer/pkg/util"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: uc}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request, config config.EnvConfig) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST 메서드만 허용됩니다", http.StatusMethodNotAllowed)
		return
	}
	var req request.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청 본문입니다", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	// Usecase를 호출하여 로그인 로직 실행 및 토큰 발급
	user, err := h.usecase.Login(r.Context(), req)
	if err != nil {
		http.Error(w, "이메일 또는 비밀번호가 잘못되었습니다", http.StatusUnauthorized)
		return
	}
	// 수정된 함수를 호출하고 에러를 확인합니다.
	// if err := util.ProvideTokenCookie(w, config, user.ID, user.Name); err != nil {
	// 	log.Printf("쿠키 제공 실패: %v\n", err)
	// 	http.Error(w, "서버 내부 오류로 토큰을 발급하지 못했습니다", http.StatusInternalServerError)
	// 	return // 에러 발생 시 핸들러를 확실히 종료합니다.
	// }
	tokenString, err := util.ProvideToken(w, config, user.ID, user.Name)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
	// json.NewEncoder(w).Encode(response.LoginResponse{User: *user})
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST 메서드만 허용됩니다", http.StatusMethodNotAllowed)
		return
	}

	var req request.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청 본문입니다", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	success, _ := h.usecase.Register(r.Context(), req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK 상태 코드 반환

	json.NewEncoder(w).Encode(response.UserRegisterResponse{Success: success})

}

func (h *UserHandler) FindUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST 메서드만 허용됩니다", http.StatusMethodNotAllowed)
		return
	}
	var req request.UserFindByIdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청 본문입니다", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	user, err := h.usecase.FindById(r.Context(), req.Id)
	if err != nil {
		http.Error(w, "유저를 찾지 못했습니다.", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK 상태 코드 반환

	json.NewEncoder(w).Encode(response.UserFindByIdReponse{User: user})
}

func (h *UserHandler) FindUserByToken(w http.ResponseWriter, r *http.Request) {
	claims, ok := appcontext.GetUserClaims(r.Context())
	if !ok {
		// 미들웨어를 통과했다면 이 에러는 거의 발생하지 않지만, 안전장치로 둡니다.
		http.Error(w, "Could not retrieve user info from context", http.StatusInternalServerError)
		return
	}
	user, err := h.usecase.FindById(r.Context(), claims.Subject)
	if err != nil {
		http.Error(w, "유저를 찾지 못했습니다.", http.StatusUnauthorized)
		return
	}
	response := map[string]string{
		"username": user.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
