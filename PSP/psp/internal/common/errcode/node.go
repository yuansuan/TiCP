package errcode

import (
	"google.golang.org/grpc/codes"
)

// Node service error codes from 18001 ~ 19000
const (
	ErrNodeQueryFailed     codes.Code = 18001
	ErrNodeFailGet         codes.Code = 18002
	ErrNodeFailOperate     codes.Code = 18003
	ErrOperateTypeNotExist codes.Code = 18004
	ErrNodeFailCoreNum     codes.Code = 18005
	ErrNodeNotExist        codes.Code = 18006
	ErrSchedulerTypeNotSet codes.Code = 18007
	ErrNodeListFail        codes.Code = 18008
	ErrClusterInfoFail     codes.Code = 18009
	ErrResourceInfoFail    codes.Code = 18010
	ErrJobStatusInfoFail   codes.Code = 18011
)
const (
	MsgNodeFailGet         = "获取节点信息失败"
	MsgNodeListFail        = "获取节点列表失败"
	MsgNodeFailOperate     = "更改节点状态失败"
	MsgOperateTypeNotExist = "操作类型不存在"
	MsgNodeFailCoreNum     = "获取节点总核数失败"
	MsgNodeNotExist        = "节点不存在"
	MsgSchedulerTypeNotSet = "调度器类型未设置"
	MsgClusterInfoFail     = "获取集群信息失败"
	MsgResourceInfoFail    = "获取资源信息失败"
	MsgJobStatusInfoFail   = "获取作业状态信息失败"
)
