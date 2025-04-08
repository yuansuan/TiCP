package module

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/accountcashvoucher/add"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/accountcashvoucher/get"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/accountcashvoucher/list"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/accountcashvoucher/statusmodify"
)

func TestAccountCashVoucherGetByID(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &get.Request{
			AccountCashVoucherID: "4TdLQ4gXrj2",
		}

		accountResp, err := AccountCashVoucherGetByID(Ctx, model, "3N5XGJTRUG9")

		t.Logf("%v", accountResp)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, accountResp)
	})
}
func TestAdd(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &add.Request{
			CashVoucherID: "4TaQtdB845C",
			AccountIDs:    "3Mnj8NRiBYL,3MsuiH9ZpSG",
		}

		err := Add(Ctx, model, "3N5XGJTRUG9")

		assert.Nil(t, err)
	})
}

func TestList(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &list.Request{
			AccountID: "3Mnczvenxr3",
			StartTime: "2023-06-18 11:31:22",
			EndTime:   "2023-06-30 11:31:16",
			PageIndex: 1,
			PageSize:  10,
		}

		accountResp, err := List(Ctx, model, "3N5XGJTRUG9")

		t.Logf("%v", accountResp)
		if err != nil {
			t.Fatal(err)
		}
		assert.Greater(t, len(accountResp.AccountCashVouchers), 0)
	})
}

func TestStatusModify(t *testing.T) {
	initEev()

	t.Run("success", func(t *testing.T) {
		model := &statusmodify.Request{
			AccountCashVoucherID: "4TdLQ4gXrj2",
			Status:               "DISABLED",
		}

		err := StatusModify(Ctx, model, "3N5XGJTRUG9")

		assert.Nil(t, err)
	})
}
