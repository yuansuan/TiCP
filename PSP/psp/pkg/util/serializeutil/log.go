package serializeutil

import (
	"encoding/json"
	"fmt"
)

const ErrTraceLogToStringFailed = "##ErrTraceLogToString##"

func GetStringForTraceLog(value any) string {
	content := GetString(value, ErrTraceLogToStringFailed)
	if ErrTraceLogToStringFailed != content {
		return content
	}
	return fmt.Sprintf("%+v", value)
}

func GetString(value any, msg string) string {
	content, err := json.Marshal(value)
	if err != nil {
		content = []byte(msg)
	}
	return string(content)
}
