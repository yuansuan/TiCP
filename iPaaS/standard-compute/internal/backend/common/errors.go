package common

import _errors "errors"

var (
	// 没有查到节点
	ErrNoNodesAvailable = _errors.New("get cpu usage: 没查到节点")

	// 根据j.Queue读到的配置文件里 节点的CPU数量为0
	ErrZeroCPUsPerNode = _errors.New("get cpu usage: 配置文件中节点CPU数量为0")

	// 分配给作业的CPU核数为0
	ErrZeroAllocCPU = _errors.New("get cpu usage: 分配给作业的CPU核数为0")

	// 未能获取所有节点的CPU使用率
	ErrWrongCPUUsage = _errors.New("get cpu usage: 未能获取所有节点的CPU使用率")
)
