package models

import (
	"encoding/json"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

type Account struct {
	Id              snowflake.ID `json:"id" xorm:"not null pk default 0 comment('账号ID') BIGINT(20)"`
	CustomerId      snowflake.ID `json:"customer_id" xorm:"not null default 0 comment('客户ID') BIGINT(20)"`
	RealCustomerId  snowflake.ID `json:"real_customer_id" xorm:"not null default 0 comment('实名认证客户ID') BIGINT(20)"`
	Name            string       `json:"name" xorm:"not null default '' comment('账户名称，个人为手机号') VARCHAR(255)"`
	Currency        string       `json:"currency" xorm:"not null default 'CNY' comment('币种') VARCHAR(8)"`
	AccountBalance  int64        `json:"account_balance" xorm:"not null default 0 comment('账户余额（赠余额+普通余额-冻结金额）') BIGINT(20)"`
	FreezedAmount   int64        `json:"freezed_amount" xorm:"not null default 0 comment('冻结金额') BIGINT(20)"`
	NormalBalance   int64        `json:"normal_balance" xorm:"not null default 0 comment('普通余额（可提现金额）') BIGINT(20)"`
	AwardBalance    int64        `json:"award_balance" xorm:"not null default 0 comment('赠送余额') BIGINT(20)"`
	WithdrawEnabled int          `json:"withdraw_enabled" xorm:"not null default 1 comment('是否允许提现') TINYINT(1)"`
	CreditQuota     int64        `json:"credit_quota" xorm:"not null default 0 comment('授信额度') BIGINT(20)"`
	Status          int          `json:"status" xorm:"not null default 1 comment('状态:0,已删除;1,正常') TINYINT(2)"`
	AccountType     int32        `json:"account_type" xorm:"not null default 1 comment('账号类型:1,企业; 2,个人') TINYINT(2)"`
	IsFreeze        bool         `json:"is_freeze" xorm:"not null default 0 comment('账户是否被人工冻结') BOOLEAN"`
	CreateTime      time.Time    `json:"create_time" xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime      time.Time    `json:"update_time" xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') DATETIME"`
}

func (a *Account) TableName() string {
	return "account"
}

func (a *Account) calAccountBalance() {
	a.AccountBalance = a.NormalBalance + a.AwardBalance - a.FreezedAmount
}

func (a *Account) Add(normal, award int64) {
	defer a.calAccountBalance()

	a.NormalBalance += normal
	a.AwardBalance += award
}

// Reduce Reduce
// return (normalReduced, awardReduced)
func (a *Account) Reduce(amount int64) (normalReduced int64, awardReduced int64) {
	defer a.calAccountBalance()

	reduceToZero := func(balance, amount int64) (int64, int64) {
		if balance < 0 {
			return balance, 0
		}
		if balance > amount {
			return balance - amount, amount
		}
		return 0, balance
	}

	// 优先使用 NormalBalance
	// AwardBalance 不能透支
	// NormalBalance 可透支（业务控制是否透支，资金不够不提交业务单，或撤销已提交业务单 ）
	tmp := amount
	a.NormalBalance, normalReduced = reduceToZero(a.NormalBalance, tmp)
	tmp = tmp - normalReduced

	a.AwardBalance, awardReduced = reduceToZero(a.AwardBalance, tmp)
	tmp = tmp - awardReduced

	a.NormalBalance = a.NormalBalance - tmp
	normalReduced += tmp

	return normalReduced, awardReduced
}

// Freeze Freeze + delta
func (a *Account) Freeze(delta int64) {
	defer a.calAccountBalance()

	a.FreezedAmount = a.FreezedAmount + delta
}

// Unfreeze Unfreeze amount
func (a *Account) Unfreeze(amount int64) {
	a.Freeze(-amount)
}

// IsOverCreditQuota IsOverCreditQuota
func (a *Account) IsOverCreditQuota() bool {
	return a.AwardBalance < -a.CreditQuota
}

func (a *Account) String() string {
	bs, _ := json.Marshal(a)
	return string(bs)
}
