package response

import "github.com/Zepelown/Go_WebServer/pkg/domain/payload/dto"

type PostLoadOneReponse struct {
	Post    *dto.PostItem
	Comment []*dto.CommentItem
}
