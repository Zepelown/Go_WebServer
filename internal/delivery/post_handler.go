package handler

import (
	"encoding/json"
	"net/http" // 웹 서버 기능

	// "strings" // 더 이상 필요하지 않으므로 제거할 수 있습니다.

	"github.com/Zepelown/Go_WebServer/internal/usecase"
	"github.com/Zepelown/Go_WebServer/pkg/appcontext"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/request"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/response"
)

type PostHandler struct {
	postUsecase    usecase.PostUsecase
	userUsecase    usecase.UserUsecase
	commentUsecase usecase.CommentUsecase
}

func NewPostHandler(
	postUC usecase.PostUsecase,
	userUC usecase.UserUsecase,
	commentUC usecase.CommentUsecase) *PostHandler {
	return &PostHandler{
		postUsecase:    postUC,
		userUsecase:    userUC,
		commentUsecase: commentUC,
	}
}

// LoadAllPosts 핸들러 (HTTP 메서드 확인 제거)
func (h *PostHandler) LoadAllPosts(w http.ResponseWriter, r *http.Request) {
	// main.go에서 "GET"으로 라우팅되므로 메서드 확인이 필요 없습니다.
	defer r.Body.Close()

	posts, err := h.postUsecase.LoadAllPosts(r.Context())
	if err != nil {
		http.Error(w, "불러오기가 실패하였습니다", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.PostLoadAllReponse{Posts: posts})
}

// WritePost 핸들러 (HTTP 메서드 확인 제거)
func (h *PostHandler) WritePost(w http.ResponseWriter, r *http.Request) {
	// main.go에서 "POST"로 라우팅되므로 메서드 확인이 필요 없습니다.
	var req request.WritePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청 본문입니다", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	claims, ok := appcontext.GetUserClaims(r.Context())
	if !ok {
		http.Error(w, "Could not retrieve user info from context", http.StatusInternalServerError)
		return
	}

	id, err := h.postUsecase.WritePost(r.Context(), req, claims.Subject)
	if err != nil {
		http.Error(w, "잘못된 형식입니다", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.WritePostResponse{ID: id})
}

// LoadOnePost 핸들러 (r.PathValue 사용 및 메서드 확인 제거)
func (h *PostHandler) LoadOnePost(w http.ResponseWriter, r *http.Request) {
	// main.go에서 "GET"으로 라우팅되므로 메서드 확인이 필요 없습니다.
	defer r.Body.Close()

	// 변경된 부분: r.PathValue를 사용하여 경로 변수 "id"를 가져옵니다.
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "게시글 ID가 필요합니다.", http.StatusBadRequest)
		return
	}

	post, err := h.postUsecase.LoadPost(r.Context(), id)
	if err != nil {
		http.Error(w, "불러오기가 실패하였습니다", http.StatusUnauthorized)
		return
	}

	comments, err := h.commentUsecase.LoadAllCommentsInPost(r.Context(), id)
	if err != nil {
		http.Error(w, "댓글 불러오기가 실패하였습니다", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.PostLoadOneReponse{Post: post, Comment: comments})
}

// WriteComment 핸들러 (수정 필요 없음 - 이미 올바르게 작성됨)
func (h *PostHandler) WriteComment(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("postId")
	if postId == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	var req request.WriteCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청 본문입니다", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	claims, ok := appcontext.GetUserClaims(r.Context())
	if !ok {
		http.Error(w, "Could not retrieve user info from context", http.StatusInternalServerError)
		return
	}

	id, err := h.commentUsecase.WriteComment(r.Context(), claims.Subject, postId, req)
	if err != nil {
		http.Error(w, "댓글 작성이 실패하였습니다", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"id": id,
	})
}
