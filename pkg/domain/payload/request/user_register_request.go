package request

type UserRegisterRequest struct {
	Email string `json:"email"`
	Password   string `json:"password"`
	Name  string `json:"name"`
}
