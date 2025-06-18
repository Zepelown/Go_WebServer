package entity

import (
	"github.com/Zepelown/Go_WebServer/pkg/domain/payload/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Content string             `bson:"content"`
	Date    string             `bson:"date"`
	UserId  string             `bson:"userId"`
	PostId  string             `bson:"postId"`
}

func (e *Comment) CommentToCommentItem(user *User) dto.CommentItem {
	var username string
	if user == nil || user.Name == "" {
		username = "알 수 없는 사용자" // 혹은 "탈퇴한 사용자" 등
	} else {
		username = user.Name
	}
	return dto.CommentItem{
		Id:       e.ID.Hex(),
		Content:  e.Content,
		Date:     e.Date,
		Username: username,
	}
}
