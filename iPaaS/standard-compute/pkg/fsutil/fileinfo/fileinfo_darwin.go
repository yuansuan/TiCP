//go:build darwin
// +build darwin

package fileinfo

import (
	"encoding/json"
	"time"
)

// darwin 表示MacOS的额外文件信息
type darwin struct {
	UID        int64 `json:"Uid"`
	GID        int64 `json:"Gid"`
	AccessTime struct {
		Second     int64 `json:"Sec"`
		NanoSecond int64 `json:"Nsec"`
	} `json:"Atimespec"`
	CreateTime struct {
		Second     int64 `json:"Sec"`
		NanoSecond int64 `json:"Nsec"`
	} `json:"Ctimespec"`
}

// loadSys 加载额外的系统文件信息
func loadSys(data []byte, fi *extSysFileInfo) {
	var sys darwin
	if err := json.Unmarshal(data, &sys); err == nil {
		fi.uid = sys.UID
		fi.gid = sys.GID
		fi.atime = time.Unix(sys.AccessTime.Second, sys.AccessTime.NanoSecond)
		fi.ctime = time.Unix(sys.CreateTime.Second, sys.CreateTime.NanoSecond)
	}
}
