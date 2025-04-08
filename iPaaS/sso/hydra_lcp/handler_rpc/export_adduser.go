package handler_rpc

import (
	"context"
	"fmt"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/util"
	"github.com/yuansuan/ticp/iPaaS/sso/protos"
	"github.com/yuansuan/ticp/iPaaS/sso/protos/platform/idgen"
)

// AddUser 在SSO中添加一个用户
func (h *HydraLcpService) AddUser(ctx context.Context, req *protos.AddUserReq) (*protos.AddUserResp, error) {
	if !util.IsPhoneValid(req.Phone) {
		return nil, status.Error(consts.ErrHydraLcpPhoneInvalidate, "invalid phone")
	}
	if !util.IsPasswordValid(req.Password) {
		return nil, status.Error(consts.ErrHydraLcpPwdInvalidate, "invalid password")
	}

	_, err := h.userSrv.GetID(ctx, "", req.Phone, "")
	if err == nil {
		return nil, status.Errorf(consts.ErrHydraLcpUserExist, "user already exists")
	}

	sts, ok := status.FromError(err)
	if !ok || sts.Code() != consts.ErrHydraLcpDBUserNotExist {
		return nil, status.Errorf(consts.Unknown, "unknown error occurs")
	}

	gen, err := h.Idgen.GenerateID(ctx, &idgen.GenRequest{})
	if err != nil {
		logging.Default().Warnw("unable to generate user id", "error", err, "req", req)
		return nil, status.Errorf(consts.Unknown, "unable to generate user id")
	}

	if len(req.UserChannel) == 0 {
		req.UserChannel = "UNKNOWN"
	}
	if len(req.UserSource) == 0 {
		req.UserSource = "未知来源"
	}

	err = h.userSrv.Add(ctx, &models.SsoUser{
		Ysid:          gen.Id,
		Phone:         req.Phone,
		WechatUnionId: "",
		UserChannel:   req.UserChannel,
		UserSource:    req.UserSource,
	}, req.Password)
	if err != nil {
		logging.Default().Warnw("unable to create the user", "error", err, "req", req)
		return nil, status.Errorf(consts.Unknown, "unable to create the user")
	}

	return &protos.AddUserResp{UserId: snowflake.ID(gen.Id).String()}, nil
}

func (h *HydraLcpService) AddUser2(ctx context.Context, req *protos.AddUserReq) (*protos.AddUserResp, error) {
	if !util.IsPasswordValid(req.Password) {
		return nil, status.Error(consts.ErrHydraLcpPwdInvalidate, "invalid password")
	}
	gen, err := h.Idgen.GenerateID(ctx, &idgen.GenRequest{})
	if err != nil {
		logging.Default().Errorf("unable to generate user id", "error", err, "req", req)
		return nil, status.Errorf(consts.Unknown, "unable to generate user id")
	}

	if len(req.UserChannel) == 0 {
		req.UserChannel = "UNKNOWN"
	}

	err = h.userSrv.Add(ctx, &models.SsoUser{
		Ysid:        gen.Id,
		Phone:       req.Phone,
		UserChannel: req.UserChannel,
		Name:        req.Name,
		Company:     req.CompanyName,
		Email:       req.Email,
	}, req.Password)

	if err != nil {
		sta, ok := status.FromError(err)
		if ok && sta.Code() == consts.ErrHydraLcpDBDuplicatedEntry {
			logging.Default().Warnw("User existed", "request", req)
			return nil, status.Errorf(consts.ErrHydraLcpUserExist, fmt.Sprintf("%s existed", req.Phone))
		}
		logging.Default().Warnw("unable to create the user", "error", err, "req", req)
		return nil, status.Errorf(consts.Unknown, "unable to create the user")
	}
	userResp := &protos.AddUserResp{
		UserId: snowflake.ID(gen.Id).String(),
	}
	return userResp, nil
}
