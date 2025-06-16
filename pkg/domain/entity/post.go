package entity

import (
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
	Date     string             `bson:"date"`
	Category string             `bson:"category"`
	UserId   string             `bson:"userId"`
}

func (p *Post) PostToPostItem(user *User) dto.PostItem {
	var username string
	if user == nil || user.Name == "" {
		username = "알 수 없는 사용자" // 혹은 "탈퇴한 사용자" 등
	} else {
		username = user.Name
	}
	return dto.PostItem{
		Id:       p.ID.Hex(),
		Title:    p.Title,
		Content:  p.Content,
		Date:     p.Date,
		Category: p.Category,
		Username: username,
	}
}
