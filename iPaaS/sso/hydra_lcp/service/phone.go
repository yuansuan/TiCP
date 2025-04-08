package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/monitor"
	util2 "github.com/yuansuan/ticp/common/go-kit/gin-boot/util"

	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/common"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/config"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	txcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"
)

// PhoneService PhoneService
type PhoneService struct {
	userDao  *dao.UserDao
	captcha  *CaptchaService
	smsType  string
	signName string

	// aliyun
	endpoint  string
	accessID  string
	accessKey string

	// for monyun cloud
	monyunEndpoint string
	monyunKeyYS    string
	monyunKeyZS    string

	// for tencent cloud
	templateID         string
	smsSdkAppid        string
	endpointTencent    string
	endpointSmsTencent string
	regionSmsTencent   string
	secretID           string
	secretKey          string
	signYS             string
	signZS             string
}

var randPhone *rand.Rand

const (
	// SignYS 签名 远算云
	SignYS string = "远算云"
	// SignZS 签名 智算未来
	SignZS string = "智算未来"
)

const (
	TencentSms string = "tencent"
	AliyunSms  string = "aliyun"
)

const (
	PhoneCode            string = "PhoneCode"
	SignupSuccess        string = "SignupSuccess"
	SignupSuccessForCDCS string = "SignupSuccessForCDCS"
	SignupSuccessForT3   string = "SignupSuccessForT3"
	ApplyTerminalSucess  string = "ApplyTerminalSucess"
	ApplyTerminalFailed  string = "ApplyTerminalFailed"
	ProjectConsumeLimit  string = "ProjectConsumeLimit"
)

// NewPhoneSrv NewPhoneSrv
func NewPhoneSrv() *PhoneService {
	smsConfig := config.Custom.SmsConfig
	phoneService := &PhoneService{
		userDao:   dao.NewUserDao(),
		captcha:   NewCaptcha(),
		smsType:   smsConfig.SmsType,
		endpoint:  smsConfig.Aliyun.Endpoint,
		accessID:  smsConfig.Aliyun.AccessKeyID,
		accessKey: smsConfig.Aliyun.AccessKeySecret,
		signName:  smsConfig.Aliyun.SignName,

		// monyun  梦云签名与发送账号绑定，所以只能通过切换账号来切换签名
		monyunEndpoint: "http://api01.monyun.cn:7901/sms/v2/std/single_send",
		monyunKeyYS:    "aa4529e9853908ba4a0a9316a79f8334",
		monyunKeyZS:    "1b1bce1972c3f298bc735b4780ff5616",

		// tencent cloud
		smsSdkAppid:        "1400290847",
		endpointTencent:    "https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=%v&random=%v",
		endpointSmsTencent: "sms.tencentcloudapi.com",
		regionSmsTencent:   "ap-shanghai",
		secretID:           "AKIDIL8iMtiBvOcaIlohAF482A7pHNluvUuC",
		secretKey:          "xf7fy1qQurdQxEF0f6J6LcdOuJSMB2Yo",
	}
	if phoneService.smsType == "" {
		phoneService.smsType = TencentSms
	}
	return phoneService
}

func (s *PhoneService) setVerificationCodeID(ctx context.Context, phone string, captcha string) error {
	cache := boot.MW.DefaultCache()
	logger := logging.GetLogger(ctx)
	err := cache.Put(common.RedisPrefixPhoneToCodeID, phone, captcha)

	if err != nil {
		logger.Warnf("[set code exception] save captcha to phone failed, err is ", err.Error())
		return status.Error(consts.ErrHydraLcpRedisFailed, err.Error())
	}

	err = cache.Put(common.RedisPrefixPhoneToLastSendTime, phone, time.Now().Format(time.RFC3339))
	if err != nil {
		logger.Warnf("[set code exception] save captcha generated time to phone failed, err is ", err.Error())
		return status.Error(consts.ErrHydraLcpRedisFailed, err.Error())
	}

	return nil
}

