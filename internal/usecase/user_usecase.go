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
}

// userUsecase 구조체는 인터페이스의 실제 구현체입니다.
// 비공개(소문자)로 선언하여 외부에서 직접 접근하지 못하게 합니다.
type userUsecase struct {
	// 데이터베이스 접근을 위한 리포지토리 '인터페이스'에 의존합니다.
	// 실제 DB가 MySQL인지, 메모리인지 등은 전혀 알 필요가 없습니다.
	repo repository.UserRepository
}

// NewUserUsecase는 userUsecase의 생성자 함수입니다.
// main.go에서 리포지토리 구현체를 주입받아 usecase 인스턴스를 생성합니다.
// 반환 타입은 '인터페이스'로 하여, 호출하는 쪽이 구현체에 의존하지 않도 합니다.
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

var ErrUserAlreadyExists = errors.New("이미 존재하는 이메일입니다")
