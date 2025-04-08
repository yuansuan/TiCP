//go:build linux || freebsd
// +build linux freebsd

package fileinfo

import (
	"encoding/json"
	"time"
)

// unix 表示*nix系统的额外文件信息
type unix struct {
	UID        int64 `json:"Uid"`
	GID        int64 `json:"Gid"`
	AccessTime struct {
		Second     int64 `json:"Sec"`
		NanoSecond int64 `json:"Nsec"`
	} `json:"Atim"`
	CreateTime struct {
		Second     int64 `json:"Sec"`
		NanoSecond int64 `json:"Nsec"`
	} `json:"Ctim"`
}

// loadSys 加载额外的系统文件信息
func loadSys(data []byte, fi *extSysFileInfo) {
	var sys unix
	if err := json.Unmarshal(data, &sys); err == nil {
		fi.uid = sys.UID
		fi.gid = sys.GID
		fi.atime = time.Unix(sys.AccessTime.Second, sys.AccessTime.NanoSecond)
		fi.ctime = time.Unix(sys.CreateTime.Second, sys.CreateTime.NanoSecond)
	}
}