// SendVerificationCode SendVerificationCode
func (s *PhoneService) SendVerificationCode(ctx context.Context, phone string, sign string) error {
	logger := logging.GetLogger(ctx)
	logger.Infof("[send code] send code to %v", phone)
	cache := boot.MW.DefaultCache()

	var lastSendTime string
	// check whether get the captcha already within a minute
	_, ok := cache.Get(common.RedisPrefixPhoneToLastSendTime, phone, &lastSendTime)
	if ok {
		t, err := time.Parse(time.RFC3339, lastSendTime)
		if err != nil {
			logger.Warnf("[send code exception] parse last send time %v failed", lastSendTime)
			return status.Errorf(consts.ErrHydraLcpFailedToParseTime, err.Error())
		}
		if !time.Now().After(t.Add(time.Minute)) {
			logger.Warnf("[send code exception] code has been send to phone %v within a minute", phone)
			return status.Errorf(consts.ErrHydraLcpWait, "code has been send to phone %v within a minute", phone)
		}
	}

	// create a new captcha
	captcha := fmt.Sprintf("%06v", randPhone.Int31n(1000000))
	var err error

	// set redis key
	err = s.setVerificationCodeID(ctx, phone, captcha)
	if err != nil {
		return err
	}

	// send message
	err = s.SendSms(ctx, phone, []string{captcha}, common.SMSTemplatePhoneCode, sign)
	if err != nil {
		return err
	}

	logger.Infof("[send code] code sent to %v successful", phone)
	return nil
}

