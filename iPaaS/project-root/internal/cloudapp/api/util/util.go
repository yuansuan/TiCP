package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"strings"
	"time"
	"unicode"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	baseschema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/state"
	internal_common "github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/hashid"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

const (
	PrivateIPKey         = "private_ip"
	UserIdKeyInHeader    = "x-ys-user-id"
	RequestIdKeyInHeader = "x-ys-request-id"
)

var InvalidArgumentErrResp = response.WrapErrorResp(common.InvalidArgumentErrorCode, "invalidArgument")

func GetState(c *gin.Context) (*state.State, error) {
	stateI, exist := c.Get("state")
	if !exist {
		return nil, fmt.Errorf("ginCtx['state'] not exist")
	}

	s, ok := stateI.(*state.State)
	if !ok {
		return nil, fmt.Errorf("ginCtx['state'] cannot convert to *state.State")
	}

	return s, nil
}

func ParseSnowflakeIds(snowflakeIdsStr string) ([]snowflake.ID, error) {
	snowflakeIdsTemp := strings.Split(snowflakeIdsStr, ",")

	res := make([]snowflake.ID, 0)
	for _, s := range snowflakeIdsTemp {
		snowflakeId, err := snowflake.ParseString(s)
		if err != nil {
			return nil, fmt.Errorf("parse string %s to snowflake id failed, %w", s, err)
		}

		res = append(res, snowflakeId)
	}

	return res, nil
}

func GetUserId(c *gin.Context) (snowflake.ID, error) {
	userIdStr := c.GetHeader(UserIdKeyInHeader)
	if userIdStr == "" {
		return 0, fmt.Errorf("%s is empty in HTTP Header", UserIdKeyInHeader)
	}

	userId, err := snowflake.ParseString(userIdStr)
	if err != nil {
		return 0, fmt.Errorf("parse userId failed, %w", err)
	}

	return userId, nil
}

func StringInSlice(s string, list []string) bool {
	for _, v := range list {
		if s == v {
			return true
		}
	}

	return false
}

const pathSeparator = "/"

func IsAbsPath(path string) bool {
	return strings.HasPrefix(path, pathSeparator)
}

func ContainsChinese(text string) bool {
	for _, r := range text {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}

func MustParseToSnowflakeIds(list []string) []snowflake.ID {
	res := make([]snowflake.ID, 0)

	for _, v := range list {
		res = append(res, snowflake.MustParseString(v))
	}

	return res
}

func IsAccountInArrears(account *baseschema.AccountDetail) bool {
	return account.IsOverdrawn || account.AccountBalance+account.CreditQuotaAmount < 0
}

func IsAccountFrozen(account *baseschema.AccountDetail) bool {
	return account.FrozenStatus
}

func GetNonexistentList(checkList []snowflake.ID, existList []snowflake.ID) []snowflake.ID {
	nonexistentList := make([]snowflake.ID, 0)
	for _, i := range checkList {
		if isIdInSlice(i, existList) {
			continue
		}

		nonexistentList = append(nonexistentList, i)
	}

	return nonexistentList
}

func isIdInSlice(id snowflake.ID, list []snowflake.ID) bool {
	for _, v := range list {
		if id == v {
			return true
		}
	}
	return false
}

func GenShareUsername(userId snowflake.ID, num int) (string, error) {
	var res []string
	for i := 0; i < num; i++ {
		// 防止生成一样的username
		time.Sleep(time.Nanosecond)
		name, err := hashid.Encode(userId)
		if err != nil {
			return "", fmt.Errorf("encode userId to shareUsername failed, %w", err)
		}
		res = append(res, name)
	}
	return strings.Join(res, ","), nil
}

func GetDefaultMountPaths(isLinux bool) *map[string]string {
	defaultDest := "X:"
	if isLinux {
		defaultDest = "/mnt/data"
	}
	return &map[string]string{
		"": defaultDest,
	}
}

func CheckAccountInPostSession(logger *logging.Logger, s *state.State, userID snowflake.ID) (snowflake.ID, error) {
	accountDetail, err := s.OpenAPIClient.GetAccountByUserId(userID)
	if err != nil {
		if errors.Is(err, internal_common.ErrInvalidUserID) {
			logger.Infof("invalid user ID: %d", userID)
			return 0, fmt.Errorf("invalid user ID: %d", userID)
		} else if errors.Is(err, internal_common.ErrAccountResponse) {
			logger.Infof("invalid account response, userId %s", userID)
			return 0, fmt.Errorf("invalid account response, userId %s", userID)
		}
		logger.Warnf("get account by userId %s failed, %v", userID, err)
		return 0, internal_common.ErrRequestAccountByYSID
	}

	// accountId必须有且合法
	accountId, err := snowflake.ParseString(accountDetail.AccountID)
	if err != nil {
		logger.Infof("parse accountId %s to snowflakeId failed, %v", accountDetail.AccountID, err)
		return 0, errors.WithMessagef(internal_common.ErrInvalidAccountId, "parse accountId %s to snowflakeId failed", accountDetail.AccountID)
	}

	// 账户欠费
	if IsAccountInArrears(accountDetail) {
		logger.Infof("account has not enough balance, userId %s", userID)
		return 0, errors.WithMessagef(internal_common.ErrInvalidAccountStatusNotEnoughBalance, "account has not enough balance")
	}

	// 账户被冻结
	if IsAccountFrozen(accountDetail) {
		logger.Infof("account has been frozen, userId %s", userID)
		return 0, errors.WithMessagef(internal_common.ErrInvalidAccountStatusFrozen, "account has been fronzen")
	}

	return accountId, nil
}
