package rpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/util"
)

func (srv *GRPCService) GetRole(ctx context.Context, req *rbac.RoleID) (*rbac.Role, error) {
	roleInfo, err := srv.RoleService.GetRole(ctx, req.Id)
	if err != nil {
		return &rbac.Role{}, err
	}
	return util.ToGRPC.Role(roleInfo), nil
}

func (srv *GRPCService) GetRoles(ctx context.Context, req *rbac.RoleIDs) (*rbac.Roles, error) {
	roleInfoList, err := srv.RoleService.GetRoles(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return &rbac.Roles{
		Roles: util.ToGRPC.Roles(roleInfoList),
		Total: int64(len(roleInfoList)),
	}, nil
}

func (srv *GRPCService) GetRoleByName(ctx context.Context, req *rbac.RoleName) (*rbac.Role, error) {
	roleInfo, err := srv.RoleService.GetRoleByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return util.ToGRPC.Role(roleInfo), nil
}

// AddRolePerms AddRolePerms
func (srv *GRPCService) AddRolePerms(ctx context.Context, req *rbac.RolePerms) (*empty.Empty, error) {

	err := srv.RoleService.AddRolePerms(ctx, req.Role.Id, req.Perms)

	if err != nil {
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}

// UpdateRolePerms UpdateRolePerms
func (srv *GRPCService) UpdateRolePerms(ctx context.Context, req *rbac.RolePerms) (*empty.Empty, error) {
	err := srv.RoleService.UpdateRolePerms(ctx, req.Role.Id, req.Perms)
	if err != nil {
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}

// RemoveRolePerms RemoveRolePerms
func (srv *GRPCService) RemoveRolePerms(ctx context.Context, req *rbac.RolePerms) (*empty.Empty, error) {

	err := srv.RoleService.RemoveRolePerms(ctx, req.Role.Id, req.Perms)
	if err != nil {
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}

// InternalAddRolePerms InternalAddRolePerms
func (srv *GRPCService) InternalAddRolePerms(ctx context.Context, req *rbac.RolePerms) (*empty.Empty, error) {
	err := srv.RoleService.InternalAddRolePerms(ctx, req.Role.Id, req.Perms)

	if err != nil {
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}

// InternalRemoveRolePerms InternalRemoveRolePerms
func (srv *GRPCService) InternalRemoveRolePerms(ctx context.Context, req *rbac.RolePerms) (*empty.Empty, error) {
	err := srv.RoleService.InternalRemoveRolePerms(ctx, req.Role.Id, req.Perms)

	if err != nil {
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}

// AddObjectRoles AddObjectRoles
func (srv *GRPCService) AddObjectRoles(ctx context.Context, req *rbac.ObjectRoles) (*empty.Empty, error) {

	err := srv.RoleService.AddObjectRoles(ctx, req)

	if err != nil {
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}

// ListObjectsRoles ListObjectsRoles
func (srv *GRPCService) ListObjectsRoles(ctx context.Context, req *rbac.ListObjectsRolesReq) (*rbac.ListObjectsRolesResp, error) {
	return srv.RoleService.ListObjectsRoles(ctx, req.Ids, req.NeedImplicitRoles)
}

// UpdateObjectRoles UpdateObjectRoles
func (srv *GRPCService) UpdateObjectRoles(ctx context.Context, req *rbac.ObjectRoles) (*empty.Empty, error) {
	err := srv.RoleService.UpdateObjectRoles(ctx, req)
	if err != nil {
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}

// RemoveObjectRoles RemoveObjectRoles
func (srv *GRPCService) RemoveObjectRoles(ctx context.Context, req *rbac.ObjectRoles) (*empty.Empty, error) {

	err := srv.RoleService.RemoveObjectRoles(ctx, req)
	if err != nil {
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}

func (srv *GRPCService) GetRoleByObjectID(ctx context.Context, req *rbac.ObjectID) (*rbac.Roles, error) {
	roleInfoList, err := srv.RoleService.GetRoleByObjectID(ctx, req)
	if err != nil {
		return &rbac.Roles{}, err
	}

	util.ToGRPC.Roles(roleInfoList)

	return &rbac.Roles{
		Roles: util.ToGRPC.Roles(roleInfoList),
		Total: int64(len(roleInfoList)),
	}, nil

}

func (srv *GRPCService) AddRole(ctx context.Context, req *rbac.AddRoleReq) (*empty.Empty, error) {
	role := &model.Role{
		Name:    req.Name,
		Comment: req.Comment,
		Type:    consts.RoleTypeCustom,
	}

	err := srv.RoleService.AddRole(ctx, role)

	if err != nil {
		return nil, err
	}

	if len(req.Perms) > 0 {
		err = srv.RoleService.AddRolePerms(ctx, role.Id, req.Perms)
		if err != nil {
			return nil, err
		}
	}
	return &empty.Empty{}, nil
}

func (srv *GRPCService) UpdateRole(ctx context.Context, req *rbac.UpdateRoleReq) (*empty.Empty, error) {

	err := srv.RoleService.UpdateRole(ctx, &model.Role{
		Id:      req.Id,
		Name:    req.Name,
		Comment: req.Comment,
	})
	if err != nil {
		return nil, err
	}
	err = srv.RoleService.UpdateRolePerms(ctx, req.Id, req.Perms)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
func (srv *GRPCService) DelRole(ctx context.Context, req *rbac.RoleID) (*empty.Empty, error) {
	_, err := srv.RoleService.GetRole(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	err = srv.RoleService.DeleteRole(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (srv *GRPCService) SetLdapUserDefRole(ctx context.Context, req *rbac.RoleID) (*empty.Empty, error) {
	err := srv.RoleService.SetLdapUserDefRole(ctx, req.Id)

	if err != nil {
		return nil, status.Error(errcode.ErrUserLDAPDefRoleFailed, "")
	}

	return &empty.Empty{}, nil
}
