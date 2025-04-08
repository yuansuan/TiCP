package guacamole

import (
	"fmt"
	"strings"
)

const (
	invalidProtocolMsg = "0.;"
)

func BuildProtocol(msgType string, args ...string) string {
	if len(args) == 0 {
		return invalidProtocolMsg
	}

	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("%d.%s,", len(msgType), msgType))

	for i, arg := range args {
		buf.WriteString(fmt.Sprintf("%d.%s", len(arg), arg))

		// 最后一个arg，添加 ";" 结尾
		if i == len(args)-1 {
			buf.WriteString(";")
		} else { // 非最后一个arg，添加 "," 作为分割
			buf.WriteString(",")
		}
	}

	return buf.String()
}
