package handler_rpc

import (
	"context"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/iPaaS/sso/protos/platform/idgen"

	"golang.org/x/oauth2"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"
	hydra_lcp "github.com/yuansuan/ticp/iPaaS/sso/protos"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/service"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/util"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// HydraLcpService HydraLcpService
type HydraLcpService struct {
	userSrv   *service.UserService
	ssoConfig *oauth2.Config
	*util.HydraConfig

	Idgen idgen.IdGenClient `grpc_client_inject:"idgen"`

	hydra_lcp.UnimplementedHydraLcpServiceServer
}

// GetUserInfo GetUserInfo
func (h *HydraLcpService) GetUserInfo(ctx context.Context, req *hydra_lcp.GetUserInfoReq) (resp *hydra_lcp.UserInfo, err error) {
	ysid, err := parseYsidFromBase58(ctx, req.Ysid)
	if nil != err {
		return nil, err
	}

	user := models.SsoUser{Ysid: ysid}
	ok, err := h.userSrv.Get(ctx, &user)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Errorf(consts.ErrHydraLcpDBUserNotExist, "user does not exist")
	}
	return h.userSrv.ModelToProtoUserInfo(&user), nil
}

func (h *HydraLcpService) ListUsers(ctx context.Context, req *hydra_lcp.ListUserReq) (resp *hydra_lcp.UserInfoList, err error) {
	userList, total, err := h.userSrv.List(ctx, req.Page.Index, req.Page.Size, req.Name)
	if err != nil {
		return nil, err
	}

	// assemble data
	var allUserInfo []*hydra_lcp.UserInfo
	for _, user := range userList {
		allUserInfo = append(allUserInfo, h.userSrv.ModelToProtoUserInfo(user))
	}

	return &hydra_lcp.UserInfoList{UserInfo: allUserInfo, Total: total}, nil
}

// UpdateName UpdateName
func (h *HydraLcpService) UpdateName(ctx context.Context, req *hydra_lcp.UserInfoReq) (resp *hydra_lcp.UserInfo, err error) {
	if "" == req.Param {
		logging.GetLogger(ctx).Info("param name is empty")
		return nil, status.Error(consts.ErrHydraLcpNameEmpty, "param name is empty")
	}

	ysid, err := parseYsidFromBase58(ctx, req.Ysid)
	if nil != err {
		return nil, err
	}

	info, err := h.userSrv.UpdateName(ctx, models.SsoUser{Ysid: ysid, Name: req.Param})
	return h.userSrv.ModelToProtoUserInfo(&info), err
}

// QueryInfoByPhoneNumber QueryInfoByPhoneNumber
func (h *HydraLcpService) QueryInfoByPhoneNumber(ctx context.Context, req *hydra_lcp.QueryInfoByPhoneNumberReq) (resp *hydra_lcp.UserInfo, err error) {
	if "" == req.PhoneNumber {
		logging.GetLogger(ctx).Info("param phone number is empty")
		return nil, status.Error(consts.ErrHydraLcpPhoneEmpty, "param phone number is empty")
	}

	info := models.SsoUser{Phone: req.PhoneNumber}
	ok, err := h.userSrv.Get(ctx, &info)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Errorf(consts.ErrHydraLcpDBUserNotExist, "user does not exist")
	}
	return h.userSrv.ModelToProtoUserInfo(&info), nil
}

// CheckPassword CheckPassword
func (h *HydraLcpService) CheckPassword(ctx context.Context, req *hydra_lcp.UserInfoReq) (*empty.Empty, error) {
	if "" == req.Param {
		logging.GetLogger(ctx).Info("param password is empty")
		return nil, status.Error(consts.ErrHydraLcpPasswordEmpty, "param password is empty")
	}

	ysid, err := parseYsidFromBase58(ctx, req.Ysid)
	if nil != err {
		return nil, err
	}

	_, err = h.userSrv.VerifyPasswordByUserID(ctx, ysid, req.Param)

	return &empty.Empty{}, err
}

func (h *HydraLcpService) CheckPassword2(ctx context.Context, req *hydra_lcp.CheckPasswordReq) (*hydra_lcp.UserInfo, error) {
	if req.Password == "" {
		logging.GetLogger(ctx).Info("param password is empty")
		return nil, status.Error(consts.ErrHydraLcpPasswordEmpty, "param password is empty")
	}
	var id int64
	var err error
	if req.Ysid != "" {
		parseID, parseErr := parseYsidFromBase58(ctx, req.Ysid)
		if parseErr != nil {
			return nil, parseErr
		}
		id, err = h.userSrv.VerifyPasswordByUserID(ctx, parseID, req.Password)
		if err != nil {
			return nil, err
		}
	} else if req.Phone != "" {
		id, err = h.userSrv.VerifyPasswordByPhone(ctx, req.Phone, req.Password)
		if err != nil {
			return nil, err
		}
	} else if req.Email != "" {
		id, err = h.userSrv.VerifyPasswordByEmail(ctx, req.Email, req.Password)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, status.Error(consts.ErrHydraLcpBadRequest, "ys id, phone, email all are empty")
	}
	user := &models.SsoUser{Ysid: id}
	_, err = h.userSrv.Get(ctx, user)
	if err != nil {
		return nil, err
	}
	return h.userSrv.ModelToProtoUserInfo(user), nil
}

func parseYsidFromBase58(ctx context.Context, ysid string) (int64, error) {
	logger := logging.GetLogger(ctx)
	// check ysid empty
	if "" == ysid {
		logger.Info("failed to decrypt ysid, ysid is empty")
		return -1, status.Error(consts.ErrHydraLcpYsidEmpty, "ysid is empty")
	}
	// decrypt ysid
	return snowflake.MustParseString(ysid).Int64(), nil
}

// InitGRPCServer InitGRPCServer
func InitGRPCServer(drv *http.Driver) {
	s, err := boot.GRPC.DefaultServer()
	if err != nil {
		panic(err)
	}
	hander := &HydraLcpService{
		userSrv:     service.NewUserSrv(),
		HydraConfig: util.GetHydraConfig(),
	}

	grpc_boot.InjectAllClient(hander)
	hydra_lcp.RegisterHydraLcpServiceServer(s.Driver(), hander)

}
