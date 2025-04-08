package rpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/util"
)

// AddPermission AddPermission
func (srv *GRPCService) AddPermission(ctx context.Context, res *rbac.Resource) (*rbac.Resource, error) {

	// this three type only add by init sql when install
	if res.ResourceType == common.PermissionResourceTypeSystem ||
		res.ResourceType == common.PermissionResourceTypeInternal ||
		res.ResourceType == common.PermissionResourceTypeApi {
		return nil, status.Error(errcode.ErrRBACCantAddInternalPerm, "")
	}

	model := util.FromGPRC.Resource(res)

	model, err := srv.PermService.AddResource(ctx, model)
	if err != nil {
		return nil, err
	}
	res.Id = model.Id

	return res, nil
}

// GetPermission GetPermission
func (srv *GRPCService) GetPermission(ctx context.Context, req *rbac.PermissionID) (*rbac.Resource, error) {
	result, err := srv.PermService.GetResource(ctx, req.Id)
	if err != nil {
		return &rbac.Resource{}, err
	}
	return util.ToGRPC.Resource(result), nil
}

// GetPermission GetPermission
func (srv *GRPCService) GetPermissions(ctx context.Context, req *rbac.PermissionIDs) (*rbac.Permissions, error) {
	resList, err := srv.PermService.GetResources(ctx, req.Ids)
	if resList != nil {
		return &rbac.Permissions{
			Perms: util.ToGRPC.Resources(resList),
			Total: int64(len(resList)),
		}, nil
	}

	return nil, err
}

// GetResourcePerm GetResourcePerm
func (srv *GRPCService) GetResourcePerm(ctx context.Context, request *rbac.ResourceIdentity) (*rbac.Resource, error) {
	result, notFound, err := srv.PermService.FindResourceByNameOrExId(ctx, request)
	if err != nil {
		return &rbac.Resource{}, nil
	}

	if len(notFound) > 0 {
		return &rbac.Resource{}, status.Error(errcode.ErrRBACPermissionNotFound, errcode.NoSuchPermission)
	}

	if len(result) != 1 {
		return &rbac.Resource{}, status.Error(errcode.ErrRBACUnknown, errcode.RBACUnknown)
	}

	return util.ToGRPC.Resource(result[0]), nil
}

func (srv *GRPCService) UpdatePermission(ctx context.Context, res *rbac.Resource) (*empty.Empty, error) {
	err := srv.PermService.UpdateResource(ctx, res.Id, util.FromGPRC.Resource(res))
	return &empty.Empty{}, err
}

// DeletePermission DeletePermission
func (srv *GRPCService) DeletePermission(ctx context.Context, req *rbac.PermissionID) (*empty.Empty, error) {
	err := srv.PermService.DeleteResource(ctx, req.Id)
	return &empty.Empty{}, err
}

// ListPermission ListPermission
func (srv *GRPCService) ListPermission(ctx context.Context, listQuery *rbac.ListQuery) (*rbac.Permissions, error) {
	resources, total, err := srv.PermService.ListResource(ctx, util.FromGPRC.ListQuery(listQuery))
	if err != nil {
		return &rbac.Permissions{}, err
	}
	return &rbac.Permissions{Perms: util.ToGRPC.Resources(resources), Total: total}, nil
}

// ListObjectResources ListObjectResources
func (srv *GRPCService) ListObjectResources(ctx context.Context, request *rbac.ListObjectResourcesRequest) (*rbac.Permissions, error) {
	resList, err := srv.PermService.ListObjectResources(ctx, request.Id, request.ResourceType)
	if resList != nil {
		return &rbac.Permissions{
			Perms: util.ToGRPC.Resources(resList),
			Total: int64(len(resList)),
		}, nil
	}

	return nil, err
}

func (srv *GRPCService) ListObjectPermissions(ctx context.Context, objectID *rbac.ObjectID) (*rbac.Permissions, error) {
	resList, err := srv.PermService.ListObjectPermissions(ctx, objectID)
	if err != nil {
		return nil, err
	}
	return &rbac.Permissions{
		Perms: util.ToGRPC.Resources(resList),
		Total: int64(len(resList)),
	}, nil
}

func (srv *GRPCService) CheckResourcesPerm(ctx context.Context, req *rbac.CheckResourcesPermRequest) (*rbac.PermCheckResponse, error) {
	return srv.PermService.CheckResourcesPerm(ctx, req)
}

func (srv *GRPCService) CheckSelfPermissions(ctx context.Context, req *rbac.SimpleResources) (*rbac.PermCheckResponse, error) {
	err := srv.PermService.CheckSelfPermissions(ctx, req.Resources...)
	if err == status.Error(errcode.ErrRBACNoPermission, errcode.NoPermission) {
		return &rbac.PermCheckResponse{Pass: false}, nil
	}
	if err != nil {
		return nil, err
	}
	return &rbac.PermCheckResponse{Pass: true}, nil
}

func (srv *GRPCService) GetObjectsByResource(ctx context.Context, req *rbac.ResourceID) (*rbac.ObjectIDs, error) {
	objs, err := srv.PermService.GetObjectsByResource(ctx, req)
	if err != nil {
		return nil, err
	}
	return objs, nil
}
