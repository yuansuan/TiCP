//go:build !linux && !darwin && !freebsd
// +build !linux,!darwin,!freebsd

package fileinfo

// loadSys 其他系统暂时没有支持
func loadSys(_ []byte, _ *extSysFileInfo) {}
