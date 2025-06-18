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

type CommentRepository interface {
	Save(ctx context.Context, comment *entity.Comment) (id string, error error)
	GetAllCommentInPost(ctx context.Context, postId string) ([]*entity.Comment, error)
}
type mongoCommentRepository struct {
	collection *mongo.Collection
}

func NewMongoCommentRepository(collection *mongo.Collection) CommentRepository {
	return &mongoCommentRepository{
		collection: collection,
	}
}

func (r *mongoCommentRepository) Save(ctx context.Context, comment *entity.Comment) (id string, error error) {
	result, err := r.collection.InsertOne(ctx, comment)
	if err != nil {
		return "", err
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		// 타입 단언에 실패한 경우, 즉 ID가 ObjectID가 아닌 경우에 대한 처리입니다.
		return "", fmt.Errorf("failed to convert insertedID to primitive.ObjectID")
	}

	// ObjectID를 16진수 문자열로 변환하여 반환합니다.
	return oid.Hex(), nil

}

func (r *mongoCommentRepository) GetAllCommentInPost(ctx context.Context, postId string) ([]*entity.Comment, error) {
	// 	filter := bson.M{"email": email}

	// err := r.collection.FindOne(ctx, filter).Decode(&user)
	cursor, err := r.collection.Find(ctx, bson.M{"postId": postId})
	if err != nil {
		log.Printf("Error decoding documents: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var comments []*entity.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		log.Printf("Error decoding documents: %v", err)
		return nil, err
	}

	return comments, nil
}
