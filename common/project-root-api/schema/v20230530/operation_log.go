package v20230530

import "time"

type OperationLog struct {
	Id            int64     `json:"id"`
	UserId        string    `json:"userID"`
	FileName      string    `json:"fileName"`
	SrcPath       string    `json:"srcPath"`
	DestPath      string    `json:"destPath"`
	FileType      string    `json:"fileType"`
	OperationType string    `json:"operationType"`
	Size          string    `json:"size"`
	CreateTime    time.Time `json:"createTime"`
}