func init() {
	randPhone = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// VerifyCode VerifyCode
func (s *PhoneService) VerifyCode(ctx context.Context, phone string, code string) error {
	if debugCode := config.Custom.SmsConfig.DebugCode; debugCode != "" && debugCode == code {
		return nil
	}

	logger := logging.GetLogger(ctx)
	logger.Infof("[verify code] verify captcha for phone %v", phone)
	cache := boot.MW.DefaultCache()
	var captchaID string
	_, ok := cache.Get(common.RedisPrefixPhoneToCodeID, phone, &captchaID)
	if !ok {
		logger.Infof("[verify code exception] captcha for phone %v is not set", phone)
		return status.Errorf(consts.ErrHydraLcpCaptchaVerifyFailed, "captcha for phone %v is not set", phone)
	}

	if code != captchaID {
		logger.Infof("[verify code exception] verify captcha for phone %v fail", phone)
		return status.Errorf(consts.ErrHydraLcpCaptchaVerifyFailed, "verify captcha for phone %v fail", phone)
	}
	// clear phone code cache after validate success
	cache.Delete(common.RedisPrefixPhoneToCodeID, phone)
	cache.Delete(common.RedisPrefixPhoneToLastSendTime, phone)
	logger.Infof("code cache for phone %v cleared", phone)
	return nil
}

// VerifyImageCaptcha VerifyImageCaptcha
func (s *PhoneService) VerifyImageCaptcha(ctx context.Context, captchaID string, captchaContent string) error {
	logger := logging.GetLogger(ctx)
	logger.Infof("[verify code] verify image captcha for id %v", captchaID)
	cache := boot.MW.DefaultCache()
	var captchaIndex string
	_, ok := cache.Get(common.RedisPrefixImageCaptchaIDToIndex, captchaID, &captchaIndex)
	if !ok {
		logger.Warnf("[verify code exception] image captcha for id %v is not set", captchaID)
		return status.Errorf(consts.ErrHydraLcpCaptchaVerifyFailed, "image captcha for id %v is not set", captchaID)
	}

	err := s.captcha.ValidateCaptcha(captchaIndex, captchaContent)
	if err != nil {
		logger.Warnf("[verify code] verify image captcha for id %v fail", captchaID)
		return status.Errorf(consts.ErrHydraLcpCaptchaVerifyFailed, "verify image captcha for id %v fail", captchaID)
	}

	return nil
}

// SendSms 发送短信，多通道顺序尝试发送
// param参数的值，Menyun通道字符串长度最大为5， 超过将发送失败
func (s *PhoneService) SendSms(ctx context.Context, phone string, param []string, templateKey string, sign string) error {
	logger := logging.GetLogger(ctx)

	tplConfig, exist := config.Custom.SmsConfig.TemplateMap[templateKey]

	// 短信模板不存在
	if !exist {
		return status.Errorf(consts.ErrHydraLcpSendSmsTplNotFound, "sms template not found")
	}

	if len(param) != tplConfig.ParamCount {
		// 参数数目错误
		return status.Errorf(consts.ErrHydraLcpSendSmsParamCount, "the count of paramter is wrong")
	}

	// generate phone message content
	var p []interface{} = make([]interface{}, len(param))
	for k := range param {
		p[k] = param[k]
	}
	content := fmt.Sprintf(tplConfig.MonYunTpl, p...)

	if os.Getenv("DONT_SEND_SMS") == "Yes" {
		// 开启不发送短信功能；只记日志，不真实发送
		logger.Infof("[send sms] phone number: %v, content : %v, sign: %v", phone, content, sign)
		return nil
	}

	// send message
	switch s.smsType {
	case TencentSms:
		if err := s.sendSmsByTencentCloud(ctx, common.PhoneCodeChina+phone, param, tplConfig.TencentTID, sign); err != nil {
			logger.Warnf("[send code exception] tencent cloud send sms captcha failed, use monyun instead, phone is %v, err is %v", phone, err)
			// add tencent cloud send sms failed counter
			_ = boot.Monitor.AddCounter(common.HydraLcpSmsSend, 1, []*monitor.Label{
				{
					Name:  common.ServiceProvider,
					Value: "tencent_sms_service",
				},
				{
					Name:  common.SendResult,
					Value: "fail",
				},
			})
			err = s.sendSmsByMonyun(ctx, phone, content, sign)
			if err != nil {
				logger.Warnf("[send code exception] monyun cloud send sms captcha failed, phone is %v, err is %v", phone, err)
				// add monyun cloud send sms failed counter
				_ = boot.Monitor.AddCounter(common.HydraLcpSmsSend, 1, []*monitor.Label{
					{
						Name:  common.ServiceProvider,
						Value: "monyun_sms_service",
					},
					{
						Name:  common.SendResult,
						Value: "fail",
					},
				})
				return err
			}
		}
	case AliyunSms:
		err := s.sendSmsByAliyun(ctx, common.PhoneCodeChina+phone, param, tplConfig.AliyunTID, templateKey, sign)
		if err != nil {
			logger.Warnf("[send code exception] aliyun cloud send sms captcha failed, phone is %v, err is %v", phone, err)
			// add aliyun cloud send sms failed counter
			_ = boot.Monitor.AddCounter(common.HydraLcpSmsSend, 1, []*monitor.Label{
				{
					Name:  common.ServiceProvider,
					Value: "aliyun_sms_service",
				},
				{
					Name:  common.SendResult,
					Value: "fail",
				},
			})
			return err
		}
	default:
		logging.GetLogger(ctx).Warnf("[send code exception] sms type is not set")
		return status.Errorf(consts.Unknown, "sms type is not match")
	}

	_ = boot.Monitor.Add(common.HydraLcpSmsSend, 1, []*monitor.Label{
		{
			Name:  common.SendResult,
			Value: "success",
		},
	})

	return nil
}

func (s *PhoneService) sendSmsByAliyun(ctx context.Context, phone string, param []string, templateID, templateKey string, sign string) error {
	logging.GetLogger(ctx).Infof("sendSmsByAliyun params: phone: %v, param: %v, templateID: %v, templateKey: %v, sign: %v", phone, param, templateID, templateKey, sign)
	client, err := dysmsapi.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(s.accessID),
		AccessKeySecret: tea.String(s.accessKey),
		Endpoint:        tea.String(s.endpoint),
	})
	if err != nil {
		logging.GetLogger(ctx).Error("[send code exception] failed to create aliyun sms client, err is %v", err.Error())
		return err
	}
	singName := sign
	if s.signName != "" {
		singName = s.signName
	}

	templateParams, err := getAliyunTemplateParams(ctx, templateKey, param)
	if err != nil {
		logging.GetLogger(ctx).Error("[send code exception] failed to get aliyun template params, err is %v", err.Error())
		return err
	}
	sendSmsRequest := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String(singName),
		TemplateCode:  tea.String(templateID),
		TemplateParam: tea.String(templateParams),
	}
	res, err := client.SendSms(sendSmsRequest)
	if err != nil || *res.StatusCode != http.StatusOK {
		logging.GetLogger(ctx).Warnf("[send code exception] failed to send sms message to %v, err is %v", phone, err.Error())
		return err
	}
	return nil
}

func getAliyunTemplateParams(ctx context.Context, templateKey string, param []string) (string, error) {
	switch templateKey {
	case PhoneCode:
		return "{\"code\":\"" + param[0] + "\"}", nil
	case ApplyTerminalSucess:
		return "{\"name\":\"" + param[0] + "\"}", nil
	case ApplyTerminalFailed:
		return "{\"name\":\"" + param[0] + "\"}", nil
	case ProjectConsumeLimit:
		return "{\"admin\":\"" + param[0] + "\",\"project\":\"" + param[1] + "\"}", nil
	case SignupSuccess, SignupSuccessForCDCS, SignupSuccessForT3:
		return "", nil
	default:
		logging.GetLogger(ctx).Warnf("[send code exception] unknown template key %v", templateKey)
		return "", fmt.Errorf("unknown template key %v", templateKey)
	}
}

