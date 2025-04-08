package module

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/config"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/billlist"
	accountCreate "github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/create"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/creditadd"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/creditquotamodify"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/frozenmodify"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/idget"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/idreduce"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/paymentfreezeunfreeze"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/paymentreduce"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/ysidget"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/ysidreduce"
)

func TestAccountGetByID(t *testing.T) {
	initEev()
	t.Run("success", func(t *testing.T) {
		model := &idget.Request{
			AccountID: "3Mncb4Tfy6w",
		}

		accountResp, err := AccountGetByID(Ctx, model, "4ySNxwqt5jC")

		t.Logf("%v", accountResp)
		if err != nil {
			t.Fatal(err)
		}
		assert.Greater(t, len(accountResp.AccountID), 0)
	})

}

func TestAccountGetByYsID(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &ysidget.Request{
			UserID: "3MncaVbPff5",
		}

		accountResp, err := AccountGetByYsID(Ctx, model, "4ySNxwqt5jC")

		t.Logf("%v", accountResp)
		if err != nil {
			t.Fatal(err)
		}
		assert.Greater(t, len(accountResp.AccountID), 0)
	})

}

func TestAccountIDReduce(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &idreduce.Request{
			AccountID:             "4yzGH3ajKab",
			Comment:               "test",
			TradeID:               "4JQhLmCkcXn", //
			Amount:                10,
			AccountCashVoucherIDs: "4TdLQ4gXrj2",
		}

		accountResp, err := AccountIDReduce(Ctx, model, "4ySNxwqt5jC")

		t.Logf("%v", accountResp)
		if err != nil {
			t.Fatal(err)
		}
		assert.Greater(t, len(accountResp.AccountID), 0)
	})

}

func TestAccountYsIDReduce(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &ysidreduce.Request{
			UserID:  "3MncaVbPff5",
			Comment: "test",
			TradeID: "4JQhLmCkcXn", //
			Amount:  50,
		}

		accountResp, err := AccountYsIDReduce(Ctx, model, "4ySNxwqt5jC")

		t.Logf("%v", accountResp)
		if err != nil {
			t.Fatal(err)
		}
		assert.Greater(t, len(accountResp.AccountID), 0)
	})

}

func TestBillList(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &billlist.Request{
			AccountID: "4yzGH3ajKab",
			StartTime: "2023-06-18 11:31:22",
			EndTime:   "2023-06-20 11:31:16",
			PageIndex: 1,
			PageSize:  10,
		}

		accountResp, err := BillList(Ctx, model, "4ySNxwqt5jC")

		t.Logf("%v", accountResp)
		if err != nil {
			t.Fatal(err)
		}
		assert.Greater(t, len(accountResp.AccountBills), 0)
	})

}

func TestCreate(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		middleware.Init(config.Conf, logging.Default())
		model := &accountCreate.Request{
			AccountName: "test062n", //
			UserID:      "3MncaVbPff5",
			AccountType: 1,
		}

		accountResp, err := Create(Ctx, model, "4ySNxwqt5jC")

		t.Logf("%v", accountResp)
		if err != nil {
			t.Fatal(err)
		}
		assert.Greater(t, len(accountResp.AccountID), 0)
	})

}

func TestCreditAdd(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &creditadd.Request{
			AccountID:          "3Mncb4Tfy6w",
			DeltaAwardBalance:  100,
			DeltaNormalBalance: 200,
			TradeId:            "vdxniHsp",
			Comment:            "测试充值",
		}

		accountResp, err := CreditAdd(Ctx, model, "4ySNxwqt5jC")

		t.Logf("%v", accountResp)
		if err != nil {
			t.Fatal(err)
		}
		assert.Greater(t, len(accountResp.AccountID), 0)
	})

}

func TestCreditQuotaModify(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &creditquotamodify.Request{
			AccountID:         "4yzGH3ajKab",
			CreditQuotaAmount: 1000,
		}

		modify, err := CreditQuotaModify(Ctx, model, "4ySNxwqt5jC")

		t.Logf("%v", modify)
		if err != nil {
			t.Fatal(err)
		}
		assert.Greater(t, len(modify.AccountID), 0)
	})

}

func TestFrozenModify(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &frozenmodify.Request{
			AccountID:   "4yzGH3ajKab",
			FrozenState: false,
		}

		err := FrozenModify(Ctx, model, "4ySNxwqt5jC")

		if err != nil {
			t.Fatal(err)
		}
		assert.Nil(t, err)
	})

}

func TestPaymentFreezeUnfreeze(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &paymentfreezeunfreeze.Request{
			AccountID: "4yzGH3ajKab",
			Amount:    301,
			Comment:   "test",
			TradeID:   "4BiFxMJWvuC",
			IsFreezed: false,
		}

		_, err := PaymentFreezeUnfreeze(Ctx, model, "4ySNxwqt5jC")

		if err != nil {
			t.Fatal(err)
		}
		assert.Nil(t, err)
	})

}

func TestPaymentFreeze(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &paymentfreezeunfreeze.Request{
			AccountID: "4yzGH3ajKab",
			Amount:    301,
			Comment:   "test",
			TradeID:   "4BiFxMJWvuC",
			IsFreezed: true,
		}

		_, err := PaymentFreezeUnfreeze(Ctx, model, "4ySNxwqt5jC")

		if err != nil {
			t.Fatal(err)
		}
		assert.Nil(t, err)
	})

}

func TestPaymentReduce(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &paymentreduce.Request{
			AccountID: "4yzGH3ajKab",
			Comment:   "test",
			TradeID:   "4JQhLmCkcYd",
		}

		accountResp, err := PaymentReduce(Ctx, model, "4ySNxwqt5jC")

		t.Logf("%v", accountResp)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, accountResp.AccountID)
	})

}
