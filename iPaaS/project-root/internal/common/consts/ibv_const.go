package consts

import (
	"google.golang.org/grpc/codes"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

//TODO

// 3D云应用相关错误码
const (
	// ErrVisSoftwareDesiredGPU 所选择的软件需要一个GPU
	ErrVisSoftwareDesiredGPU codes.Code = 280001 + iota
	// ErrVisBadAvailableZone 无效的可用区
	ErrVisBadAvailableZone = 280008
	// ErrVisInternalError 创建会话失败: 内部错误
	ErrVisInternalError // 280003
	// ErrVisHardwareNotFound 硬件不存在
	ErrVisHardwareNotFound
	// ErrVisBadRequest 前端传参错误
	ErrVisBadRequest
	// ErrVisSoftwareNotFound 软件不存在
	ErrVisSoftwareNotFound // 280005
	// ErrVisComboBeyondValidDate 3d云应用套餐不在有效期，禁止开启会话
	ErrVisComboBeyondValidDate
	// ErrVisComboRemainTimeInsufficient 3d云应用套餐剩余时间不足，禁止开启会话
	ErrVisComboRemainTimeInsufficient
	// ErrVisComboNodeInsufficient  开启会话的个数不能超过套餐中规定的节点个数
	ErrVisComboNodeInsufficient
	// ErrVisInsufficientMachine 资源不足
	ErrVisInsufficientMachine
	// ErrVisStartInstanceFailed 创建实例失败: 腾讯云错误
	ErrVisStartInstanceFailed
)

const (
	// ZSWLProductID 智算未来产品ID 1260511618517025002("3VGTYLv6FhW")
	ZSWLProductID snowflake.ID = snowflake.ID(1260511618517025002)
)

const (
	Unknown codes.Code = 2

	InvalidArgument codes.Code = 3
)

const (
	InvalidParam codes.Code = 240006
)
