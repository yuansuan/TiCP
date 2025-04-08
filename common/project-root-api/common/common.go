package common

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// ParseResp 解析http响应
func ParseResp(resp *http.Response, data interface{}) error {
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := json.Unmarshal(body, data)
		if err != nil {
			return errors.Errorf("json unmarshal error: %v\nhttp: %v, body: %v", err, resp.Status, string(body))
		}
		return errors.Errorf("http: %v, body: %v", resp.Status, string(body))
	}
	return json.NewDecoder(resp.Body).Decode(data)
}
