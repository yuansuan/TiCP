package version

import (
	"flag"
	"os"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

var version = ""
var branch = ""
var commitId = ""
var buildTime = ""

func ShouldLogVersion() bool {
	if len(os.Args) == 2 && os.Args[1] == "version" {
		return true
	}

	// avoid flag redefined problem with glog from grpc
	cmdline := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	show := cmdline.Bool("v", false, "show the current version of the application")

	err := cmdline.Parse(os.Args[1:])
	return err == nil && *show
}

// LogVersion LogVersion
func LogVersion() {
	logging.Default().Infof("Version: %s, Branch: %s, Build: %s, Build time: %s\n",
		version, branch, commitId, buildTime)
}
