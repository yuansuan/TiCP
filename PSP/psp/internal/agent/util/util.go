package util

import (
	stdlog "log"
	"path/filepath"
	"strconv"
	"strings"
)

func GetLogFileArg(maxSize, maxNum string) (outSize int, outNum int) {
	outSize, err := strconv.Atoi(maxSize)
	if err != nil {
		outSize = 50
	}
	outNum, err = strconv.Atoi(maxNum)
	if err != nil {
		outNum = 20
	}
	return outSize, outNum
}

func CheckFilePath(path string) bool {
	if path == "" {
		stdlog.Println("log_path is required")
		return false
	}

	isAbs := filepath.IsAbs(path)
	if !isAbs {
		stdlog.Println("log_path must be absolute path")
		return false
	}

	_, fileName := filepath.Split(path)
	if !strings.Contains(fileName, ".") {
		stdlog.Println("log_path must be a full path")
		return false
	}

	return true
}
