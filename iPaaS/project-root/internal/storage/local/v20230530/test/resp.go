package test

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func ParseResp(resp *http.Response, data interface{}) error {
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := json.Unmarshal(body, data)
		if err != nil {
			return err
		}
		return errors.Errorf("http: %v, body: %v", resp.Status, string(body))
	}
	return json.NewDecoder(resp.Body).Decode(data)
}
