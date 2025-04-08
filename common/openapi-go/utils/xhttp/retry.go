package xhttp

import (
	"net/http"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

const DefaultRetryInterval = 5 * time.Second

func DefaultRetryCondition(resp *http.Response, err error) bool {
	if err != nil {
		if resp != nil && resp.Request != nil {
			logging.Default().Warnf("Error: %v, URL: %s, Host: %s", err, resp.Request.URL.String(), resp.Request.Host)
		} else {
			logging.Default().Warnf("Error: %v, URL: unknown, Host: unknown", err)
		}
		return true
	}
	if resp != nil && resp.StatusCode >= http.StatusInternalServerError {
		logging.Default().Warnf("StatusCode: %d, URL: %s, Host: %s", resp.StatusCode, resp.Request.URL.String(), resp.Request.Host)
		return true
	}
	return false
}
