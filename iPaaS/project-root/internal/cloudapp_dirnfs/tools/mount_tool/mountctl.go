package main

import (
	"os"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_dirnfs/tools/mount_tool/cmd"
)

// 文件挂载工具
// CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build mountctl.go
// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build mountctl.go
func main() {

	command := cmd.NewCSPCtlCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
