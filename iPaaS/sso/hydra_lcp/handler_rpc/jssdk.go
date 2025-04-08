package handler_rpc

import (
	"context"

	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/service"

	hydra_lcp "github.com/yuansuan/ticp/iPaaS/sso/protos"
)

// GetJsConfig 获取JSSDK的页面配置（签名）
func (h *HydraLcpService) GetJsConfig(_ context.Context, in *hydra_lcp.GetJsConfigReq) (*hydra_lcp.GetJsConfigResp, error) {
	if len(in.Uri) == 0 {
		return nil, status.Error(consts.InvalidArgument, "uri required")
	}

	js := service.GetOfficialAccount().GetJs()
	config, err := js.GetConfig(in.GetUri())
	if err != nil {
		return nil, status.Errorf(consts.Unknown, "wechat error: %s", err)
	}

	return &hydra_lcp.GetJsConfigResp{
		AppId:     config.AppID,
		Timestamp: config.Timestamp,
		NonceStr:  config.NonceStr,
		Signature: config.Signature,
	}, nil
}

// GetJsTicket 获取JSSDK的临时票据
func (h *HydraLcpService) GetJsTicket(context.Context, *hydra_lcp.GetJsTicketReq) (*hydra_lcp.GetJsTicketResp, error) {
	js := service.GetOfficialAccount().GetJs()

	ak, err := js.GetAccessToken()
	if err != nil {
		return nil, status.Errorf(consts.Unknown, "wechat error: %s", err)
	}

	ticket, err := js.GetTicket(ak)
	if err != nil {
		return nil, status.Errorf(consts.Unknown, "wechat error: %s", err)
	}

	return &hydra_lcp.GetJsTicketResp{
		Ticket: ticket,
	}, nil
}
