package user

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
	v1 "github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/controller/v1"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/rpc"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
	hydra_lcp "github.com/yuansuan/ticp/iPaaS/sso/protos"
	"google.golang.org/grpc/status"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	"time"
)

func AdminGetUser(c *gin.Context) {
	userID := strings.TrimSpace(c.Param("userId"))

	u, err := rpc.GetInstance().GetUser(c, userID)
	if u == nil && err == nil {
		// user not exist
		common.ErrorRespWithAbort(c, http.StatusNotFound, v1.UserNotFound, "user not exist")
	}

	if err != nil {
		common.InternalServerError(c, "")
		return
	}

	user := &iam_api.AdminGetUserResponse{
		UserId: u.Ysid,
		Name:   u.Name,
		Phone:  u.Phone,
	}
	common.SuccessResp(c, user)
}

func AdminListUser(c *gin.Context) {
	pageOffset := c.Query("PageOffset")
	pageSize := c.Query("PageSize")
	name := c.Query("Name")

	if pageOffset == "" || pageSize == "" {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "pageOffset or pageSize is empty")
		return
	}
	// to int64
	offset, err := strconv.ParseInt(pageOffset, 10, 64)
	if err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "pageOffset is invalid")
		return
	}
	size, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "pageSize is invalid")
		return
	}
	if offset == 0 {
		offset = 1
	}

	users, err := rpc.GetInstance().ListUsers(c, offset, size, name)
	if err != nil {
		common.InternalServerError(c, "")
		return
	}

	us := make([]*iam_api.User, 0, len(users.UserInfo))
	for _, u := range users.UserInfo {
		t := protoTimeToTime(u.CreateTime)
		us = append(us, &iam_api.User{
			Ysid:            u.Ysid,
			Name:            u.Name,
			Email:           u.Email,
			Phone:           u.Phone,
			RealName:        u.RealName,
			UserName:        u.UserName,
			DisplayUserName: u.DisplayUserName,
			UserChannel:     u.UserChannel,
			UserSource:      u.UserSource,
			UserReferer:     u.UserReferer,
			CreateTime:      t,
		})
	}
	res := &iam_api.AdminListUserByNameResponse{
		Users: us,
		Total: users.Total,
	}
	common.SuccessResp(c, res)
}

func protoTimeToTime(t *timestamppb.Timestamp) time.Time {
	return time.Unix(t.Seconds, int64(t.Nanos))
}

func AdminAddUser(c *gin.Context) {
	req := &iam_api.AdminAddUserRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, err.Error())
		return
	}

	if req.Password == "" || req.Phone == "" {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "name or phone is empty")
		return
	}

	user := &hydra_lcp.AddUserReq{
		Phone:       req.Phone,
		Password:    req.Password,
		UserSource:  "ysadmin",
		UserChannel: "project-root",
	}

	err, userID := rpc.GetInstance().AddUser(c, user)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			if s.Code() == consts.ErrHydraLcpPwdInvalidate {
				common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "password invalidate")
				return
			}
			if s.Code() == consts.ErrHydraLcpUserExist {
				common.ErrorRespWithAbort(c, http.StatusConflict, v1.AlreadyExists, "user exist")
				return
			}
		}
		common.InternalServerError(c, "")
		return
	}

	res := &iam_api.AdminAddUserResponse{
		UserId: userID,
	}

	common.SuccessResp(c, res)
}

func AdminUpdateUser(c *gin.Context) {
	userID := strings.TrimSpace(c.Param("userId"))

	req := &iam_api.AdminUpdateUserRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, err.Error())
		return
	}

	if req.Name == "" {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "name is empty")
		return
	}

	user := &hydra_lcp.UserInfoReq{
		Ysid:  userID,
		Param: req.Name,
	}

	err := rpc.GetInstance().UpdateUser(c, user)
	if err != nil {
		s, ok := status.FromError(err)
		if ok && s.Code() == consts.ErrHydraLcpDBUserNotExist {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.UserNotFound, "user not exist")
			return
		}
		if ok && s.Code() == consts.ErrHydraLcpLackYsID {
			common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "userID is invalid")
			return
		}
		common.InternalServerError(c, "")
		return
	}

	common.SuccessResp(c, nil)
}
