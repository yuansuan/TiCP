package approve

//
import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dto"
)

var (
	// go不支持把一个泛型接口作为map的value类型，只能用interface{}了
	approveImpls = make(map[ApproveTypeImplEnum]IApproveService[any])
	once         sync.Once
)

func GetImpl(approveType dto.ApproveType) IApproveService[any] {
	lazyInit()
	return approveImpls[convertType(approveType)]
}

type ApproveTypeImplEnum string

const (
	UserApproveType ApproveTypeImplEnum = "USER"
	RoleApproveType ApproveTypeImplEnum = "ROLE"
	UNKNOWN         ApproveTypeImplEnum = "UNKNOWN"
)

const (
	SignTypeID   string = "ID"
	SignTypeName string = "NAME"
)

func lazyInit() {
	once.Do(func() {
		approveImpls[UserApproveType] = NewUserApproveImpl()
		approveImpls[RoleApproveType] = NewRoleApproveImpl()
	})
}

func convertType(approveType dto.ApproveType) ApproveTypeImplEnum {
	switch approveType {
	case dto.ApproveTypeAddUser, dto.ApproveTypeDelUser, dto.ApproveTypeEditUser, dto.ApproveTypeEnableUser, dto.ApproveTypeDisableUser:
		return UserApproveType
	case dto.ApproveTypeAddRole, dto.ApproveTypeDelRole, dto.ApproveTypeEditRole, dto.ApproveTypeSetLdapDefRole:
		return RoleApproveType
	default:
		return UNKNOWN
	}
}

type IApproveService[T any] interface {
	// PreCheck 前置检查
	PreCheck(ctx context.Context, req *dto.ApplyApproveRequest) error
	// BuildContent 构建前端显示文案
	BuildContent(ctx context.Context, approveInfo T) (string, error)
	// GenSign 生成签名
	GenSign(ctx context.Context, req *dto.ApplyApproveRequest) (string, error)
	// BuildApproveInfo 构建审批信息
	BuildApproveInfo(ctx *gin.Context, req *dto.ApplyApproveRequest) (approveInfo T, err error)
	// CheckNecessary 检查审批通过数据落地必要条件, 返参 ready:是否可以执行数据落地 notReadyErr:不满足条件的错误信息 err: 通用错误信息
	CheckNecessary(ctx context.Context, approveType dto.ApproveType, approveInfo T) (ready bool, notReadyErr error, err error)
	// AfterPass 审批通过数据落地等后处理
	AfterPass(ctx context.Context, approveInfo T) error
	// ParseObject 审批信息数据库json字符串转结构体
	ParseObject(jsonString string) (T, error)
}
