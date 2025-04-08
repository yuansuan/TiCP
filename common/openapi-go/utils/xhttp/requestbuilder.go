package xhttp

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/yuansuan/ticp/common/openapi-go/utils"
)

type RequestBuilder struct {
	req *http.Request
	err error

	formOnce   sync.Once
	formData   *bytes.Buffer
	formWriter *multipart.Writer

	bodyOnce sync.Once
	body     url.Values

	query url.Values
}

func (rb *RequestBuilder) Method(method string) *RequestBuilder {
	rb.req.Method = method
	return rb
}

func (rb *RequestBuilder) URI(uri string) *RequestBuilder {
	rb.req.URL, rb.err = url.Parse(uri)

	return rb
}

func (rb *RequestBuilder) AddHeader(key, val string) *RequestBuilder {
	rb.req.Header.Add(key, val)
	return rb
}

func (rb *RequestBuilder) Json(v interface{}) *RequestBuilder {
	rb.req.Header.Set("Content-Type", "application/json")

	var buf bytes.Buffer
	rb.err = json.NewEncoder(&buf).Encode(v)
	// 使用json.NewEncoder()会在末尾加上一个换行符，这里去掉，否则会导致签名错误
	body := strings.TrimRight(buf.String(), "\n")
	rb.req.Body = io.NopCloser(bytes.NewBufferString(body))

	// call http2 frequency may call error below
	// cannot retry err [http2: Transport received Server's graceful shutdown GOAWAY] after Request.Body was written; define Request.GetBody to avoid this error
	// issue like https://github.com/connectrpc/connect-go/issues/541, https://github.com/stripe/stripe-go/pull/711
	rb.req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBufferString(body)), nil
	}
	return rb
}

func (rb *RequestBuilder) JsonCond(cond bool, v interface{}) *RequestBuilder {
	if cond {
		return rb.Json(v)
	}
	return rb
}

func (rb *RequestBuilder) BytesBody(v []byte) *RequestBuilder {
	rb.req.Body = io.NopCloser(bytes.NewBuffer(v))
	return rb
}

func (rb *RequestBuilder) Body(reader io.Reader) *RequestBuilder {
	rb.req.Body = io.NopCloser(reader)
	return rb
}

func (rb *RequestBuilder) AddQuery(key, value string) *RequestBuilder {
	rb.query.Add(key, value)

	return rb
}

func (rb *RequestBuilder) AddQueryCond(cond bool, key, val string) *RequestBuilder {
	if cond {
		return rb.AddQuery(key, val)
	}
	return rb
}

func (rb *RequestBuilder) AddForm(key, value string) *RequestBuilder {
	rb.initForm()
	rb.err = rb.formWriter.WriteField(key, value)

	return rb
}

func (rb *RequestBuilder) AddFormFileSlice(key string, data []byte) *RequestBuilder {
	rb.initForm()
	w, err := rb.formWriter.CreateFormFile(key, "file_slice")
	if err != nil {
		rb.err = err
		return rb
	}
	_, rb.err = w.Write(data)

	return rb
}

func (rb *RequestBuilder) AddFormCond(cond bool, key, val string) *RequestBuilder {
	if cond {
		return rb.AddForm(key, val)
	}
	return rb
}

func (rb *RequestBuilder) AddBody(key, val string) *RequestBuilder {
	rb.initBody()
	rb.body.Add(key, val)

	return rb
}

func (rb *RequestBuilder) AddBodyCond(cond bool, key, val string) *RequestBuilder {
	if cond {
		return rb.AddBody(key, val)
	}

	return rb
}

func (rb *RequestBuilder) Ref(v interface{}) *RequestBuilder {
	rv := reflect.Indirect(reflect.ValueOf(v))
	if rv.Kind() == reflect.Struct {
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			fv, ft := rv.Field(i), rt.Field(i)
			if k, ok := ft.Tag.Lookup("xquery"); ok {
				rb.AddQuery(k, utils.Stringify(fv.Interface()))
			} else if k, ok = ft.Tag.Lookup("xform"); ok {
				rb.AddForm(k, utils.Stringify(fv.Interface()))
			} else if k, ok = ft.Tag.Lookup("xbody"); ok {
				rb.AddBody(k, utils.Stringify(fv.Interface()))
			} else if k, ok = ft.Tag.Lookup("xheader"); ok {
				rb.AddHeader(k, utils.Stringify(fv.Interface()))
			}
		}
	}

	return rb
}

func (rb *RequestBuilder) Build() (*http.Request, error) {
	if rb.err != nil {
		return nil, rb.err
	}

	qs := rb.req.URL.Query()
	for k, v := range rb.query {
		qs[k] = v
	}
	rb.req.URL.RawQuery = qs.Encode()

	if rb.formData != nil {
		if err := rb.formWriter.Close(); err != nil {
			return nil, err
		}
		rb.req.Header.Set("Content-Type", rb.formWriter.FormDataContentType())
		rb.req.Body = io.NopCloser(rb.formData)
	} else if rb.body != nil {
		rb.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rb.req.Body = io.NopCloser(strings.NewReader(rb.body.Encode()))
	}

	return rb.req, nil
}

func (rb *RequestBuilder) initForm() {
	rb.formOnce.Do(func() {
		rb.formData = new(bytes.Buffer)
		rb.formWriter = multipart.NewWriter(rb.formData)
	})
}

func (rb *RequestBuilder) initBody() {
	rb.bodyOnce.Do(func() {
		rb.body = make(url.Values)
	})
}

func NewRequestBuilder() *RequestBuilder {
	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		panic(err)
	}

	return &RequestBuilder{
		req:   req,
		query: make(url.Values),
	}
}

func NewPostRequestBuilder() *RequestBuilder {
	return NewRequestBuilder().Method(http.MethodPost)
}

func NewPutRequestBuilder() *RequestBuilder {
	return NewRequestBuilder().Method(http.MethodPut)
}

func NewDeleteRequestBuilder() *RequestBuilder {
	return NewRequestBuilder().Method(http.MethodDelete)
}

func NewPatchRequestBuilder() *RequestBuilder {
	return NewRequestBuilder().Method(http.MethodPatch)
}
