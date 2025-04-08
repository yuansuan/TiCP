package service

import (
	"context"
	"encoding/json"
	"os"
	"time"

	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"google.golang.org/grpc/status"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/config"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"
	hydra_lcp "github.com/yuansuan/ticp/iPaaS/sso/protos"
)

// OffiaccountBindingService OffiaccountBindingService
type OffiaccountBindingService struct {
	offiaccountBindingDao *dao.OffiaccountBindingDao
}

// NewOffiaccountBindingSrv NewOffiaccountBindingSrv
func NewOffiaccountBindingSrv() *OffiaccountBindingService {
	offiaccountBindingDao := dao.NewOffiaccountBindingDao()
	return &OffiaccountBindingService{
		offiaccountBindingDao: offiaccountBindingDao,
	}
}

// AddOffiaccountJobSubscription AddOffiaccountJobSubscription
func (offiAcctSrv *OffiaccountBindingService) AddOffiaccountJobSubscription(ctx context.Context, offiAcctBinding *models.OffiaccountBinding) (err error) {
	err = offiAcctSrv.offiaccountBindingDao.AddOffiaccountJobSubscription(ctx, offiAcctBinding)
	if err != nil {
		return err
	}

	return nil
}

// ModelToProtoTime convert time.Time to proto timestamp
func ModelToProtoTime(time *time.Time) *timestamp.Timestamp {
	return &timestamp.Timestamp{
		Seconds: time.Unix(),
		Nanos:   int32(time.Nanosecond()),
	}
}

// ProtoTimeToTime convert proto timestamp to time.Time
func ProtoTimeToTime(t *timestamp.Timestamp) time.Time {
	return time.Unix(t.Seconds, int64(t.Nanos))
}

// GetOffiaccountBinding GetOffiaccountBinding service
func (offiAcctSrv *OffiaccountBindingService) GetOffiaccountBinding(ctx context.Context, offiAcctBinding *models.OffiaccountBinding) (bool, error) {
	return offiAcctSrv.offiaccountBindingDao.GetOffiaccountBinding(ctx, offiAcctBinding)
}

// GetOffiaccountBalanceSubscriptions GetOffiaccountBalanceSubscriptions
func (offiAcctSrv *OffiaccountBindingService) GetOffiaccountBalanceSubscriptions(ctx context.Context, req *dao.OffiaccountBindingListInput) ([]*models.OffiaccountBinding, int64, error) {
	return offiAcctSrv.offiaccountBindingDao.GetOffiaccountBalanceSubscriptions(ctx, req)
}

// GetOffiaccountBalanceSubsByPlatform GetOffiaccountBalanceSubsByPlatform service
func (offiAcctSrv *OffiaccountBindingService) GetOffiaccountBalanceSubsByPlatform(ctx context.Context, platform string, companyID snowflake.ID) ([]string, error) {
	return offiAcctSrv.offiaccountBindingDao.GetOffiaccountBalanceSubsByPlatform(ctx, platform, companyID)
}

// GetOffiaccountTopupSubscriptions GetOffiaccountTopupSubscriptions
func (offiAcctSrv *OffiaccountBindingService) GetOffiaccountTopupSubscriptions(ctx context.Context, companyID snowflake.ID) ([]string, error) {
	return offiAcctSrv.offiaccountBindingDao.GetOffiaccountTopupSubscriptions(ctx, companyID)
}

// AddOffiaccountBalanceSubscription AddOffiaccountBalanceSubscription
func (offiAcctSrv *OffiaccountBindingService) AddOffiaccountBalanceSubscription(ctx context.Context, offiaccountBinding *models.OffiaccountBinding) (int64, error) {
	return offiAcctSrv.offiaccountBindingDao.AddOffiaccountBalanceSubscription(ctx, offiaccountBinding)
}

// UpdateActivation Update subscription and notification status
func (offiAcctSrv *OffiaccountBindingService) UpdateActivation(ctx context.Context, offiAcctBinding *models.OffiaccountBinding) (int64, error) {
	num, err := offiAcctSrv.offiaccountBindingDao.UpdateActivation(ctx, offiAcctBinding)
	if err != nil {
		return num, err
	}

	return num, nil
}

