# 超算中心管理

超算中心管理, 直连超算, 可查看HPC剩余资源以及在HPC集群上执行shell命令

使用方法:

ysadmin hpc

可用子命令:

- cmd   在HPC集群上执行shell命令
- freeresource   获取HPC剩余资源

## 在HPC集群上执行shell命令

在HPC集群上执行shell命令

使用方法:

ysadmin hpc cmd [cmd] [flags]

示例:

- 在HPC集群上执行ls -l命令, 区域为az-zhigu
  - ysadmin hpc cmd "ls -l" -Z az-zhigu
- 在HPC集群上执行ls -l命令, 超时时间为20秒, 区域为az-zhigu
  - ysadmin hpc cmd "ls -l" -Z az-zhigu -T 20

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -T, --timeout |  int |   执行命令的超时时间 \(default 10\) |
| -Z, --zone |  string |   区域 \(必填\) |

## 获取HPC剩余资源

获取HPC剩余资源

使用方法:

ysadmin hpc freeresource [flags]

示例:

- 获取HPC剩余资源, 区域为az-zhigu
  - ysadmin hpc freeresource -Z az-zhigu

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -Z, --zone |  string |   区域 \(必填\) |

