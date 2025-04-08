package signer

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/yuansuan/ticp/common/openapi-go/credential"
	"github.com/yuansuan/ticp/common/openapi-go/utils"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhash"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xio"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xurl"
)

// Signer represents a signer for user request
type Signer struct {
	cred *credential.Credential

	debug func(s string)
}

// Sign generate a signature for user input (Form: [key1=val1key2=val2...] + AppSecret)
func (s *Signer) Sign(values map[string]interface{}) (*Fingerprint, error) {
	var sb strings.Builder
	utils.SortEach(values, func(key string, val interface{}) {
		sb.WriteString(fmt.Sprintf("%s=%v", key, val))
	})
	srcRaw := sb.String()
	sb.WriteString(s.cred.GetAppSecret())
	s.debug("Signature Source: " + sb.String())

	return &Fingerprint{
		AppKey:    s.cred.GetAppKey(),
		SourceStr: srcRaw,
		Signature: xhash.MD5(sb.String()),
	}, nil
}

func (s *Signer) SignHttp(req *http.Request) (*Fingerprint, error) {
	values := req.URL.Query()

	if req.Body != nil && req.Body != http.NoBody {
		raw := req.Header.Get("Content-Type")
		if raw == "" {
			raw = "application/octet-stream"
		}

		ct, _, err := mime.ParseMediaType(raw)
		if err != nil {
			return nil, errors.Wrap(err, "parse request")
		}

		var dup io.ReadCloser
		req.Body, dup, err = xio.DupReader(req.Body)
		if err != nil {
			return nil, errors.Wrap(err, "duplicate request body")
		}

		if req.Body != nil {
			switch ct {
			case "application/x-www-form-urlencoded":
				if err = req.ParseForm(); err != nil {
					return nil, errors.Wrap(err, "parse application/x-www-form-urlencoded")
				}
				values = xurl.AppendValues(values, req.PostForm)
			case "multipart/form-data":
				// 10 MB is a lot of data
				if err = req.ParseMultipartForm(int64(10 << 20)); err != nil {
					return nil, errors.Wrap(err, "parse multipart/form-data")
				}
				values = xurl.AppendValues(values, req.MultipartForm.Value)
			case "application/json":
				h := sha1.New()
				if _, err = io.Copy(h, req.Body); err != nil {
					return nil, errors.Wrap(err, "read request body")
				}
				values = xurl.AppendValues(values, url.Values{
					"_body": []string{hex.EncodeToString(h.Sum(nil))},
				})
			}

			req.Body = dup
		}
	}

	params := make(map[string]interface{})
	for key, val := range values {
		if key != "Signature" {
			params[key] = ""
			if len(val) != 0 {
				params[key] = val[0]
			}
		}
	}

	return s.Sign(params)
}

type Fingerprint struct {
	AppKey    string `json:"AppKey"`
	Signature string `json:"Signature"`
	SourceStr string
}

func (f *Fingerprint) AsQuery(req *http.Request) *http.Request {
	qs := req.URL.Query()
	qs.Set("AppKey", f.AppKey)
	qs.Set("Signature", f.Signature)

	req.URL.RawQuery = qs.Encode()
	return req
}

type Option func(s *Signer) error

func WithDebugPrinter(debug func(string)) Option {
	return func(s *Signer) error {
		s.debug = debug
		return nil
	}
}

func NewSigner(cred *credential.Credential, options ...Option) (*Signer, error) {
	s := &Signer{cred: cred, debug: silent}
	for _, opt := range options {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func silent(string) {}
