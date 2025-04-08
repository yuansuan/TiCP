package rpc

import (
	"context"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	licenseinfo "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_info"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/license"
)

// QueueLicenseTypeList 获取license type列表
func (s *GRPCService) QueueLicenseTypeList(ctx context.Context, in *license.QueueLicenseTypeListRequest) (*license.QueueLicenseTypeListResponse, error) {
	logger := logging.GetLogger(ctx)

	licManagerList, err := s.localAPI.Client.License.ListLicenseManager()
	if err != nil {
		logger.Errorf("list license manager list err: %v", err)
		return nil, err
	}

	//在内存中过滤
	var resultItems []*license.LicenseType
	for _, item := range licManagerList.Data.Items {
		if len(item.LicenseInfos) == 0 {
			continue
		}
		licenseManager := &license.LicenseType{
			Id:           item.Id,
			TypeName:     item.AppType,
			LicenceValid: FindLicenseValid(ctx, item.LicenseInfos),
		}
		resultItems = append(resultItems, licenseManager)
	}

	return &license.QueueLicenseTypeListResponse{
		LicenseTypes: resultItems,
	}, nil
}

func FindLicenseValid(ctx context.Context, licenseInfos []*licenseinfo.GetLicenseInfoResponseData) bool {
	logger := logging.GetLogger(ctx)
	if len(licenseInfos) == 0 {
		logger.Infof("license is null")
		return false
	}
	for _, lic := range licenseInfos {
		if !LicenceValid(ctx, lic.BeginTime, lic.EndTime) {
			return false
		}
	}
	return true
}

func LicenceValid(ctx context.Context, beginTime string, endTime string) bool {
	logger := logging.GetLogger(ctx)
	if beginTime == "" || endTime == "" {
		return false
	}

	layout := "2006-01-02 15:04:05"
	begin, err := time.Parse(layout, beginTime)
	if err != nil {
		logger.Errorf("license begin time err: %v ,begin:%v", err, begin)
		return false
	}
	end, err := time.Parse(layout, endTime)
	if err != nil {
		logger.Errorf("license end time err: %v ,end: %v", err, end)
		return false
	}
	beginTimestamp := begin.Unix()
	endTimestamp := end.Unix()
	currentUnixTime := time.Now().Unix()
	if currentUnixTime >= beginTimestamp && currentUnixTime < endTimestamp {
		return true
	}
	return false
}
