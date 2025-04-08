package fsutil

import (
	"io/fs"
	"time"
)

type linkFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
	sys     any
}

func (fi *linkFileInfo) Name() string {
	return fi.name
}

func (fi *linkFileInfo) Size() int64 {
	return fi.size
}

func (fi *linkFileInfo) Mode() fs.FileMode {
	return fi.mode
}

func (fi *linkFileInfo) ModTime() time.Time {
	return fi.modTime
}

func (fi *linkFileInfo) IsDir() bool {
	return fi.isDir
}

func (fi *linkFileInfo) Sys() any {
	return fi.sys
}
