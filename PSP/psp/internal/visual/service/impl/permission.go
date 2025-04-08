package impl

import (
	"context"
	"fmt"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/service/client"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/reflectutil"
)

func (s *VisualService) AddVisualPermission(ctx context.Context, visualResourceType string, resource any) error {
	resourceID, resourceName, err := s.getResourceInfo(ctx, resource, false)
	if err != nil {
		return err
	}

	permission, err := client.GetInstance().RBAC.Permission.AddPermission(ctx, &rbac.Resource{
		ResourceName: resourceName,
		DisplayName:  resourceName,
		ResourceType: visualResourceType,
		ResourceId:   resourceID,
		Custom:       common.ENABLE_CUSTOM,
		Action:       common.ResourceActionNONE,
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

	tracelog.Info(ctx, fmt.Sprintf("add visual type: [%v] permission, resoureID: [%v], resourceName: [%v]", visualResourceType, resourceID, resourceName))

	return nil
}

func (s *VisualService) UpdateVisualPermission(ctx context.Context, visualResourceType string, originResource, resource any) error {
	resourceID, resourceName, err := s.getResourceInfo(ctx, originResource, true)
	if err != nil {
		return err
	}

	permission, err := client.GetInstance().RBAC.Permission.GetResourcePerm(ctx, &rbac.ResourceIdentity{
		Identity: &rbac.ResourceIdentity_Name{
			Name: &rbac.ResourceName{
				Type: visualResourceType,
				Name: resourceName,
			},
		},
	})
	if err != nil {
		return err
	}

	_, resourceName, err = s.getResourceInfo(ctx, resource, false)
	if err != nil {
		return err
	}

	permission.ResourceName = resourceName
	permission.DisplayName = resourceName
	_, err = client.GetInstance().RBAC.Permission.UpdatePermission(ctx, permission)
	if err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("update visual type: [%v] permission, resoureID: [%v], resourceName: [%v]", visualResourceType, resourceID, resourceName))

	return nil
}

func (s *VisualService) DeleteVisualPermission(ctx context.Context, visualResourceType string, resource any) error {
	resourceID, resourceName, err := s.getResourceInfo(ctx, resource, false)
	if err != nil {
		return err
	}

	permission, err := client.GetInstance().RBAC.Permission.GetResourcePerm(ctx, &rbac.ResourceIdentity{
		Identity: &rbac.ResourceIdentity_Name{
			Name: &rbac.ResourceName{
				Type: visualResourceType,
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

	tracelog.Info(ctx, fmt.Sprintf("update visual type: [%v] permission, resoureID: [%v], resourceName: [%v]", visualResourceType, resourceID, resourceName))

	return nil
}

func (s *VisualService) getResourceInfo(ctx context.Context, resource any, origin bool) (int64, string, error) {
	var err error

	resourceID, resourceName := int64(0), ""
	switch v := resource.(type) {
	case *model.Hardware:
		if origin {
			return int64(v.ID), v.Name, nil
		}
		hardware, exist, errTmp := s.hardwareDao.GetHardware(ctx, v.ID, "")
		if exist && hardware != nil {
			resourceID, resourceName = int64(hardware.ID), hardware.Name
		}
		err = errTmp
	case *model.Software:
		if origin {
			return int64(v.ID), v.Name, nil
		}
		software, exist, errTmp := s.softwareDao.GetSoftware(ctx, v.ID, "")
		if exist && software != nil {
			resourceID, resourceName = int64(software.ID), software.Name
		}
		err = errTmp
	default:
		return resourceID, resourceName, fmt.Errorf("resource: [%v] type not support", reflectutil.GetStructName(v))
	}

	return resourceID, resourceName, err
}
