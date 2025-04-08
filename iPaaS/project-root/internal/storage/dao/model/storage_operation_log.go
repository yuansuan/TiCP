package model

import (
	"time"
)

type StorageOperationLog struct {
	Id            int64     `json:"id" xorm:"pk autoincr comment('主键id') BIGINT(20)"`
	UserId        string    `json:"user_id" xorm:"not null index(idx_user_id) VARCHAR(255)"`
	FileName      string    `json:"file_name" xorm:"not null VARCHAR(255)"`
	SrcPath       string    `json:"src_path" xorm:"not null MEDIUMTEXT"`
	DestPath      string    `json:"dest_path" xorm:"not null MEDIUMTEXT"`
	FileType      string    `json:"file_type" xorm:"not null comment('文件类型, 可选值: file-普通文件, folder-文件夹, batch-批量操作') VARCHAR(20)"`
	OperationType string    `json:"operation_type" xorm:"not null comment('操作类型, 可选值: upload-上传, download-下载, delete-删除, move-移动, mkdir-添加文件夹, compress-压缩, copy-拷贝, copy_range-指定范围拷贝, create-创建, link-链接, read_at-读, write_at-写') VARCHAR(20)"`
	Size          string    `json:"size" xorm:"not null VARCHAR(20)"`
	IsDeleted     int       `json:"is_deleted" xorm:"not null INT(2)"`
	CreateTime    time.Time `json:"create_time" xorm:"not null DATETIME"`
}
