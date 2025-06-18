package dto

type CommentItem struct {
	Id       string `json:"id"`
	Content  string `json:"content"`
	Date     string `json:"date"`
	Username string `json:"username"`
}
