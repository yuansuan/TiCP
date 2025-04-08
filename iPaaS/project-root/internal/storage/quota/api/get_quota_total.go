package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/storage/quota/api"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	storageQuotaService "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/storageQuota"
	"math"
)

func (q *Quota) GetStorageQuotaTotal(ctx *gin.Context) {
	userID, _, err := q.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "GetStorageQuotaTotal", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	storageUsageTotal, err := storageQuotaService.GetStorageQuotaTotal(ctx, q.Engine, q.StorageQuotaDao)
	if err != nil {
		msg := "Error getting storage quota total, err: " + err.Error()
		logger.Infof(msg)
		common.InternalServerError(ctx, "Error getting storage quota total")
		return
	}

	common.SuccessResp(ctx, &api.GetStorageQuotaResponseData{
		StorageUsage: math.Round(storageUsageTotal*100) / 100,
		StorageLimit: float64(q.MaxSystemStorageLimit),
	})

}
