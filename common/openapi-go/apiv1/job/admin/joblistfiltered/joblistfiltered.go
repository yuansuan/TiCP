package joblistfiltered

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	jlf "github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/joblistfiltered"
	"net/http"
	"time"
)

type Option func(req *jlf.Request) error

type API func(options ...Option) (*jlf.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*jlf.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		builder := xhttp.NewRequestBuilder().URI("/admin/jobs/filtered")
		builder.AddQuery("PageSize", utils.Stringify(*req.PageSize))
		builder.AddQuery("PageOffset", utils.Stringify(*req.PageOffset))
		builder.AddQuery("WithDelete", utils.Stringify(req.WithDelete))
		builder.AddQuery("IsSystemFailed", utils.Stringify(req.IsSystemFailed))
		if req.JobState != "" {
			builder.AddQuery("JobState", req.JobState)
		}
		if req.UserID != "" {
			builder.AddQuery("UserID", req.UserID)
		}
		if req.Zone != "" {
			builder.AddQuery("Zone", req.Zone)
		}
		if req.Name != "" {
			builder.AddQuery("Name", req.Name)
		}
		if req.JobID != "" {
			builder.AddQuery("JobID", req.JobID)
		}
		if req.AppID != "" {
			builder.AddQuery("AppID", req.AppID)
		}
		if req.AccountID != "" {
			builder.AddQuery("AccountID", req.AccountID)
		}
		if req.FileSyncState != "" {
			builder.AddQuery("FileSyncState", req.FileSyncState)
		}
		if req.StartTime.IsZero() == false {
			builder.AddQuery("StartTime", req.StartTime.Format(time.RFC3339))
		}
		if req.EndTime.IsZero() == false {
			builder.AddQuery("EndTime", req.EndTime.Format(time.RFC3339))
		}

		resolver := hc.Prepare(builder)

		ret := new(jlf.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*jlf.Request, error) {
	offset := int64(0)
	size := int64(100)
	req := &jlf.Request{}
	req.PageOffset = &offset
	req.PageSize = &size

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

func (api API) PageOffset(pageOffset int64) Option {
	return func(req *jlf.Request) error {
		req.PageOffset = &pageOffset
		return nil
	}
}

func (api API) PageSize(size int64) Option {
	return func(req *jlf.Request) error {
		req.PageSize = &size
		return nil
	}
}

func (api API) JobState(jobState string) Option {
	return func(req *jlf.Request) error {
		req.JobState = jobState
		return nil
	}
}

func (api API) Zone(zone string) Option {
	return func(req *jlf.Request) error {
		req.Zone = zone
		return nil
	}
}

func (api API) Name(name string) Option {
	return func(req *jlf.Request) error {
		req.Name = name
		return nil
	}
}

func (api API) JobID(jobId string) Option {
	return func(req *jlf.Request) error {
		req.JobID = jobId
		return nil
	}
}

func (api API) UserID(userId string) Option {
	return func(req *jlf.Request) error {
		req.UserID = userId
		return nil
	}
}

func (api API) AppID(appId string) Option {
	return func(req *jlf.Request) error {
		req.AppID = appId
		return nil
	}
}

func (api API) AccountID(accountID string) Option {
	return func(req *jlf.Request) error {
		req.AccountID = accountID
		return nil
	}
}

func (api API) FileSyncState(syncState string) Option {
	return func(req *jlf.Request) error {
		req.FileSyncState = syncState
		return nil
	}
}

func (api API) StartTime(startTime time.Time) Option {
	return func(req *jlf.Request) error {
		req.StartTime = startTime
		return nil
	}
}

func (api API) EndTime(endTime time.Time) Option {
	return func(req *jlf.Request) error {
		req.EndTime = endTime
		return nil
	}
}

func (api API) WithDelete(withDelete bool) Option {
	return func(req *jlf.Request) error {
		req.WithDelete = withDelete
		return nil
	}
}

func (api API) IsSystemFailed(isSystemFailed bool) Option {
	return func(req *jlf.Request) error {
		req.IsSystemFailed = isSystemFailed
		return nil
	}
}
