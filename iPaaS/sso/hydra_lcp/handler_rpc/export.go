package handler_rpc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/yuansuan/ticp/iPaaS/sso/protos/platform/idgen"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	ginUtil "github.com/yuansuan/ticp/common/go-kit/gin-boot/util"

	"github.com/silenceper/wechat/v2/officialaccount/menu"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"

	hydra_lcp "github.com/yuansuan/ticp/iPaaS/sso/protos"
	"github.com/yuansuan/ticp/iPaaS/sso/protos/platform/ptype"

	"github.com/yuansuan/ticp/iPaaS/sso/protos/platform/company"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/rpc"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/common"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/config"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/service"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/util"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// HydraLcpService HydraLcpService
type HydraLcpService struct {
	userSrv        *service.UserService
	phoneSvr       *service.PhoneService
	ldapSvr        *service.LdapService
	offiaccountSrv *service.OffiaccountBindingService
	ssoConfig      *oauth2.Config
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

// GetUserInfoBatch GetUserInfoBatch
func (h *HydraLcpService) GetUserInfoBatch(ctx context.Context, req *hydra_lcp.GetUserInfoBatchReq) (resp *hydra_lcp.UserInfoBatch, err error) {
	if len(req.Ysid) == 0 {
		return &hydra_lcp.UserInfoBatch{}, nil
	}

	// convert ysid from base58 string to int64
	var ysid []int64
	for _, v := range req.Ysid {
		id, err := parseYsidFromBase58(ctx, v)
		if nil != err {
			return nil, err
		}
		ysid = append(ysid, id)
	}

	// get all user in batch mode
	userInfo, err := h.userSrv.GetBatch(ctx, ysid)
	if err != nil {
		return nil, err
	}

	if len(userInfo) == 0 {
		return &hydra_lcp.UserInfoBatch{}, nil
	}

	// assemble data
	var allUserInfo []*hydra_lcp.UserInfo
	for _, user := range userInfo {
		allUserInfo = append(allUserInfo, h.userSrv.ModelToProtoUserInfo(user))
	}

	return &hydra_lcp.UserInfoBatch{UserInfo: allUserInfo}, nil
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

// UpdatePhoneNumber UpdatePhoneNumber
func (h *HydraLcpService) UpdatePhoneNumber(ctx context.Context, req *hydra_lcp.UpdatePhoneNumberReq) (resp *hydra_lcp.UserInfo, err error) {
	if "" == req.PhoneNumberNew {
		logging.GetLogger(ctx).Info("param phone number is empty")
		return nil, status.Error(consts.ErrHydraLcpPhoneEmpty, "param phone number is empty")
	}

	if h.userSrv.CheckUserExists(ctx, models.SsoUser{Phone: req.PhoneNumberNew}) == nil {
		return nil, status.Errorf(consts.ErrHydraLcpPhoneExist, "phone number(%v) exists", req.PhoneNumberNew)
	}

	err = h.phoneSvr.VerifyCode(ctx, req.PhoneNumberNew, req.Captcha)
	if err != nil {
		return nil, err
	}

	ysid, err := parseYsidFromBase58(ctx, req.Ysid)
	if nil != err {
		return nil, err
	}

	info, err := h.userSrv.UpdatePhone(ctx, models.SsoUser{Ysid: ysid, Phone: req.PhoneNumberNew})
	return h.userSrv.ModelToProtoUserInfo(&info), err
}

// UpdateEmail UpdateEmail
func (h *HydraLcpService) UpdateEmail(ctx context.Context, req *hydra_lcp.UserInfoReq) (resp *hydra_lcp.UserInfo, err error) {
	ysid, err := parseYsidFromBase58(ctx, req.Ysid)
	if nil != err {
		return nil, err
	}

	info, err := h.userSrv.UpdateEmail(ctx, models.SsoUser{Ysid: ysid, Email: req.Param})
	return h.userSrv.ModelToProtoUserInfo(&info), err
}

// UpdateWechatInfo UpdateWechatInfo
func (h *HydraLcpService) UpdateWechatInfo(ctx context.Context, req *hydra_lcp.WechatInfoReq) (resp *hydra_lcp.UserInfo, err error) {
	ysid, err := parseYsidFromBase58(ctx, req.Ysid)
	if nil != err {
		return nil, err
	}

	info, err := h.userSrv.UpdateWechatInfo(ctx, models.SsoUser{
		Ysid: ysid, WechatUnionId: req.WechatUnionId, WechatOpenId: req.WechatOpenId, WechatNickName: req.WechatNickName})
	return h.userSrv.ModelToProtoUserInfo(&info), err
}

// UpdateRealName UpdateRealName
func (h *HydraLcpService) UpdateRealName(ctx context.Context, req *hydra_lcp.UserInfoReq) (resp *hydra_lcp.UserInfo, err error) {
	ysid, err := parseYsidFromBase58(ctx, req.Ysid)
	if nil != err {
		return nil, err
	}

	info, err := h.userSrv.UpdateRealName(ctx, models.SsoUser{Ysid: ysid, RealName: req.Param})
	return h.userSrv.ModelToProtoUserInfo(&info), err
}

// UpdateHeadimg UpdateRealName
func (h *HydraLcpService) UpdateHeadimg(ctx context.Context, req *hydra_lcp.UserInfoReq) (resp *hydra_lcp.UserInfo, err error) {
	ysid, err := parseYsidFromBase58(ctx, req.Ysid)
	if nil != err {
		return nil, err
	}

	info, err := h.userSrv.UpdateHeadimg(ctx, models.SsoUser{Ysid: ysid, HeadimgUrl: req.Param})
	return h.userSrv.ModelToProtoUserInfo(&info), err
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

// UpdatePassword UpdatePassword
func (h *HydraLcpService) UpdatePassword(ctx context.Context, req *hydra_lcp.UpdatePasswordReq) (*empty.Empty, error) {
	if req.PasswordNew == "" || req.Captcha == "" {
		logging.GetLogger(ctx).Info("password or captcha is empty")
		return nil, status.Error(consts.ErrHydraLcpPasswordEmpty, "password or captcha is empty")
	}

	ysid, err := parseYsidFromBase58(ctx, req.Ysid)
	if nil != err {
		return nil, err
	}
	// get user info
	var userInfo models.SsoUser
	userInfo.Ysid = ysid
	ok, err := h.userSrv.Get(ctx, &userInfo)
	if err != nil {
		return nil, err
	}
	if !ok {
		logging.GetLogger(ctx).Infof("user with ID: %v id not exist", ysid)
		return nil, status.Error(consts.ErrHydraLcpDBUserNotExist, "user not exist")
	}

	// verify phone code
	err = h.phoneSvr.VerifyCode(ctx, userInfo.Phone, req.Captcha)
	if err != nil {
		return nil, err
	}

	err = h.userSrv.UpdatePwd(ctx, models.SsoUser{Ysid: ysid}, req.PasswordNew)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// UpdatePasswordByOld UpdatePassword
func (h *HydraLcpService) UpdatePasswordByOld(ctx context.Context, req *hydra_lcp.UpdatePasswordByOldReq) (*empty.Empty, error) {
	if "" == req.OldPwd {
		logging.GetLogger(ctx).Info("old password is empty")
		return nil, status.Error(consts.ErrHydraLcpPasswordEmpty, "old password is empty")
	}

	ysid, err := parseYsidFromBase58(ctx, req.YsId)
	if nil != err {
		return nil, err
	}

	_, err = h.userSrv.VerifyPasswordByUserID(ctx, ysid, req.OldPwd)
	if err != nil {
		logging.GetLogger(ctx).Info("does not match the original password ")
		fmt.Println("err", err)
		return nil, status.Error(consts.ErrHydraLcpPasswordEmpty, "does not match the original password")
	}

	err = h.userSrv.UpdatePwd(ctx, models.SsoUser{Ysid: ysid}, req.NewPwd)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// AddUserFromPhone  通过手机号添加用户
func (h *HydraLcpService) AddUserFromPhone(ctx context.Context, req *hydra_lcp.AddUserFromPhoneReq) (*hydra_lcp.UserInfo, error) {
	phone := strings.TrimSpace(req.Phone)
	// 手机号验证
	if m, _ := regexp.MatchString("^1[0-9]{10}$", phone); !m {
		return nil, status.Error(consts.ErrHydraLcpPhoneInvalidate, "phone invalidate")
	}

	userModel := &models.SsoUser{}
	userModel.Phone = phone
	//	get userID, if user not exist, add User
	ok, err := h.userSrv.Get(ctx, userModel)
	if err != nil {
		return nil, err
	}
	if ok {
		return h.userSrv.ModelToProtoUserInfo(userModel), nil
	}
	// 用户不存在，创建用户
	id, err := rpc.GenID(ctx)
	if err != nil {
		return nil, err
	}

	userModel.Ysid = id.Int64()
	err = h.userSrv.Add(ctx, userModel, util.RandomString(20))
	if err != nil {
		return nil, err
	}

	// 发送短信
	sign := service.SignYS
	if req.ProductId == consts.ZSWLProductID.String() {
		sign = service.SignZS
	}

	doNotSend := false
	if req.CompanyId != "" {
		companyID := snowflake.MustParseString(req.CompanyId)
		logger := logging.GetLogger(ctx)
		company, err := rpc.GetInstance().Company.Company.GetCompany(ctx, &company.GetCompanyRequest{
			CompanyId: req.CompanyId,
		})
		if err != nil {
			doNotSend = true
			logger.Errorf("[AddUserFromPhone] get company error: %v, companyID: %s", err, companyID)
		} else if company.Domain != "" {
			t, ok := config.Custom.SmsConfig.DomainMap[company.Domain]
			if ok {
				doNotSend = true
				err = h.phoneSvr.SendSms(ctx, phone, []string{}, t.Key, sign)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	if !doNotSend {
		err = h.phoneSvr.SendSms(ctx, phone, []string{}, common.SMSTemplateSignupSuccess, sign)
		if err != nil {
			return nil, err
		}
	}

	return h.userSrv.ModelToProtoUserInfo(userModel), nil
}

// SendSms 发送短信
func (h *HydraLcpService) SendSms(ctx context.Context, req *hydra_lcp.SendSmsReq) (*hydra_lcp.SendSmsResp, error) {
	// 默认签名
	sign := service.SignYS
	if req.ProductId != "" { // 参数优先级最高
		if req.ProductId == consts.ZSWLProductID.String() {
			sign = service.SignZS
		}
	} else {
		// metadata中优先级次之
		productIDStr, _ := ginUtil.GetInMetadata(ctx, "product_id")
		if productIDStr != "" {
			if productIDStr == consts.ZSWLProductID.String() {
				sign = service.SignZS
			}
		}
	}

	// 发送短信
	err := h.phoneSvr.SendSms(ctx, req.Phone, req.Param, req.TplName, sign)
	if err != nil {
		return nil, err
	}

	// 成功
	return &hydra_lcp.SendSmsResp{
		IsSucceed: true,
	}, nil
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

// GetOffiaccountBinding 获取公众号绑定(已激活)数据
func (h *HydraLcpService) GetOffiaccountBinding(ctx context.Context, req *hydra_lcp.GetOffiaccountBindingReq) (*hydra_lcp.OffiaccountBinding, error) {
	offiacctBinding := models.OffiaccountBinding{
		Platform:              req.Platform,
		WechatOpenid:          req.WechatOpenid,
		NotificationType:      req.NotificationType,
		NotificationActivated: 1,
	}
	// 远算云平台
	if req.Platform == models.OffiaccountBindingPlatformCloud {
		if req.NotificationType == models.OffiaccountBindingNotificationTypeJob {
			// 作业通知才验证userID
			offiacctBinding.UserId = snowflake.MustParseString(req.UserId)
		}
		// 远算云平台对企业用户进行验证
		if req.CompanyId != "" {
			companyUserRelation := models.CompanyUserRelation{
				UserId:    snowflake.MustParseString(req.UserId),
				CompanyId: snowflake.MustParseString(req.CompanyId),
				Status:    1,
			}
			session := boot.MW.DefaultSession(ctx)
			defer session.Close()
			exist, err := session.Get(&companyUserRelation)
			if err != nil {
				return nil, err
			}
			if exist {
				offiacctBinding.CompanyId = snowflake.MustParseString(req.CompanyId)
			} else {
				return nil, status.Error(consts.InvalidParam, "company_id is invalid")
			}
		}
	}
	ok, err := h.offiaccountSrv.GetOffiaccountBinding(ctx, &offiacctBinding)
	if err != nil {
		return nil, err
	}
	if !ok {
		return h.offiaccountSrv.ModelToProtoOffiaccountBinding(&models.OffiaccountBinding{}), nil
	}
	return h.offiaccountSrv.ModelToProtoOffiaccountBinding(&offiacctBinding), nil
}

// GetOffiaccountBalanceSubList 公众号余额通知关注列表 内部账号使用。
func (h *HydraLcpService) GetOffiaccountBalanceSubList(ctx context.Context, req *hydra_lcp.GetOffiaccountBalanceSubListReq) (*hydra_lcp.GetOffiaccountBalanceSubListResp, error) {
	logger := logging.GetLogger(ctx)

	companyListReq := &dao.OffiaccountBindingListInput{
		Platform:  req.Platform,
		UserID:    snowflake.MustParseString(req.UserId),
		CompanyID: snowflake.MustParseString(req.CompanyId),
		PageIdx:   req.Page.Index,
		PageSize:  req.Page.Size,
	}

	bindingList, total, err := h.offiaccountSrv.GetOffiaccountBalanceSubscriptions(ctx, companyListReq)
	if err != nil {
		return nil, status.Error(codes.Code(consts.Unknown), err.Error())
	}

	logger.Infof("[GetOffiaccountBalanceSubscriptions] total %v", total)
	var response = &hydra_lcp.GetOffiaccountBalanceSubListResp{
		PageCtx: &ptype.PageCtx{
			Index: req.Page.Index,
			Size:  req.Page.Size,
			Total: total,
		},
	}

	for _, binding := range bindingList {
		response.OffiaccountBindings = append(response.OffiaccountBindings, h.offiaccountSrv.ModelToProtoOffiaccountBinding(binding))
	}
	return response, nil
}

// DeactivateOffiaccountNotification DeactivateOffiaccountNotification
func (h *HydraLcpService) DeactivateOffiaccountNotification(ctx context.Context, req *hydra_lcp.DeactivateOffiaccountNotificationReq) (*hydra_lcp.DeactivateOffiaccountNotificationResp, error) {
	bindingData := models.OffiaccountBinding{
		Platform:              req.Platform,
		UserId:                snowflake.MustParseString(req.UserId),
		WechatOpenid:          req.WechatOpenid,
		NotificationType:      req.NotificationType,
		NotificationActivated: 0,
		DeactivateTime:        time.Now(),
	}
	if req.Platform == models.OffiaccountBindingPlatformCloud && req.CompanyId != "" {
		bindingData.CompanyId = snowflake.MustParseString(req.CompanyId)
	}
	num, err := h.offiaccountSrv.Deactivation(ctx, &bindingData)
	if err != nil {
		return nil, err
	}

	return &hydra_lcp.DeactivateOffiaccountNotificationResp{UpdateNum: num}, nil
}

// SendOffiaccountNotification SendOffiaccountNotification
func (h *HydraLcpService) SendOffiaccountNotification(ctx context.Context, req *hydra_lcp.SendOffiaccountNotificationReq) (*hydra_lcp.SendOffiaccountNotificationResp, error) {
	logger := logging.GetLogger(ctx)
	notiType := req.NotificationType
	if notiType == "" {
		notiType = "job"
	}
	notiContent := req.NotificationContent
	offiacctBinding := models.OffiaccountBinding{
		UserId:                snowflake.MustParseString(req.UserId),
		NotificationType:      notiType,
		NotificationActivated: 1,
	}
	ok, err := h.offiaccountSrv.GetOffiaccountBinding(ctx, &offiacctBinding)
	if err != nil {
		return nil, err
	}
	if !ok {
		logger.Infof("[SendOffiaccountNotification] 未关注微信公众号，不进行通知 %v", ok)
		return &hydra_lcp.SendOffiaccountNotificationResp{NotificationId: 0}, nil
	}
	officialAccount := service.GetOfficialAccount()
	var remark string
	templateColor := "#173177"

	os.Getenv("OFFIACCOUNT_MSG_TEMPLATE_JOB")
	templateID := os.Getenv("OFFIACCOUNT_MSG_TEMPLATE_JOB")
	if templateID == "" {
		templateID = config.Custom.Offiaccount.MsgTemplate.Job
	}
	remark = notiContent.Remark
	if remark == "" {
		remark = "感谢您的使用！"
	}

	msgID, err := officialAccount.GetTemplate().Send(&message.TemplateMessage{
		ToUser:     offiacctBinding.WechatOpenid,
		TemplateID: templateID,
		Data: map[string]*message.TemplateDataItem{
			"first": {
				Value: notiContent.First,
				Color: templateColor,
			},
			"keyword1": {
				Value: notiContent.Keyword1,
				Color: templateColor,
			},
			"keyword2": {
				Value: notiContent.Keyword2,
				Color: templateColor,
			},
			"keyword3": {
				Value: notiContent.Keyword3,
				Color: templateColor,
			},
			"keyword4": {
				Value: time.Now().Local().Format("2006年01月02日 15:04"),
				Color: templateColor,
			},
			"remark": {
				Value: remark,
				Color: templateColor,
			},
		},
	})
	logger.Infof("[SendJobNotice] send notice msgID %v", msgID)
	if err != nil {
		logger.Errorf("[SendJobNotice] send notice error %v", err)
	}
	return &hydra_lcp.SendOffiaccountNotificationResp{NotificationId: msgID}, nil
}

// SendOffiaccountBalanceNotifications SendOffiaccountBalanceNotifications
func (h *HydraLcpService) SendOffiaccountBalanceNotifications(ctx context.Context, req *hydra_lcp.SendOffiaccountBalanceNotificationsReq) (*empty.Empty, error) {
	logger := logging.GetLogger(ctx)

	notiType := req.NotificationType
	notiContent := req.NotificationContent
	templateColor := "#173177"
	remark := notiContent.Remark
	if remark == "" {
		remark = "感谢您的使用！"
	}
	tempData := &message.TemplateMessage{
		// TemplateID: templateID,
		Data: map[string]*message.TemplateDataItem{
			"first": {
				Value: notiContent.First,
				Color: templateColor,
			},
			"keyword1": {
				Value: notiContent.Keyword1,
				Color: templateColor,
			},
			"remark": {
				Value: remark,
				Color: templateColor,
			},
		},
	}

	switch notiType {
	case "balance":
		templateID := os.Getenv("OFFIACCOUNT_MSG_TEMPLATE_BALANCE")
		if templateID == "" {
			templateID = config.Custom.Offiaccount.MsgTemplate.Balance
		}
		tempData.TemplateID = templateID
		tempData.Data["keyword2"] = &message.TemplateDataItem{
			Value: notiContent.Keyword2,
			Color: templateColor,
		}
	case "topup":
		templateID := os.Getenv("OFFIACCOUNT_MSG_TEMPLATE_TOPUP")
		if templateID == "" {
			templateID = config.Custom.Offiaccount.MsgTemplate.Topup
		}
		tempData.TemplateID = templateID
		tempData.Data["keyword2"] = &message.TemplateDataItem{
			Value: time.Now().Local().Format("2006年01月02日 15:04"),
			Color: templateColor,
		}
		tempData.Data["keyword3"] = &message.TemplateDataItem{
			Value: notiContent.Keyword3,
			Color: templateColor,
		}
	}
	platform := req.Platform

	var subList []string
	var err error
	// 余额通知->客户/客服
	if notiType == models.OffiaccountBindingNotificationTypeBalance {
		subList, err = h.offiaccountSrv.GetOffiaccountBalanceSubsByPlatform(ctx, platform, snowflake.MustParseString(req.CompanyId))
		if err != nil {
			return nil, err
		}
		logger.Infof("[GetOffiaccountBalanceSubsByPlatform] subList %v", subList)
		logger.Infof("[GetOffiaccountBalanceSubsByPlatform] sub count %v", len(subList))
	} else {
		// 充值通知->客户
		subList, err = h.offiaccountSrv.GetOffiaccountTopupSubscriptions(ctx, snowflake.MustParseString(req.CompanyId))
		if err != nil {
			return nil, err
		}
		logger.Infof("[GetOffiaccountTopupSubscriptions] subList %v", subList)
		logger.Infof("[GetOffiaccountTopupSubscriptions] sub count %v", len(subList))
	}

	officialAccount := service.GetOfficialAccount()

	for _, openID := range subList {
		tempData.ToUser = openID
		msgID, err := officialAccount.GetTemplate().Send(tempData)
		if err != nil {
			logger.Errorf("[SendOffiaccountBalanceNotifications] send notice error %v", err)
		}
		logger.Infof("[SendOffiaccountBalanceNotifications] msgID %v", msgID)
	}
	return &empty.Empty{}, nil
}

// SendOffiaccountVisJobNotification SendOffiaccountVisJobNotification
func (h *HydraLcpService) SendOffiaccountVisJobNotification(ctx context.Context, req *hydra_lcp.SendOffiaccountVisJobNotificationReq) (*empty.Empty, error) {
	logger := logging.GetLogger(ctx)

	notiType := req.NotificationType
	notiContent := req.NotificationContent
	templateColor := "#173177"
	remark := notiContent.Remark
	if remark == "" {
		remark = "如忘记关闭，请前往平台主动关闭3D云应用"
	}
	tempData := &message.TemplateMessage{
		// TemplateID: templateID,
		Data: map[string]*message.TemplateDataItem{
			"first": {
				Value: notiContent.First,
				Color: templateColor,
			},
			"keyword1": {
				Value: notiContent.Keyword1,
				Color: templateColor,
			},
			"keyword2": {
				Value: notiContent.Keyword2,
				Color: templateColor,
			},
			"keyword3": {
				Value: notiContent.Keyword3,
				Color: templateColor,
			},
			"remark": {
				Value: remark,
				Color: templateColor,
			},
		},
	}
	templateID := os.Getenv("OFFIACCOUNT_MSG_TEMPLATE_VIS_JOB")
	if templateID == "" {
		templateID = config.Custom.Offiaccount.MsgTemplate.VisJob
	}
	tempData.TemplateID = templateID
	offiacctBinding := models.OffiaccountBinding{
		UserId:                snowflake.MustParseString(req.UserId),
		NotificationType:      notiType,
		NotificationActivated: 1,
	}
	ok, err := h.offiaccountSrv.GetOffiaccountBinding(ctx, &offiacctBinding)
	if err != nil {
		return nil, err
	}
	if !ok {
		logger.Warnf("[SendOffiaccountVisJobNotification] 未关注微信公众号，不进行通知 %v", ok)
		return &empty.Empty{}, nil
	}
	tempData.ToUser = offiacctBinding.WechatOpenid

	officialAccount := service.GetOfficialAccount()
	msgID, err := officialAccount.GetTemplate().Send(tempData)
	if err != nil {
		logger.Errorf("[SendOffiaccountVisJobNotification] error %v", err)
		return nil, err
	}
	logger.Infof("[SendOffiaccountVisJobNotification] msgID %v", msgID)
	return &empty.Empty{}, nil
}

// AddJobToNotify AddJobToNotify
func (h *HydraLcpService) AddJobToNotify(ctx context.Context, req *hydra_lcp.AddJobToNotifyReq) (*hydra_lcp.JobToNotify, error) {
	logger := logging.GetLogger(ctx)

	ok, err := h.offiaccountSrv.GetOffiaccountBinding(ctx, &models.OffiaccountBinding{
		UserId:                snowflake.MustParseString(req.UserId),
		CompanyId:             snowflake.MustParseString(req.CompanyId),
		NotificationType:      models.OffiaccountBindingNotificationTypeJob,
		NotificationActivated: 1,
	})
	if err != nil {
		return nil, err
	}
	// Stop creating JobToNitify record if no Wechat Official Account bound!
	if !ok {
		return nil, status.Errorf(consts.WechatOffiaccountNotSubscribed, "尚未绑定公众号，不能开启作业通知！")
	}
	jobNotifyModel := models.JobToNotify{
		Id:         snowflake.MustParseString(req.Id),
		UserId:     snowflake.MustParseString(req.UserId),
		JobId:      snowflake.MustParseString(req.JobId),
		CreateTime: time.Now(),
	}
	num, err := h.offiaccountSrv.AddJobToNotify(ctx, &jobNotifyModel)
	if err != nil {
		return nil, err
	}
	logger.Infof("[AddJobToNotify] affected lines %v", num)

	return h.offiaccountSrv.ModelToProtoJobNotify(&jobNotifyModel), nil
}

// GetJobToNotify GetJobToNotify
func (h *HydraLcpService) GetJobToNotify(ctx context.Context, req *hydra_lcp.GetJobToNotifyReq) (*hydra_lcp.JobToNotify, error) {
	jobNotifyModel := models.JobToNotify{
		UserId: snowflake.MustParseString(req.UserId),
		JobId:  snowflake.MustParseString(req.JobId),
	}
	ok, err := h.offiaccountSrv.GetJobToNotify(ctx, &jobNotifyModel)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Errorf(consts.JobNotificationSwitchOff, "通知未开启")
	}
	return h.offiaccountSrv.ModelToProtoJobNotify(&jobNotifyModel), nil
}

// SetJobToNotifyStatus SetJobToNotifyStatus
func (h *HydraLcpService) SetJobToNotifyStatus(ctx context.Context, req *hydra_lcp.SetJobToNotifyStatusReq) (*hydra_lcp.SetJobToNotifyStatusResp, error) {
	num, err := h.offiaccountSrv.SetJobToNotifyStatus(ctx, &models.JobToNotify{
		UserId:     snowflake.MustParseString(req.UserId),
		JobId:      snowflake.MustParseString(req.JobId),
		Status:     req.Status,
		UpdateTime: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &hydra_lcp.SetJobToNotifyStatusResp{UpdateNum: num}, nil
}

// GetOffiaccountAK Get Official Account access token from central server.
func (h *HydraLcpService) GetOffiaccountAK(ctx context.Context, req *hydra_lcp.GetOffiaccountAKReq) (*hydra_lcp.GetOffiaccountAKResp, error) {
	if req.AppId != os.Getenv("OFFIACCOUNT_APP_ID") {
		return nil, status.Errorf(consts.WechatOffiaccountAppIDInvalid, "AppId is invalid")
	}
	ak, err := service.GetOfficialAccount().GetAccessToken()
	if err != nil {
		return nil, err
	}
	return &hydra_lcp.GetOffiaccountAKResp{AccessToken: ak}, nil
}

// SaveOffiaccountMenu SaveOffiaccountMenu
func (h *HydraLcpService) SaveOffiaccountMenu(ctx context.Context, req *hydra_lcp.SaveOffiaccountMenuReq) (*empty.Empty, error) {
	logger := logging.GetLogger(ctx)

	// 尝试发布到微信公众号 失败则停止数据写入
	if req.Publish {
		var buttons []*menu.Button
		err := json.Unmarshal([]byte(req.Button), &buttons)
		if err != nil {
			return nil, err
		}
		err = service.GetOfficialAccount().GetMenu().SetMenu(buttons)
		if err != nil {
			logger.Errorf("[SaveOffiaccountMenu] create menu failed %v", err)
			return nil, status.Error(consts.OffiaccountMenuInvalid, err.Error())
		}
	}

	menuData := models.OffiaccountMenu{
		Id:    snowflake.MustParseString(req.Id),
		AppId: req.AppId,
		// 存储前先转码，兼容更多特殊字符
		Button:    base64.StdEncoding.EncodeToString([]byte(req.Button)),
		CreatorId: snowflake.MustParseString(req.CreatorId),
	}
	if req.Id != "" {
		menuData.Id = snowflake.MustParseString(req.Id)
		menuData.UpdateTime = time.Now()
		err := h.offiaccountSrv.UpdateMenu(ctx, &menuData)
		if err != nil {
			logger.Errorf("[SaveOffiaccountMenu] UpdateMenu failed %v", err)
			return nil, err
		}
	} else {
		menuID, err := rpc.GenID(ctx)
		if err != nil {
			return nil, err
		}
		menuData.Id = menuID
		menuData.CreateTime = time.Now()
		err = h.offiaccountSrv.InsertMenu(ctx, &menuData)
		if err != nil {
			logger.Errorf("[SaveOffiaccountMenu] InsertMenu failed %v", err)
			return nil, err
		}
	}
	logger.Infof("[SaveOffiaccountMenu] menuData %v", menuData)

	return &empty.Empty{}, nil
}

// GetOffiaccountMenu GetOffiaccountMenu
func (h *HydraLcpService) GetOffiaccountMenu(ctx context.Context, req *hydra_lcp.GetOffiaccountMenuReq) (*hydra_lcp.OffiaccountMenu, error) {
	menu, err := h.offiaccountSrv.GetMenu(ctx, req.AppId)
	if err != nil {
		return nil, err
	}
	// 使用前先解码
	rawButtonData, err := base64.StdEncoding.DecodeString(menu.Button)
	if err != nil {
		return nil, err
	}
	menu.Button = string(rawButtonData[:])
	return h.offiaccountSrv.ModelToProtoMenu(menu), nil
}

// AddOffiaccountReplyRule 新增自动回复规则，默认不启用
func (h *HydraLcpService) AddOffiaccountReplyRule(ctx context.Context, req *hydra_lcp.AddOffiaccountReplyRuleReq) (*empty.Empty, error) {
	// 首先检查规则有效性
	ruleErr := h.offiaccountSrv.ValidateAutoRule(ctx, req.Keywords, req.ReplyList)
	if ruleErr != nil {
		return nil, ruleErr
	}
	ruleID, err := rpc.GenID(ctx)
	if err != nil {
		return nil, err
	}
	ruleData := models.OffiaccountReplyRule{
		Id:         ruleID,
		RuleName:   req.RuleName,
		Keywords:   req.Keywords,
		ReplyList:  req.ReplyList,
		ReplyMode:  req.ReplyMode,
		IsActive:   "no",
		CreatorId:  snowflake.MustParseString(req.UserId),
		CreateTime: time.Now(),
	}
	err = h.offiaccountSrv.InsertAutoReplyRule(ctx, &ruleData)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// GetOffiaccountReplyRule GetOffiaccountReplyRule
func (h *HydraLcpService) GetOffiaccountReplyRule(ctx context.Context, req *hydra_lcp.GetOffiaccountReplyRuleReq) (*hydra_lcp.OffiaccountReplyRule, error) {
	rule, err := h.offiaccountSrv.GetAutoReplyRule(ctx, snowflake.MustParseString(req.RuleId).Int64())
	if err != nil {
		return nil, err
	}
	return h.offiaccountSrv.ModelToProtoReplyRule(rule), nil
}

// UpdateOffiaccountReplyRule UpdateOffiaccountReplyRule
func (h *HydraLcpService) UpdateOffiaccountReplyRule(ctx context.Context, req *hydra_lcp.OffiaccountReplyRule) (*empty.Empty, error) {
	// 检查规则有效性
	ruleErr := h.offiaccountSrv.ValidateAutoRule(ctx, req.Keywords, req.ReplyList)
	if ruleErr != nil {
		return nil, ruleErr
	}
	rule := &models.OffiaccountReplyRule{
		Id:         snowflake.MustParseString(req.Id),
		RuleName:   req.RuleName,
		Keywords:   req.Keywords,
		ReplyList:  req.ReplyList,
		ReplyMode:  req.ReplyMode,
		CreatorId:  snowflake.MustParseString(req.CreatorId),
		UpdateTime: time.Now(),
	}
	err := h.offiaccountSrv.UpdateAutoReplyRule(ctx, rule)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// GetOffiaccountReplyRuleList GetOffiaccountReplyRuleList 分页
func (h *HydraLcpService) GetOffiaccountReplyRuleList(ctx context.Context, req *hydra_lcp.GetOffiaccountReplyRuleListReq) (*hydra_lcp.GetOffiaccountReplyRuleListResp, error) {
	logger := logging.GetLogger(ctx)

	ruleInput := dao.AutoReplyRuleListInput{
		PageIdx:  req.Page.Index,
		PageSize: req.Page.Size,
	}
	ruleList, total, err := h.offiaccountSrv.GetKeywordReplyRuleList(ctx, &ruleInput)
	if err != nil {
		return nil, status.Error(codes.Code(consts.Unknown), err.Error())
	}

	logger.Infof("[GetOffiaccountReplyRuleList] total %v", total)
	var response = &hydra_lcp.GetOffiaccountReplyRuleListResp{
		PageCtx: &ptype.PageCtx{
			Index: req.Page.Index,
			Size:  req.Page.Size,
			Total: total,
		},
	}

	for _, autoRule := range ruleList {
		response.OffiaccountReplyRules = append(response.OffiaccountReplyRules, h.offiaccountSrv.ModelToProtoReplyRule(autoRule))
	}
	return response, nil
}

// SwitchOffiaccountReplyRule 启用/停用关键词回复规则
func (h *HydraLcpService) SwitchOffiaccountReplyRule(ctx context.Context, req *hydra_lcp.SwitchOffiaccountReplyRuleReq) (*hydra_lcp.OffiaccountReplyRule, error) {
	rule, err := h.offiaccountSrv.SwitchAutoReplyRule(ctx, snowflake.MustParseString(req.RuleId).Int64())
	if err != nil {
		return nil, err
	}
	return h.offiaccountSrv.ModelToProtoReplyRule(rule), nil
}

// SaveOffiaccountSubGeneralReply 保存(创建+更新)公众号订阅自动回复/收到消息回复
func (h *HydraLcpService) SaveOffiaccountSubGeneralReply(ctx context.Context, req *hydra_lcp.SaveOffiaccountSubGeneralReplyReq) (*empty.Empty, error) {
	if req.ReplyMode == "" {
		return nil, status.Errorf(consts.OffiaccountReplyRuleInvalidReplyMode, "必须提供reply_mode")
	}
	ruleName := "被关注回复"
	if req.ReplyMode == models.OffiaccountAutoReplyModeGeneral {
		ruleName = "收到消息回复"
	}
	ruleData := models.OffiaccountReplyRule{
		RuleName:  ruleName,
		ReplyList: req.ReplyList,
		ReplyMode: req.ReplyMode,
		IsActive:  "yes",
		CreatorId: snowflake.MustParseString(req.UserId),
	}
	ok, _, err := h.offiaccountSrv.GetSubGeneralReply(ctx, req.ReplyMode)
	if err != nil {
		return nil, err
	}
	if ok {
		ruleData.UpdateTime = time.Now()
		err = h.offiaccountSrv.UpdateSubscriptionReply(ctx, &ruleData)
		if err != nil {
			return nil, err
		}
	} else {
		ruleID, err := rpc.GenID(ctx)
		if err != nil {
			return nil, err
		}
		ruleData.Id = ruleID
		ruleData.CreateTime = time.Now()
		err = h.offiaccountSrv.InsertAutoReplyRule(ctx, &ruleData)
		if err != nil {
			return nil, err
		}
	}
	return &empty.Empty{}, nil
}

// GetOffiaccountSubGeneralReply 获取公众号订阅自动回复/收到消息回复
func (h *HydraLcpService) GetOffiaccountSubGeneralReply(ctx context.Context, req *hydra_lcp.GetOffiaccountSubGeneralReplyReq) (*hydra_lcp.OffiaccountReplyRule, error) {
	replyMode := req.ReplyMode
	if replyMode == "" {
		replyMode = models.OffiaccountAutoReplyModeSubscribe
	}
	_, subReply, err := h.offiaccountSrv.GetSubGeneralReply(ctx, replyMode)

	if err != nil {
		return nil, err
	}
	return h.offiaccountSrv.ModelToProtoReplyRule(subReply), nil
}

// DeleteOffiaccountSubGeneralReply 删除公众号订阅自动回复/收到消息回复
func (h *HydraLcpService) DeleteOffiaccountSubGeneralReply(ctx context.Context, req *hydra_lcp.DeleteOffiaccountSubGeneralReplyReq) (*empty.Empty, error) {
	if req.ReplyMode == "" {
		return nil, status.Errorf(consts.OffiaccountReplyRuleInvalidReplyMode, "必须提供reply_mode")
	}
	err := h.offiaccountSrv.DeleteSubGerenalReply(ctx, req.ReplyMode)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// InitGRPCServer InitGRPCServer
func InitGRPCServer(drv *http.Driver) {
	s, err := boot.GRPC.DefaultServer()
	if err != nil {
		panic(err)
	}

	// config ldap from env
	// if env LDAP_STARTUP exists and its value is "yes", use ldap. if not, can't use ldap
	vv, startup := os.LookupEnv("LDAP_STARTUP")
	var ldapSrv *service.LdapService
	if startup && vv == "yes" {
		dsn := os.Getenv("LDAP_DSN")
		ldapSrv = service.NewLdapService(dsn)
	}

	hander := &HydraLcpService{
		userSrv:        service.NewUserSrv(),
		phoneSvr:       service.NewPhoneSrv(),
		ldapSvr:        ldapSrv,
		offiaccountSrv: service.NewOffiaccountBindingSrv(),
		HydraConfig:    util.GetHydraConfig(),
	}

	grpc_boot.InjectAllClient(hander)
	hydra_lcp.RegisterHydraLcpServiceServer(s.Driver(), hander)

}

// SuperVerificationCodeForOms 后台超级验证码
func (h *HydraLcpService) SuperVerificationCodeForOms(ctx context.Context, req *hydra_lcp.SuperVerificationCodeForOmsReq) (*hydra_lcp.SuperVerificationCodeForOmsResp, error) {
	phone := req.Phone

	code, err := h.phoneSvr.SuperVerificationCode(ctx, phone)
	if err != nil {
		return nil, err
	}

	return &hydra_lcp.SuperVerificationCodeForOmsResp{Code: code}, nil
}
