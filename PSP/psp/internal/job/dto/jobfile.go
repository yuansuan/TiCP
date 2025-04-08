package dto

// CreateTempDirRequest 创建作业临时目录请求数据
type CreateTempDirRequest struct {
	ComputeType string `json:"compute_type"` // 计算类型(local, cloud)
}

// CreateTempDirResponse 创建作业临时目录响应数据
type CreateTempDirResponse struct {
	Path string `json:"path"` // 创建作业的临时目录
}

// GetWorkSpaceResponse 获取工作空间响应数据
type GetWorkSpaceResponse struct {
	Name string `json:"name"` // 工作空间名称
}
