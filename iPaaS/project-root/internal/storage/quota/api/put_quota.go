package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/quota/admin"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	storageQuotaService "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/storageQuota"
	"math"
	"net/http"
	"time"
)

func (q *Quota) PutStorageQuota(ctx *gin.Context) {
	userID, _, err := q.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "PutStorageQuota", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	userIDStr := ctx.Param(UserIDKey)
	if userIDStr == "" {
		msg := fmt.Sprintf("Error parsing userID: %s,error: %v", userIDStr, err.Error())
		logger.Infof(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidUserID, "invalid user id")
		return
	}
	id, err := snowflake.ParseString(userIDStr)
	if err != nil {
		msg := fmt.Sprintf("Error parsing userID: %s,error: %v", userIDStr, err.Error())
		logger.Infof(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidUserID, "invalid user id")
		return
	}

	request := &admin.PutStorageQuotaRequest{}
	if err := ctx.BindJSON(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, err.Error())
		return
	}

	if request.StorageLimit <= 0 || request.StorageLimit > float64(q.MaxUserStorageLimit) {
		msg := fmt.Sprintf("invalid storage limit: %v", math.Round(request.StorageLimit*100)/100)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidStorageLimit, msg)
		return
	}

	exists, _, err := storageQuotaService.GetStorageQuotaInfo(ctx, q.Engine, q.StorageQuotaDao, &model.StorageQuota{UserId: id})
	if err != nil {
		msg := fmt.Sprintf("Error getting storage quota: %v", err)
		logger.Errorf(msg)
		common.InternalServerError(ctx, "Error getting storage quota")
		return
	}
	if !exists {
		//insert storage quota
		storageQuota := &model.StorageQuota{
			UserId:       id,
			StorageLimit: request.StorageLimit,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		err = storageQuotaService.InsertStorageQuotaInfo(ctx, q.Engine, q.StorageQuotaDao, storageQuota)
		if err != nil {
			msg := fmt.Sprintf("Error inserting storage quota: %v", err)
			logger.Errorf(msg)
			common.InternalServerError(ctx, "Error inserting storage quota")
			return
		}
	} else {
		//update storage quota
		storageQuota := &model.StorageQuota{
			UserId:       id,
			StorageLimit: request.StorageLimit,
			UpdateTime:   time.Now(),
		}
		err = storageQuotaService.UpdateStorageQuotaInfo(ctx, q.Engine, q.StorageQuotaDao, id, storageQuota)
		if err != nil {
			msg := fmt.Sprintf("Error updating storage quota: %v", err)
			logger.Errorf(msg)
			common.InternalServerError(ctx, "Error updating storage quota")
			return
		}
	}

	common.SuccessResp(ctx, nil)

}
