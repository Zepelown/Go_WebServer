package usecase

import (
	"context"
	"errors"
	"log"

	"github.com/Zepelown/Go_WebServer/internal/repository"
	"github.com/Zepelown/Go_WebServer/pkg/domain/entity"
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/request"
)

type UserUsecase interface {
	Login(ctx context.Context, request request.UserLoginRequest) (*entity.User, error)
	Register(ctx context.Context, request request.UserRegisterRequest) (bool, error)
	FindById(ctx context.Context, id string) (*entity.User, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) UserUsecase {
	return &userUsecase{
		repo: r,
	}
}

func (uc *userUsecase) Login(ctx context.Context, request request.UserLoginRequest) (*entity.User, error) {
	user, error := uc.repo.FindByEmail(ctx, request.Email)
	if error != nil || user == nil {
		return nil, error
	}
	return user, nil
}

func (uc *userUsecase) Register(ctx context.Context, request request.UserRegisterRequest) (bool, error) {
	exists, err := uc.repo.IsExistEmail(ctx, request.Email)
	if err != nil {
		log.Println("이메일 존재 여부 확인 중 DB 에러 발생:", err)
		return false, err // DB 에러를 그대로 위로 전달
	}

	// 3. 이메일이 이미 존재하는지 확인합니다.
	if exists {
		return false, ErrUserAlreadyExists // "이미 존재하는 이메일입니다" 에러 반환
	}

	err = uc.repo.Save(ctx, &entity.User{Email: request.Email, Password: request.Password, Name: request.Name})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (uc *userUsecase) FindById(ctx context.Context, id string) (*entity.User, error) {
	user, err := uc.repo.FindById(ctx, id)
	if err != nil {
		log.Println("유저를 찾지 못했습니다.", err)
		return nil, err // DB 에러를 그대로 위로 전달
	}
	return user, nil
}

var ErrUserAlreadyExists = errors.New("이미 존재하는 이메일입니다")
