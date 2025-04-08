package v20230530

// swagger:model CommonResponse
type Response struct {
	// error code of response
	//
	// required: false
	ErrorCode string `json:"ErrorCode"`
	// error message of response
	//
	// required: false
	ErrorMsg string `json:"ErrorMsg"`

	RequestID string `json:"RequestID"`
}
