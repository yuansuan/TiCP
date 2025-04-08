package openstack

import (
	"encoding/json"
	"strings"

	_servers "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud/openstack/compute/servers"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
)

type instance struct {
	raw *_servers.Server
}

func newInstance(server *_servers.Server) *instance {
	return &instance{raw: server}
}

// ID 实例的ID
func (i *instance) ID() string {
	if i == nil || i.raw == nil {
		return ""
	}

	return i.raw.ID
}

// Raw 返回实例的原始数据
func (i *instance) Raw() string {
	if i == nil || i.raw == nil {
		return ""
	}

	rawData, err := json.Marshal(i.raw)
	if err != nil {
		return ""
	}

	return string(rawData)
}

var vmStateMap = map[string]models.InstanceStatus{
	"ACTIVE":  models.InstanceRunning,
	"BUILD":   models.InstancePending,
	"REBOOT":  models.InstanceRebooting,
	"ERROR":   models.InstanceLaunchFailed,
	"STOPPED": models.InstanceStopped,
}

var taskStateMap = map[string]models.InstanceStatus{
	"DELETING":       models.InstanceTerminating,
	"REBOOTING":      models.InstanceRebooting,
	"STOPPING":       models.InstanceStopping,
	"STARTING":       models.InstanceStarting,
	"REBOOT_STARTED": models.InstanceRebooting,
	"POWERING-OFF":   models.InstanceStopping,
	"POWERING-ON":    models.InstanceStarting,
}

// Status 实例当前状态
func (i *instance) Status() models.InstanceStatus {
	if i == nil || i.raw == nil {
		return models.InstanceInvalid
	}

	vmState, vmStateExist := vmStateMap[strings.ToUpper(i.raw.VmState)]
	taskState, taskStateExist := taskStateMap[strings.ToUpper(i.raw.TaskState)]
	if !vmStateExist && !taskStateExist {
		return models.InstanceInvalid
	}

	// 正处于某个任务运行中状态，此时vmState也会存在，优先返回taskState对应的vis_ibv内部定义的状态
	if taskStateExist {
		return taskState
	}

	return vmState
}

func (i *instance) String() string {
	return i.Raw()
}
