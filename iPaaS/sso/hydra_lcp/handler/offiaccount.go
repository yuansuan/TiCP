package handler

import (
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/basic"
	"github.com/silenceper/wechat/v2/officialaccount/message"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/rpc"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/service"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// OffiaccountCallback 微信公众号服务器配置URL 仅限微信公众号(服务号)后台使用
// 需同时支持GET POST
// @GET /api/offiaccount/callback
// @POST /api/offiaccount/callback
func (h *Handler) OffiaccountCallback(c *gin.Context) {
	logger := logging.GetLogger(c)
	officialAccount := service.GetOfficialAccount()
	// 传入request和responseWriter
	server := officialAccount.GetServer(c.Request, c.Writer)
	// 测试时启用此项
	// server.SkipValidate(true)

	//设置接收消息的处理方法
	server.SetMessageHandler(
		func(msg message.MixMessage) *message.Reply {
			logger.Infof("[SetMessageHandler] msg %v", msg)
			var reply *message.Reply
			msgType := string(msg.CommonToken.MsgType)
			if msgType == "event" {
				eventKey := msg.EventKey
				var userID snowflake.ID
				var notificationType, companyID string
				platform := models.OffiaccountBindingPlatformCloud
				wechatOpenid := string(msg.CommonToken.FromUserName)
				logger.Infof("[SetMessageHandler] wechatOpenid %v", wechatOpenid)
				offiaccountBindingQuery := models.OffiaccountBinding{WechatOpenid: wechatOpenid}

				if msg.Event == message.EventSubscribe || msg.Event == message.EventScan {
					logger.Infof("[SetMessageHandler] subscribed %v", msg.CommonToken)
					if eventKey != "" {
						splits := strings.Split(eventKey, "_")
						switch msg.Event {
						case message.EventSubscribe:
							if splits[1] == "" {
								return reply
							}
							userID = snowflake.MustParseString(splits[1])
							notificationType = splits[2]
							if len(splits) >= 4 {
								platform = splits[3]
							}
							if len(splits) >= 5 {
								companyID = splits[4]
							}
						case message.EventScan:
							userID = snowflake.MustParseString(splits[0])
							notificationType = splits[1]
							if len(splits) >= 3 {
								platform = splits[2]
							}
							if len(splits) >= 4 {
								companyID = splits[3]
							}
						}
					}
					offiaccountBindingQuery.Platform = platform
					offiaccountBindingQuery.NotificationType = notificationType

					logger.Infof(" platform %v userID %v notificationType %v companyID %v", platform, userID, notificationType, companyID)
					if userID.String() != "" && platform == models.OffiaccountBindingPlatformCloud {
						if notificationType == models.OffiaccountBindingNotificationTypeJob {
							// 当且仅当cloud绑定作业通知时，才按照UserID查询
							offiaccountBindingQuery.UserId = userID
						}
					}
					if companyID != "" {
						offiaccountBindingQuery.CompanyId = snowflake.MustParseString(companyID)
					}
					exist, err := h.offiacctBindingSrv.GetOffiaccountBinding(c, &offiaccountBindingQuery)
					if err != nil {
						logger.Errorf("[SetMessageHandler] GetOffiaccountBinding %v", err)
					}
					logger.Infof("GetOffiaccountBinding exist %v", offiaccountBindingQuery)
					if !exist {
						// 获取微信用户额外信息
						userInfo, err := officialAccount.GetUser().GetUserInfo(wechatOpenid)
						if err != nil {
							logger.Errorf("[SetMessageHandler] GetUserInfo", err)
							// ref https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
							// 40164	调用接口的IP地址不在白名单中，请在接口IP白名单中进行设置。
							// *errors.errorString=get access_token error : errcode=40164 , errormsg=invalid ip 106.54.234.106 ipv6 ::ffff:106.54.234.106, not in whitelist
							return reply
						}
						logger.Infof("userInfo %v", userInfo)
						bindingID, err := rpc.GenID(c)
						if err != nil {
							logger.Errorf("GenID error ", err)
						}
						bindingData := models.OffiaccountBinding{
							Id:                    bindingID,
							Platform:              platform,
							WechatOpenid:          wechatOpenid,
							WechatUnionid:         userInfo.UnionID,
							WechatNickname:        userInfo.Nickname,
							WechatHeadimgurl:      userInfo.Headimgurl,
							WechatLanguage:        userInfo.Language,
							UserGender:            userInfo.Sex,
							UserCity:              userInfo.City,
							NotificationType:      notificationType,
							NotificationActivated: 1,
							SubscribeScene:        userInfo.SubscribeScene,
							IsSubscribed:          1,
							SubscribeTime:         time.Now(),
							ActivateTime:          time.Now(),
						}
						if companyID != "" {
							bindingData.CompanyId = snowflake.MustParseString(companyID)
						}
						if platform == models.OffiaccountBindingPlatformOMS {
							bindingData.CompanyIds = "all"
						}
						if notificationType == models.OffiaccountBindingNotificationTypeJob {
							bindingData.UserId = userID
							logger.Infof("bindingData %v ", bindingData)
							// 插入作业通知绑定
							err = h.offiacctBindingSrv.AddOffiaccountJobSubscription(c, &bindingData)
							if err != nil {
								logger.Errorf("[SetMessageHandler] insert %v", err)
							}
						} else {
							bindingData.CreateBy = userID
							// 插入余额通知绑定 & For CAE部门渠道订阅
							num, err := h.offiacctBindingSrv.AddOffiaccountBalanceSubscription(c, &bindingData)
							if err != nil {
								logger.Errorf("[AddOffiaccountBalanceSubscription or websub] insert %v", err)
							}
							logger.Infof("[AddOffiaccountBalanceSubscription] num %v", num)
						}
					} else {
						updateData := models.OffiaccountBinding{
							Id:                    offiaccountBindingQuery.Id,
							Platform:              platform,
							WechatOpenid:          wechatOpenid,
							NotificationType:      notificationType,
							NotificationActivated: 1,
							IsSubscribed:          1,
							SubscribeTime:         time.Now(),
							ActivateTime:          time.Now(),
						}
						if platform == models.OffiaccountBindingPlatformCloud && notificationType == models.OffiaccountBindingNotificationTypeJob {
							updateData.UserId = userID
						}
						// 更新绑定
						num, err := h.offiacctBindingSrv.UpdateActivation(c, &updateData)
						if err != nil {
							logger.Errorf("[SetMessageHandler] subscribe UpdateActivation %v", err)
						}
						logger.Infof("[SetMessageHandler] subscribe UpdateActivation %v", num)
					}
					reply = h.subscribeReply(c)
				}
				if msg.Event == message.EventUnsubscribe {
					logger.Infof("[SetMessageHandler] unsubscribed")
					// 解除绑定
					num, err := h.offiacctBindingSrv.Unsubscribe(c, &models.OffiaccountBinding{
						WechatOpenid: wechatOpenid,
						// NotificationType:      "job",
						NotificationActivated: 0,
						IsSubscribed:          0,
						UnsubscribeTime:       time.Now(),
					})
					if err != nil {
						logger.Errorf("[SetMessageHandler] unsubscribe %v", err)
					}
					logger.Infof("[SetMessageHandler] unsubscribe num %v", num)
				}

			} else if msgType == string(message.MsgTypeText) {
				rules := h.autoReply(c, msgType, msg.Content)
				for _, rule := range rules {
					switch rule.Type {
					case "news":
						newsList := rule.NewsInfo.List
						articles := make([]*message.Article, len(newsList))
						logger.Infof(" newsList    %v", newsList)
						for index, news := range newsList {
							logger.Infof("[autoReply] news %v", news.ContentURL)
							articles[index] = message.NewArticle(news.Title, news.Digest, news.CoverURL, news.ContentURL)
						}

						newsMsg := message.NewNews(articles)
						logger.Infof("[autoReply] newsMsg %v", newsMsg)
						reply = &message.Reply{MsgType: message.MsgTypeNews, MsgData: newsMsg}
						logger.Infof("8 reply ===>> %v", reply)

					case "text":
						text := message.NewText(rule.Content)
						logger.Infof("[autoReply] text %v", text)
						reply = &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
					case "img":
						imgMsg := message.NewImage(rule.MediaID)
						reply = &message.Reply{MsgType: message.MsgTypeImage, MsgData: imgMsg}
					case "voice":
						voiceMsg := message.NewVoice(rule.MediaID)
						reply = &message.Reply{MsgType: message.MsgTypeVoice, MsgData: voiceMsg}
					case "video":
						if rule.MediaID != "" {
							videoMsg := message.NewVideo(rule.MediaID, rule.Title, rule.Description)
							logger.Infof("[autoReply] videoMsg %v %v %v", rule.MediaID, rule.Title, rule.Description)

							reply = &message.Reply{MsgType: message.MsgTypeVideo, MsgData: videoMsg}
						}
					}
				}

			} else {
				// 确保在任何情况下正确回复微信服务器!
				reply = &message.Reply{MsgType: message.MsgTypeNews, MsgData: "success"}
			}
			return reply

		})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		logger.Infof("[OffiaccountCallback] server error %v", err)
		return
	}
	//发送回复的消息
	server.Send()
}

