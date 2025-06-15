package request

type WritePostRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Date     string `json:"date"`
	Category string `json:"category"`
	UserId   string `json:"userId"`
}
