package consts

import (
	"google.golang.org/grpc/codes"
)

// from 100001 to 100100
const (
	// 获取区域列表错误
	ErrGetZoneList codes.Code = 300001 + iota
	// 获取区域信息错误
	ErrGetZoneInfo // 300002
	// 区域不存在
	ErrZoneNotExist // 300003
)
