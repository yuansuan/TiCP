# CAD商业格式转换管理

管理CAD商业格式转换作业及配额

使用方法:

ysadmin CADConverter

可用子命令:

- list   查询CAD转换作业列表
- create   创建作业
- cancel   取消CAD转换作业
- delete   删除CAD转换作业
- quota-list   所有用户的CAD转换配额列表
- quota-update   更新指定用户的CAD转换配额

## 查询CAD转换作业列表

查询CAD转换作业列表

使用方法:

ysadmin CADConverter list [flags]

示例:

- 列出所有CAD转换作业(默认1000条)
    - ysadmin CADConverter list
- 查询所有运行中的CAD转换作业, 偏移量为0, 限制条数为1000
    - ysadmin CADConverter list -O 0 -S 1000 -s Running
- 查询CAD转换作业ID为 52hgnbSTGWd 的作业信息
    - ysadmin CADConverter list -J 52hgnbSTGWd

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -O, --pageOffset |  int |  分页偏移量 (默认0) |
| -S, --pageSize |  int |  分页大小 (默认1000) |
| -s, --state |  string |  作业状态 |
| -J, --jobIDs |  stringSlice |  作业ID列表，多个用逗号分隔 |

## 创建作业

创建CAD转换作业

使用方法:

ysadmin CADConverter create [flags]

示例:

- 创建CAD转换作业
    - ysadmin CADConverter create -I /test/aaa -T .stp
- 创建CAD转换作业，作业完成6000秒后删除作业
    - ysadmin CADConverter create -I /test/aaa -T .stp -A 6000

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --input |  string |  输入路径 (必填) |
| -T, --target |  string |  目标转换格式 (必填) |
| -A, --autoDelete |  int |  作业完成x秒自动删除，0不删除 (默认0) |

## 取消CAD转换作业

取消指定作业ID的CAD转换作业

使用方法:

ysadmin CADConverter cancel [flags]

示例:

- 取消作业ID为 52hgnbSTGWd 的CAD转换作业
    - ysadmin CADConverter cancel -J 52hgnbSTGWd

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -J, --jobID |  string |  作业ID (必填) |

## 删除CAD转换作业

删除指定作业ID的CAD转换作业

使用方法:

ysadmin CADConverter delete [flags]

示例:

- 删除作业ID为 52hgnbSTGWd 的CAD转换作业
    - ysadmin CADConverter delete -J 52hgnbSTGWd

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -J, --jobID |  string |  作业ID (必填) |

## 所有用户的CAD转换配额列表

查看所有用户的CAD转换配额信息

使用方法:

ysadmin CADConverter quota-list

示例:

- 查看所有用户的CAD转换配额列表
    - ysadmin CADConverter quota-list

## 更新指定用户的CAD转换配额

修改用户的CAD转换配额数量

使用方法:

ysadmin CADConverter quota-update [flags]

示例:

- 更新用户 4TiSsZonTa3 的CAD转换配额
    - ysadmin CADConverter quota-update -U 4TiSsZonTa3 -N 1000

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -N, --number |  int |  配额数量 (必填) |
| -U, --userID |  string |  用户ID (必填) |