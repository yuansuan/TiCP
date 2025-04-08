package errcode

import (
	"google.golang.org/grpc/codes"
)

// Storage service error codes from 15001 to 16000
const (
	ErrFilePathTooLong             codes.Code = 15001
	ErrFilePathInvalidChar         codes.Code = 15002
	ErrFileNotExist                codes.Code = 15003
	ErrFileNoPermission            codes.Code = 15004
	ErrFileGetInfo                 codes.Code = 15005
	ErrFileArgs                    codes.Code = 15006
	ErrFileNameWasTaken            codes.Code = 15007
	ErrFileAlreadyExist            codes.Code = 15008
	ErrFileFailToBindRequest       codes.Code = 15009
	ErrFileCacheFail               codes.Code = 15010
	ErrFileFailToGetSlice          codes.Code = 15011
	ErrFileFailToOpenSlice         codes.Code = 15012
	ErrFileSliceIncompleted        codes.Code = 15013
	ErrFileSameFile                codes.Code = 15014
	ErrFileList                    codes.Code = 15015
	ErrFileFailRemove              codes.Code = 15016
	ErrFileFailMkdir               codes.Code = 15017
	ErrFileFailCopy                codes.Code = 15018
	ErrFileSamePath                codes.Code = 15019
	ErrFileFailGet                 codes.Code = 15020
	ErrFileFailRename              codes.Code = 15021
	ErrFileFailMove                codes.Code = 15022
	ErrFileFailRead                codes.Code = 15023
	ErrFileDownloadPermission      codes.Code = 15024
	ErrFileUploadInit              codes.Code = 15025
	ErrFileUploadSlice             codes.Code = 15026
	ErrFileDownload                codes.Code = 15027
	ErrFilePreDownload             codes.Code = 15028
	ErrFileUnknownUser             codes.Code = 15029
	ErrFilePaas                    codes.Code = 15030
	ErrFileUserHomeDirNotExist     codes.Code = 15031
	ErrFileBatchDownload           codes.Code = 15032
	ErrFileBatchDownloadPermission codes.Code = 15033
	ErrFileSync                    codes.Code = 15034
	ErrFileNameValidate            codes.Code = 15035
	ErrFileFailCover               codes.Code = 15036
	ErrFileDirAlreadyExist         codes.Code = 15037
	ErrFileUploadHPCFile           codes.Code = 15038
	ErrFileUploadHPCFileEmpty      codes.Code = 15039
	ErrFileUploadHPCFileNotExist   codes.Code = 15040
	ErrFileCancelUploadHPCFile     codes.Code = 15041
	ErrFileResumeUploadHPCFile     codes.Code = 15042
	ErrFileGetUploadHPCFile        codes.Code = 15043
	ErrFileMoveItself              codes.Code = 15044
	ErrFileCopyFileNotExist        codes.Code = 15045
	ErrFileFailLink                codes.Code = 15046
	ErrFileShareFailed             codes.Code = 15050
	ErrFileShareNotExist           codes.Code = 15051
	ErrFileShareSpaceNotExist      codes.Code = 15052
	ErrFileShareGetFailed          codes.Code = 15053
	ErrFileShareRecordGetFailed    codes.Code = 15054
	ErrFileShareRecordUpdateFailed codes.Code = 15055
	ErrrFileOverwriteParent        codes.Code = 15056
	ErrFileFailReadFolder          codes.Code = 15057
	ErrFileCompress                codes.Code = 15058
	ErrFileGetCompressTasks        codes.Code = 15059
	ErrFileSubmitCompressTask      codes.Code = 15060

	ErrFileOperLogAdd    codes.Code = 15070
	ErrFileOperLogGet    codes.Code = 15071
	ErrFileOperLogExport codes.Code = 15072
)

var StorageCodeMsg = map[codes.Code]string{
	ErrFileUnknownUser:             "未知用户",
	ErrFileUploadInit:              "文件上传失败",
	ErrFileUploadSlice:             "文件上传失败",
	ErrFilePreDownload:             "文件预下载失败",
	ErrFileDownload:                "文件下载失败",
	ErrFileFailCopy:                "文件复制失败",
	ErrFileNotExist:                "文件不存在",
	ErrFileUserHomeDirNotExist:     "用户目录不存在",
	ErrFileList:                    "文件列表查询失败",
	ErrFileFailGet:                 "文件信息查询失败",
	ErrFileFailRename:              "文件重命名失败",
	ErrFileFailMove:                "文件移动失败",
	ErrFileFailMkdir:               "创建文件夹失败",
	ErrFileFailRemove:              "文件删除失败",
	ErrFileFailRead:                "查看文件失败",
	ErrFileNoPermission:            "文件权限不足",
	ErrFileBatchDownload:           "文件批量下载失败",
	ErrFileSync:                    "文件同步失败",
	ErrFileNameValidate:            "文件名存在特殊字符",
	ErrFileUploadHPCFile:           "上传hpc文件失败",
	ErrFileUploadHPCFileEmpty:      "上传文件列表为空",
	ErrFileAlreadyExist:            "文件已存在",
	ErrFileFailCover:               "文件覆盖失败",
	ErrFileUploadHPCFileNotExist:   "上传任务不存在",
	ErrFileCancelUploadHPCFile:     "取消上传任务失败",
	ErrFileResumeUploadHPCFile:     "恢复上传任务失败",
	ErrFileGetUploadHPCFile:        "查询上传任务失败",
	ErrFileDirAlreadyExist:         "文件夹已存在",
	ErrFileSamePath:                "目标路径不能与源文件路径一致",
	ErrFileMoveItself:              "无法将文件夹移动到其子目录中",
	ErrFileCopyFileNotExist:        "复制任务不存在",
	ErrFileFailLink:                "文件链接失败",
	ErrFileShareFailed:             "分享失败",
	ErrFileShareNotExist:           "分享文件已取消",
	ErrFileShareSpaceNotExist:      "共享空间不存在",
	ErrFileShareGetFailed:          "共享文件获取失败",
	ErrFileShareRecordUpdateFailed: "共享文件处理状态修改失败",
	ErrFileShareRecordGetFailed:    "文件分享记录获取失败",
	ErrFileOperLogAdd:              "添加文件操作日志失败",
	ErrFileOperLogExport:           "导出文件操作日志失败",
	ErrrFileOverwriteParent:        "不能覆盖文件的父目录",
	ErrFileFailReadFolder:          "不能查看读取文件夹内容",
	ErrFileCompress:                "压缩文件失败",
	ErrFileGetCompressTasks:        "压缩任务获取失败",
	ErrFileSubmitCompressTask:      "压缩任务提交失败",
}
