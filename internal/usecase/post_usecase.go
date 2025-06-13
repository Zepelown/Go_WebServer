package usecase

import (
	"context"

	"github.com/Zepelown/Go_WebServer/internal/repository"
	"github.com/Zepelown/Go_WebServer/pkg/domain/entity"
)

type PostUsecase interface {
	LoadAllPosts(ctx context.Context) ([]*entity.Post, error)
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
