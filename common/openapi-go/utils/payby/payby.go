package payby

import (
	"fmt"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/openapi-go/credential"
	"github.com/yuansuan/ticp/common/openapi-go/utils/signer"
	"strconv"
	"strings"
)

type PayBy struct {
	accessKeyID  string
	accessSecret string
	resourceTag  string
	timestamp    int64
	signature    string
}

func NewPayBy(accessKeyID, accessSecret, resourceTag string, timestamp int64) (*PayBy, error) {
	payBy := &PayBy{
		accessKeyID:  accessKeyID,
		accessSecret: accessSecret,
		resourceTag:  strings.ReplaceAll(resourceTag, ":", "+"),
		timestamp:    timestamp,
	}

	sign, err := payBy.Sign()
	if err != nil {
		logging.Default().Errorf("accessKeyID: %s generate sign err: %v", accessKeyID, err)
		return payBy, err
	}
	payBy.signature = sign

	return payBy, nil
}

func ParseToken(str string) (*PayBy, error) {
	if len(str) == 0 {
		return nil, fmt.Errorf("request string is empty")
	}

	// "appKeyID:resourceTagID:timestamp:signature"
	sArray := strings.Split(str, ":")
	if len(sArray) != 4 {
		return nil, fmt.Errorf("parse payBy token invalid: %s", str)
	}

	timestamp, err := strconv.ParseInt(sArray[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse payBy token timestamp %s invalid, err: %v", sArray[2], err)
	}

	return &PayBy{
		accessKeyID: sArray[0],
		resourceTag: sArray[1],
		timestamp:   timestamp,
		signature:   sArray[3],
	}, nil
}

func (p *PayBy) Token() string {
	if p.accessKeyID == "" || p.accessSecret == "" || p.resourceTag == "" {
		return ""
	}

	timestamp := strconv.FormatInt(p.timestamp, 10)
	token := strings.Join([]string{p.accessKeyID, p.resourceTag, timestamp, p.signature}, ":")
	return token
}

func (p *PayBy) SignEqualTo(target *PayBy) bool {
	if target == nil {
		return false
	}

	return p.signature == target.signature
}

func (p *PayBy) Sign() (string, error) {
	m := make(map[string]interface{})
	m["access_key_id"] = p.accessKeyID
	m["resource_tag"] = p.resourceTag
	m["time_stamp"] = p.timestamp

	tokenSigner, err := signer.NewSigner(credential.NewCredential(p.accessKeyID, p.accessSecret))
	if err != nil {
		return "", err
	}

	sign, err := tokenSigner.Sign(m)
	if err != nil {
		return "", err
	}

	return sign.Signature, nil
}

func (p *PayBy) GetAccessKeyID() string {
	return p.accessKeyID
}

func (p *PayBy) GetTimestamp() int64 {
	return p.timestamp
}

func (p *PayBy) GetResourceTag() string {
	return p.resourceTag
}

func (p *PayBy) GetSignature() string {
	return p.signature
}
