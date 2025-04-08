package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	openys "github.com/yuansuan/ticp/common/openapi-go"
	"github.com/yuansuan/ticp/common/openapi-go/credential"
	iamClient "github.com/yuansuan/ticp/common/project-root-iam/iam-client"
)

// Cmds 所有的命令
var Cmds []*cobra.Command

// RegisterCmd 注册命令
func RegisterCmd(command *cobra.Command) {
	Cmds = append(Cmds, command)
}

// _YsClient ysadmin的openapi client
var _YsClient *openys.Client

// InitYSClient 初始化ysadmin的openapi client
func InitYSClient() {
	_YsClient = NewYSClient()
	if _YsClient == nil {
		panic("New Ys Client Fail")
	}
}

// GetYsClient 获取ysadmin的openapi client
func GetYsClient() *openys.Client {
	if _YsClient == nil {
		InitYSClient()
	}
	return _YsClient
}

// NewYSClient 创建ysadmin的openapi client
func NewYSClient() *openys.Client {
	fmt.Printf("current env: %s [%s]\n\n", color.New(color.FgGreen).Add(color.Bold).SprintFunc()(Cfg.CurrentEnvironment), Cfg.Environments[Cfg.CurrentEnvironment].Endpoint)

	CheckComputeConfigExit()
	c, err := openys.NewClient(credential.NewCredential(CurrentCfg.ComputeAccessKeyID, CurrentCfg.ComputeAccessKeySecret), openys.WithBaseURL(CurrentCfg.Endpoint))
	if err != nil {
		fmt.Println("new client fail：", err)
		return nil
	}
	return c
}

// NewIamAdminClient 创建iam admin的openapi client
func NewIamAdminClient() *iamClient.IamAdminClient {
	CheckIamAdminConfigExit()
	return iamClient.NewAdminClient(CurrentCfg.IamAdminEndpoint, CurrentCfg.IamAdminAccessKeyID, CurrentCfg.IamAdminAccessKeySecret, "")
}

// NewIamClient 创建iam的openapi client
func NewIamClient() *iamClient.IamClient {
	fmt.Printf("current env: %s [%s]\n\n", color.New(color.FgGreen).Add(color.Bold).SprintFunc()(Cfg.CurrentEnvironment), Cfg.Environments[Cfg.CurrentEnvironment].Endpoint)

	CheckComputeConfigExit()
	return iamClient.NewClient(CurrentCfg.Endpoint, CurrentCfg.ComputeAccessKeyID, CurrentCfg.ComputeAccessKeySecret)
}

// PrintResp 打印响应
func PrintResp(resp interface{}, err error, errTitle string) {
	if err != nil {
		fmt.Printf("%s Error:\n%s\n", errTitle, err.Error())
		return
	}
	if resp == nil {
		return
	}
	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		fmt.Println("MarshalError: \n", err.Error())
		return
	}
	fmt.Println(string(data))
}

// ReadAndUnmarshal 读取文件并解析
func ReadAndUnmarshal(path string, saved interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Read Error: %s, Path: %s\n", err.Error(), path)
		os.Exit(3)
	}

	err = json.Unmarshal(data, saved)
	if err != nil {
		fmt.Printf("unmarshal file data fail: %s, Path: %s\n", err.Error(), path)
		os.Exit(3)
	}
	return nil
}

// BaseOptions 基础参数
type BaseOptions struct {
	Id        string
	State     string
	Zone      string
	Offset    int64
	Limit     int64
	StartTime string
	EndTime   string
	JsonFile  string
	All       bool
}

// AddBaseOptions 添加基础参数
func (o *BaseOptions) AddBaseOptions(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "创建资源的json格式文件")
	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "资源id")
	cmd.Flags().StringVarP(&o.State, "state", "S", "", "作业状态")
	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "", "区域")
	cmd.Flags().Int64VarP(&o.Offset, "offset", "O", 0, "offset")
	cmd.Flags().Int64VarP(&o.Limit, "limit", "L", 1000, "limit")
	cmd.Flags().BoolVarP(&o.All, "all", "", false, "所有的条目")
	cmd.Flags().StringVarP(&o.StartTime, "start_time", "", "", "开始时间")
	cmd.Flags().StringVarP(&o.EndTime, "end_time", "", "", "结束时间")
}

// ValidateArgsCount 验证参数个数
func ValidateArgsCount(cmd *cobra.Command, args []string, expected int) {
	if len(args) < expected {
		cmd.Help()
		os.Exit(2)
	}
}

// 存储类型
const (
	HPCType   = "hpc"
	CloudType = "cloud"
)

// NewStorageClient 创建存储openapi client
func NewStorageClient(zone string, t string, akID string, akSecret string) (*openys.Client, error) {
	endpoint, err := getStorageEndpoint(zone, t)
	if err != nil {
		return nil, err
	}
	c, err := openys.NewClient(credential.NewCredential(akID, akSecret), openys.WithBaseURL(endpoint))
	if err != nil {
		fmt.Println("new storage client fail: ", err)
		return nil, err
	}
	return c, nil
}

// GetStorageClient 获取存储openapi client
func GetStorageClient(zone, t string) *openys.Client {
	printZoneAndType(zone, t)
	var akID, akSecret string
	if t == CloudType {
		CheckStorageConfigExit()
		akID, akSecret = CurrentCfg.StorageAccessKeyID, CurrentCfg.StorageAccessKeySecret
	} else {
		CheckComputeConfigExit()
		akID, akSecret = CurrentCfg.ComputeAccessKeyID, CurrentCfg.ComputeAccessKeySecret
	}

	c, err := NewStorageClient(zone, t, akID, akSecret)
	if err != nil {
		os.Exit(1)
	}
	return c
}

func printZoneAndType(zone, storageType string) {
	blue := color.New(color.FgBlue).SprintFunc()
	fmt.Printf("Zone: %s, Type: %s\n\n", blue(zone), blue(storageType))
}

func getStorageEndpoint(zone string, t string) (string, error) {
	if zone == "local" {
		return localEndPoint, nil
	}
	res, err := GetYsClient().Job.ZoneList()
	if err != nil {
		fmt.Printf("List Zones Fail, Error: %s\n", err.Error())
		return "", err
	}
	if _, ok := res.Data.Zones[zone]; !ok {
		fmt.Printf("Zone %s not Found\n", zone)
		return "", errors.New("zone not found")
	}
	endpoint := res.Data.Zones[zone].StorageEndpoint
	if t == "hpc" {
		endpoint = res.Data.Zones[zone].HPCEndpoint
	}
	if endpoint == "" {
		zones, _ := json.Marshal(res.Data)
		fmt.Printf("this storage endpoint not exists in this zone, Zones: %s\n", string(zones))
		return "", errors.New("EmptyEndpoint")
	}
	return endpoint, nil
}
