package module

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/config"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/env"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/cashvoucher/add"
	availabilityModify "github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/cashvoucher/availabilitymodify"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/cashvoucher/get"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/cashvoucher/list"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

var (
	Driver = "mysql"
	DsName = "lambdacal:1234yskj@tcp(0.0.0.0:3306)/platform_dev?charset=utf8&parseTime=true&loc=Local"
	DB     *xorm.Engine // 创建xorm引擎
	DBErr  error
	Sess   *xorm.Session
	Ctx    context.Context
	Log    *zap.SugaredLogger
)

func initEev() error {
	//初始化db
	DB, DBErr = xorm.NewEngine(Driver, DsName)
	if DBErr != nil {
		return DBErr
	}
	Sess := DB.NewSession().MustLogSQL(true)
	Ctx = with.KeepSession(context.Background(), Sess)
	modeStr := env.ModeName(env.Env.Mode)
	//初始化配置文件
	configFilePath := "../" + config.ConfigDir + string(os.PathSeparator) + modeStr + ".yml"
	config.InitConfig(configFilePath)
	middleware.Init(config.Conf, logging.Default())
	return nil
}
func TestAddCashVoucher(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &add.Request{
			CashVoucherName: "cashVoucherNameTestUnit0703",
			Amount:          10000,
			ExpiredType:     1,
			AbsExpiredTime:  "2025-01-01 23:59:59",
			RelExpiredTime:  0,
			Comment:         "commentTest",
		}

		voucher, err := AddCashVoucher(Ctx, model, "4N8E9pem4b1")
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, voucher)
	})
}
func TestAvailabilityModify(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &availabilityModify.Request{
			CashVoucherID:      "4TaE3Jvxi1U",
			AvailabilityStatus: "AVAILABLE",
		}

		err := AvailabilityModify(Ctx, model, "4N8E9pem4b1")

		assert.Nil(t, err)
	})
}

func TestGetCashVoucherByID(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &get.Request{
			CashVoucherID: "4TaE3Jvxi1U",
		}

		accountResp, err := GetCashVoucherByID(Ctx, model, "4N8E9pem4b1")

		t.Logf("%v", accountResp)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, accountResp)
	})
}

func TestListCashVoucher(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &list.Request{
			AvailabilityStatus: "AVAILABLE",
			StartTime:          "2023-06-18 11:31:22",
			EndTime:            "2023-06-30 11:31:16",
			PageIndex:          1,
			PageSize:           10,
		}

		accountResp, err := ListCashVoucher(Ctx, model, "4N8E9pem4b1")

		t.Logf("%v", accountResp)
		if err != nil {
			t.Fatal(err)
		}
		assert.Greater(t, len(accountResp.CashVouchers), 0)
	})
}
