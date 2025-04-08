package dto

import (
	"context"
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/util/lockutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type UploadRequest struct {
	UploadID  string `form:"upload_id"`
	Path      string `form:"path"`
	FileSize  int64  `form:"file_size"`
	Offset    int64  `form:"offset"`
	SliceSize int64  `form:"slice_size"`
	Finish    bool   `form:"finish"`
	Cross     bool   `form:"cross"`
	IsCloud   bool   `form:"is_cloud"`
	UserName  string `form:"user_name"`
}

type PreUploadRequest struct {
	Path     string `form:"path"`
	FileSize int64  `form:"file_size"`
	Cross    bool   `form:"cross"`
	IsCloud  bool   `form:"is_cloud"`
	UserName string `form:"user_name"`
}

type PreUploadResponse struct {
	UploadId string `json:"upload_id"`
}

type File struct {
	Name      string `json:"name"`
	Mode      string `json:"mode"`
	Size      int64  `json:"size"`
	MDate     int64  `json:"mdate"`
	Type      string `json:"type"`
	IsDir     bool   `json:"isdir"`
	IsSymLink bool   `json:"issymlink"`
	Path      string `json:"path"`
	IsText    bool   `json:"istext"`
	SubFile   []File `json:"sub_file"`
	OnlyRead  bool   `json:"only_read"`
}

type BatchDownloadPreRequest struct {
	FilePaths  []string `json:"file_paths" form:"file_paths"`   // 文件夹/文件路径(单个文件直接全路径)
	FileName   string   `json:"file_name" form:"file_name"`     // 文件名称(以.zip结尾, 但是单个文件的时候直接原始文件名)
	Cross      bool     `json:"cross" form:"cross"`             // 是否跨越用户目录
	IsCompress bool     `json:"is_compress" form:"is_compress"` // 是否压缩
	IsCloud    bool     `json:"is_cloud" form:"is_cloud"`       // 是否云端
	UserName   string   `json:"user_name" form:"user_name"`     // 用户名称
}

type BatchDownloadPreResponse struct {
	Token string `json:"token"` // 批量下载 Token
}

type CompressResponse struct {
	CompressID string `json:"compress_id"`
	TargetPath string `json:"target_path"`
	IsCloud    bool   `json:"is_cloud"`
}

type CompressTask struct {
	CompressId string
	Status     int8
	IsCloud    bool `json:"is_cloud"`
}

type DownloadCache struct {
	FilePaths  []string
	FileName   string
	Cross      bool
	IsCompress bool
	IsCloud    bool
	UserName   string
}

type CopyRequest struct {
	DstPath      string   `json:"dst_path" `
	Overwrite    bool     `json:"overwrite" `
	Cross        bool     `json:"cross"`
	IsCloud      bool     `json:"is_cloud" `
	UserName     string   `json:"user_name"`
	SrcFilePaths []string `json:"src_file_paths"` // 上传文件路径
	SrcDirPaths  []string `json:"src_dir_paths"`  // 上传文件夹路径
	CurrentPath  string   `json:"current_path"`   // 当前所在路径
}

type LinkRequest struct {
	DstPath      string   `json:"dst_path" `
	Overwrite    bool     `json:"overwrite" `
	Cross        bool     `json:"cross"`
	IsCloud      bool     `json:"is_cloud" `
	UserName     string   `json:"user_name"`
	SrcFilePaths []string `json:"src_file_paths"` // 上传文件路径
	SrcDirPaths  []string `json:"src_dir_paths"`  // 上传文件夹路径
	FilterPaths  []string `json:"filter_paths"`   // 过滤文件路径
	CurrentPath  string   `json:"current_path"`   // 当前所在路径
}

type SymLinkRequest struct {
	DstPath   string `json:"dst_path" `
	Overwrite bool   `json:"overwrite" `
	Cross     bool   `json:"cross"`
	IsCloud   bool   `json:"is_cloud" `
	UserName  string `json:"user_name"`
	SrcPath   string `json:"src_path"` // 上传文件路径
}

type MoveRequest struct {
	SrcPaths  []string `form:"src_paths" json:"src_paths"`
	DstPath   string   `form:"dst_path" json:"dst_path"`
	Overwrite bool     `form:"overwrite"  json:"overwrite"`
	Cross     bool     `form:"cross"  json:"cross"`
	IsCloud   bool     `form:"is_cloud" json:"is_cloud"`
	UserName  string   `form:"user_name" json:"user_name"`
}

type ListRequest struct {
	Path             string      `json:"path"`
	Page             *xtype.Page `json:"page"`
	Cross            bool        `json:"cross"`
	ShowHideFile     bool        `json:"show_hide_file"`
	IsCloud          bool        `json:"is_cloud"`
	UserName         string      `json:"user_name"`
	FilterRegexpList []string    `json:"filter_regexp_list"`
}

type ListResponse struct {
	Name     string `json:"name"`
	Mode     string `json:"mode"`
	Size     int64  `json:"size"`
	MDate    int64  `json:"m_date"`
	Type     string `json:"type"`
	IsDir    bool   `json:"is_dir"`
	Path     string `json:"path"`
	IsText   bool   `json:"is_text"`
	OnlyRead bool   `json:"only_read"`
}

type ListRecurRequest struct {
	Paths            []string `json:"paths"`
	Cross            bool     `json:"cross"`
	IsCloud          bool     `json:"is_cloud"`
	UserName         string   `json:"user_name"`
	ShowHideFile     bool     `json:"show_hide_file"`
	FilterRegexpList []string `json:"filter_regexp_list"`
}

type GetRequest struct {
	Paths    []string `json:"paths"`
	Cross    bool     `json:"cross"`
	IsCloud  bool     `json:"is_cloud"`
	UserName string   `json:"user_name"`
}

type RenameRequest struct {
	Path      string `json:"path"`
	NewPath   string `json:"newpath"`
	Overwrite bool   `json:"overwrite"`
	Cross     bool   `json:"cross"`
	IsCloud   bool   `json:"is_cloud"`
	UserName  string `json:"user_name"`
}

type CreateDirRequest struct {
	Path     string `json:"path"`
	Cross    bool   `json:"cross"`
	IsCloud  bool   `json:"is_cloud"`
	UserName string `json:"user_name"`
}

type RemoveRequest struct {
	Paths    []string `json:"paths"`
	Cross    bool     `json:"cross"`
	IsCloud  bool     `json:"is_cloud"`
	UserName string   `json:"user_name"`
}

type CompressRequest struct {
	SrcPaths []string `form:"src_paths" json:"src_paths"`
	DstPath  string   `form:"dst_path" json:"dst_path"`
	BasePath string   `form:"base_path" json:"base_path"`
	IsCloud  bool     `json:"is_cloud"`
}

type CompressTasksRequest struct {
	IsCloud bool `json:"is_cloud"`
}

type ReadRequest struct {
	Path     string `json:"path"`
	Offset   int64  `json:"offset"`
	Len      int64  `json:"len"`
	Cross    bool   `json:"cross"`
	IsCloud  bool   `json:"is_cloud"`
	UserName string `json:"user_name"`
}

type BatchDownloadRequest struct {
	Token   string `json:"token" form:"token"`       // 批量下载token
	IsCloud bool   `json:"is_cloud" form:"is_cloud"` // 是否云端
}

type HpcDownloadRequest struct {
	SrcFilePaths []string `json:"src_file_paths"` // 云端文件路径
	SrcDirPaths  []string `json:"src_dir_paths"`  // 云端文件夹路径
	DestDirPath  string   `json:"dest_dir_path"`  // 目标文件夹路径
	CurrentPath  string   `json:"current_path"`   // 当前路径
	Overwrite    bool     `json:"overwrite"`      // 是否覆盖
	UserName     string   `json:"user_name"`      // 用户名
}

type QueryHPCFileRequest struct {
	SrcFilePaths []string `json:"src_file_paths"` // 上传文件路径
	SrcDirPaths  []string `json:"src_dir_paths"`  // 上传文件夹路径
	UserName     string   `json:"user_name"`      // 用户名
}

type QueryHPCFileTaskResponse struct {
	DirTasks  []string
	FileTasks map[string]*HPCUploadTask
}

type UploadHPCRequest struct {
	SrcFilePaths []string `json:"src_file_paths"` // 上传文件路径
	SrcDirPaths  []string `json:"src_dir_paths"`  // 上传文件夹路径
	DestDirPath  string   `json:"dest_dir_path"`  // 目标文件夹路径
	CurrentPath  string   `json:"current_path"`   // 当前所在路径
	Overwrite    bool     `json:"overwrite"`      // 是否覆盖
	UserName     string   `json:"user_name"`      // 用户名
	Cross        bool     `json:"cross"`          // 是否跨越用户目录
}

type UploadState int8

const (
	UploadStateFailure   = UploadState(1) // 失败
	UploadStateUploading = UploadState(2) // 上传中
	UploadStatePending   = UploadState(3) // 等待
	UploadStateSuccess   = UploadState(4) // 成功
	UploadStateCancel    = UploadState(5) // 取消
)

type HPCUploadTask struct {
	FileName    string             `json:"file_name"`    // 文件名
	SrcPath     string             `json:"src_path"`     // 文件路径
	DestPath    string             `json:"dest_path"`    // 目标路径
	TotalSize   int64              `json:"total_size"`   // 总大小
	CurrentSize int64              `json:"current_size"` // 已上传大小
	State       UploadState        `json:"state"`        // 状态 0-未开始 -1-失败 1-上传成功 2-上传中
	Cancel      context.CancelFunc `json:"-"`
	CancelCtx   context.Context    `json:"-"`
	ErrMsg      string             `json:"err_msg"` // 错误信息
}

func (task *HPCUploadTask) CancelTask(taskKey string) {
	for {
		// 锁只是保险用，如果redis不可用，不阻塞业务
		successFlag, err := lockutil.TryLock(taskKey)
		if err != nil {
			successFlag = true
		}
		if successFlag {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	task.Cancel()
	lockutil.UnLock(taskKey)
}

func (task *HPCUploadTask) Update(taskKey string, state UploadState, errMsg string) {
	for {
		// 锁只是保险用，如果redis不可用，不阻塞业务
		successFlag, err := lockutil.TryLock(taskKey)
		if err != nil {
			successFlag = true
		}
		if successFlag {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	if state != 0 {
		task.State = state
	}
	if strutil.IsNotEmpty(errMsg) {
		task.ErrMsg = errMsg
	}
	lockutil.UnLock(taskKey)
}

type HPCUploadTaskResponse struct {
	UploadID    string      `json:"upload_id"`    // id
	FileName    string      `json:"file_name"`    // 文件名
	SrcPath     string      `json:"src_path"`     // 文件路径
	DestPath    string      `json:"dest_path"`    // 目标路径
	TotalSize   int64       `json:"total_size"`   // 总大小
	CurrentSize int64       `json:"current_size"` // 已上传大小
	State       UploadState `json:"state"`        // 状态 0-未开始 -1-失败 1-上传成功 2-上传中
	ErrMsg      string      `json:"err_msg"`      // 错误信息
}

type CancelUploadHPCFileTaskRequest struct {
	TaskKey  string `json:"task_key"`
	UploadID string `json:"upload_id"`
}

type ResumeUploadHPCFileTaskRequest struct {
	TaskKey  string `json:"task_key"`
	UploadID string `json:"upload_id"`
}

type CopyState int8

const (
	CopyStateFailure = CopyState(1) // 失败
	CopyStateCopying = CopyState(2) // 复制中
	CopyStateSuccess = CopyState(3) // 成功
)

type CompressState int8

const (
	CompressStateFailure     = CompressState(1)
	CompressStateCompressing = CompressState(2)
	CompressStateSuccess     = CompressState(3)
)

func (e CompressState) String() string {
	switch e {
	case CompressStateFailure:
		return "failed"
	case CompressStateCompressing:
		return "compressing"
	case CompressStateSuccess:
		return "success"
	default:
		return ""
	}
}

const (
	CompressDuration = time.Hour * 24
)

type ShareType int8

const (
	ShareTypeLink = ShareType(1)
	ShareTypeCopy = ShareType(2)
)

type GenerateShareRequest struct {
	ShareUserList []string  `json:"share_user_list"` // 分享用户
	ShareFilePath string    `json:"share_file_path"` // 分享文件路径
	ShareType     ShareType `json:"share_type"`      // 分享方式,1-复制 2-硬链接
	//ExpireTime    int64    `json:"expire_time"` // 到期时间
}

type ShareFileInfo struct {
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	MDate     int64  `json:"mdate"`
	Type      string `json:"type"`
	IsDir     bool   `json:"isdir"`
	Path      string `json:"path"`
	ShareType int8   `json:"share_type"`
}

type AddFileOperLogRequest struct {
	FileName    string `json:"file_name"`
	FileType    int8   `json:"file_type"`
	OperateType int8   `json:"operate_type"`
	StorageSize int64  `json:"storage_size"`
	FilePath    string `json:"file_path"`
}

type GetRecordListRequest struct {
	Page   *xtype.Page `json:"page"`
	Filter string      `json:"filter"`
}

type FileOperLogListRequest struct {
	Page      *xtype.Page `json:"page"`
	FileName  string      `json:"file_name"`
	StartTime int64       `json:"start_time"`
	EndTime   int64       `json:"end_time"`
}

type FileOperLogListResponse struct {
	Page *xtype.PageResp    `json:"page"`
	List []*FileOperLogInfo `json:"list"`
}

type FileOperLogInfo struct {
	ID          string       `json:"id"`
	FileName    string       `json:"file_name"`
	FilePath    string       `json:"file_path"`
	FileType    FileTypeEnum `json:"file_type"`
	OperateType OpTypeEnum   `json:"operate_type"`
	StorageSize int64        `json:"storage_size"`
	OperateTime time.Time    `json:"operate_time"`
}

type ShareRecordListResponse struct {
	Page *xtype.PageResp    `json:"page"`
	List []*ShareRecordInfo `json:"list"`
}

type ShareRecordInfo struct {
	Id        string    `json:"id"`
	Content   string    `json:"content"`    // 内容
	ShareTime time.Time `json:"share_time"` // 分享时间
	State     int8      `json:"state"`      // 1-未处理 2-已处理
	FilePath  string    `json:"file_path"`
	Owner     string    `json:"owner"`
	ShareType int8      `json:"share_type"`
}

type UpdateRecordRequest struct {
	RecordIDs []string `json:"record_ids"` // 消息ID列表
}

type FileTypeEnum int8

const (
	// 未知
	FILE_UNKNOWN FileTypeEnum = 0
	//普通文件
	FILE FileTypeEnum = 1
	//文件夹
	FOLDER FileTypeEnum = 2
	//批量操作
	BATCH FileTypeEnum = 3
)

func (e FileTypeEnum) String() string {
	switch e {
	case FILE:
		return "文件"
	case FOLDER:
		return "文件夹"
	case BATCH:
		return "批量操作"
	default:
		return "未知"
	}
}

type OpTypeEnum int8

const (
	// OP_UNKNOWN 未知
	OP_UNKNOWN OpTypeEnum = 0
	// UPLOAD 上传
	UPLOAD OpTypeEnum = 1
	// DOWNLOAD 下载
	DOWNLOAD OpTypeEnum = 2
	// DELETE 删除
	DELETE OpTypeEnum = 3
	// RENAME 重命名
	RENAME OpTypeEnum = 4
	// ADD_FOLDER 添加文件夹
	ADD_FOLDER OpTypeEnum = 5
	// OPEN 公开
	OPEN OpTypeEnum = 6
	// SHARE 分享
	SHARE OpTypeEnum = 7
)

func (e OpTypeEnum) String() string {
	switch e {
	case UPLOAD:
		return "上传"
	case DOWNLOAD:
		return "下载"
	case DELETE:
		return "删除"
	case RENAME:
		return "重命名"
	case ADD_FOLDER:
		return "新建文件夹"
	case OPEN:
		return "公开"
	case SHARE:
		return "分享"
	default:
		return "未知"
	}
}
