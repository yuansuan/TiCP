package rpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	userMgr "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

// Exist Exist
func (s *GRPCService) Exist(ctx context.Context, in *userMgr.UserIdentity) (*userMgr.UserExistResponse, error) {
	exist, err := s.UserService.Exist(ctx, snowflake.MustParseString(in.Id).Int64())
	if err != nil {
		return nil, err
	}

	return &userMgr.UserExistResponse{
		Exist: exist,
	}, nil
}

// Get get user
func (s *GRPCService) Get(ctx context.Context, in *userMgr.UserIdentity) (*userMgr.UserObj, error) {
	user, err := s.UserService.Get(ctx, snowflake.MustParseString(in.Id).Int64())
	if err != nil {
		return nil, err
	}
	return util.ToGRPC.UserObj(user), nil
}

// GetIncludeDeleted get user
func (s *GRPCService) GetIncludeDeleted(ctx context.Context, in *userMgr.UserIdentity) (*userMgr.UserObj, error) {
	user, err := s.UserService.GetIncludeDeleted(ctx, snowflake.MustParseString(in.Id).Int64())
	if err != nil {
		return nil, err
	}
	return util.ToGRPC.UserObj(user), nil
}

// GetIdByName GetIdByName
func (s *GRPCService) GetIdByName(ctx context.Context, in *userMgr.NameCondRequest) (*userMgr.UserIdentity, error) {
	id, err := s.UserService.GetIdByName(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return &userMgr.UserIdentity{
		Id: snowflake.ID(id).String(),
	}, nil
}

// GetIdByName GetIdByName
func (s *GRPCService) GetUserByName(ctx context.Context, in *userMgr.NameCondRequest) (*userMgr.UserObj, error) {
	user, err := s.UserService.GetUserByName(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return util.ToGRPC.UserObj(user), nil
}

// BatchGetUser Batch Get User
func (s *GRPCService) BatchGetUser(ctx context.Context, in *userMgr.UserIdentities) (*userMgr.BatchUsersResponse, error) {
	users, err := s.UserService.BatchGetUser(ctx, util.FromGPRC.UserIDs(in.UserIdentities))

	return &userMgr.BatchUsersResponse{
		Success: true,
		Total:   int64(len(users)),
		UserObj: util.ToGRPC.UserObjs(users),
	}, err

}

// GetAllUserName GetAllUserName
func (s *GRPCService) GetAllUserName(ctx context.Context, _ *userMgr.GetAllUserRequest) (*userMgr.GetAllUserResponse, error) {
	users, err := s.UserService.GetAllUser(ctx)
	if err != nil {
		return nil, err
	}

	userNameList := make([]string, 0)

	for _, user := range users {
		userNameList = append(userNameList, user.Name)
	}

	return &userMgr.GetAllUserResponse{
		Names: userNameList,
	}, err
}

func (s *GRPCService) AddUserWithRole(ctx context.Context, req *userMgr.AddUserWithRoleRequest) (*userMgr.UserIdentity, error) {
	userId, err := s.UserService.AddUserWithRole(ctx, dto.UserAddRequest{
		Name:          req.Name,
		Password:      req.Password,
		Email:         req.Email,
		Mobile:        req.Mobile,
		RealName:      req.RealName,
		Roles:         req.RoleIds,
		EnableOpenapi: req.EnableOpenapi,
	})

	return &userMgr.UserIdentity{
		Id: snowflake.ID(userId).String(),
	}, err
}

func (s *GRPCService) UpdateUser(ctx context.Context, req *userMgr.UpdateUserRequest) (*empty.Empty, error) {
	err := s.UserService.Update(ctx, model.User{
		Id:            snowflake.MustParseString(req.Id).Int64(),
		Email:         req.Email,
		Mobile:        req.Mobile,
		EnableOpenapi: req.EnableOpenapi,
	})

	if err != nil {
		return nil, err
	}

	_, err = client.GetInstance().Role.UpdateObjectRoles(ctx, &rbac.ObjectRoles{
		Id: &rbac.ObjectID{
			Id:   req.Id,
			Type: rbac.ObjectType_USER,
		},
		Roles: req.RoleIds,
	})

	return &empty.Empty{}, err
}

func (s *GRPCService) DelUser(ctx context.Context, req *userMgr.UserIdentity) (*empty.Empty, error) {
	userID := snowflake.MustParseString(req.Id).Int64()
	user, err := s.UserService.Get(ctx, userID)

	if err != nil || user == nil {
		return nil, status.Error(errcode.ErrUserNotExist, "")
	}

	err = s.UserService.Delete(ctx, userID)

	if err != nil {
		return nil, err
	}

	client.GetInstance().Role.RemoveObjectRoles(ctx, &rbac.ObjectRoles{
		Id: &rbac.ObjectID{
			Id:   req.Id,
			Type: rbac.ObjectType_USER,
		},
	})

	return &empty.Empty{}, err
}

func (s *GRPCService) EnableUser(ctx context.Context, req *userMgr.EnableUserRequest) (*empty.Empty, error) {
	var err error
	if req.Enable {
		err = s.UserService.ActiveUser(ctx, snowflake.MustParseString(req.Id).Int64())
	} else {
		err = s.UserService.InactiveUser(ctx, snowflake.MustParseString(req.Id).Int64())

	}

	return &empty.Empty{}, err
}

func (s *GRPCService) GetUserRoleNames(ctx context.Context, req *userMgr.UserIdentity) (*userMgr.GetUserRoleNamesResponse, error) {
	return &userMgr.GetUserRoleNamesResponse{
		RoleNames: s.UserService.GetUserRoleNames(ctx, req.Id),
	}, nil
}
