package storage

import (
	"github.com/gin-gonic/gin"
)

type Storage interface {
	// Ls 用于分页list文件夹，每页最大1000，如果list出来的文件数量小于请求时指定的PageSize，说明分页list已经结束，不再提供递归list的接口。
	Ls(ctx *gin.Context)
	// Mkdir 用于创建文件夹，必须'/'开头，路径里不能带'../'。
	Mkdir(ctx *gin.Context)
	// Mv 移动一个文件或文件夹。
	Mv(ctx *gin.Context)
	// Rm 删除一个文件或文件夹。
	Rm(ctx *gin.Context)
	// Stat 获取一个文件或文件夹的信息。
	Stat(ctx *gin.Context)
	// UploadInit 用于初始化一个文件上传，返回一个uploadID，后续的上传slice都需要带上这个uploadID。
	UploadInit(ctx *gin.Context)
	// UploadSlice 用于上传一个文件的一个分片。
	UploadSlice(ctx *gin.Context)
	// UploadComplete 用于完成一个文件的上传。
	UploadComplete(ctx *gin.Context)
	// UploadFile 用于直接上传一个1G以内的文件。
	UploadFile(ctx *gin.Context)
	// Download 用于下载一个文件。
	Download(ctx *gin.Context)
	// BatchDownload 用于批量下载文件，返回一个zip文件。
	BatchDownload(ctx *gin.Context)
	// Realpath 根据相对路径转成绝对路径
	Realpath(ctx *gin.Context)
	// Copy 复制一个文件/目录(包括子目录和文件)
	Copy(ctx *gin.Context)
	// CopyRange 复制一个文件的一部分到另一个文件
	CopyRange(ctx *gin.Context)
	// Link 创建一个链接文件
	Link(ctx *gin.Context)
	// Truncate 改变文件大小
	Truncate(ctx *gin.Context)
	// WriteAt 随机写文件
	WriteAt(ctx *gin.Context)
	// ReadAt 随机读文件
	ReadAt(ctx *gin.Context)
	// Create 创建一个文件
	Create(ctx *gin.Context)
	// CompressStart 开始压缩文件
	CompressStart(ctx *gin.Context)
	// CompressStatus 查看压缩状态
	CompressStatus(ctx *gin.Context)
	// CompressCancel 取消压缩任务
	CompressCancel(ctx *gin.Context)
}
