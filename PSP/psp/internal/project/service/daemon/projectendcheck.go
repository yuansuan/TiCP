package daemon

import (
	"context"
	"fmt"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	pbnotice "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/notice"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/service/client"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

// ProjectCheck ...
type ProjectCheck struct {
	projectDao dao.ProjectDao
	sid        *snowflake.Node
	rpc        *client.GRPC
}

// NewProjectCheck ...
func NewProjectCheck() (*ProjectCheck, error) {
	projectDao, err := dao.NewProjectDao()
	if err != nil {
		return nil, err
	}

	node, err := snowflake.GetInstance()
	if err != nil {
		return nil, err
	}

	rpc := client.GetInstance()

	impl := &ProjectCheck{
		projectDao: projectDao,
		sid:        node,
		rpc:        rpc,
	}

	return impl, nil
}

// ProjectCheckStart 项目结束检查
func (loader *ProjectCheck) ProjectCheckStart() {
	go loader.ProjectCheckTicker()
	go loader.UpdateProjectStatus()
}

// ProjectCheckTicker 检查
func (loader *ProjectCheck) ProjectCheckTicker() {
	ctx := context.Background()
	logger := logging.GetLogger(ctx)

	// 根据配置启用
	projectCheck := config.GetConfig().ProjectCheck
	if !projectCheck.Enable {
		logger.Infof("check project end routine has disabled")
		return
	}

	currentTime := time.Now()
	dailyTime := projectCheck.DailyTime
	if dailyTime == 0 {
		dailyTime = 9 // 默认9点
	}
	minTime := projectCheck.MinTime

	tomorrow9AM := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+1, dailyTime, minTime, 0, 0, currentTime.Location())
	duration := tomorrow9AM.Sub(currentTime)

	timer := time.NewTimer(duration)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if err := loader.ProjectCheck(); err != nil {
				logger.Errorf("project end check err: %v", err)
			}
			// 重新计算下个9点的定时区间
			currentTime := time.Now()
			tomorrow9AM := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+1, dailyTime, minTime, 0, 0, currentTime.Location())
			duration := tomorrow9AM.Sub(currentTime)
			timer.Reset(duration)
		}
	}
}

// ProjectCheck 检查
func (loader *ProjectCheck) ProjectCheck() error {
	ctx := context.Background()
	logger := logging.GetLogger(ctx)

	start, end := getTomorrow()

	//1.查询明日 要结束的项目
	projects, err := loader.projectDao.GetProjectListByTimePeriod(ctx, start, end)
	if err != nil {
		return err
	}

	//2.项目结束，发送提示消息
	for _, project := range projects {
		content := fmt.Sprintf("项目[编号:%v 项目名:%v] %v", project.Id, project.ProjectName, consts.EndProjectContent)
		msg := &pbnotice.WebsocketMessage{
			UserId:  project.ProjectOwner.String(),
			Type:    common.ProjectEventType,
			Content: content,
		}

		if _, err := loader.rpc.Notice.SendWebsocketMessage(ctx, msg); err != nil {
			logger.Errorf("project end send err: %v", err)
		}
	}

	return nil
}

// UpdateProjectStatus 更新项目状态
func (loader *ProjectCheck) UpdateProjectStatus() {
	ctx := context.Background()
	logger := logging.GetLogger(ctx)

	// 根据配置启用
	syncData := config.GetConfig().ProjectCheck
	if !syncData.Enable {
		logger.Infof("update project status routine has disabled")
		return
	}

	timerDuration := time.Second * time.Duration(syncData.Interval)
	timer := time.NewTimer(timerDuration)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if err := loader.UpdateStatus(); err != nil {
				logger.Errorf("project end check err: %v", err)
			}
			timer.Reset(timerDuration)
		}
	}
}

// UpdateStatus 更新状态
func (loader *ProjectCheck) UpdateStatus() error {
	ctx := context.Background()
	logger := logging.GetLogger(ctx)

	//1.将已经结束的项目 状态更新为 Completed
	err := loader.projectDao.UpdateProjectStatus(ctx, 0, time.Now().Unix(), common.ProjectCompleted)
	if err != nil {
		logger.Errorf("update project status err: %v", err)
		return err
	}

	//2.将已经开始的项目 状态更新为 Running
	err = loader.projectDao.UpdateProjectStatus(ctx, time.Now().Unix(), 0, common.ProjectRunning)
	if err != nil {
		logger.Errorf("update project status err: %v", err)
		return err
	}

	return nil
}

func getTomorrow() (start, end int64) {
	// 获取当前时间
	now := time.Now()

	// 获取明天的日期
	tomorrow := now.AddDate(0, 0, 1)

	// 构建明天的起始时间和结束时间
	startOfDay := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
	endOfDay := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 23, 59, 59, 999999999, tomorrow.Location())

	return startOfDay.Unix(), endOfDay.Unix()
}