// Unsubscribe Unsubscribe Wechat Official Account
func (offiAcctSrv *OffiaccountBindingService) Unsubscribe(ctx context.Context, offiAcctBinding *models.OffiaccountBinding) (int64, error) {
	num, err := offiAcctSrv.offiaccountBindingDao.Unsubscribe(ctx, offiAcctBinding)
	if err != nil {
		return 0, err
	}

	return num, nil
}

// Deactivation Deactivation
func (offiAcctSrv *OffiaccountBindingService) Deactivation(ctx context.Context, offiAcctBinding *models.OffiaccountBinding) (int64, error) {
	num, err := offiAcctSrv.offiaccountBindingDao.Deactivation(ctx, offiAcctBinding)
	if err != nil {
		return num, err
	}

	return num, nil
}

// ModelToProtoOffiaccountBinding ModelToProtoOffiaccountBinding
func (offiAcctSrv *OffiaccountBindingService) ModelToProtoOffiaccountBinding(m *models.OffiaccountBinding) *hydra_lcp.OffiaccountBinding {
	return &hydra_lcp.OffiaccountBinding{
		Id:                    m.Id.String(),
		Platform:              m.Platform,
		UserId:                m.UserId.String(),
		CompanyId:             m.CompanyId.String(),
		CompanyIds:            m.CompanyIds,
		WechatOpenid:          m.WechatOpenid,
		WechatUnionid:         m.WechatUnionid,
		WechatNickname:        m.WechatNickname,
		WechatHeadimgurl:      m.WechatHeadimgurl,
		WechatLanguage:        m.WechatLanguage,
		UserGender:            m.UserGender,
		UserCity:              m.UserCity,
		NotificationType:      m.NotificationType,
		NotificationActivated: m.NotificationActivated,
		IsSubscribed:          m.IsSubscribed,
		SubscribeScene:        m.SubscribeScene,
		SubscribeTime:         ModelToProtoTime(&m.SubscribeTime),
		UnsubscribeTime:       ModelToProtoTime(&m.UnsubscribeTime),
		ActivateTime:          ModelToProtoTime(&m.SubscribeTime),
		DeactivateTime:        ModelToProtoTime(&m.UnsubscribeTime),
	}
}

// GetOfficialAccount 获取微信公众号实例
func GetOfficialAccount() *officialaccount.OfficialAccount {
	wc := wechat.NewWechat()
	oaConf := config.Custom.Offiaccount

	appID := os.Getenv("OFFIACCOUNT_APP_ID")
	if appID == "" {
		appID = oaConf.AppID
	}
	appSecret := os.Getenv("OFFIACCOUNT_APP_SECRET")
	if appSecret == "" {
		appSecret = oaConf.AppSecret
	}
	cfgToken := os.Getenv("OFFIACCOUNT_CFG_TOKEN")
	if cfgToken == "" {
		cfgToken = oaConf.CfgToken
	}
	// aesKey := os.Getenv("OFFIACCOUNT_CFG_TOKEN")
	// if aesKey == "" {
	// aesKey = oaConf.EncodingAESKey
	// }
	redisOptions := boot.MW.DefaultRedis().Options()
	redisOpts := &cache.RedisOpts{
		Host:     redisOptions.Addr,
		Password: redisOptions.Password,
	}
	cfg := &offConfig.Config{
		AppID:     appID,
		AppSecret: appSecret,
		Token:     cfgToken,
		// EncodingAESKey: aesKey,
		Cache: cache.NewRedis(redisOpts),
	}
	officialAccount := wc.GetOfficialAccount(cfg)
	return officialAccount
}

// GetOffiaccountReplyRules GetOffiaccountReplyRules 不分页
func (offiAcctSrv *OffiaccountBindingService) GetOffiaccountReplyRules(ctx context.Context) ([]*models.OffiaccountReplyRule, error) {
	aRules, err := offiAcctSrv.offiaccountBindingDao.GetKeywordReplyRules(ctx)
	if err != nil {
		return nil, err
	}
	return aRules, nil
}

// AddJobToNotify AddJobToNotify
func (offiAcctSrv *OffiaccountBindingService) AddJobToNotify(ctx context.Context, toNotifyJob *models.JobToNotify) (int64, error) {
	num, err := offiAcctSrv.offiaccountBindingDao.AddJobToNotify(ctx, toNotifyJob)
	if err != nil {
		return num, err
	}
	return num, nil
}

