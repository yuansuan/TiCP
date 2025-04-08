package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/quota/admin"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	storageQuotaService "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/storageQuota"
	"math"
	"net/http"
)

func (q *Quota) ListStorageQuota(ctx *gin.Context) {
	userID, _, err := q.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "ListStorageQuota", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &admin.ListStorageQuotaRequest{}
	if err := ctx.BindQuery(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, err.Error())
		return
	}

	if request.PageOffset < 0 {
		msg := fmt.Sprintf("invalid page offset, page offset should be greater than or equal to 0, got: %v", request.PageOffset)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPageOffset, msg)
		return
	}

	if request.PageSize < 1 || request.PageSize > 1000 {
		msg := fmt.Sprintf("invalid page size, page size should be in [1, 1000], got: %v", request.PageSize)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPageSize, msg)
		return
	}

	res, err, nextMarker, total := storageQuotaService.ListStorageQuotaInfo(ctx, q.Engine, q.StorageQuotaDao, request.PageOffset, request.PageSize)
	if err != nil {
		msg := fmt.Sprintf("Error list storage quota: %v", err)
		logger.Errorf(msg)
		common.InternalServerError(ctx, "Error list storage quota")
		return
	}

	storageQuota := make([]*admin.StorageQuota, 0, len(res))
	for _, v := range res {
		storageQuotaLimit := v.StorageLimit
		if storageQuotaLimit == 0 {
			storageQuotaLimit = float64(q.DefaultUserStorageLimit)
		}
		storageQuota = append(storageQuota, &admin.StorageQuota{
			UserID:       v.UserId.String(),
			StorageUsage: math.Round(v.StorageUsage*100) / 100,
			StorageLimit: math.Round(storageQuotaLimit*100) / 100,
		})
	}

	common.SuccessResp(ctx, &admin.ListStorageQuotaResponseData{
		Total:        total,
		NextMarker:   nextMarker,
		StorageQuota: storageQuota,
	})

}
