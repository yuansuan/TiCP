package dto

// Request ...
type Request struct {
	Start int64 `form:"start"` //开始时间，毫秒单位
	End   int64 `form:"end"`   //结束时间，毫秒单位
}

// ClusterResponse ...
type ClusterResponse struct {
	ClusterInfo *ClusterInfo  `json:"clusterInfo"` //集群信息
	NodeList    []*NodeDetail `json:"nodeList"`    //节点信息
	Disks       *Disk         `json:"disks"`       //磁盘信息
}

// ClusterInfo ...
type ClusterInfo struct {
	ClusterName      string `json:"clusterName"`      //集群名字
	Cores            int    `json:"cores"`            //总核数
	UsedCores        int    `json:"usedCores"`        // 已使用核数
	FreeCores        int    `json:"freeCores"`        // 空闲核数
	AvailableNodeNum int    `json:"availableNodeNum"` // 可用节点数
	TotalNodeNum     int    `json:"totalNodeNum"`     // 总节点数
}

// Disk ...
type Disk struct {
	Data   []map[string]interface{} `json:"data"`   //磁盘信息
	Fields []string                 `json:"fields"` //磁盘挂在路径
}

// ResourceResponse CPU平均利用率/磁盘IO速率/内存平均利用率
type ResourceResponse struct {
	MetricCpuUtAvg []*ValueStruct `json:"metric_cpu_ut_avg"` //CPU 平均利用率
	MetricIoUtAvg  []*ValueStruct `json:"metric_io_ut_avg"`  //磁盘IO速率
	MetricMemUtAvg []*ValueStruct `json:"metric_mem_ut_avg"` //内存平均利用率
}

// ValueStruct ...
type ValueStruct struct {
	Values []*Value `json:"d"` //值集合
	Name   string   `json:"n"` //名字
}

type Value struct {
	T int64   `json:"t"` //时间戳
	V float64 `json:"v"` //值
}

// JobResponse 作业状态
type JobResponse struct {
	JobResLatest []*JobStatusValue `json:"jobResLatest"` //当前作业状态
	JobResRange  []*JobStatusValue `json:"jobResRange"`  //作业状态记录
}
type JobResLatest struct {
	Suspended int `json:"suspended"`
	Pending   int `json:"pending"`
	Running   int `json:"running"`
	Exited    int `json:"exited"`
	Done      int `json:"done"`
}
type JobStatusValue struct {
	JobCount  int64  `json:"job_count"` //作业数量
	Status    string `json:"status"`    //作业状态
	Timestamp int64  `json:"timestamp"` //时间戳
}
