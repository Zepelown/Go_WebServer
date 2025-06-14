package handler

import (
	"encoding/json"
	"net/http" // 웹 서버 기능
	"strings"

	"github.com/Zepelown/Go_WebServer/internal/usecase"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/request"
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

func (h *PostHandler) WritePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST 메서드만 허용됩니다", http.StatusMethodNotAllowed)
		return
	}
	var req request.WritePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청 본문입니다", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	id, err := h.usecase.WritePost(r.Context(), req)
	if err != nil {
		http.Error(w, "잘못된 형식입니다", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK 상태 코드 반환
	json.NewEncoder(w).Encode(response.WritePostResponse{ID: id})
}

func (h *PostHandler) LoadOnePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "GET 메서드만 허용됩니다", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	id := strings.TrimPrefix(r.URL.Path, "/post/")
	if id == "" {
		// ID가 없는 경우, 예를 들어 /post/ 로만 요청이 온 경우
		http.Error(w, "게시글 ID가 필요합니다.", http.StatusBadRequest)
		return
	}

	post, err := h.usecase.LoadPost(r.Context(), id)
	if err != nil {
		http.Error(w, "불러오기가 실패하였습니다", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK 상태 코드 반환

	json.NewEncoder(w).Encode(response.PostLoadOneReponse{Post: post})
}
