package fileinfo

import (
	"encoding/json"
	"os"
	"time"
)

// FileInfo 是一个扩展文件系统的文件信息
type FileInfo interface {
	os.FileInfo

	UID() int64
	GID() int64
	AccessTime() time.Time
	CreateTime() time.Time
}

// extSysFileInfo 保存系统文件信息
type extSysFileInfo struct {
	os.FileInfo

	uid   int64
	gid   int64
	atime time.Time
	ctime time.Time
}

// UID 所有者ID
func (fi *extSysFileInfo) UID() int64 {
	return fi.uid
}

// GID 所有组ID
func (fi *extSysFileInfo) GID() int64 {
	return fi.gid
}

// AccessTime 最后一次访问时间
func (fi *extSysFileInfo) AccessTime() time.Time {
	return fi.atime
}

// CreateTime 文件创建时间
func (fi *extSysFileInfo) CreateTime() time.Time {
	return fi.ctime
}

// New 读取文件信息并返回扩展的文件信息
func New(fi os.FileInfo) FileInfo {
	nfi := &extSysFileInfo{FileInfo: fi}
	if sys := fi.Sys(); sys != nil {
		if data, err := json.Marshal(sys); err == nil {
			loadSys(data, nfi)
		}
	}
	return nfi
}
