# 计算应用管理工具

计算应用管理工具, 用于管理计算应用

使用方法:

ysadmin application

可用子命令:

- add   添加计算应用
- delete   删除计算应用
- example   创建计算应用示例文件
- get   获取计算应用
- list   列出计算应用
- publish   计算应用发布/取消发布
- put   更新计算应用
- quota   计算应用配额管理工具

## 添加计算应用

添加计算应用, 指定一个.json文件作为参数文件, 可使用'ysadmin application example'生成参考参数文件

使用方法:

ysadmin application add -F file.json [flags]

示例:

- 添加应用, 参数文件为app.json
  - ysadmin application add -F app.json
- 添加应用, 参数文件为app.json, command文件为script.sh
  - ysadmin application add -F app.json --sh script.sh

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   JSON 文件路径 |
| --sh |  string |   脚本文件路径,若指定则会覆盖json文件中的command字段 |

## 删除计算应用

删除计算应用

使用方法:

ysadmin application delete -I id [flags]

示例:

- 删除应用ID为 52XkfrGM9vE 的应用
  - ysadmin application delete -I 52XkfrGM9vE

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   应用 ID \(必填\) |

## 创建计算应用示例文件

创建计算应用示例文件

使用方法:

ysadmin application example

示例:

- 创建应用示例文件
  - ysadmin application example

## 获取计算应用

获取计算应用

使用方法:

ysadmin application get -I id [flags]

示例:

- 获取应用ID为 52XkfrGM9vE 的应用
  - ysadmin application get -I 52XkfrGM9vE
- 获取应用ID为 52XkfrGM9vE 的应用, 并输出到当前目录下的 52XkfrGM9vE.json 文件中
  - ysadmin application get -I 52XkfrGM9vE -O

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   应用 ID \(必填\) |
| -O, --output |  |   输出到当前目录下的 \[appID\].json 文件中, 用于put命令使用 |

## 列出计算应用

列出计算应用, 可以指定用户 ID 查询该用户有配额的应用列表

使用方法:

ysadmin application list [flags]

示例:

- 列出所有应用
  - ysadmin application list
- 列出用户ID为 52XmLrdCkHb 的用户有配额的应用列表
  - ysadmin application list -U 52XmLrdCkHb

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -U, --user_id |  string |   用户 ID, 查询该用户有配额的应用列表 |

## 计算应用发布/取消发布

计算应用发布/取消发布

使用方法:

ysadmin application publish -I id [-D] [flags]

示例:

- 发布应用ID为 52XkfrGM9vE 的应用
  - ysadmin application publish -I 52XkfrGM9vE
- 取消发布应用ID为 52XkfrGM9vE 的应用
  - ysadmin application publish -I 52XkfrGM9vE -D

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -D, --down |  |   取消发布 |
| -I, --id |  string |   应用 ID \(必填\) |

## 更新计算应用

更新计算应用, 是对应用的全量更新, 未指定的字段会被置空, 指定一个.json文件作为参数文件, 可通过'ysadmin application get -O'获取当前应用的参数文件

使用方法:

ysadmin application put -F file.json -I id [flags]

示例:

- 更新应用ID为 52XkfrGM9vE 的应用, 参数文件为app.json
  - ysadmin application put -F app.json -I 52XkfrGM9vE
- 更新应用ID为 52XkfrGM9vE 的应用, 参数文件为app.json, command文件为script.sh
  - ysadmin application put -F app.json -I 52XkfrGM9vE --sh script.sh

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   JSON 文件路径 |
| -I, --id |  string |   应用 ID \(必填\) |
| --sh |  string |   脚本文件路径,若指定则会覆盖json文件中的command字段 |

## 计算应用配额管理工具

计算应用配额管理工具, 用于管理计算应用配额, 配额是指用户对应用的使用权限

使用方法:

ysadmin application quota

可用子命令:

- add   添加计算应用配额
- delete   删除计算应用配额
- get   获取计算应用配额

### 添加计算应用配额

添加计算应用配额

使用方法:

ysadmin application quota add -I app_id -U user_id [flags]

示例:

- 为用户ID为 52XmLrdCkHb 的用户添加应用ID为 52XkfrGM9vE 的应用配额
  - ysadmin application quota add -I 52XkfrGM9vE -U 52XmLrdCkHb

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   应用 ID \(必填\) |
| -U, --user_id |  string |   用户 ID \(必填\) |

### 删除计算应用配额

删除计算应用配额

使用方法:

ysadmin application quota delete -I app_id -U user_id [flags]

示例:

- 删除用户ID为 52XmLrdCkHb 的用户的应用ID为 52XkfrGM9vE 的应用配额
  - ysadmin application quota delete -I 52XkfrGM9vE -U 52XmLrdCkHb

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   应用 ID \(必填\) |
| -U, --user_id |  string |   用户 ID \(必填\) |

### 获取计算应用配额

获取计算应用配额

使用方法:

ysadmin application quota get -I app_id -U user_id [flags]

示例:

- 查看用户ID为 52XmLrdCkHb 的用户是否有应用ID为 52XkfrGM9vE 的应用配额
  - ysadmin application quota get -I 52XkfrGM9vE -U 52XmLrdCkHb

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   应用 ID \(必填\) |
| -U, --user_id |  string |   用户 ID \(必填\) |

