package service

import (
	"context"
	"mime/multipart"

	"github.com/gin-gonic/gin"

	"github.com/yuansuan/ticp/PSP/psp/internal/storage/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type FileService interface {
	// PreUpload 预上传
	PreUpload(ctx context.Context, req *dto.PreUploadRequest) (*dto.PreUploadResponse, error)
	// Upload 分片上传
	Upload(ctx *gin.Context, req *dto.UploadRequest, slice []*multipart.FileHeader) error
	// BatchDownload 批量下载
	BatchDownload(ctx *gin.Context, userName, token string, isCloud bool) error
	// Compress 压缩文件
	Compress(ctx *gin.Context, req *dto.CompressRequest) (*dto.CompressResponse, error)
	//CompressSubmit 压缩任务是否可以提交
	CompressSubmit(ctx *gin.Context) (bool, error)
	//CompressTasks 压缩任务的信息
	CompressTasks(ctx *gin.Context) ([]*dto.CompressTask, error)
	// BatchDownloadPre 批量预下载
	BatchDownloadPre(ctx *gin.Context, req *dto.BatchDownloadPreRequest) (*dto.BatchDownloadPreResponse, error)
	// Exist 文件是否存在
	Exist(ctx context.Context, userName string, cross bool, paths ...string) bool
	// CreateDir 创建文件夹
	CreateDir(ctx context.Context, userName, path string, cross bool) error
	// Remove 删除文件
	Remove(ctx context.Context, userName, path string, cross bool) error
	// Rename 文件重命名
	Rename(ctx context.Context, userName, path string, newPath string, overWrite bool, cross bool) error
	// Move 移动文件
	Move(ctx *gin.Context, userName string, cross, overwrite bool, dstdir string, srcpaths ...string) error
	// Get 查询单文件详情
	Get(ctx context.Context, userName, path string, cross bool) (file *dto.File, err error)
	// Read 读取文件内容
	Read(ctx context.Context, userName, path string, offset int64, len int64, cross bool) ([]byte, error)
	// List 文件列表
	List(ctx context.Context, userName, dir string, page *xtype.Page, cross bool, showHideFile bool, filterRegexps []string) (files []*dto.File, err error)
	// Write 写入文件内容
	Write(ctx *gin.Context, userName string, path string, fileSize int64, offset int64, sliceData []byte) error
	// Copy 复制文件到destPath目录
	Copy(ctx *gin.Context, req *dto.CopyRequest) (string, error)
	// ListOfRecur 递归的查询文件列表
	ListOfRecur(ctx context.Context, name string, paths []string, cross bool, showHideFile bool, filterRegexp []string) ([]*dto.File, error)
	// Mv 移动文件夹(支持重命名)
	Mv(ctx context.Context, name string, cross, overwrite bool, srcpath string, dstpath string) error
	// HpcDownload 下载文件到我的文件
	HpcDownload(ctx *gin.Context, req *dto.HpcDownloadRequest) error
	// Realpath 获取文件真实路径
	Realpath(ctx context.Context, relativePath string) (string, error)
	// QueryUploadHPCFileTask 查询需要上传的hpc文件夹与文件
	QueryUploadHPCFileTask(ctx context.Context, srcFilePaths, srcDirPaths []string, userName string, cross bool) (*dto.QueryHPCFileTaskResponse, error)
	// UploadHPCFile 上传HPC文件
	UploadHPCFile(ctx context.Context, req *dto.UploadHPCRequest, needUploadHPCFile *dto.QueryHPCFileTaskResponse) (string, error)
	// CancelUploadHPCFileTask 取消上传HPC文件任务
	CancelUploadHPCFileTask(ctx *gin.Context, name string, id string) error
	// GetUploadHPCFileTask 查询上传HPC文件任务查询上传HPC文件任务列表列表
	GetUploadHPCFileTask(ctx context.Context, taskKey string) ([]*dto.HPCUploadTaskResponse, error)
	// ResumeUploadHPCFileTask 恢复上传HPC文件任务
	ResumeUploadHPCFileTask(ctx *gin.Context, key string, id string) error
	// AbortAllTask 终止指定任务集的所有任务
	AbortAllTask(ctx *gin.Context, key string) error
	// GetCopyStatus 获取复制文件的状态
	GetCopyStatus(ctx *gin.Context, key string) (dto.CopyState, error)
	// HardLink 硬链接，如果传入文件夹，会先创建好目标文件夹，再给其中的文件做硬链接
	HardLink(ctx context.Context, req *dto.LinkRequest) error
	// Link 链接文件(夹)
	Link(ctx context.Context, req *dto.LinkRequest) (err error)
	// GenerateShareLink  生成分享链接并自动发送消息给被分享者
	GenerateShareLink(ctx *gin.Context, req dto.GenerateShareRequest) error
	// GetRecordList 获取分享记录列表
	GetRecordList(ctx *gin.Context, userID int64, req dto.GetRecordListRequest) (*dto.ShareRecordListResponse, error)
	// GetShareFile 获取分享文件
	GetShareFile(ctx *gin.Context, id snowflake.ID) (*dto.ShareFileInfo, error)
	// GetShareCount 获取分享消息数量
	GetShareCount(ctx *gin.Context, userId int64, share int8) (int64, error)
	// CheckUserHomePath 检查用户家目录是否创建
	CheckUserHomePath(ctx context.Context, userName string, cross bool) error
	// CheckSharePath  检查共享路径是否创建
	CheckSharePath(ctx context.Context, userName string, cross bool) error
	// UpdateRecordState 修改共享记录状态
	UpdateRecordState(ctx *gin.Context, userId int64, ids []string) error
	// SymLink 创建软连接
	SymLink(ctx context.Context, d *dto.SymLinkRequest) error
	// ReadAll 分享记录全部已处理
	ReadAll(ctx *gin.Context, id int64) error
	// CheckOnlyReadDir 检查文件是否是只读的(不能移动/删除/重命名)
	CheckOnlyReadDir(ctx *gin.Context, dir string, files []*dto.File) ([]*dto.File, error)
	// GetFileInfoByStat 查看单个文件信息
	GetFileInfoByStat(ctx context.Context, userName, path string, cross bool) (file *dto.File, err error)
}
