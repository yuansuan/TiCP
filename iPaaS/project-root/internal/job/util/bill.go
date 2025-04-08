package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/openapi-go/utils/payby"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
	iam_client "github.com/yuansuan/ticp/common/project-root-iam/iam-client"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/account"
)

func CheckAccountAndMerchandiseInCreateJob(logger *logging.Logger, userID snowflake.ID, appID snowflake.ID) (snowflake.ID, error) {

	accountId, err := CheckAccountInCreateJob(logger, userID)
	if err != nil {
		return 0, err
	}

	return accountId, nil
}

func CheckAccountInCreateJob(logger *logging.Logger, userID snowflake.ID) (snowflake.ID, error) {
	accountDetail, err := account.Client().GetAccountByUserId(userID)
	if err != nil {
		if errors.Is(err, common.ErrInvalidUserID) {
			logger.Infof("invalid user ID: %d", userID)
			return 0, common.ErrInvalidUserID
		} else if errors.Is(err, common.ErrAccountResponse) {
			logger.Infof("invalid account response, userId %s", userID)
			return 0, common.ErrInvalidPayBy
		}
		logger.Warnf("get account by userId %s failed, %v", userID, err)
		return 0, fmt.Errorf("get account by userId %s failed", userID)
	}

	// accountId必须有且合法
	accountId, err := snowflake.ParseString(accountDetail.AccountID)
	if err != nil {
		logger.Infof("parse accountId %s to snowflakeId failed, %v", accountDetail.AccountID, err)
		return 0, common.ErrInvalidAccountId
	}

	// 账户欠费
	if IsAccountInArrears(accountDetail) {
		logger.Infof("account has not enough balance, userId %s", userID)
		return 0, common.ErrInvalidAccountStatusNotEnoughBalance
	}

	// 账户被冻结
	if IsAccountFrozen(accountDetail) {
		logger.Infof("account has been frozen, userId %s", userID)
		return 0, common.ErrInvalidAccountStatusFrozen
	}

	return accountId, nil
}

func CheckPayByAccount(logger *logging.Logger, payBy string) (snowflake.ID, error) {
	if payBy == "" {
		return snowflake.ID(0), fmt.Errorf("invalid payBy param")
	}

	reqPayBy, err := payby.ParseToken(payBy)
	if err != nil {
		return snowflake.ID(0), errors.WithMessagef(common.ErrInvalidPayBy, "invalid payby")
	}

	accessKeyID := reqPayBy.GetAccessKeyID()
	timestamp := reqPayBy.GetTimestamp()
	resourceTag := reqPayBy.GetResourceTag()
	// 获取accountId
	iamClient := iam_client.NewClient(config.GetConfig().OpenAPIEndpoint, config.GetConfig().AK, config.GetConfig().AS)
	resp, err := iamClient.GetSecret(&iam_api.GetSecretRequest{
		AccessKeyId: accessKeyID,
	})
	if err != nil {
		if strings.Contains(err.Error(), "secret not found") {
			logger.Infof("accessKeyId not found. accessKeyId: %s", reqPayBy.GetAccessKeyID())
			return snowflake.ID(0), err
		}

		logger.Warnf("get secret %s from iam server failed %s", accessKeyID, err.Error())
		return snowflake.ID(0), err
	}

	// 校验token签名是否一致
	newPayBy, err := payby.NewPayBy(accessKeyID, resp.AccessKeySecret, resourceTag, timestamp)
	if err != nil {
		logger.Infof("generate payBy failed, payBy: %+v, err: %v", newPayBy, err)
		return snowflake.ID(0), errors.Errorf("gen payby error!")
	}

	if !newPayBy.SignEqualTo(reqPayBy) {
		logger.Infof("payBy sign check failed, request payBy params: %+v, generate new sign: %+v", reqPayBy, newPayBy)
		return snowflake.ID(0), errors.WithMessagef(common.ErrInvalidArgumentPayBySignature, "accessKeyID:%s invalid Token", accessKeyID)
	}

	// token 是否过期
	if time.Now().Sub(time.UnixMilli(timestamp)) > time.Minute*5 {
		logger.Infof("payBy token expired, token: %s", payBy)
		return snowflake.ID(0), errors.WithMessagef(common.ErrPayByTokenExpire, "payBy token expired, accessKeyID: %s", accessKeyID)
	}

	payByUserId, err := snowflake.ParseString(resp.YSId)
	if err != nil {
		logger.Errorf("parse userId %s from snowflake failed %s", resp.YSId, err.Error())
		return snowflake.ID(0), err
	}

	payByAccountId, err := CheckAccountInCreateJob(logger, payByUserId)
	if err != nil {
		logger.Infof("check account failed, ysid: %s, err: %s", resp.YSId, err.Error())
		return snowflake.ID(0), err
	}

	return payByAccountId, nil
}
