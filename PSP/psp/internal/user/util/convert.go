package util

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type toGRPC struct{}

var ToGRPC = toGRPC{}

func (toGRPC) UserObj(userModel *model.User) *user.UserObj {
	if userModel == nil {
		return nil
	}

	return &user.UserObj{
		Id:         snowflake.ID(userModel.Id).String(),
		Name:       userModel.Name,
		Email:      userModel.Email,
		Mobile:     userModel.Mobile,
		Enabled:    userModel.Enabled,
		CreatedAt:  userModel.CreatedAt.UnixNano() / int64(time.Millisecond),
		IsInternal: userModel.IsInternal,
		RealName:   "",
	}
}

func (toGRPC) UserObjs(users []*model.User) []*user.UserObj {
	result := make([]*user.UserObj, 0, len(users))
	for _, role := range users {
		result = append(result, ToGRPC.UserObj(role))
	}
	return result
}

type fromGRPC struct{}

var FromGPRC = fromGRPC{}

func (fromGRPC) UserID(identity *user.UserIdentity) int64 {
	if identity == nil {
		return 0
	}

	return snowflake.MustParseString(identity.Id).Int64()
}

func (fromGRPC) UserIDs(identitys []*user.UserIdentity) []int64 {
	result := make([]int64, 0, len(identitys))
	for _, role := range identitys {
		result = append(result, FromGPRC.UserID(role))
	}
	return result
}

type toDTO struct{}

var ToDTO = toDTO{}

func (toDTO) UserObj(userModel *model.User, roleIds []int64) *dto.UserInfo {
	if userModel == nil {
		return nil
	}

	return &dto.UserInfo{
		Id:            snowflake.ID(userModel.Id).String(),
		Name:          userModel.Name,
		Email:         userModel.Email,
		Mobile:        userModel.Mobile,
		Enabled:       userModel.Enabled,
		CreatedAt:     userModel.CreatedAt.UnixNano() / int64(time.Millisecond),
		IsInternal:    userModel.IsInternal,
		RealName:      userModel.RealName,
		EnableOpenapi: userModel.EnableOpenapi,
		Roles:         roleIds,
	}
}
