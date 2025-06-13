package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
	Date     string             `bson:"date"`
	Category string             `bson:"category"`
	UserId   string             `bson:"userId"`
}
