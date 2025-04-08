package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"github.com/dchest/captcha"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/common"
)

// CaptchaService CaptchaService
type CaptchaService struct {
	store *customStore
}

// NewCaptcha NewCaptcha
func NewCaptcha() *CaptchaService {
	store := &customStore{}
	captcha.SetCustomStore(store)
	return &CaptchaService{
		store: store,
	}
}

// image is base64 encoded
func (srv *CaptchaService) CreateImageCaptcha(ctx context.Context, width int, height int, imageCaptchaID string) (string, image string, err error) {
	index := captcha.New()

	buf := new(bytes.Buffer)

	logger := logging.GetLogger(ctx)
	err = captcha.WriteImage(buf, index, width, height)
	if err != nil {
		logger.Warnf("[create image captcha exception] fail to create image captcha: %v", err)
		return "", "", status.Errorf(consts.ErrHydraLcpCaptchaFailed, "fail to create image captcha: %v", err)
	}

	image = base64.StdEncoding.EncodeToString(buf.Bytes())

	// save image captcha request, id -> index -> content
	cache := boot.MW.DefaultCache()
	err = cache.Put(common.RedisPrefixImageCaptchaIDToIndex, imageCaptchaID, index)
	if err != nil {
		logger.Warnf("[create image captcha exception] fail to save image captcha: %v", err)
		return "", "", status.Errorf(consts.ErrHydraLcpRedisFailed, "fail to save image captcha: %v", err)
	}

	return index, image, nil
}

// CreateDigitCaptcha CreateDigitCaptcha
func (srv *CaptchaService) CreateDigitCaptcha() (id string, digits string) {
	// new captcha id and captcha context, and save captcha id
	id = captcha.New()
	// get captcha context by captcha id
	b := srv.store.Get(id, false)
	// convert []byte to string
	for _, i := range b {
		digits = digits + fmt.Sprint(i)
	}
	return id, digits
}

// ValidateCaptcha ValidateCaptcha
func (srv *CaptchaService) ValidateCaptcha(id string, digits string) error {
	logging.Default().Info(">>>>", id, ">>>>>", digits)
	match := captcha.VerifyString(id, digits)
	if match {
		return nil
	}

	return status.Error(consts.ErrHydraLcpCaptchaVerifyFailed, "captcha digits doesn't match")
}

// Reload Reload
func (srv *CaptchaService) Reload(id string) error {
	exist := captcha.Reload(id)
	if exist {
		return nil
	}

	return status.Error(consts.ErrHydraLcpCaptchaFailed, "captcha id doesn't exist")
}

type customStore struct {
}

// Set sets the digits for the captcha id.
func (s *customStore) Set(id string, digits []byte) {
	c := boot.MW.DefaultCache()
	c.Put(common.RedisPrefixCodeIDToCaptcha, id, digits)
}

// Get returns stored digits for the captcha id. Clear indicates
// whether the captcha must be deleted from the store.
func (s *customStore) Get(id string, clear bool) (digits []byte) {
	c := boot.MW.DefaultCache()

	_, _ = c.Get(common.RedisPrefixCodeIDToCaptcha, id, &digits)

	if clear {
		c.Delete(common.RedisPrefixCodeIDToCaptcha, id)
	}

	return digits
}
