package consts

import (
	"google.golang.org/grpc/codes"
)

// from 80001 to 90000
const (
	// 服务内部错误 数据库连接 redis连接 等等
	ErrServerInternal codes.Code = 80001
	// 参数验证错误
	ErrInvalidParams codes.Code = 80002
	// 缺少参数
	ErrMissingParams codes.Code = 80003

	// 商品已存在
	ErrMerchandiseExist codes.Code = 80051

	// 80101 请求服务超时
	ErrTimeout codes.Code = 80101
	// 80102 交易已成功
	ErrTradeAlreadySucceeded codes.Code = 80102
	// 80103 有未完成交易,请先关闭
	ErrPendingTradeExists codes.Code = 80103
	// 80104 并发操作，请重试
	ErrConcurrencyModification codes.Code = 80104
	// 80105 支付方式不支持
	ErrPayTypeNotSupport codes.Code = 80105
	// 80106 product_id未注册
	ErrProductIDNotRegister codes.Code = 80106

	// 80121 支付单不存在
	ErrTradeNotFound codes.Code = 80121

	// 充值单id不存在
	ErrCreditIDNotExist codes.Code = 80201
	// 活动id不存在
	ErrActivityIDNotExist codes.Code = 80202

	// 账户id不存在
	ErrAccountIDNotExist codes.Code = 80301
	// 账户没有信息被更新
	ErrAccountNotUpdated codes.Code = 80302
	// 资金操作不支持
	ErrFundOperateNotSupport codes.Code = 80303

	// 发送账户变更事件失败
	ErrSendAccountEventToKafka codes.Code = 80401
	// 读取账户变更事件失败
	ErrReadAccountEventFromKafka codes.Code = 80402
	// 读取IDM同步事件失败
	ErrReadIDMEventFromKafka codes.Code = 80403

	// 账单不存在
	ErrBillingNotFound codes.Code = 80501
	// 退款单已存在
	ErrRefundAlreadyExist codes.Code = 80521
	// 退款单不存在
	ErrRefundNotFound codes.Code = 80522
	// 只能取消本人提交的退款单
	ErrRefundCancelOnlyBySelf codes.Code = 80523
	// 退款单非待审状态，不能取消
	ErrRefundNotCancelOnTodo codes.Code = 80524
	// 不能审批本人提交的退款申请
	ErrRefundNotApproveBySelf codes.Code = 80525
	// 请不要重复审批
	ErrRefundAlreadyApproved codes.Code = 80526

	// 审批流程已存在
	ErrApprovalProcessExist          codes.Code = 80601
	ErrApprovalProcessNotExist       codes.Code = 80602
	ErrApprovalProcessIsBinding      codes.Code = 80603
	ErrApprovalProcessInvalid        codes.Code = 80604
	ErrApprovalProcessUnbound        codes.Code = 80605
	ErrApprovalProcessAlreadyBinding codes.Code = 80606

	// MerchandisePrice
	ErrMerchandisePriceNotExist codes.Code = 80701

	// 当前license正在使用，不能取消
	ErrLicenseIsInUse codes.Code = 80801
)
