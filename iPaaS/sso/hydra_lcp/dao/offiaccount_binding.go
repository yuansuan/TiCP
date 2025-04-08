package dao

import (
	"context"
	"time"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	cuModel "github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
)

// Keywords 公众号自动回复规则关键词
type Keywords struct {
	Type      string `json:"type"`
	Keyword   string `json:"content"`
	MatchMode string `json:"match_mode"`
}

// KeywordReplyInfo 关键词回复消息体
type KeywordReplyInfo struct {
	Type        string   `json:"type"`
	Content     string   `json:"content"`
	MediaID     string   `json:"media_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	NewsInfo    NewsInfo `json:"news_info"`
}

// NewsInfo NewsInfo
type NewsInfo struct {
	List []Article `json:"list"`
}

// Article Article 图文
type Article struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	Digest     string `json:"digest"`
	ShowCover  int8   `json:"show_cover"`
	CoverURL   string `json:"cover_url"`
	ContentURL string `json:"content_url"`
	SourceURL  string `json:"source_url"`
}

// OffiaccountBindingDao OffiaccountBindingDao
type OffiaccountBindingDao struct {
}

// NewOffiaccountBindingDao NewOffiaccountBindingDao
func NewOffiaccountBindingDao() *OffiaccountBindingDao {
	return &OffiaccountBindingDao{}
}

// AddOffiaccountJobSubscription AddOffiaccountJobSubscription 仅对企业账户
func (b *OffiaccountBindingDao) AddOffiaccountJobSubscription(ctx context.Context, offiaccountBinding *models.OffiaccountBinding) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	if offiaccountBinding.UserId.Int64() > 0 && offiaccountBinding.CompanyId.Int64() > 0 {
		companyUserRelation := cuModel.CompanyUserRelation{
			UserId: offiaccountBinding.UserId,
			Status: 1,
		}
		// Check if current user belongs to any company
		exist, err := session.Get(&companyUserRelation)
		if err != nil {
			return err
		}
		if !exist {
			return status.Errorf(consts.InvalidParam, "AddOffiaccountJobSubscription, companyId not exist")
		}
	}

	_, err := session.Nullable("user_id", "create_by", "update_by", "company_id", "company_ids", "wechat_unionid", "wechat_headimgurl", "wechat_language").Insert(offiaccountBinding)
	if err != nil {
		return status.Error(consts.ErrHydraLcpDBOpFail, err.Error())
	}

	return nil
}

// GetOffiaccountBinding GetOffiaccountBinding
func (b *OffiaccountBindingDao) GetOffiaccountBinding(ctx context.Context, offiaccountBinding *models.OffiaccountBinding) (bool, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	ok, err := session.Get(offiaccountBinding)
	if err != nil {
		return false, status.Errorf(consts.ErrHydraLcpDBOpFail, "failed to get offiaccountBinding, err: %v", err.Error())
	}
	return ok, nil
}

// OffiaccountBindingListInput 绑定列表参数
type OffiaccountBindingListInput struct {
	UserID    snowflake.ID `json:"user_id"`
	CompanyID snowflake.ID `json:"company_id"`
	// Could be cloud or oms
	Platform string `json:"platform"`
	PageIdx  int64  `json:"page_idx"`
	PageSize int64  `json:"page_size"`
}

// GetOffiaccountBalanceSubscriptions Get Official Account subscriptions of account balance
// userid 参数弃用
func (b *OffiaccountBindingDao) GetOffiaccountBalanceSubscriptions(ctx context.Context, req *OffiaccountBindingListInput) (bindingList []*models.OffiaccountBinding, total int64, err error) {

	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()
	if req.Platform == "" {
		req.Platform = models.OffiaccountBindingPlatformCloud
	}
	// 远算云 余额绑定 需对应企业ID
	if req.Platform == models.OffiaccountBindingPlatformCloud {
		if req.CompanyID.Int64() > 0 {
			sess.Where("company_id = ?", req.CompanyID)
		}
	}

	sess.Where("platform = ?", req.Platform)
	sess.Where("notification_activated = ?", 1)
	sess.Where("notification_type = ?", models.OffiaccountBindingNotificationTypeBalance)
	sess.Where("is_subscribed = ?", 1)

	limitSize, limitOffset := int(req.PageSize), int((req.PageIdx-1)*req.PageSize)

	total, err = sess.Limit(limitSize, limitOffset).OrderBy("subscribe_time desc").FindAndCount(&bindingList)

	return bindingList, total, err
}

// GetOffiaccountBalanceSubsByPlatform GetOffiaccountBalanceSubsByPlatform
func (b *OffiaccountBindingDao) GetOffiaccountBalanceSubsByPlatform(ctx context.Context, platform string, companyID snowflake.ID) (openIDs []string, err error) {
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()
	sess.Where("platform = ? and notification_type = ? and notification_activated = ?", platform, models.OffiaccountBindingNotificationTypeBalance, 1)
	if platform == models.OffiaccountBindingPlatformCloud {
		// 余额通知->客户
		sess.Where("company_id = ?", companyID)
	} else if platform == models.OffiaccountBindingPlatformOMS {
		// 余额通知->客服
		// Todo companyIds将来可能允许编辑 需利用compandyID进行过滤
		sess.Where("company_ids = ?", "all")
	}
	err = sess.Table("offiaccount_binding").Cols("wechat_openid").Find(&openIDs)
	if err != nil {
		return nil, err
	}
	return openIDs, nil
}

// GetOffiaccountTopupSubscriptions GetOffiaccountTopupSubscriptions
func (b *OffiaccountBindingDao) GetOffiaccountTopupSubscriptions(ctx context.Context, companyID snowflake.ID) (openIDs []string, err error) {
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()
	sess.Where("platform = ? and notification_type = ? and notification_activated = ?", models.OffiaccountBindingPlatformCloud, models.OffiaccountBindingNotificationTypeBalance, 1)
	sess.Where("company_id = ?", companyID)
	err = sess.Table("offiaccount_binding").Cols("wechat_openid").Find(&openIDs)
	if err != nil {
		return nil, err
	}
	return openIDs, nil
}

// AddOffiaccountBalanceSubscription 添加公众号余额订阅
func (b *OffiaccountBindingDao) AddOffiaccountBalanceSubscription(ctx context.Context, offiaccountBinding *models.OffiaccountBinding) (int64, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	if offiaccountBinding.Platform == models.OffiaccountBindingPlatformCloud {
		if offiaccountBinding.UserId.Int64() > 0 && offiaccountBinding.CompanyId.Int64() > 0 {
			companyUserRelation := cuModel.CompanyUserRelation{
				UserId:    offiaccountBinding.UserId,
				CompanyId: offiaccountBinding.CompanyId,
				Status:    1,
			}
			// Check if current user belongs to the company
			exist, err := session.Get(&companyUserRelation)
			if err != nil {
				return 0, err
			}
			if !exist {
				// offiaccountBinding.CompanyId = 0
				return 0, status.Errorf(consts.InvalidParam, "failed to AddOffiaccountBalanceSubscription, err: %v", err.Error())
			}
		}
	}

	num, err := session.Insert(offiaccountBinding)
	if err != nil {
		return num, status.Error(consts.ErrHydraLcpDBOpFail, err.Error())
	}
	return num, nil
}

// UpdateActivation UpdateActivation
func (b *OffiaccountBindingDao) UpdateActivation(ctx context.Context, offiaccountBinding *models.OffiaccountBinding) (int64, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	if offiaccountBinding.Id.Int64() > 0 {
		session.Where("id = ?", offiaccountBinding.Id)
	}
	session.Where("platform = ?", offiaccountBinding.Platform)
	if offiaccountBinding.UserId.Int64() > 0 {
		session.Where("user_id = ?", offiaccountBinding.UserId)
	}
	session.Where("wechat_openid = ?", offiaccountBinding.WechatOpenid)
	session.Where("notification_type = ?", offiaccountBinding.NotificationType)
	num, err := session.Cols("notification_activated", "is_subscribed", "subscribe_time", "activate_time").Update(offiaccountBinding)
	if err != nil {
		return 0, status.Error(consts.Unknown, err.Error())
	}

	return num, err
}

// Unsubscribe 解除对应公众号的所有绑定
func (b *OffiaccountBindingDao) Unsubscribe(ctx context.Context, offiaccountBinding *models.OffiaccountBinding) (int64, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	session.Where("wechat_openid = ?", offiaccountBinding.WechatOpenid)
	// session.Where("notification_type = ?", offiaccountBinding.NotificationType)
	num, err := session.Cols("notification_activated", "is_subscribed", "unsubscribe_time").Update(offiaccountBinding)
	if err != nil {
		return 0, status.Error(consts.Unknown, err.Error())
	}

	return num, err
}

// Deactivation Deactivation
func (b *OffiaccountBindingDao) Deactivation(ctx context.Context, offiaccountBinding *models.OffiaccountBinding) (int64, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	if offiaccountBinding.Platform != "" {
		session.Where("platform = ?", offiaccountBinding.Platform)
	}
	if offiaccountBinding.UserId.String() != "" {
		session.Where("user_id = ?", offiaccountBinding.UserId)
	}
	if offiaccountBinding.CompanyId.String() != "" {
		session.Where("company_id = ?", offiaccountBinding.CompanyId)
	}
	session.Where("wechat_openid = ?", offiaccountBinding.WechatOpenid)
	session.Where("notification_type = ?", offiaccountBinding.NotificationType)

	num, err := session.Cols("notification_activated", "deactivate_time").Update(offiaccountBinding)
	if err != nil {
		return 0, status.Error(consts.Unknown, err.Error())
	}

	return num, err
}

// GetKeywordReplyRules GetKeywordReplyRules
func (b *OffiaccountBindingDao) GetKeywordReplyRules(ctx context.Context) ([]*models.OffiaccountReplyRule, error) {
	var aRules []*models.OffiaccountReplyRule
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	session.Where("is_active = ?", "yes")
	err := session.Where("reply_mode = ? or reply_mode = ?", "random_one", "reply_all").Find(&aRules)
	if err != nil {
		return nil, status.Errorf(consts.ErrHydraLcpDBOpFail, "failed to get aRules, err: %v", err.Error())
	}
	return aRules, nil
}

// AddJobToNotify AddJobToNotify
func (b *OffiaccountBindingDao) AddJobToNotify(ctx context.Context, toNotifyJob *models.JobToNotify) (int64, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	num, err := session.Nullable("update_time").Insert(toNotifyJob)
	if err != nil {
		return 0, status.Error(consts.Unknown, err.Error())
	}

	return num, nil
}

// GetJobToNotify GetJobToNotify
func (b *OffiaccountBindingDao) GetJobToNotify(ctx context.Context, toNotifyJob *models.JobToNotify) (bool, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	ok, err := session.Get(toNotifyJob)
	if err != nil {
		return ok, status.Errorf(consts.Unknown, "failed to get toNotifyJob, err: %v", err.Error())
	}
	return ok, nil
}

// UpdateJobToNotifyStatus UpdateJobToNotifyStatus
func (b *OffiaccountBindingDao) UpdateJobToNotifyStatus(ctx context.Context, toNotifyJob *models.JobToNotify) (int64, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	session.Where("job_id = ?", toNotifyJob.JobId)
	num, err := session.Cols("status", "update_time").Update(toNotifyJob)
	if err != nil {
		return num, status.Errorf(consts.ErrHydraLcpDBOpFail, "failed to update toNotifyJob status, err: %v", err.Error())
	}
	return num, nil
}

// InsertMenu InsertMenu
func (b *OffiaccountBindingDao) InsertMenu(ctx context.Context, menuData *models.OffiaccountMenu) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	ok, err := session.Get(&models.OffiaccountMenu{AppId: menuData.AppId})
	if err != nil {
		return err
	}
	if ok {
		return status.Errorf(consts.OffiaccountMenuDup, "Menu exists, cannot insert")
	}
	_, err = session.Insert(menuData)
	if err != nil {
		return err
	}
	return nil
}

// UpdateMenu UpdateMenu
func (b *OffiaccountBindingDao) UpdateMenu(ctx context.Context, menuData *models.OffiaccountMenu) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	_, err := session.ID(menuData.Id).Cols("button", "update_time").Update(menuData)
	if err != nil {
		return err
	}
	return nil
}

// GetMenu GetMenu
func (b *OffiaccountBindingDao) GetMenu(ctx context.Context, appID string) (*models.OffiaccountMenu, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	menu := &models.OffiaccountMenu{}
	ok, err := session.Where("app_id = ?", appID).Get(menu)
	if err != nil {
		return nil, status.Errorf(consts.Unknown, "failed to get menu, err: %v", err.Error())
	}
	if !ok {
		return nil, status.Errorf(consts.OffiaccountMenuNotExits, "No menu exists")
	}
	return menu, nil
}

// InsertAutoReplyRule InsertAutoReplyRule
func (b *OffiaccountBindingDao) InsertAutoReplyRule(ctx context.Context, rule *models.OffiaccountReplyRule) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	_, err := session.Insert(rule)
	if err != nil {
		return err
	}
	return nil
}

// GetAutoReplyRule GetAutoReplyRule
func (b *OffiaccountBindingDao) GetAutoReplyRule(ctx context.Context, ruleID int64) (*models.OffiaccountReplyRule, error) {

	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()

	rule := &models.OffiaccountReplyRule{}
	ok, err := sess.Where("id = ?", ruleID).Get(rule)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Errorf(consts.OffiaccountReplyRuleNotExists, "This rule not exists")
	}

	return rule, nil
}

// UpdateAutoReplyRule UpdateAutoReplyRule
func (b *OffiaccountBindingDao) UpdateAutoReplyRule(ctx context.Context, rule *models.OffiaccountReplyRule) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	ok, err := session.Get(&models.OffiaccountReplyRule{Id: rule.Id})
	if err != nil {
		return err
	}
	if !ok {
		return status.Errorf(consts.OffiaccountReplyRuleNotExists, "This rule not exists")
	}

	_, err = session.ID(rule.Id).Update(rule)
	if err != nil {
		return err
	}
	return nil
}

// AutoReplyRuleListInput AutoReplyRuleListInput
type AutoReplyRuleListInput struct {
	PageIdx  int64 `json:"page_idx"`
	PageSize int64 `json:"page_size"`
}

// GetKeywordReplyRuleList GetKeywordReplyRuleList 分页
func (b *OffiaccountBindingDao) GetKeywordReplyRuleList(ctx context.Context, req *AutoReplyRuleListInput) (ruleList []*models.OffiaccountReplyRule, total int64, err error) {

	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()

	sess.Where("reply_mode = ? or reply_mode = ?", "reply_all", "random_one")

	limitSize, limitOffset := int(req.PageSize), int((req.PageIdx-1)*req.PageSize)

	total, err = sess.Limit(limitSize, limitOffset).OrderBy("create_time desc").FindAndCount(&ruleList)

	return ruleList, total, err
}

// SwitchAutoReplyRule SwitchAutoReplyRule
func (b *OffiaccountBindingDao) SwitchAutoReplyRule(ctx context.Context, ruleID int64) (*models.OffiaccountReplyRule, error) {
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()

	oldRule := &models.OffiaccountReplyRule{}
	ok, err := sess.Where("id = ?", ruleID).Get(oldRule)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Errorf(consts.OffiaccountReplyRuleNotExists, "This rule not exists")
	}
	status := oldRule.IsActive
	if status == "yes" {
		status = "no"
	} else {
		status = "yes"
	}
	oldRule.IsActive = status
	oldRule.UpdateTime = time.Now()
	_, err = sess.Where("id = ?", ruleID).Cols("is_active", "update_time").Update(oldRule)
	if err != nil {
		return nil, err
	}
	return oldRule, nil
}

// GetSubGeneralReply GetSubGeneralReply
func (b *OffiaccountBindingDao) GetSubGeneralReply(ctx context.Context, replyMode string) (bool, *models.OffiaccountReplyRule, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	subReply := &models.OffiaccountReplyRule{}
	ok, err := session.Where("reply_mode = ? and is_active = ?", replyMode, "yes").Get(subReply)
	if err != nil {
		return ok, nil, status.Errorf(consts.Unknown, "failed to get sub reply, err: %v", err.Error())
	}
	return ok, subReply, nil
}

// UpdateSubscriptionReply UpdateSubscriptionReply
func (b *OffiaccountBindingDao) UpdateSubscriptionReply(ctx context.Context, subReply *models.OffiaccountReplyRule) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	_, err := session.Where("reply_mode = ?", subReply.ReplyMode).Update(subReply)
	if err != nil {
		return status.Errorf(consts.Unknown, "failed to update subreply, err: %v", err.Error())
	}
	return nil
}

// DeleteSubGerenalReply DeleteSubGerenalReply
func (b *OffiaccountBindingDao) DeleteSubGerenalReply(ctx context.Context, replyMode string) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	subReply := &models.OffiaccountReplyRule{ReplyMode: replyMode}
	_, err := session.Delete(subReply)
	if err != nil {
		return status.Errorf(consts.Unknown, "failed to get sub reply, err: %v", err.Error())
	}
	return nil
}
