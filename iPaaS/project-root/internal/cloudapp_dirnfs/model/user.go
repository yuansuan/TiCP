package model

type AddUserRequest struct {
	Username      string `uri:"username"`
	SubPath       string `json:"sub_path"`
	Password      string `json:"password"`
	ExcludeUserID bool   `json:"exclude_user_id"`
}

type DeleteUserRequest struct {
	Username string `uri:"username"`
}

type BaseResponse struct {
	ErrorMessage string `json:"error_message"`
	RequestId    string `json:"request_id"`
}
