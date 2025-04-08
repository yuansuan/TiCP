//go:build windows

package onclose

import "os"

const (
	envPath = "C:\\Windows\\ys\\agent_env"
)

func Clean() {
	_ = os.RemoveAll(envPath)
}
