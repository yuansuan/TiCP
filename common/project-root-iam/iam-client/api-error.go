package iam_client

import (
	"net/http"

	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
)

// 404 error
var ErrNotFound = ErrorResponse{
	Message: "secret not found in Cache",
	Code:    "AccessKeyNotFound",
	Status:  http.StatusNotFound,
}

type ErrorResponse struct {
	Message string `json:"Message"`
	Code    string `json:"Code"`
	Status  int    `json:"Status"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func httpRespToErrorResponse(code int, data *iam_api.BasicResponse) error {
	err := ErrorResponse{
		Message: data.ErrorMessage,
		Code:    data.ErrorCode,
		Status:  code,
	}
	return err
}
