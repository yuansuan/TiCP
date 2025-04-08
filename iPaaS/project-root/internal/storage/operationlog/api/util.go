package api

import (
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"time"
)

func ToResponseOperationLogs(operationLogs []*model.StorageOperationLog) []*v20230530.OperationLog {
	var res []*v20230530.OperationLog
	for _, operationLog := range operationLogs {
		res = append(res, ToResponseOperationLog(operationLog))
	}
	return res
}

func ToResponseOperationLog(operationLog *model.StorageOperationLog) *v20230530.OperationLog {
	return &v20230530.OperationLog{
		Id:            operationLog.Id,
		UserId:        operationLog.UserId,
		FileName:      operationLog.FileName,
		SrcPath:       operationLog.SrcPath,
		DestPath:      operationLog.DestPath,
		FileType:      operationLog.FileType,
		OperationType: operationLog.OperationType,
		Size:          operationLog.Size,
		CreateTime:    operationLog.CreateTime,
	}
}

func UnixToTime(unixTimestamp int64) time.Time {
	return time.Unix(unixTimestamp, 0)
}

var validFileTypes = []string{commoncode.FILE, commoncode.FOLDER, commoncode.Batch}

func CheckFileType(fileType string) bool {
	for _, validType := range validFileTypes {
		if fileType == validType {
			return true
		}
	}
	return false
}

var validOperationTypes = []string{commoncode.FOLDER, commoncode.DOWNLOAD, commoncode.DELETE, commoncode.UPLOAD,
	commoncode.MOVE, commoncode.MKDIR, commoncode.COPY,
	commoncode.COPY_RANGE, commoncode.CREATE, commoncode.LINK,
	commoncode.READ_AT, commoncode.WRITE_AT}

func CheckOperationType(operationType string) bool {
	for _, validType := range validOperationTypes {
		if operationType == validType {
			return true
		}
	}
	return false
}

func CheckTime(timestamp int64) bool {
	if timestamp < 0 {
		return false
	}
	tm := time.Unix(timestamp, 0)
	return tm.Unix() == timestamp
}
