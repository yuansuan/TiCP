package dao

import (
	"context"

	"github.com/qiniu/qmgo"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"go.mongodb.org/mongo-driver/bson"
)

type residualDaoImpl struct {
	coll *qmgo.Collection
}

func NewResidualDaoImpl(coll *qmgo.Collection) ResidualDao {
	return &residualDaoImpl{
		coll: coll,
	}
}

// GetJobResidual 获取作业残差图
func (r *residualDaoImpl) GetJobResidual(ctx context.Context, jobId snowflake.ID) (*models.Residual, error) {
	residual := models.Residual{}
	err := r.coll.Find(ctx, bson.M{"job_id": jobId}).One(&residual)
	return &residual, err
}

// GetUnfinishedResidual 获取未完成的残差图
func (r *residualDaoImpl) GetUnfinishedResidual(ctx context.Context) ([]*models.Residual, error) {
	var residuals []*models.Residual
	err := r.coll.Find(ctx, bson.M{"finished": false}).All(&residuals)
	return residuals, err
}

// InsertResidual 插入残差图
func (r *residualDaoImpl) InsertResidual(ctx context.Context, residual *models.Residual) error {
	_, err := r.coll.InsertOne(ctx, residual)
	return err
}

// UpdateResidual 更新残差图
func (r *residualDaoImpl) UpdateResidual(ctx context.Context, residual *models.Residual) error {
	err := r.coll.UpdateOne(ctx, bson.M{"id": residual.ID}, bson.M{"$set": residual})
	return err
}
