package response

import "github.com/Zepelown/Go_WebServer/pkg/domain/payload/dto"

type PostLoadAllReponse struct {
	Posts []*dto.PostItem
}
