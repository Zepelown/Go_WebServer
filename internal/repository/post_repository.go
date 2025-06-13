package repository

import (
	"context"
	"log"

	"github.com/Zepelown/Go_WebServer/pkg/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository interface {
	Save(ctx context.Context, post *entity.Post) error
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

func (r *mongoPostRepository) Save(ctx context.Context, post *entity.Post) error {
	_, err := r.collection.InsertOne(ctx, post)
	return err
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
