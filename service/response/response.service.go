package responseService

import (
	"github.com/sonhnguyen/pcchecker/model/response"
)

func ResponseError(status int, msg error, code string) responseModel.Error {
	var errorResponse responseModel.Error
	errorResponse.Code = code
	errorResponse.Message = msg.Error()
	errorResponse.Status = status
	return errorResponse
}
