//go:build linux

package onclose

import "os"

const (
	envPath = "/etc/ys-agent/agent_env"
)

func Clean() {
	_ = os.RemoveAll(envPath)
}
