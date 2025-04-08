package filemode

import "os"

const (
	// Directory 表示文件夹的默认权限配置(rwx.r-x.r-x)
	Directory os.FileMode = 0755

	// RegularFile 表示文件的默认权限配置(rw-.r--.r--)
	RegularFile os.FileMode = 0644
)
