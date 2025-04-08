package v20230530

import (
	"time"
)

type Merchandise struct {
	Id            string       `json:"Id,omitempty"`
	Name          string       `json:"Name,omitempty"`
	ChargeType    ChargeType   `json:"ChargeType,omitempty"`
	UnitPrice     float64      `json:"UnitPrice,omitempty"`    // 单位/元
	QuantityUnit  string       `json:"QuantityUnit,omitempty"` // 示例：核时
	Formula       string       `json:"Formula,omitempty"`      // 计费公式，当前仅支持 $UnitPrice * $Quantity
	YSProduct     string       `json:"YSProduct,omitempty"`
	OutResourceId string       `json:"OutResourceId,omitempty"`
	PublishState  PublishState `json:"PublishState,omitempty"`
	Description   string       `json:"Description,omitempty"`
}

type ChargeType string

const (
	PrePaid  ChargeType = "PrePaid"
	PostPaid ChargeType = "PostPaid"
)

func (ct ChargeType) IsValid() bool {
	switch ct {
	case PrePaid:
		return true
	case PostPaid:
		return true
	default:
		return false
	}
}

type ChargeParams struct {
	ChargeType *ChargeType `json:"ChargeType,omitempty"` // 计费类型 [ PrePaid | PostPaid ]，不填默认为 PostPaid
	PeriodType *string     `json:"PeriodType,omitempty"` // 计费类型为 PrePaid 时必填，表示付费单位，可填[ hour | day | month ]
	PeriodNum  *int        `json:"PeriodNum,omitempty"`  // 计费类型为 PrePaid 时必填，表示付费单位数量
}

type PublishState string

const (
	PublishStateUp   PublishState = "Up"
	PublishStateDown PublishState = "Down"
)

type SpecialPrice struct {
	MerchandiseId string  `json:"MerchandiseId,omitempty"`
	AccountId     string  `json:"AccountId,omitempty"`
	UnitPrice     float64 `json:"UnitPrice,omitempty"` // 单位/元
}

type Order struct {
	Id            string     `json:"Id,omitempty"`
	MerchandiseId string     `json:"MerchandiseId,omitempty"`
	AccountId     string     `json:"AccountId,omitempty"`
	Quantity      float64    `json:"Quantity,omitempty"`
	Comment       string     `json:"Comment,omitempty"`
	ChargeType    ChargeType `json:"ChargeType,omitempty"`
	CreateTime    time.Time  `json:"CreateTime,omitempty"`
}
