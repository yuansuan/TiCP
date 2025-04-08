package filemode

import "os"

const (
	// Directory 表示文件夹的默认权限配置(rwx.r-x.r-x)
	Directory os.FileMode = 0755
)
