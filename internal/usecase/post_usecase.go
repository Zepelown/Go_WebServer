package usecase

import (
	"context"

	"github.com/Zepelown/Go_WebServer/internal/repository"
	"github.com/Zepelown/Go_WebServer/pkg/domain/entity"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/dto"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/request"
)

type PostUsecase interface {
	LoadAllPosts(ctx context.Context) ([]*dto.PostItem, error)
	WritePost(ctx context.Context, request request.WritePostRequest) (id string, err error)
	LoadPost(ctx context.Context, id string) (*dto.PostItem, error)
}

type postUsecase struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

func NewPostUsecase(postRepo repository.PostRepository, userRepo repository.UserRepository) PostUsecase {
	return &postUsecase{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

func (u *postUsecase) LoadAllPosts(ctx context.Context) ([]*dto.PostItem, error) {
	posts, err := u.postRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	postItems := make([]*dto.PostItem, 0, len(posts))

	for _, post := range posts {
		user, _ := u.userRepo.FindById(ctx, post.UserId)
		postItem := post.PostToPostItem(user)
		postItems = append(postItems, &postItem)
	}
	return postItems, nil
}

func (u *postUsecase) WritePost(ctx context.Context, request request.WritePostRequest) (id string, err error) {
	id, err = u.postRepo.Save(ctx, &entity.Post{
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

func (u *postUsecase) LoadPost(ctx context.Context, id string) (*dto.PostItem, error) {
	post, err := u.postRepo.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	user, _ := u.userRepo.FindById(ctx, post.UserId)

	postItem := post.PostToPostItem(user)

	return &postItem, nil
}
