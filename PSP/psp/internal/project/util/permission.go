package util

import (
	"context"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/service/client"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

// 是否是项目管理权限角色或者是系统管理员
func CheckProjectAdminRole(ctx context.Context, loginUserID snowflake.ID, checkPmPermission bool) (bool, error) {
	// 系统管理员权限校验
	sysManageResourceIdentity := &rbac.ResourceIdentity{
		Identity: &rbac.ResourceIdentity_Name{
			Name: &rbac.ResourceName{
				Type: common.PermissionResourceTypeSystem,
				Name: common.ResourceSysManagerName,
			},
		},
	}

	permRequest := rbac.CheckResourcesPermRequest{
		Id:        &rbac.ObjectID{Id: loginUserID.String()},
		Resources: []*rbac.ResourceIdentity{sysManageResourceIdentity},
	}
	checkResp, err := client.GetInstance().Rbac.CheckResourcesPerm(ctx, &permRequest)
	if err != nil {
		return false, err
	}

	// 只要检查系统管理员权限
	if !checkPmPermission && !checkResp.Pass {
		return false, nil
	}

	// 校验项目管理权限
	if checkPmPermission {
		projectManageResourceIdentity := &rbac.ResourceIdentity{
			Identity: &rbac.ResourceIdentity_Name{
				Name: &rbac.ResourceName{
					Type: common.PermissionResourceTypeSystem,
					Name: common.ResourceProjectName,
				},
			},
		}

		pmRequest := rbac.CheckResourcesPermRequest{
			Id:        &rbac.ObjectID{Id: loginUserID.String()},
			Resources: []*rbac.ResourceIdentity{projectManageResourceIdentity},
		}
		checkPMResp, err := client.GetInstance().Rbac.CheckResourcesPerm(ctx, &pmRequest)
		if err != nil {
			return false, err
		}

		if !checkPMResp.Pass {
			return false, nil
		}
	}

	return true, nil
}
