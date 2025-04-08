package common

import (
	"fmt"
	"strings"
)

const (
	yrnPartition string = "ys"
)

type YRN struct {
	Partition    string
	Service      string
	Region       string
	AccountID    string
	ResourceType string
	ResourceID   string
}

func (y *YRN) String() string {
	return fmt.Sprintf("yrn:%s:%s:%s:%s:%s/%s", y.Partition, y.Service, y.Region, y.AccountID, y.ResourceType, y.ResourceID)
}

func ParseYRN(arnStr string) (yrn *YRN, err error) {
	ps := strings.SplitN(arnStr, ":", 6)
	if len(ps) != 6 ||
		ps[0] != "yrn" {
		err = fmt.Errorf("Invalid YRN string format, Like yrn:partition:service:region:account-id:resource-id")
		return
	}

	if ps[1] != yrnPartition {
		err = fmt.Errorf("Invalid YRN - bad partition field")
		return
	}

	res := strings.SplitN(ps[5], "/", 2)
	if len(res) != 2 {
		err = fmt.Errorf("Invalid YRN - resource does not contain a \"/\"")
		return
	}

	if res[0] == "" {
		err = fmt.Errorf("Invalid YRN - missing resource-type field")
		return
	}

	if res[1] == "" {
		err = fmt.Errorf("Invalid YRN - missing resource-id field")
		return
	}

	yrn = &YRN{
		Partition:    ps[1],
		Service:      ps[2],
		Region:       ps[3],
		AccountID:    ps[4],
		ResourceType: res[0],
		ResourceID:   res[1],
	}
	return
}
