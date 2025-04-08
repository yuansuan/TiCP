package handler_rpc

import (
	"context"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/service"

	hydra_lcp "github.com/yuansuan/ticp/iPaaS/sso/protos"
)

// GetExternalUserInfo GetExternalUserInfo
func (h *HydraLcpService) GetExternalUserInfo(ctx context.Context, req *hydra_lcp.GetExternalUserInfoReq) (resp *hydra_lcp.GetExternalUserInfoReply, err error) {
	ysid, err := parseYsidFromBase58(ctx, req.Ysid)
	if nil != err {
		return nil, err
	}

	user := models.SsoExternalUser{Ysid: ysid}
	err = service.ExternalUserService.Get(ctx, &user)
	if err != nil {
		return nil, err
	}

	if user.UserName == "" {
		return &hydra_lcp.GetExternalUserInfoReply{}, err
	}

	return &hydra_lcp.GetExternalUserInfoReply{
		Ysid:     snowflake.ID(user.Ysid).String(),
		UserName: user.UserName,
	}, err
}

// CheckExternalUserExist CheckExternalUserExist
func (h *HydraLcpService) CheckExternalUserExist(ctx context.Context, req *hydra_lcp.CheckExternalUserExistReq) (resp *hydra_lcp.CheckExternalUserExistReply, err error) {

	exist, err := h.ldapSvr.SearchUser(req.UserName)

	return &hydra_lcp.CheckExternalUserExistReply{
		IsExist: exist,
	}, err

}
