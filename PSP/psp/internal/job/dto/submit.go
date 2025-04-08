package dto

import "github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"

// JobSubmitRequest 作业提交请求数据
type JobSubmitRequest struct {
	AppID     string   `json:"app_id"`     // 应用 ID
	ProjectID string   `json:"project_id"` // 项目 ID
	MainFiles []string `json:"main_files"` // 主文件
	QueueName string   `json:"queue_name"` // 队列名称
	WorkDir   *WorkDir `json:"work_dir"`   // 工作目录
	Fields    []*Field `json:"fields"`     // 参数信息列表
}

// JobSubmitResponse 作业提交请求数据
type JobSubmitResponse struct{}

type ResubmitRequest struct {
	JobId string `json:"job_id"` // 作业 ID
}

type Extension struct {
	AppType  string `json:"app_type"`  // 应用类型
	UploadId string `json:"upload_id"` // 上传 ID
}

type ResubmitResponse struct {
	Param     *SubmitParam `json:"param"`     // 作业提交参数
	Extension *Extension   `json:"extension"` // 扩展参数
}

type SubmitParam struct {
	AppID     string       `json:"app_id"`     // 应用 ID
	ProjectID string       `json:"project_id"` // 项目 ID
	UserID    snowflake.ID `json:"user_id"`    // 用户 ID
	UserName  string       `json:"user_name"`  // 用户名称
	QueueName string       `json:"queue_name"` // 队列名称
	MainFiles []string     `json:"main_files"` // 主文件
	WorkDir   *WorkDir     `json:"work_dir"`   // 工作目录
	Fields    []*Field     `json:"fields"`     // 参数信息列表
	IsOpenApi bool         `json:"is_openapi"` // 是否为oepnapi提交
}

// WorkDir 工作目录信息
type WorkDir struct {
	Path   string `json:"path"`    // 工作目录
	IsTemp bool   `json:"is_temp"` // 是否临时目录
}

// Field 作业提交表单的字段信息
type Field struct {
	ID     string   `json:"id"`     // 参数 ID
	Type   string   `json:"type"`   // 参数类型
	Value  string   `json:"value"`  // 单个参数值
	Values []string `json:"values"` // 多个参数值
}

type JobSubmitParamInfo struct {
	Param *SubmitParam `json:"param"` // 作业提交参数
	Dirs  []string     `json:"dirs"`  // 上传的目录列表
	Files []string     `json:"files"` // 上传的文件列表
}
