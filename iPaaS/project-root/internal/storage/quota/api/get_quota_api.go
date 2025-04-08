package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/quota/api"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	storageQuotaService "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/storageQuota"
	"math"
	"net/http"
)

func (q *Quota) GetStorageQuotaAPI(ctx *gin.Context) {
	userID, _, err := q.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "GetStorageQuotaAPI", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	if userID == "" {
		msg := fmt.Sprintf("Error parsing userID: %v,got: %v", err, userID)
		logger.Infof(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidUserID, "invalid user id")
		return
	}
	id, err := snowflake.ParseString(userID)
	if err != nil {
		msg := fmt.Sprintf("Error parsing userID: %v,got: %v", err, userID)
		logger.Infof(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidUserID, "invalid user id")
		return
	}

	exists, storageQuota, err := storageQuotaService.GetStorageQuotaInfo(ctx, q.Engine, q.StorageQuotaDao, &model.StorageQuota{UserId: id})
	if err != nil {
		msg := fmt.Sprintf("Error getting storage quota: %v", err)
		logger.Errorf(msg)
		common.InternalServerError(ctx, "Error getting storage quota")
		return
	}
	if !exists {
		msg := fmt.Sprintf("Storage quota not found for user: %v", id)
		logger.Infof(msg)
		common.ErrorResp(ctx, http.StatusNotFound, commoncode.StorageQuotaNotFound, msg)
		return
	}

	storageQuotaLimit := storageQuota.StorageLimit
	if storageQuotaLimit == 0 {
		storageQuotaLimit = float64(q.DefaultUserStorageLimit)
	}

	common.SuccessResp(ctx, &api.GetStorageQuotaResponseData{
		StorageUsage: math.Round(storageQuota.StorageUsage*100) / 100,
		StorageLimit: math.Round(storageQuotaLimit*100) / 100,
	})

}
