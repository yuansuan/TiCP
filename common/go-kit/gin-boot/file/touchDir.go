package file

import (
	"os"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"
)

// TouchDir TouchDir
func TouchDir(dir string) {
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			util.ChkErr(err)
		}
	}
}
