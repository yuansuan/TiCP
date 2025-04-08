package fsutil

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/fsutil/filemode"
)

func TestMkdirAuto(t *testing.T) {
	root = &Filesystem{
		realpath: "/",
		Fs:       afero.NewBasePathFs(afero.NewOsFs(), "/"),
	}
	testDirPath := "/tmp/test_tmp_dir_" + time.Now().String()

	if err := MkdirAuto(testDirPath, filemode.Directory); err != nil {
		t.Fatal(err)
	}

	os.Remove(testDirPath)

}
