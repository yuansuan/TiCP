package dto

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

// ReadMessageRequest 消息设置已读列表请求数据
type ReadMessageRequest struct {
	MessageIDs []string `json:"message_ids"` // 消息ID列表
}

// MessageCountRequest 用户消息统计请求数据
type MessageCountRequest struct {
	State int `form:"state"` // 消息状态：1未读；2已读；
}

// MessageCountResponse 消息统计请求数据
type MessageCountResponse struct {
	Total int64 `json:"total"` // 消息总数
}

// MessageListRequest 消息列表请求数据
type MessageListRequest struct {
	Page      *xtype.Page      `json:"page"`       // 分页信息
	OrderSort *xtype.OrderSort `json:"order_sort"` // 排序条件
	Filter    *MessageFilter   `json:"filter"`     // 过滤条件
}

// MessageListResponse 消息列表返回数据
type MessageListResponse struct {
	Page     *xtype.PageResp `json:"page"`     // 分页信息
	Messages []*Message      `json:"messages"` // 消息列表
}

// Message 消息信息
type Message struct {
	ID         string    `json:"id"`          // 消息ID
	UserID     string    `json:"user_id"`     // 用户ID
	Type       string    `json:"type"`        // 消息类型
	Content    string    `json:"content"`     // 消息内容
	State      int       `json:"state"`       // 消息状态
	CreateTime time.Time `json:"create_time"` // 创建时间
	UpdateTime time.Time `json:"update_time"` // 修改时间
}

// MessagePage 消息分页
type MessagePage struct {
	UserID    int64            `json:"user_id"`
	Page      *xtype.Page      `json:"page"`
	OrderSort *xtype.OrderSort `json:"order_sort"`
	Filter    *MessageFilter   `json:"filter"`
}

// MessageFilter 消息过滤器
type MessageFilter struct {
	State   int    `json:"state"`   // 消息状态：1未读；2已读；
	Content string `json:"content"` // 消息内容
}

type WebsocketMessage struct {
	UserId  string `json:"user_id"` // 发送目标用户ID
	Type    string `json:"type"`    // 消息类型
	Content string `json:"content"` // 消息内容
}
