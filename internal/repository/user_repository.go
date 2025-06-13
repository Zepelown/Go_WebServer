package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	entity "github.com/Zepelown/Go_WebServer/pkg/domain/entity"
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	IsExistEmail(ctx context.Context, email string) (bool, error)
}

type mongoUserRepository struct {
	collection *mongo.Collection // 특정 컬렉션에 대한 포인터
}

func NewMongoUserRepository(collection *mongo.Collection) UserRepository {
	return &mongoUserRepository{
		collection: collection,
	}
}

func (r *mongoUserRepository) Save(ctx context.Context, user *entity.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *mongoUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	// bson.M은 map[string]interface{}과 유사한 타입으로, 쿼리 문서를 만들 때 사용합니다.
	filter := bson.M{"email": email}

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepository) IsExistEmail(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"email": email}
	err := r.collection.FindOne(ctx, filter).Err()
	if err == nil {
		return true, nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}

	return false, err
}
