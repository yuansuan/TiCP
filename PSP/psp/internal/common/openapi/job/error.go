package job

import (
	"errors"
	"fmt"

	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

var (
	// ErrJobIDEmpty 作业ID不能为空
	ErrJobIDEmpty = errors.New("the job id is empty")

	// ErrSnapshotPathIsNotExist 云图路径不存在
	ErrSnapshotPathIsNotExist = errors.New("the snapshot path is not exist")

	// ErrJobIsNotExist 作业不存在
	ErrJobIsNotExist = errors.New("the job is not exist")

	// ErrResidualIsNotExist 残差图不存在
	ErrResidualIsNotExist = errors.New("the residual is not exist")

	// ErrSnapshotsIsNotExist 云图集不存在
	ErrSnapshotsIsNotExist = errors.New("the snapshots is not exist")

	// ErrSnapshotIsNotExist 云图不存在
	ErrSnapshotIsNotExist = errors.New("the snapshot is not exist")

	// ErrPageIndexInvalid 分页索引值不合规
	ErrPageIndexInvalid = errors.New("page index must greater than '0'")

	// ErrPageSizeInvalid 分页大小不合规
	ErrPageSizeInvalid = fmt.Errorf("page size cannot greater than '%v'", xtype.MaxPageSize)
)
