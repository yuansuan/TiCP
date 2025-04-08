package openapivisual

import (
	"errors"
)

var (
	// ErrSessionIDEmpty 会话ID不能为空
	ErrSessionIDEmpty = errors.New("the session id is empty")
	// ErrSessionIDsStrEmpty 会话IDs字符串不能为空
	ErrSessionIDsStrEmpty = errors.New("the session ids is empty")
	// ErrHardwareIDEmpty 硬件ID不能为空
	ErrHardwareIDEmpty = errors.New("the hardware id is empty")
	// ErrSoftwareIDEmpty 软件ID不能为空
	ErrSoftwareIDEmpty = errors.New("the software id is empty")
	// ErrExistReasonEmpty 退出原因不能为空
	ErrExistReasonEmpty = errors.New("the exist reason is empty")

	// ErrHardwareNameEmpty 硬件名称不能为空
	ErrHardwareNameEmpty = errors.New("the hardware name is empty")
	// ErrHardwareInstanceTypeEmpty 硬件实例类型不能为空
	ErrHardwareInstanceTypeEmpty = errors.New("the hardware instance type is empty")
	// ErrHardwareInstanceFamilyEmpty 硬件实例族不能为空
	ErrHardwareInstanceFamilyEmpty = errors.New("the hardware instance family is empty")
	// ErrHardwareZoneEmpty 硬件可用区不能为空
	ErrHardwareZoneEmpty = errors.New("the hardware zone is empty")

	// ErrSoftwareNameEmpty 软件名称不能为空
	ErrSoftwareNameEmpty = errors.New("the software name is empty")
	// ErrSoftwarePlatformEmpty 软件平台不能为空
	ErrSoftwarePlatformEmpty = errors.New("the software platform is empty")
	// ErrSoftwareImageIDEmpty 软件镜像ID不能为空
	ErrSoftwareImageIDEmpty = errors.New("the software image id is empty")

	// ErrRemoteAppIDEmpty 远程应用ID不能为空
	ErrRemoteAppIDEmpty = errors.New("the remote app id is empty")
	// ErrRemoteAppNameEmpty 远程应用名称不能为空
	ErrRemoteAppNameEmpty = errors.New("the remote app name is empty")
	// ErrRemoteAppBaseURLEmpty 远程应用基础URL不能为空
	ErrRemoteAppBaseURLEmpty = errors.New("the remote app base url is empty")
)
