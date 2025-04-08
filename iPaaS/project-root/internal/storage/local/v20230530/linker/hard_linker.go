package linker

import (
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	"os"
	"path/filepath"
)

type HardLink struct {
}

func (hl *HardLink) Link(sourcePath string, destPath string) error {
	// 创建目录
	dirPath := filepath.Dir(destPath)
	if err := os.MkdirAll(dirPath, filemode.Directory); err != nil {
		return errors.Wrap(err, "HardLink")
	}

	// 硬件连源文件
	err := os.Link(sourcePath, destPath)
	if err != nil {
		return errors.Wrap(err, "HardLink")
	}

	return nil
}
