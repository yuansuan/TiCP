package fsutil

import (
	"fmt"
	"strings"
)

// InvalidChars 包含：三种引号，分号
const (
	InvalidChars       = "" //"\";'`"
	TmpBucketUploading = ".uploading"
	TmpFileCompress    = ".compress"
)

// IsSafePath 判断非法路径
func IsSafePath(path string) (bool, string) {
	if path == "" {
		return false, "path can not be empty"
	}
	if !strings.HasPrefix(path, "/") {
		return false, "path must start with /"
	}
	if strings.Contains(fmt.Sprintf("/%s/", path), "/../") {
		return false, "path can not contain .."
	}
	if strings.ContainsAny(path, InvalidChars) {
		return false, "path contains invalid char"
	}
	return true, ""
}

// ValidateUserIDPath 验证路径 并获得userID
func ValidateUserIDPath(path string) (isValid bool, userID string, errMsg string) {
	if path == "" {
		return false, "", "path can not be empty"
	}
	if strings.Contains(path, "..") {
		return false, "", "path can not contain .."
	}
	if strings.ContainsAny(path, InvalidChars) {
		return false, "", "path contains invalid char"
	}
	parts := strings.Split(path, "/")
	if !strings.HasPrefix(path, "/") || len(parts) < 2 || parts[1] == "" {
		return false, "", "path must start with /{userID}"
	}
	userID = parts[1]
	return true, userID, ""
}

// TrimPrefix ...
func TrimPrefix(path, prefix string) string {
	if strings.HasPrefix(path, prefix) {
		return strings.TrimPrefix(path, prefix)
	}
	return path
}

// TmpBucketUploadingJoin ...
func TmpBucketUploadingJoin(uploadID string) string {
	return fmt.Sprintf("%s/%s", TmpBucketUploading, uploadID)
}

// TmpFileCompressJoin ...
func TmpFileCompressJoin(fileName string) string {
	return fmt.Sprintf("%s/%s", TmpFileCompress, fileName)
}
