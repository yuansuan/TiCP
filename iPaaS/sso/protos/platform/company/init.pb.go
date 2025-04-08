package company

import grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"

var _ grpc_boot.ServerType

func init() {
	grpc_boot.RegisterClient("company", NewCompanyApplyServiceClient)
	grpc_boot.RegisterClient("company", NewTrialApplyServiceClient)
	grpc_boot.RegisterClient("company", NewInnerTrialServiceClient)
	grpc_boot.RegisterClient("company", NewCompanyServiceClient)
	grpc_boot.RegisterClient("company", NewCompanyCheckServiceClient)
	grpc_boot.RegisterClient("company", NewCompanyMerchandiseServiceClient)
	grpc_boot.RegisterClient("company", NewCompanyUserConfigServiceClient)
	grpc_boot.RegisterClient("company", NewCompanyConfigServiceClient)
	grpc_boot.RegisterClient("company", NewDepartmentServiceClient)
	grpc_boot.RegisterClient("company", NewPermissionServiceClient)
	grpc_boot.RegisterClient("company", NewProjectServiceClient)
	grpc_boot.RegisterClient("company", NewProjectReadOnlyServiceClient)
	grpc_boot.RegisterClient("company", NewRoleServiceClient)
	grpc_boot.RegisterClient("company", NewUserServiceClient)
}