// GetJobToNotify GetJobToNotify
func (offiAcctSrv *OffiaccountBindingService) GetJobToNotify(ctx context.Context, toNotifyJob *models.JobToNotify) (bool, error) {
	ok, err := offiAcctSrv.offiaccountBindingDao.GetJobToNotify(ctx, toNotifyJob)
	if err != nil {
		return ok, err
	}
	return ok, nil
}

// SetJobToNotifyStatus SetJobToNotifyStatus
func (offiAcctSrv *OffiaccountBindingService) SetJobToNotifyStatus(ctx context.Context, toNotifyJob *models.JobToNotify) (int64, error) {
	num, err := offiAcctSrv.offiaccountBindingDao.UpdateJobToNotifyStatus(ctx, toNotifyJob)
	if err != nil {
		return num, err
	}
	return num, nil
}

// InsertMenu InsertMenu
func (offiAcctSrv *OffiaccountBindingService) InsertMenu(ctx context.Context, menu *models.OffiaccountMenu) error {
	return offiAcctSrv.offiaccountBindingDao.InsertMenu(ctx, menu)
}

// UpdateMenu UpdateMenu
func (offiAcctSrv *OffiaccountBindingService) UpdateMenu(ctx context.Context, menu *models.OffiaccountMenu) error {
	return offiAcctSrv.offiaccountBindingDao.UpdateMenu(ctx, menu)
}

// GetMenu GetMenu
func (offiAcctSrv *OffiaccountBindingService) GetMenu(ctx context.Context, appID string) (*models.OffiaccountMenu, error) {
	return offiAcctSrv.offiaccountBindingDao.GetMenu(ctx, appID)
}

// ModelToProtoJobNotify ModelToProtoJobNotify
func (offiAcctSrv *OffiaccountBindingService) ModelToProtoJobNotify(m *models.JobToNotify) *hydra_lcp.JobToNotify {
	return &hydra_lcp.JobToNotify{
		Id:         m.Id.String(),
		UserId:     m.UserId.String(),
		JobId:      m.JobId.String(),
		Status:     m.Status,
		CreateTime: ModelToProtoTime(&m.CreateTime),
		UpdateTime: ModelToProtoTime(&m.UpdateTime),
	}
}

// ModelToProtoMenu ModelToProtoMenu
func (offiAcctSrv *OffiaccountBindingService) ModelToProtoMenu(m *models.OffiaccountMenu) *hydra_lcp.OffiaccountMenu {
	return &hydra_lcp.OffiaccountMenu{
		Id:         m.Id.String(),
		AppId:      m.AppId,
		Button:     m.Button,
		CreatorId:  m.CreatorId.String(),
		CreateTime: ModelToProtoTime(&m.CreateTime),
		UpdateTime: ModelToProtoTime(&m.UpdateTime),
	}
}

// ModelToProtoReplyRule ModelToProtoReplyRule
func (offiAcctSrv *OffiaccountBindingService) ModelToProtoReplyRule(m *models.OffiaccountReplyRule) *hydra_lcp.OffiaccountReplyRule {
	return &hydra_lcp.OffiaccountReplyRule{
		Id:         m.Id.String(),
		RuleName:   m.RuleName,
		Keywords:   m.Keywords,
		ReplyList:  m.ReplyList,
		ReplyMode:  m.ReplyMode,
		IsActive:   m.IsActive,
		CreatorId:  m.CreatorId.String(),
		CreateTime: ModelToProtoTime(&m.CreateTime),
		UpdateTime: ModelToProtoTime(&m.UpdateTime),
	}
}

// InsertAutoReplyRule InsertAutoReplyRule
func (offiAcctSrv *OffiaccountBindingService) InsertAutoReplyRule(ctx context.Context, rule *models.OffiaccountReplyRule) error {
	return offiAcctSrv.offiaccountBindingDao.InsertAutoReplyRule(ctx, rule)
}

// GetAutoReplyRule GetAutoReplyRule service
func (offiAcctSrv *OffiaccountBindingService) GetAutoReplyRule(ctx context.Context, ruleID int64) (*models.OffiaccountReplyRule, error) {
	return offiAcctSrv.offiaccountBindingDao.GetAutoReplyRule(ctx, ruleID)
}

