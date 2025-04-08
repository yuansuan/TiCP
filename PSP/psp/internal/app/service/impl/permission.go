package impl

import (
	"context"
	"fmt"

	"github.com/yuansuan/ticp/PSP/psp/internal/app/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/reflectutil"
)

func AddAppPermission(ctx context.Context, appResourceType string, service, resource any) error {
	resourceID, resourceName, err := getResourceInfo(ctx, service, resource, false)
	if err != nil {
		return err
	}

	permission, err := client.GetInstance().RBAC.Permission.AddPermission(ctx, &rbac.Resource{
		ResourceName: resourceName,
		DisplayName:  resourceName,
		ResourceType: appResourceType,
		ResourceId:   resourceID,
		Action:       common.ResourceActionNONE,
		Custom:       common.ENABLE_CUSTOM,
	})
	if err != nil {
		return err
	}

	_, err = client.GetInstance().RBAC.Role.InternalAddRolePerms(ctx, &rbac.RolePerms{
		Role: &rbac.RoleID{
			Id: int64(rbac.RoleType_ROLE_SUPER_ADMIN),
		},
		Perms: []int64{permission.Id},
	})
	if err != nil {
		return err
	}

	return nil
}

func DeleteAppPermission(ctx context.Context, appResourceType string, service, resource any) error {
	_, resourceName, err := getResourceInfo(ctx, service, resource, false)
	if err != nil {
		return err
	}

	permission, err := client.GetInstance().RBAC.Permission.GetResourcePerm(ctx, &rbac.ResourceIdentity{
		Identity: &rbac.ResourceIdentity_Name{
			Name: &rbac.ResourceName{
				Type: appResourceType,
				Name: resourceName,
			},
		},
	})
	if err != nil {
		return err
	}

	if permission.ResourceName == "" {
		return nil
	}

	_, err = client.GetInstance().RBAC.Permission.DeletePermission(ctx, &rbac.PermissionID{Id: permission.Id})
	if err != nil {
		return err
	}

	return nil
}

func getResourceInfo(ctx context.Context, service, resource any, origin bool) (int64, string, error) {
	var err error

	resourceID, resourceName := int64(0), ""
	switch v := resource.(type) {
	case *model.App:
		if origin {
			return int64(v.ID), v.Name, nil
		}
		app, exist, errTmp := getResource(ctx, service, v.ID)
		if exist && app != nil {
			resourceID, resourceName = int64(app.ID), app.Name
		}
		err = errTmp
	default:
		return resourceID, resourceName, fmt.Errorf("resource type [%v] not support", reflectutil.GetStructName(v))
	}

	return resourceID, resourceName, err
}

func getResource(ctx context.Context, service any, ID snowflake.ID) (*model.App, bool, error) {
	switch v := service.(type) {
	case *AppService:
		app := &model.App{ID: ID}
		exist, err := v.appDao.GetApp(ctx, app)
		if err != nil {
			return nil, false, err
		}
		return app, exist, nil
	case *AppLoader:
		app := &model.App{ID: ID}
		exist, err := v.appDao.GetApp(ctx, app)
		if err != nil {
			return nil, false, err
		}
		return app, exist, nil
	default:
		return nil, false, fmt.Errorf("service type [%v] not support", reflectutil.GetStructName(v))
	}
}
