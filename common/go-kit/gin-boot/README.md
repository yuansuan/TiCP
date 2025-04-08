# github.com/yuansuan/ticp/common/go-kit/gin-boot

这是一个基于gin的微服务框架，支持自动配置、REST API、gRPC、Metric监控、注解、优雅停服。

## 说明

requirement: go version >= 1.12，推荐最新稳定版本。

## 安装

### 快速生成项目

```bash
git clone ssh://vcssh@phabricator.intern.yuansuan.cn/source/helloworld.git
```

### 自定义项目导入包

导入master

```bash
go get github.com/yuansuan/ticp/common/go-kit/gin-boot@master
```

或者按tag导入（推荐，可```go mod tidy```更新）

```bash
go get github.com/yuansuan/ticp/common/go-kit/gin-boot
```

## 编译、运行项目

```bash
cd helloworld
#安装依赖包, 
#需要翻墙: 
export https_proxy=10.81.254.21:3128
export all_proxy=10.81.254.21:3128
go get -v -x ./...
#编译, 请提前安装 protoc: https://github.com/protocolbuffers/protobuf/releases/ 和 https://github.com/golang/protobuf
make build
#运行
./helloworld 
```

## 查看运行状态

```bash
#REST HTTP:
curl -isv 'http://127.0.0.1:8080/'

```

## 环境变量文件

项目根目录中 .env

```ini
Mode=[dev|test|stage|prod]
LogLevel=[debug|info|warn|error|fatal|off]
Type={任何string变量，如IDC}
```

## API用法

在IDE中键入"boot."会自动提示API的使用方式。

可参考示例: github.com/yuansuan/ticp/go-kit/examples/

```go
func main() {
boot.
//使用默认http server
Default().
Register( //注册路由策略
router.UseRoutersGenerated,
handler.InitHandlers,
handler_rpc.InitGRPCServer,
).RegisterRoutine( //注册go-routine在后台运行
handler_rpc.InitGRPCClient,
).OnShutdown( //注册退出事件
handler_rpc.OnShutdown,
).Run() //启动运行
}

func monitor(){
//添加counter
boot.Monitor.Add("show", 1, nil)
//设置counter
boot.Monitor.Set("show", 10, nil)
//设置summary/histogram
boot.Monitor.Observe("show", 10, monitor.Objectives{1: 0.1, 2:0.3}, []*monitor.Label{
{"color", "red"},
{"color", "yellow"},
})
}

func middleware(){
rows, _ := boot.MW.DefaultMysql().Query("select id,img,tag from img_tags_audit limit 4")
for rows.Next() {
_ = rows.Scan(&doc.ID, &doc.Img, &doc.Tag)
}

res, e := boot.MW.Redis("demo2").Info().Result()
}

func log(){
//支持方法有： debug[f]|info[f]|warn[f]|error[f]|fatal[f]
//支持自定义Map
boot.Logger.Fatal(logging.Map{
"panic":  r,
"detail": fmt.Sprintf("%s", logging.Stack(5)),
})
//支持格式化日志
boot.Logger.Debugf("mysql %v %v %v ", doc.ID, doc.Img, doc.Tag)

}

func grpc(){
//创建默认grpc client
clientConn, err := boot.GRPC.DefaultClient()
add := protos.NewAddClient(clientConn)

//创建grpc server，注册
s, err := boot.GRPC.DefaultServer()
protos.RegisterAddServer(s.Driver(), &server{})
}

```

## 一键生成代码

### 部分功能需要安装nodejs、python

```bash
# json to go
sh ${GOPATH}/src/github.com/yuansuan/ticp/common/go-kit/gin-boot/bin/cli.sh gen j2g [file or dir(include .json)] 
# yaml to go
sh ${GOPATH}/src/github.com/yuansuan/ticp/common/go-kit/gin-boot/bin/cli.sh gen y2g [file or dir(include .yml)] 
# proto to go
sh ${GOPATH}/src/github.com/yuansuan/ticp/common/go-kit/gin-boot/bin/cli.sh gen p2g [file or dir(include .proto)] 

```


 
