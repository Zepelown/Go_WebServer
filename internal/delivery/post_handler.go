package handler

import (
	"encoding/json"
	"net/http" // 웹 서버 기능

	"github.com/Zepelown/Go_WebServer/internal/usecase"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/response"
)

type PostHandler struct {
	usecase usecase.PostUsecase
}

func NewPostHandler(uc usecase.PostUsecase) *PostHandler {
	return &PostHandler{usecase: uc}
}

func (h *PostHandler) LoadAllPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "GET 메서드만 허용됩니다", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	posts, err := h.usecase.LoadAllPosts(r.Context())
	if err != nil {
		http.Error(w, "불러오기가 실패하였습니다", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK 상태 코드 반환

	json.NewEncoder(w).Encode(response.PostLoadAllReponse{Posts: posts})
}
