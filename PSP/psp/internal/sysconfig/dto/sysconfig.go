package dto

type ComputeTypeName struct {
	ComputeType string `json:"compute_type"` // 计算类型
	ShowName    string `json:"show_name"`    // 展示名称
}

type GetGlobalSysConfigRequest struct{}

type GetGlobalSysConfigResponse struct {
	EnableVisual bool               `json:"enable_visual"` // 开启可视化服务
	ComputeTypes []*ComputeTypeName `json:"compute_types"` // 计算类型列表
}

type GetJobBurstConfigRequest struct{}

type GetJobBurstConfigResponse struct {
	Enable    bool  `json:"enable"`    // 是否开启
	Threshold int64 `json:"threshold"` // 等待时间
}

type SetJobBurstConfigRequest struct {
	Enable    bool  `json:"enable" from:"enable"`       // 是否开启
	Threshold int64 `json:"threshold" from:"threshold"` // 等待时间
}

type SetJobBurstConfigResponse struct{}

type GetJobConfigRequest struct{}

type GetJobConfigResponse struct {
	Queue string `json:"queue"` // 队列名称
}

type SetJobConfigRequest struct {
	Queue string `json:"queue" from:"queue"` // 队列名称
}

type SetJobConfigResponse struct{}

type GetThreePersonConfigResponse struct {
	DefSafeUserID   string `json:"def_safe_user_id"`
	DefSafeUserName string `json:"def_safe_user_name"`
}

type SetThreePersonConfigRequest struct {
	DefSafeUserID   string `json:"def_safe_user_id"`
	DefSafeUserName string `json:"def_safe_user_name"`
}

type SetThreePersonConfigResponse struct{}
