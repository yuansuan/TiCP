# 作业管理

作业管理，用于提交作业，查询作业，删除作业，终止作业

使用方法:

ysadmin job

可用子命令:

- delete   删除作业
- example   创建作业示例文件
- get   查询作业
- list   查询作业列表
- submit   提交作业
- terminate   终止作业
- retransmit   重传作业结果文件

## 删除作业

删除作业, 只可删除状态在[Terminated, Completed, Failed]且回传状态在[Paused, Completed, Failed]的作业

使用方法:

ysadmin job delete [flags]

示例:

- 删除作业ID为 52hgnbSTGWd 的作业
  - ysadmin job delete -I 52hgnbSTGWd

可用Flags:

| 命令参数 |  类型  | 说明           |
| -------: | :----: | :------------- |
| -I, --id | string | 作业ID\(必填\) |

## 创建作业示例文件

创建作业示例文件

使用方法:

ysadmin job example

示例:

ysadmin job example

## 查询作业

查询作业

使用方法:

ysadmin job get [flags]

示例:

- 查询作业ID为 52hgnbSTGWd 的作业
  - ysadmin job get -I 52hgnbSTGWd

可用Flags:

| 命令参数 |  类型  | 说明           |
| -------: | :----: | :------------- |
| -I, --id | string | 作业ID\(必填\) |

## 查询作业列表

查询作业列表, 可以指定作业状态, 区域, 偏移量, 限制条数, 用户ID, 默认偏移量0, 限制条数1000

使用方法:

ysadmin job list [flags]

示例:

- 查询所有作业(默认1000条)
  - ysadmin job list
- 查询分区为az-zhigu的所有运行中的作业, 偏移量为0, 限制条数为1000
  - ysadmin job list -S running -Z az-zhigu -O 0 -L 1000
- 查询用户ID为4TiSsZonTa3的所有作业
  - ysadmin job list -U 4TiSsZonTa3

可用Flags:

|      命令参数 |  类型  | 说明                                                                                                                                 |
| ------------: | :----: | :----------------------------------------------------------------------------------------------------------------------------------- |
|   -L, --limit |  int  | 指定的条数上限, 最多1000\(default 1000\)                                                                                             |
|  -O, --offset |  int  | 指定的开始偏移量, 从0开始                                                                                                            |
|   -S, --state | string | 作业状态, 可选值: Initiated, InitiallySuspended, Pending, Running, Suspending, Suspended, Terminating, Terminated, Completed, Failed |
| -U, --user_id | string | 用户ID, 例如4TiSsZonTa3                                                                                                              |
|    -Z, --zone | string | 指定区域, 例如az-jinan, az-zhigu。可以通过ysadmin zone list查看所有的区域信息                                                        |

## 提交作业

提交作业, 指定一个.json文件作为参数文件, 可使用ysadmin job example生成参考参数文件

使用方法: 

ysadmin job submit [flags]

示例:

- 提交作业, 参数文件为job.json
  - ysadmin job submit -F job.json
- 提交作业, 参数文件为job.json, command文件为script.sh
  - ysadmin job submit -F job.json --sh script.sh

可用Flags:

|   命令参数 |  类型  | 说明                                                  |
| ---------: | :----: | :---------------------------------------------------- |
| -F, --file | string | JSON 文件路径\(必填\)                                 |
|       --sh | string | command文件路径,若指定则会覆盖json文件中的command字段 |

## 终止作业

终止作业

使用方法:

ysadmin job terminate [flags]

示例:

- 终止作业ID为 52hgnbSTGWd 的作业
  - ysadmin job terminate -I 52hgnbSTGWd

可用Flags:

| 命令参数 |  类型  | 说明           |
| -------: | :----: | :------------- |
| -I, --id | string | 作业ID\(必填\) |

## 重传作业结果文件

重传作业结果文件

使用方法:

ysadmin job retransmit [flags]

示例:

- 终止作业ID为 52hgnbSTGWd 的作业
  - ysadmin job retransmit -I 52hgnbSTGWd

可用Flags:

| 命令参数 |  类型  | 说明           |
| -------: | :----: | :------------- |
| -I, --id | string | 作业ID\(必填\) |