type qrCodeReq struct {
	UserID           string `json:"userID" binding:"required"`
	NotificationType string `json:"notificationType" binding:"required"`
	Platform         string `json:"platform"`
	CompanyID        string `json:"companyId"`
	IsPermanent      int64  `json:"isPermanent"`
}

// CreateQRCodeWithParam 生成带参二维码
//
//	{
//		 "userId": "42sT8sLhNYq",
//		 "notificationType": "balance",
//	  "platform": "oms" 格物渠道预留gw-channel-001 ~ gw-channel-050
//	  "companyId": "42vHnMB91hh" required for cloud platform
//	  "isPermanent": 1 是否永久二维码
//	}
//
// @POST /api/offiaccount/createqrcode
func (h *Handler) CreateQRCodeWithParam(c *gin.Context) {
	logger := logging.GetLogger(c)
	var req qrCodeReq
	if err := c.BindJSON(&req); err != nil {
		return
	}
	if req.Platform == "" {
		req.Platform = models.OffiaccountBindingPlatformCloud
	}
	if req.Platform == models.OffiaccountBindingPlatformOMS {
		if req.NotificationType != models.OffiaccountBindingNotificationTypeBalance {
			http.Errf(c, consts.GetQRCodeFailed, "notificationType must be 'balance' for oms")
			return
		}
	}
	logger.Infof("CreateQRCodeWithParam params %v", req)
	officialAccount := service.GetOfficialAccount()

	isPermanent := req.IsPermanent
	var expireSeconds int64
	if isPermanent != 1 {
		expireSeconds = h.conf.Offiaccount.ExpireSeconds
		if os.Getenv("OFFIACCOUNT_EXPIRE_SECONDS") != "" {
			expireSeconds, _ = strconv.ParseInt(os.Getenv("OFFIACCOUNT_EXPIRE_SECONDS"), 10, 64)
		}
	}

	// QR code will be valid in 1 min by default
	tq := &basic.Request{
		ExpireSeconds: expireSeconds,
	}
	tq.ActionName = "QR_STR_SCENE"
	if isPermanent == 1 {
		tq.ActionName = "QR_LIMIT_STR_SCENE"
	}
	tq.ActionInfo.Scene.SceneStr = req.UserID + "_" + req.NotificationType + "_" + req.Platform
	if req.CompanyID != "" {
		tq.ActionInfo.Scene.SceneStr += "_" + req.CompanyID
	}
	qrTicket, err := officialAccount.GetBasic().GetQRTicket(tq)
	logger.Infof("qrTicket %v", qrTicket)
	if err != nil {
		logger.Errorf("GetQRTicket err %v", err)
		http.Errf(c, consts.GetQRTicketFailed, "Failed to get QR ticket %v", err)
		return
	}
	qrString := basic.ShowQRCode(qrTicket)
	http.Ok(c, QRCodeRes{
		ExpireSeconds: expireSeconds,
		QrcodeURL:     qrString,
	})
}

