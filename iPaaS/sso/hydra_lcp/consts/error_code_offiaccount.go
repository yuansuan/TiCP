package consts

import (
	"google.golang.org/grpc/codes"
)

// official account err, from 240000 to 249999
const (
	InvalidParam                  codes.Code = 240006
	OffiaccountMenuNotExits       codes.Code = 240007
	OffiaccountMenuDup            codes.Code = 240008
	OffiaccountReplyRuleNotExists codes.Code = 240009
)
