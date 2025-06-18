package usecase

import (
	"context"

	"github.com/Zepelown/Go_WebServer/internal/repository"
	"github.com/Zepelown/Go_WebServer/pkg/domain/entity"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/dto"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/request"
)

type CommentUsecase interface {
	LoadAllCommentsInPost(ctx context.Context, postId string) ([]*dto.CommentItem, error)
	WriteComment(ctx context.Context, userId string, postId string, request request.WriteCommentRequest) (id string, err error)
}

type commentUsecase struct {
	postRepo    repository.PostRepository
	userRepo    repository.UserRepository
	commentRepo repository.CommentRepository
}

func NewCommentUsecase(
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
	commentRepo repository.CommentRepository) CommentUsecase {
	return &commentUsecase{
		postRepo:    postRepo,
		userRepo:    userRepo,
		commentRepo: commentRepo,
	}
}

func (u *commentUsecase) LoadAllCommentsInPost(ctx context.Context, postId string) ([]*dto.CommentItem, error) {
	comments, err := u.commentRepo.GetAllCommentInPost(ctx, postId)
	if err != nil {
		return nil, err
	}
	commentItems := make([]*dto.CommentItem, 0, len(comments))
	for _, comment := range comments {
		user, _ := u.userRepo.FindById(ctx, comment.UserId)
		commentItem := comment.CommentToCommentItem(user)
		commentItems = append(commentItems, &commentItem)
	}
	return commentItems, nil
}

func (u *commentUsecase) WriteComment(ctx context.Context, userId string, postId string, request request.WriteCommentRequest) (id string, err error) {
	id, err = u.commentRepo.Save(ctx, &entity.Comment{
		Content: request.Content,
		Date:    request.Date,
		UserId:  userId,
		PostId:  postId,
	})
	if err != nil {
		return "", err
	}
	return id, nil
}
