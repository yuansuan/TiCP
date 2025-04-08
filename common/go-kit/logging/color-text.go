/*
 * // Copyright (C) 2018 LambdaCal Inc.
 *
 */

package logging

import "net/http"

var (
	Green        = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	White        = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	Yellow       = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	Red          = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	Blue         = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	Magenta      = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	Cyan         = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	ResetColor   = string([]byte{27, 91, 48, 109})
	DisableColor = false
)

// ColorForStatus ColorForStatus
func ColorForStatus(code int) string {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return Green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return White
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return Yellow
	default:
		return Red
	}
}

// ColorForMethod ColorForMethod
func ColorForMethod(method string) string {
	switch method {
	case "GET":
		return Blue
	case "POST":
		return Cyan
	case "PUT":
		return Yellow
	case "DELETE":
		return Red
	case "PATCH":
		return Green
	case "HEAD":
		return Magenta
	case "OPTIONS":
		return White
	default:
		return ResetColor
	}
}