// QRCodeRes QRcode response
type QRCodeRes struct {
	ExpireSeconds int64  `json:"expireSeconds"`
	QrcodeURL     string `json:"qrcodeUrl"`
}

// subscribeReply 处理订阅自动回复机制
func (h *Handler) subscribeReply(c *gin.Context) *message.Reply {
	logger := logging.GetLogger(c)

	ok, subReply, err := h.offiacctBindingSrv.GetSubGeneralReply(c, models.OffiaccountAutoReplyModeSubscribe)
	if err != nil {
		logger.Errorf("[subscribeReply] GetSubGeneralReply %v", err)
	}
	var reply *message.Reply
	if !ok {
		text := message.NewText(`
Hello~欢迎加入远算云学院！
在这里你能找到关于远算创物CAD、格物CAE的使用方法、产品资料、技术技巧、进阶课程和技能培训~
点击下方菜单栏，更多知识等你发现！`)
		reply = &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	} else {
		subReplyInfo := dao.KeywordReplyInfo{}
		logger.Infof("subReply %v", subReply.ReplyList)
		json.Unmarshal([]byte(subReply.ReplyList), &subReplyInfo)
		if subReplyInfo.Type == "text" {
			text := message.NewText(subReplyInfo.Content)
			reply = &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
		}
	}
	return reply
}