// UpdateAutoReplyRule UpdateAutoReplyRule
func (offiAcctSrv *OffiaccountBindingService) UpdateAutoReplyRule(ctx context.Context, rule *models.OffiaccountReplyRule) error {
	return offiAcctSrv.offiaccountBindingDao.UpdateAutoReplyRule(ctx, rule)
}

// GetKeywordReplyRuleList GetKeywordReplyRuleList
func (offiAcctSrv *OffiaccountBindingService) GetKeywordReplyRuleList(ctx context.Context, req *dao.AutoReplyRuleListInput) ([]*models.OffiaccountReplyRule, int64, error) {
	return offiAcctSrv.offiaccountBindingDao.GetKeywordReplyRuleList(ctx, req)
}

// SwitchAutoReplyRule SwitchAutoReplyRule service
func (offiAcctSrv *OffiaccountBindingService) SwitchAutoReplyRule(ctx context.Context, ruleID int64) (*models.OffiaccountReplyRule, error) {
	return offiAcctSrv.offiaccountBindingDao.SwitchAutoReplyRule(ctx, ruleID)
}

// GetSubGeneralReply GetSubGeneralReply
func (offiAcctSrv *OffiaccountBindingService) GetSubGeneralReply(ctx context.Context, replyMode string) (bool, *models.OffiaccountReplyRule, error) {
	return offiAcctSrv.offiaccountBindingDao.GetSubGeneralReply(ctx, replyMode)
}

// UpdateSubscriptionReply UpdateSubscriptionReply
func (offiAcctSrv *OffiaccountBindingService) UpdateSubscriptionReply(ctx context.Context, rule *models.OffiaccountReplyRule) error {
	return offiAcctSrv.offiaccountBindingDao.UpdateSubscriptionReply(ctx, rule)
}

// DeleteSubGerenalReply DeleteSubGerenalReply
func (offiAcctSrv *OffiaccountBindingService) DeleteSubGerenalReply(ctx context.Context, replyMode string) error {
	return offiAcctSrv.offiaccountBindingDao.DeleteSubGerenalReply(ctx, replyMode)
}

// ValidateAutoRule ValidateAutoRule
func (offiAcctSrv *OffiaccountBindingService) ValidateAutoRule(ctx context.Context, keywordStr string, replyStr string) error {
	keywordErr := status.Errorf(consts.OffiaccountInvalidAutoRule, "无效的关键词")
	if keywordStr == "" {
		return keywordErr
	}
	var keywords []*dao.Keywords
	err := json.Unmarshal([]byte(keywordStr), &keywords)
	if err != nil {
		return status.Errorf(consts.OffiaccountInvalidAutoRule, "无效的关键词 %v", err.Error())
	}
	for _, keywordPattern := range keywords {
		if keywordPattern.Keyword == "" {
			// 忽略空keyword
			return keywordErr
		}
	}
	var replyErr error
	if replyStr == "" {
		replyErr = status.Errorf(consts.OffiaccountInvalidAutoRule, "无效的规则内容")
		return replyErr
	}
	var replyList []*dao.KeywordReplyInfo
	err = json.Unmarshal([]byte(replyStr), &replyList)
	if err != nil {
		return status.Errorf(consts.OffiaccountInvalidAutoRule, "无效的规则内容 %v", err.Error())
	}
	var replyContent string
	for _, replyItem := range replyList {
		replyContent = replyItem.Content
		switch replyItem.Type {
		case "text":
			if replyContent == "" {
				replyErr = status.Errorf(consts.OffiaccountInvalidAutoRule, "无效的文本消息")
				break
			}
		case "news":
			if len(replyItem.NewsInfo.List) == 0 {
				replyErr = status.Errorf(consts.OffiaccountInvalidAutoRule, "无效的图文消息")
				break
			}
		case "img":
			if replyItem.MediaID == "" {
				replyErr = status.Errorf(consts.OffiaccountInvalidAutoRule, "无效的图片消息")
				break
			}
		case "video":
			if replyItem.MediaID == "" {
				replyErr = status.Errorf(consts.OffiaccountInvalidAutoRule, "无效的视频消息")
				break
			}
		case "voice":
			if replyItem.MediaID == "" {
				replyErr = status.Errorf(consts.OffiaccountInvalidAutoRule, "无效的音频消息")
				break
			}
		default:
			replyErr = status.Errorf(consts.OffiaccountInvalidAutoRule, "无效的规则内容")
			break
		}
	}
	return replyErr
}
