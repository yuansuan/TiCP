package consts

import (
	"google.golang.org/grpc/codes"
)

// official account err, from 240000 to 249999
const (
	GetQRTicketFailed                    codes.Code = 240001
	WechatOffiaccountNotSubscribed       codes.Code = 240002
	JobNotificationSwitchOff             codes.Code = 240003
	WechatOffiaccountAppIDInvalid        codes.Code = 240004
	GetQRCodeFailed                      codes.Code = 240005
	InvalidParam                         codes.Code = 240006
	OffiaccountMenuNotExits              codes.Code = 240007
	OffiaccountMenuDup                   codes.Code = 240008
	OffiaccountReplyRuleNotExists        codes.Code = 240009
	OffiaccountReplyRuleInvalidReplyMode codes.Code = 240010
	OffiaccountMenuInvalid               codes.Code = 240011
	OffiaccountInvalidAutoRule           codes.Code = 240012
)
