package usecase

import (
	"context"

	"github.com/Zepelown/Go_WebServer/internal/repository"
	"github.com/Zepelown/Go_WebServer/pkg/domain/entity"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/request"
)

type PostUsecase interface {
	LoadAllPosts(ctx context.Context) ([]*entity.Post, error)
	WritePost(ctx context.Context, request request.WritePostRequest) (id string, err error)
	LoadPost(ctx context.Context, id string) (*entity.Post, error)
}

type postUsecase struct {
	repo repository.PostRepository
}

func NewPostUsecase(r repository.PostRepository) PostUsecase {
	return &postUsecase{
		repo: r,
	}
}

func (u *postUsecase) LoadAllPosts(ctx context.Context) ([]*entity.Post, error) {
	posts, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
func (u *postUsecase) WritePost(ctx context.Context, request request.WritePostRequest) (id string, err error) {
	id, err = u.repo.Save(ctx, &entity.Post{
		Title:    request.Title,
		Content:  request.Content,
		Date:     request.Date,
		Category: request.Category,
		UserId:   request.UserId,
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (u *postUsecase) LoadPost(ctx context.Context, id string) (*entity.Post, error) {
	post, err := u.repo.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return post, nil
}
