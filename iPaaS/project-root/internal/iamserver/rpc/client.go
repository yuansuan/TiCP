package rpc

import (
	"context"
	"sync"

	hydra_lcp "github.com/yuansuan/ticp/iPaaS/sso/protos"
	"github.com/yuansuan/ticp/iPaaS/sso/protos/platform/ptype"

	grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	"google.golang.org/grpc/status"
)

type Client struct {
	HydraLcp struct {
		HydraLcp hydra_lcp.HydraLcpServiceClient `grpc_client_inject:"hydra_lcp"`
	}
}

var Instance *Client

var mtx sync.Mutex

// GetInstance get instance
func GetInstance() *Client {
	mtx.Lock()
	defer mtx.Unlock()

	if Instance == nil {
		Instance = &Client{}
		grpc_boot.InjectAllClient(Instance)
	}
	return Instance
}

func (c *Client) GetUser(ctx context.Context, userID string) (*hydra_lcp.UserInfo, error) {
	reply, err := GetInstance().HydraLcp.HydraLcp.GetUserInfo(ctx, &hydra_lcp.GetUserInfoReq{
		Ysid: userID,
	})

	if err != nil {
		trace.GetLogger(ctx).Warnf("get user info failed, err: %v", err)
		if s, ok := status.FromError(err); ok && s.Code() == consts.ErrHydraLcpDBUserNotExist {
			return nil, nil
		}
		return nil, err
	}

	return reply, nil
}

func (c *Client) ListUsers(ctx context.Context, offset, limit int64, name string) (*hydra_lcp.UserInfoList, error) {
	ptypePage := &ptype.Page{
		Index: offset,
		Size:  limit,
	}
	reply, err := GetInstance().HydraLcp.HydraLcp.ListUsers(ctx, &hydra_lcp.ListUserReq{
		Page: ptypePage,
		Name: name,
	})
	if err != nil {
		trace.GetLogger(ctx).Warnf("list users failed, err: %v", err)
		return nil, err
	}
	return reply, nil
}

func (c *Client) AddUser(ctx context.Context, user *hydra_lcp.AddUserReq) (error, string) {
	reply, err := GetInstance().HydraLcp.HydraLcp.AddUser(ctx, user)
	if err != nil {
		trace.GetLogger(ctx).Warnf("add user failed, err: %v", err)
		return err, ""
	}
	return nil, reply.UserId
}

func (c *Client) AddUser2(ctx context.Context, user *hydra_lcp.AddUserReq) (*hydra_lcp.AddUserResp, error) {
	return GetInstance().HydraLcp.HydraLcp.AddUser2(ctx, user)
}

func (c *Client) UpdateUser(ctx context.Context, user *hydra_lcp.UserInfoReq) error {
	_, err := GetInstance().HydraLcp.HydraLcp.UpdateName(ctx, user)
	if err != nil {
		trace.GetLogger(ctx).Warnf("update user failed, err: %v", err)
		return err
	}
	return nil
}

func (c *Client) QueryInfoByPhoneNumber(ctx context.Context, phoneNumber, password string) (*hydra_lcp.UserInfo, error) {
	reply, err := GetInstance().HydraLcp.HydraLcp.QueryInfoByPhoneNumber(ctx, &hydra_lcp.QueryInfoByPhoneNumberReq{
		PhoneNumber: phoneNumber,
	})
	if err != nil {
		trace.GetLogger(ctx).Infof("query info by phone number failed, Phone: %s, err: %v", phoneNumber, err)
		return nil, err
	}
	_, err = GetInstance().HydraLcp.HydraLcp.CheckPassword(ctx, &hydra_lcp.UserInfoReq{
		Ysid:  reply.Ysid,
		Param: password,
	})
	return reply, err
}

func (c *Client) CheckPassword2(ctx context.Context, ysid, phone, email, password string) (*hydra_lcp.UserInfo, error) {
	reply, err := GetInstance().HydraLcp.HydraLcp.CheckPassword2(ctx, &hydra_lcp.CheckPasswordReq{
		Ysid:     ysid,
		Phone:    phone,
		Email:    email,
		Password: password,
	})
	if err != nil {
		trace.GetLogger(ctx).Infof("check password failed, err: %v", err)
		return nil, err
	}
	return reply, nil
}
