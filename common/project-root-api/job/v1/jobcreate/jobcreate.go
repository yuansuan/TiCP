package jobcreate

import (
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
// swagger:model JobCreateRequest
type Request struct {
	Name          string              `json:"Name"`                      //作业名称
	Params        Params              `json:"Params" binding:"required"` //作业参数
	Timeout       int64               `json:"Timeout"`                   //超时时间,单位为秒,默认为-1表示不设置超时
	Zone          string              `json:"Zone"`                      //分区，默认az-jinan，如果不指定分区，将会有一个调度算法选择一个zone。一个zone就是一个超算中心。
	Comment       string              `json:"Comment"`                   //备注信息
	ChargeParam   schema.ChargeParams `json:"ChargeParam"`               //计费参数
	PayBy         string              `json:"PayBy"`                     //代支付参数
	NoRound       bool                `json:"NoRound"`                   //单节点是否不进行取整,仅限内部用户使用
	AllocType     string              `json:"AllocType"`                 //CPU资源的分配方式,仅限内部用户使用
	PreScheduleID string              `json:"PreScheduleID"`             //预调度ID，如果为不为空则从预调度信息中获取资源信息等
}

// Params 作业参数
// swagger:model JobCreateParams
type Params struct {
	Application       Application       `json:"Application"`               //求解软件信息，非预调度时必填
	Resource          *Resource         `json:"Resource"`                  //作业需要的总资源，非预调度时必填
	EnvVars           map[string]string `json:"EnvVars"`                   //环境变量
	Input             *Input            `json:"Input"`                     //求解输入数据，非预调度时必填，input.Type为hpc_storage且input.Destination为空时不进行实际的数据传输
	Output            *Output           `json:"Output"`                    //存放计算结果的路径前缀，如果为空，将不进行回传。
	TmpWorkdir        bool              `json:"TmpWorkdir"`                //默认为true。如果为false，将在Input.Destination路径下执行作业。如果为true，将在Input.Destination路径前增加一个唯一的临时前缀目录。
	SubmitWithSuspend bool              `json:"SubmitWithSuspend"`         //如果为true，提交的作业状态会变成 InitiallySuspended，用户可以通过resume操作，让job进入Pending状态。默认false。
	CustomStateRule   *CustomStateRule  `json:"CustomStateRule,omitempty"` //自定义状态规则，如果为空，将使用默认规则
}

// Application 求解软件信息
// swagger:model JobCreateApplication
type Application struct {
	Command string `json:"Command"` //求解器命令行，不填的情况为提交非命令行作业。
	AppID   string `json:"AppID"`   //求解器id，可以通过ListApp接口获取
}

// Resource 作业需要的总资源
// swagger:model JobCreateResource
type Resource struct {
	Cores  *int `json:"Cores" binding:"required"` //期望的核数，实际数值可能有稍微差别
	Memory *int `json:"Memory"`                   //期望的内存数，单位为M，暂未起作用
}

// Input 求解输入数据
// swagger:model JobCreateInputFile
type Input struct {
	Type        string `json:"Type" binding:"required"`   //输入数据类型为超算存储或者远算云盒子
	Source      string `json:"Source" binding:"required"` //输入数据的路径，必须是文件夹路径，绝对路径，必须以"/"开头，不包含".."。带域名的绝对路径。不同区域的存储服务域名不同。
	Destination string `json:"Destination"`               //输入文件的目标路径，不准为空，不准以"/", "."开头。如果dest为空，默认和Source同路径。作业的工作目录
}

// Output 存放计算结果的路径前缀，如果为空，将不进行回传。
// swagger:model JobCreateOutputFile
type Output struct {
	Type          string `json:"Type"`          //结算结果路径类型，暂时只支持和input同类型。
	Address       string `json:"Address"`       //计算结果存放路径，带域名的全路径。
	NoNeededPaths string `json:"NoNeededPaths"` //正则表达式，符合规则的文件路径将不会被回传
	NeededPaths   string `json:"NeededPaths"`   //正则表达式，符合规则的文件路径将会被回传
}

// CustomStateRule 自定义状态规则
// swagger:model JobCreateCustomStateRule
type CustomStateRule struct {
	KeyStatement string `json:"KeyStatement,omitempty"` //关键字，如果作业输出中包含该关键字，则会修改状态为ResultState
	ResultState  string `json:"ResultState,omitempty"`  //包含KeyStatement关键字时，认为的结果状态，仅有两种可选择[ completed | failed ]
}

// ResultState 可选择的作业结果状态
const (
	ResultStateCompleted = "completed"
	ResultStateFailed    = "failed"
)

// Response 返回
// swagger:model JobCreateResponse
type Response struct {
	schema.Response `json:",inline"`
	Data            *Data `json:"Data,omitempty"`
}

// Data 数据
// swagger:model JobCreateData
type Data struct {
	JobID string `json:"JobID"`
}
