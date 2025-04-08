package common

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

var (
	ErrResponseNil = errors.New("response is nil")
	ErrDataNil     = errors.New("data is nil")
)

type ErrorResp struct {
	err error
	// 此时resp.Body已经读完了，内容为空
	resp *http.Response
	body []byte
}

func (e *ErrorResp) Error() string {
	if e == nil {
		return "ErrorResp is nil"
	}

	if e.resp == nil {
		return "ErrorResp.resp is nil"
	}

	var url, otherErr string
	if e.resp.Request != nil && e.resp.Request.URL != nil {
		url = e.resp.Request.URL.String()
	}
	if e.err != nil {
		otherErr = e.err.Error()
	}

	return fmt.Sprintf("statusCode: %d, status: %s, body: %s, url: %s, other err: %s",
		e.resp.StatusCode, e.resp.Status, string(e.body), url, otherErr)
}

func ParseResp(resp *http.Response, data interface{}) error {
	errorResp := &ErrorResp{
		resp: resp,
	}
	if resp == nil {
		errorResp.err = ErrResponseNil
		return errorResp
	}

	if data == nil {
		errorResp.err = ErrDataNil
		return errorResp
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		errorResp.err = fmt.Errorf("read response body failed, %w", err)
		return errorResp
	}

	if err = jsoniter.Unmarshal(body, data); err != nil {
		errorResp.body = body
		errorResp.err = fmt.Errorf("json unmarshal failed, %w", err)
		return errorResp
	}

	if resp.StatusCode != http.StatusOK {
		errorResp.body = body
		return errorResp
	}

	return nil
}
