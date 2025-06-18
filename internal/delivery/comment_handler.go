package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Zepelown/Go_WebServer/internal/usecase"
	"github.com/Zepelown/Go_WebServer/pkg/appcontext"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/request"
)

type CommentHandler struct {
	postUsecase    usecase.PostUsecase
	userUsecase    usecase.UserUsecase
	commentUsecase usecase.CommentUsecase
}

func NewCommentHandler(
	postUC usecase.PostUsecase,
	userUC usecase.UserUsecase,
	commentUC usecase.CommentUsecase) *CommentHandler {
	return &CommentHandler{
		postUsecase:    postUC,
		userUsecase:    userUC,
		commentUsecase: commentUC,
	}
}

func (h *CommentHandler) LoadAllCommentsInPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST 메서드만 허용됩니다", http.StatusMethodNotAllowed)
		return
	}
	var req request.WriteCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청 본문입니다", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	_, ok := appcontext.GetUserClaims(r.Context())
	if !ok {
		// 미들웨어를 통과했다면 이 에러는 거의 발생하지 않지만, 안전장치로 둡니다.
		http.Error(w, "Could not retrieve user info from context", http.StatusInternalServerError)
		return
	}
}
