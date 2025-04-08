package models

import (
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

type AccountBill struct {
	Id                  snowflake.ID `json:"id" xorm:"pk default 0 comment('主键ID') BIGINT(20)"`
	AccountId           snowflake.ID `json:"account_id" xorm:"not null default 0 comment('账户ID') BIGINT(20)"`
	Sign                int          `json:"sign" xorm:"not null default 0 comment('1加款;2扣款;3冻结;4解冻') TINYINT(1)"`
	Amount              int64        `json:"amount" xorm:"not null default 0 comment('操作金额') BIGINT(20)"`
	TradeType           int          `json:"trade_type" xorm:"not null default 0 comment('交易类型: 1支付；2充值；3退款；4提现') INT(11)"`
	TradeId             string       `json:"trade_id" xorm:"not null default '' comment('交易单ID') VARCHAR(128)"`
	IdempotentId        string       `json:"idempotent_id" xorm:"not null default '' comment('幂等ID') VARCHAR(255)"`
	AccountBalance      int64        `json:"account_balance" xorm:"not null default 0 comment('账户余额（赠余额+普通余额-冻结金额,操作后账户余额') BIGINT(20)"`
	FreezedAmount       int64        `json:"freezed_amount" xorm:"not null default 0 comment('冻结金额（操作后账户冻结余额）') BIGINT(20)"`
	DeltaNormalBalance  int64        `json:"delta_normal_balance" xorm:"not null default 0 comment('普通余额操作金额') BIGINT(20)"`
	DeltaAwardBalance   int64        `json:"delta_award_balance" xorm:"not null default 0 comment('赠送余额操作金额') BIGINT(20)"`
	DeltaVoucherBalance int64        `json:"delta_voucher_balance" xorm:"not null default 0 comment('代金券消费金额') BIGINT(20)"`
	AccountVoucherIds   string       `json:"account_voucher_ids" xorm:"not null default '' comment('账户代金券关联ids') VARCHAR(512)"`
	Comment             string       `json:"comment" xorm:"not null default '' comment('备注') VARCHAR(255)"`
	OutTradeId          snowflake.ID `json:"out_trade_id" xorm:"not null default 0 comment('ID') BIGINT(20)"`
	MerchandiseId       string       `json:"merchandise_id" xorm:"null default '' comment('商品ID') VARCHAR(128)"`
	MerchandiseName     string       `json:"merchandise_name" xorm:"not null default '' comment('商品名称') VARCHAR(255)"`
	UnitPrice           int64        `json:"unit_price" xorm:"not null default 0 comment('单价') BIGINT(20)"`
	PriceDes            string       `json:"price_des" xorm:"not null default '' comment('单价描述') VARCHAR(255)"`
	Quantity            float64      `json:"quantity" xorm:"not null default 0.0 comment('消耗数量') DOUBLE"`
	QuantityUnit        string       `json:"quantity_unit" xorm:"not null default '' comment('消耗数量单位描述') VARCHAR(255)"`
	ResourceId          string       `json:"resource_id" xorm:"not null default '' comment('资源id') VARCHAR(128)"`
	ProductName         string       `json:"product_name" xorm:"not null default '' comment('产品类型 pass求解作业: CloudCompute，云应用: CloudApp') VARCHAR(64)"`
	StartTime           time.Time    `json:"start_time" xorm:"null comment('扣费周期开始时间，按量付费使用') DATETIME"`
	EndTime             time.Time    `json:"end_time" xorm:"null comment('扣费周期结束时间，按量付费使用') DATETIME"`
	CreateTime          time.Time    `json:"create_time" xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime          time.Time    `json:"update_time" xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') DATETIME"`
}

func (b *AccountBill) TableName() string {
	return "account_bill"
}