func (h *Handler) autoReply(c *gin.Context, msgType string, keyword string) []dao.KeywordReplyInfo {
	logger := logging.GetLogger(c)
	aRules, err := h.offiacctBindingSrv.GetOffiaccountReplyRules(c)
	if err != nil {
		logger.Errorf("%v", err)
	}
	logger.Infof("aRules %v", aRules)
	var result []dao.KeywordReplyInfo
	for _, aRule := range aRules {
		var keywords []dao.Keywords
		json.Unmarshal([]byte(aRule.Keywords), &keywords)
		logger.Infof("keywords %v", keywords)
		matched := false
	KeyLoop:
		for _, keywordPattern := range keywords {
			if keywordPattern.Keyword == "" {
				// 忽略空keyword
				continue
			}
			switch keywordPattern.MatchMode {
			case "contain":
				isContain := strings.Contains(keyword, keywordPattern.Keyword)
				logger.Infof("1. contain keyword=%v keywordPattern=%v", keyword, keywordPattern.Keyword)

				if isContain {
					restKeyword := strings.Replace(keyword, keywordPattern.Keyword, "", -1)
					logger.Infof("2. contain restKeyword %v", restKeyword)
					if len(restKeyword) >= 0 {
						matched = true
						break KeyLoop
					}
				}
				continue

			case "equal":
				if keywordPattern.Keyword == keyword {
					matched = true
					logger.Infof("3. equal  %v", matched)
					break KeyLoop
				}
			}
		}
		logger.Infof("4. matched  %v", matched)
		if !matched {
			continue
		}
		var replyList []dao.KeywordReplyInfo
		json.Unmarshal([]byte(aRule.ReplyList), &replyList)
		logger.Infof("5. %v %v", keywords[0].Keyword, replyList[0].Type)
		if aRule.ReplyMode == "random_one" {
			rand.Seed(time.Now().Unix())
			randomIndex := rand.Intn(len(replyList))
			result = []dao.KeywordReplyInfo{replyList[randomIndex]}
		}
		if aRule.ReplyMode == "reply_all" {
			result = replyList
		}
		break
	}
	logger.Infof("6. %v", result)

	if len(result) == 0 {
		// 进入'收到消息回复'模式
		replyInfo := h.generalReply(c)
		result = append(result, replyInfo)
	}
	return result
}

// generalReply 当关键词匹配未成功时，检查是否设置了收到消息回复
func (h *Handler) generalReply(c *gin.Context) dao.KeywordReplyInfo {
	logger := logging.GetLogger(c)

	ok, genReply, err := h.offiacctBindingSrv.GetSubGeneralReply(c, models.OffiaccountAutoReplyModeGeneral)
	if err != nil {
		logger.Errorf("[generalReply] GetSubGeneralReply %v", err)
	}
	var replyInfo dao.KeywordReplyInfo
	if ok {
		json.Unmarshal([]byte(genReply.ReplyList), &replyInfo)
		logger.Infof("[generalReply] GetSubGeneralReply replyInfo %v", replyInfo)
	}
	return replyInfo
}
