// Copyright (C) 2019 LambdaCal Inc.

package osutil

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DetectFileMIMEType returns file MIME type
func DetectFileMIMEType(absPath string) (string, error) {
	f, err := os.Open(absPath)
	if err != nil {
		return "", fmt.Errorf("DetectFileMIMEType: cannot open file %v", absPath)
	}
	defer f.Close()

	// read the first 512 bytes of the file
	// because http.DetectContentType() considers at most the first 512 bytes of data
	b := make([]byte, 512)
	n, err := f.Read(b)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("DetectFileMIMEType: cannot read file %v, err=%v", absPath, err)
	}

	return http.DetectContentType(b[:n]), nil
}

// IsFileText checks whether a file is text
func IsFileText(absPath string) (bool, error) {
	t, err := DetectFileMIMEType(absPath)
	if err != nil {
		return false, err
	}

	return strings.HasPrefix(t, "text"), nil
}

// EnsureDir ensure the file full path has a directory, if not
// create it.
func EnsureDir(fileName string) error {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			return merr
		}
	}
	return nil
}
