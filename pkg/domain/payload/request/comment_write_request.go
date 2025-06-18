package request

type WriteCommentRequest struct {
	Content string `json:"content"`
	Date    string `json:"date"`
}
