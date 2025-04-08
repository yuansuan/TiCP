package util

import (
	"bufio"
	"os"
	"strings"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/config"
)

var (
	hostNameMap = make(map[string]string)
)

func GetHostNameMap() map[string]string {
	return hostNameMap
}

func InitHostNameMapping() {
	hostData := config.GetConfig().HostNameMapping
	logger := logging.Default()
	if !hostData.Enable {
		return
	}

	file, err := os.Open(hostData.Path)
	if err != nil {
		logger.Errorf("failed to open file: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			hostNameMap[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Errorf("file reading failure: %v", err)
		return
	}
}
