package handler_rpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/util"
	"github.com/yuansuan/ticp/iPaaS/sso/protos"
)

// ResetPassword 直接重置用户的密码
func (h *HydraLcpService) ResetPassword(ctx context.Context, req *protos.ResetPasswordReq) (*emptypb.Empty, error) {
	if !util.IsPasswordValid(req.NewPwd) {
		return nil, status.Error(consts.ErrHydraLcpPwdInvalidate, "invalid password")
	}

	uid, err := parseYsidFromBase58(ctx, req.Ysid)
	if nil != err {
		return nil, err
	}

	// get user info
	var userInfo models.SsoUser
	userInfo.Ysid = uid
	ok, err := h.userSrv.Get(ctx, &userInfo)
	if err != nil {
		return nil, err
	} else if !ok {
		logging.GetLogger(ctx).Infof("user with ID: %v id not exist", uid)
		return nil, status.Error(consts.ErrHydraLcpDBUserNotExist, "user not exist")
	}

	err = h.userSrv.UpdatePwd(ctx, models.SsoUser{Ysid: uid}, req.NewPwd)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
