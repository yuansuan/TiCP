package util

import (
	"os"
)

func IsFileExist(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}

	return false
}
