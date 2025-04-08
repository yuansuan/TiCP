package v20230530

import (
	"os"
	"time"
)

// swagger:model storageFileInfo
type FileInfo struct {
	Name    string    `json:"Name"`    // base name of the file
	Size    int64     `json:"Size"`    // length in bytes for regular files; system-dependent for others
	Mode    uint32    `json:"Mode"`    // file mode bits
	ModTime time.Time `json:"ModTime"` // modification time
	IsDir   bool      `json:"IsDir"`   // abbreviation for Mode().IsDir()
}

func ToRespFileInfo(file os.FileInfo) *FileInfo {
	return &FileInfo{
		Name:    file.Name(),
		Size:    file.Size(),
		Mode:    uint32(file.Mode()),
		ModTime: file.ModTime(),
		IsDir:   file.IsDir(),
	}
}

func ToRespFileInfos(fileInfos []os.FileInfo) []*FileInfo {
	result := make([]*FileInfo, len(fileInfos))
	for i, f := range fileInfos {
		result[i] = ToRespFileInfo(f)
	}
	return result
}
