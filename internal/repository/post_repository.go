package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/Zepelown/Go_WebServer/pkg/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository interface {
	Save(ctx context.Context, post *entity.Post) (id string, error error)
	GetAll(ctx context.Context) ([]*entity.Post, error)
}

type mongoPostRepository struct {
	collection *mongo.Collection
}

func NewMongoPostRepository(collection *mongo.Collection) PostRepository {
	return &mongoPostRepository{
		collection: collection,
	}
}

func (r *mongoPostRepository) Save(ctx context.Context, post *entity.Post) (string, error) {
	result, err := r.collection.InsertOne(ctx, post)
	if err != nil {
		return "", err
	}

	// InsertOne의 결과에서 InsertedID 필드를 가져옵니다.
	// 이 필드의 타입은 interface{} 이므로, 실제 타입인 primitive.ObjectID로 타입 단언(type assertion)을 해줘야 합니다.
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		// 타입 단언에 실패한 경우, 즉 ID가 ObjectID가 아닌 경우에 대한 처리입니다.
		return "", fmt.Errorf("failed to convert insertedID to primitive.ObjectID")
	}

	// ObjectID를 16진수 문자열로 변환하여 반환합니다.
	return oid.Hex(), nil
}

func (r *mongoPostRepository) GetAll(ctx context.Context) ([]*entity.Post, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		log.Printf("Error decoding documents: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var posts []*entity.Post
	if err = cursor.All(ctx, &posts); err != nil {
		log.Printf("Error decoding documents: %v", err)
		return nil, err
	}

	return posts, nil
}