func (s *PhoneService) sendSmsByMonyun(ctx context.Context, phone string, content string, sign string) error {
	logger := logging.GetLogger(ctx)
	var req = SendReq{
		Mobile: phone,
		Apikey: s.monyunKeyYS,
	}

	if sign == SignZS {
		req.Apikey = s.monyunKeyZS
	}

	gbkContent, err := util2.Utf8ToGbk([]byte(content))
	if err != nil {
		logger.Warnf("[send code exception] change message content encoding error, message content is %v, err is %v", content, err.Error())
		return status.Error(consts.ErrHydraLcpUtf8ToGbk, err.Error())
	}

	req.Content = url.QueryEscape(string(gbkContent))

	b, err := json.Marshal(req)
	if err != nil {
		logger.Warnf("[send code exception] marshal request failed, err is %v", err.Error())
		return status.Error(consts.ErrHydraLcpJSONMarshal, err.Error())
	}

	resp, err := http.Post(s.monyunEndpoint, "Content-Type: application/json", bytes.NewBuffer(b))
	if err != nil {
		logger.Warnf("[send code exception] post request failed, err is %v", err.Error())
		return status.Error(consts.ErrHydraLcpSendHTTPPostReq, err.Error())
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Warnf("[send code exception] read post response failed, err is %v", err.Error())
		return status.Error(consts.ErrHydraLcpSendHTTPPostReq, err.Error())
	}

	var r SendResp
	err = json.Unmarshal(b, &r)
	if err != nil {
		logger.Warnf("[send code exception] unmarshal response failed, err is %v", err.Error())
		return status.Error(consts.ErrHydraLcpJSONUnmarshal, err.Error())
	}

	if r.Result != 0 {
		logger.Warnf("[send code exception] failed to send sms message to %v", phone)
		return status.Errorf(consts.ErrHydraLcpFailToSendSMS, "failed to send sms message to %v", phone)
	}

	return nil
}

func (s *PhoneService) sendSmsByTencentCloud(ctx context.Context, phone string, param []string, templateID string, sign string) error {
	logger := logging.GetLogger(ctx)

	// authentication
	credential := txcommon.NewCredential(s.secretID, s.secretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = s.endpointSmsTencent
	client, _ := sms.NewClient(credential, s.regionSmsTencent, cpf)

	// form request
	req := SendReqTencentSms{
		PhoneNumberSet:   []string{phone},
		TemplateID:       templateID,
		Sign:             sign,
		TemplateParamSet: param,
		SmsSdkAppid:      s.smsSdkAppid,
	}
	params, err := json.Marshal(req)
	if err != nil {
		logger.Warnf("[send code exception] marshal sms request from json object to string failed, json object is %v, err is %v", req, err)
		return err
	}
	request := sms.NewSendSmsRequest()
	err = request.FromJsonString(string(params))
	if err != nil {
		logger.Warnf("[send code exception] read sms request failed, request is %v, err is %v", string(params), err)
		return err
	}

	// send sms
	response, err := client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logger.Warnf("[send code exception] an api error has returned: %v, sms response is %v", err, response.ToJsonString())
		return err
	}
	if err != nil {
		logger.Warnf("[send code exception] error occurs: %v, sms response is %v", err, response.ToJsonString())
		return err
	}
	return nil
}

// SuperVerificationCode 生成超级验证码
func (s *PhoneService) SuperVerificationCode(ctx context.Context, phone string) (code string, err error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("[super code] super code for %v", phone)

	captcha := fmt.Sprintf("%06v", randPhone.Int31n(1000000))

	err = s.setVerificationCodeID(ctx, phone, captcha)
	if err != nil {
		return "", err
	}

	logger.Infof("[super code] super code for %v create successful", phone)
	return captcha, nil
}

// SendReq SendReq
type SendReq struct {
	Apikey  string `json:"apikey"`
	Mobile  string `json:"mobile"`
	Content string `json:"content"`
}

// SendResp SendResp
type SendResp struct {
	Result int `json:"result"`
}

// SendReqTencentSms SendReqTencentSms
type SendReqTencentSms struct {
	PhoneNumberSet   []string `json:"PhoneNumberSet"`
	TemplateID       string   `json:"TemplateID"`
	Sign             string   `json:"Sign"`
	TemplateParamSet []string `json:"TemplateParamSet"`
	SmsSdkAppid      string   `json:"SmsSdkAppid"`
}
