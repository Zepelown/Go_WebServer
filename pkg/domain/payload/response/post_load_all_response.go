package response

import "github.com/Zepelown/Go_WebServer/pkg/domain/entity"

type PostLoadAllReponse struct {
	Posts []*entity.Post
}
