package dto

type License struct {
	Name          string `json:"name"`           // 产品名称
	Version       string `json:"version"`        // 产品版本
	Expiry        string `json:"expiry"`         // 过期时间
	MachineID     string `json:"machine_id"`     // 机器 ID
	Key           string `json:"key"`            // 加密 key
	AvailableDays int    `json:"available_days"` // 可用天数
}

type GetMachineIDRequest struct{}

type GetMachineIDResponse struct {
	ID string `json:"id"` // 机器 ID
}

type GetLicenseInfoRequest struct{}

type GetLicenseInfoResponse struct {
	Name          string `json:"name"`           // 产品名称
	Version       string `json:"version"`        // 产品版本
	Expiry        string `json:"expiry"`         // 过期时间
	AvailableDays int    `json:"available_days"` // 可用天数
}

type UpdateLicenseInfoRequest struct {
	Name      string `json:"name" form:"name"`             // 产品名称
	Version   string `json:"version" form:"version"`       // 产品版本
	Expiry    string `json:"expiry" form:"expiry"`         // 过期时间
	MachineID string `json:"machine_id" form:"machine_id"` // 机器 ID
	Key       string `json:"key" form:"key"`               // 加密 key
}

type UpdateLicenseInfoResponse struct{}
