package dto

type PostItem struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Date     string `json:"date"`
	Category string `json:"category"`
	Username string `json:"username"`
}
