package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
)

func TestConsoleLog(t *testing.T) {
	oldStderr := os.Stderr

	rfd, wfd, err := os.Pipe()
	assert.NoError(t, err)
	os.Stderr = wfd

	err = InitLogger(config.Log{
		Level:        string(InfoLevel),
		ReleaseLevel: string(DevelopmentLevel),
		UseConsole:   true,
	})
	assert.NoError(t, err)

	Info("info")
	Error("error")
	wfd.Close()

	buf := make([]byte, 1024)
	n, err := rfd.Read(buf)
	assert.NoError(t, err)
	assert.Contains(t, string(buf[:n]), "info")
	assert.Contains(t, string(buf[:n]), "error")

	os.Stdout = oldStderr
}
