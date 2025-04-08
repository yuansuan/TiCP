package openapi

import "github.com/yuansuan/ticp/PSP/psp/pkg/xtype"

type PreUploadRequest struct {
	Path        string `form:"path" json:"path" validate:"required"`
	FileSize    int64  `form:"file_size" json:"file_size" validate:"required"`
	ComputeType string `json:"compute_type" form:"compute_type" validate:"required,oneof=local cloud"` // 计算类型
	IsTemp      bool   `form:"is_temp" json:"is_temp"`
}

type PreUploadResponse struct {
	UploadId string `json:"upload_id"`
}

type UploadRequest struct {
	UploadID    string `form:"upload_id" validate:"required"`
	Path        string `form:"path" validate:"required"`
	FileSize    int64  `form:"file_size" validate:"required"`
	Offset      int64  `form:"offset"`
	SliceSize   int64  `form:"slice_size"`
	Finish      bool   `form:"finish"`
	IsTemp      bool   `form:"is_temp"`
	ComputeType string `json:"compute_type" form:"compute_type" validate:"required,oneof=local cloud"` // 计算类型
}

type ListRequest struct {
	Path             string      `json:"path" validate:"required"`
	Page             *xtype.Page `json:"page"`
	IsTemp           bool        `json:"is_temp"`
	ShowHideFile     bool        `json:"show_hide_file"`
	ComputeType      string      `json:"compute_type" form:"compute_type" validate:"required,oneof=local cloud"` // 计算类型
	FilterRegexpList []string    `json:"filter_regexp_list"`
}

type ListResponse struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	IsDir bool   `json:"is_dir"`
	Path  string `json:"path"`
}

type RemoveRequest struct {
	Paths       []string `json:"paths" validate:"required"`
	IsTemp      bool     `json:"is_temp"`
	ComputeType string   `json:"compute_type" form:"compute_type" validate:"required,oneof=local cloud"` // 计算类型
}

type BatchDownloadPreRequest struct {
	FilePaths   []string `json:"file_paths" form:"file_paths" validate:"required"`                       // 文件夹/文件路径(单个文件直接全路径)
	FileName    string   `json:"file_name" form:"file_name" validate:"required"`                         // 文件名称(以.zip结尾, 但是单个文件的时候直接原始文件名)
	IsTemp      bool     `json:"is_temp" form:"is_temp"`                                                 // 是否跨越用户目录
	IsCompress  bool     `json:"is_compress" form:"is_compress"`                                         // 是否压缩
	ComputeType string   `json:"compute_type" form:"compute_type" validate:"required,oneof=local cloud"` // 计算类型
}

type BatchDownloadRequest struct {
	Token       string `json:"token" form:"token" validate:"required"`                                 // 批量下载token
	ComputeType string `json:"compute_type" form:"compute_type" validate:"required,oneof=local cloud"` // 计算类型
}

func ConvertIsCloud(computeType string) bool {
	return computeType == "cloud"
}
