# 3D云应用，管理会话/软件/硬件/RemoteApp等资源

3D云应用，管理会话/软件/硬件/RemoteApp等资源

使用方法:

ysadmin cloudapp

可用子命令:

- hardware   硬件管理，创建/删除/更新/查询等
- remoteapp   remoteapp相关功能
- session   会话管理，创建/删除/查询等
- software   软件管理，创建/删除/更新/查询等

## 硬件管理，创建/删除/更新/查询等

硬件管理，创建/删除/更新/查询等

使用方法:

ysadmin cloudapp hardware

可用子命令:

- admin   Hardware管理员功能，增删改查等
- user   Hardware普通用户功能，查询等

### Hardware管理员功能，增删改查等

Hardware管理员功能，增删改查等

使用方法:

ysadmin cloudapp hardware admin

可用子命令:

- add   创建硬件
- add-users   添加用户
- delete   删除硬件
- delete-users   删除用户
- get   查询单个硬件
- list   批量查询硬件
- modify   修改硬件（增量）

#### 创建硬件

创建硬件

使用方法:

ysadmin cloudapp hardware admin add -F req.json [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   request json file |

#### 添加用户

添加用户
支持单个/批量,用逗号隔开 批量: add-users hardwareId1,hardwareId2 user1,user2
                     单个: add-users hardwareId userId

使用方法:

ysadmin cloudapp hardware admin add-users <hardwares> <users>

#### 删除硬件

删除硬件

使用方法:

ysadmin cloudapp hardware admin delete <hardware-id>

#### 删除用户

删除用户
支持单个/批量,用逗号隔开 批量: delete-users hardwareId1,hardwareId2 user1,user2
                     单个: delete-users hardwareId userId

使用方法:

ysadmin cloudapp hardware admin delete-users <hardwares> <users>

#### 查询单个硬件

查询单个硬件

使用方法:

ysadmin cloudapp hardware admin get <hardware-id>

#### 批量查询硬件

批量查询硬件

使用方法:

ysadmin cloudapp hardware admin list [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| --cpu |  string |   cpu |
| --gpu |  string |   gpu count |
| --mem |  string |   memory \[MB\] |
| --name |  string |   name |
| --offset |  int |   page offset |
| --size |  int |   page size \(default 1000\) |
| --zone |  string |   zone |

#### 修改硬件（增量）

修改硬件（增量）

使用方法:

ysadmin cloudapp hardware admin modify <hardware-id> -F req.json [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   request json file |

### Hardware普通用户功能，查询等

Hardware普通用户功能，查询等

使用方法:

ysadmin cloudapp hardware user

可用子命令:

- get   查询单个硬件
- list   批量查询硬件

#### 查询单个硬件

查询单个硬件

使用方法:

ysadmin cloudapp hardware user get <hardware-id>

#### 批量查询硬件

批量查询硬件

使用方法:

ysadmin cloudapp hardware user list [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| --cpu |  string |   cpu |
| --gpu |  string |   gpu count |
| --mem |  string |   memory \[MB\] |
| --name |  string |   name |
| --offset |  int |   page offset |
| --size |  int |   page size \(default 1000\) |
| --zone |  string |   zone |

## remoteapp相关功能

remoteapp相关功能

使用方法:

ysadmin cloudapp remoteapp

可用子命令:

- admin   管理员相关功能，增删改等
- user   普通用户相关功能，查询等

### 管理员相关功能，增删改等

管理员相关功能，增删改等

使用方法:

ysadmin cloudapp remoteapp admin

可用子命令:

- add   创建远程应用
- delete   创建远程应用
- modify   修改远程应用（增量）

#### 创建远程应用

创建远程应用

使用方法:

ysadmin cloudapp remoteapp admin add -F req.json [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   request json file |

#### 创建远程应用

创建远程应用

使用方法:

ysadmin cloudapp remoteapp admin delete <remoteapp-id>

#### 修改远程应用（增量）

修改远程应用（增量）

使用方法:

ysadmin cloudapp remoteapp admin modify <remoteapp-id> -F req.json [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   request json file |

### 普通用户相关功能，查询等

普通用户相关功能，查询等

使用方法:

ysadmin cloudapp remoteapp user

可用子命令:

- get   查询对应会话的远程应用信息

#### 查询对应会话的远程应用信息

查询对应会话的远程应用信息

使用方法:

ysadmin cloudapp remoteapp user get <session-id> <remoteapp-name>

## 会话管理，创建/删除/查询等

会话管理，创建/删除/查询等

使用方法:

ysadmin cloudapp session

可用子命令:

- admin   Session管理员功能，查询/关闭等
- user   Session普通用户功能，创建/关闭/删除/查询/ready等

### Session管理员功能，查询/关闭等

Session管理员功能，查询/关闭等

使用方法:

ysadmin cloudapp session admin

可用子命令:

- close   关闭会话（删除可视化机器）
- list   批量查询会话
- restart   重启
- restore   重建
- start   开机
- stop   关机

#### 关闭会话（删除可视化机器）

关闭会话（删除可视化机器）

使用方法:

ysadmin cloudapp session admin close <session-id> <reason>

#### 批量查询会话

批量查询会话

使用方法:

ysadmin cloudapp session admin list [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| --offset |  int |   page offset |
| --session-ids |  string |   "ida" or "ida,idb" |
| --size |  int |   page size \(default 1000\) |
| --status |  string |   "STARTED" or "STARTED,STARTING,CLOSING,CLOSED" |
| --user-ids |  string |   "ida" or "ida,idb" |
| --zone |  string |   zone |

#### 重启

重启

使用方法:

ysadmin cloudapp session admin restart <session-id>

#### 重建

以指定用户原会话的启动盘，重建一个新会话

使用方法:

ysadmin cloudapp session admin restore <user-id> <old-session-id>

#### 开机

开机

使用方法:

ysadmin cloudapp session admin start <session-id>

#### 关机

关机

使用方法:

ysadmin cloudapp session admin stop <session-id>

#### 执行脚本

执行脚本

使用方法：

ysadmin cloudapp session admin exec-script <session-id> --script-runner powershell --script-content-encoded <content-base64-encoded> --wait
ysadmin cloudapp session admin exec-script <session-id> --script-runner powershell --script-path /path/to/script/path --wait

### Session普通用户功能，创建/关闭/删除/查询/ready等

Session普通用户功能，创建/关闭/删除/查询/ready等

使用方法:

ysadmin cloudapp session user

可用子命令:

- close   关闭会话（删除可视化机器）
- delete   删除会话
- get   查询单个会话
- list   批量查询会话
- post   创建会话
- ready   检查会话是否ready
- restart   重启
- restore   重建
- start   开机
- stop   关机

#### 关闭会话（删除可视化机器）

关闭会话（删除可视化机器）

使用方法:

ysadmin cloudapp session user close <session-id>

#### 删除会话

删除会话

使用方法:

ysadmin cloudapp session user delete <session-id>

#### 查询单个会话

查询单个会话

使用方法:

ysadmin cloudapp session user get <session-id>

#### 批量查询会话

批量查询会话

使用方法:

ysadmin cloudapp session user list [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| --offset |  int |   page offset |
| --session-ids |  string |   "ida" or "ida,idb" |
| --size |  int |   page size \(default 1000\) |
| --status |  string |   "STARTED" or "STARTED,STARTING,CLOSING,CLOSED" |
| --zone |  string |   zone |

#### 创建会话

创建会话

使用方法:

ysadmin cloudapp session user post -F req.json [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   request json file |

#### 检查会话是否ready

检查会话是否ready

使用方法:

ysadmin cloudapp session user ready <session-id>

#### 重启

重启

使用方法:

ysadmin cloudapp session user restart <session-id>

#### 重建

从原会话的启动盘，重建一个新会话

使用方法:

ysadmin cloudapp session user restore <old-session-id>

#### 开机

开机

使用方法:

ysadmin cloudapp session user start <session-id>

#### 关机

关机

使用方法:

ysadmin cloudapp session user stop <session-id>

#### 执行脚本

执行脚本

使用方法：

ysadmin cloudapp session user exec-script <session-id> --script-runner powershell --script-content-encoded <content-base64-encoded> --wait
ysadmin cloudapp session user exec-script <session-id> --script-runner powershell --script-path /path/to/script/path --wait

#### 挂载

将某个用户存储目录挂载至会话中

使用方法：

ysadmin cloudapp session user mount <session-id> --share-directory <sub_dir> --mount-point <mount_point>

#### 解挂载

将会话中的某个已挂载的用户存储点解挂载

使用方法：

ysadmin cloudapp session user umount <session-id> --mount-point <mount_point>

## 软件管理，创建/删除/更新/查询等

软件管理，创建/删除/更新/查询等

使用方法:

ysadmin cloudapp software

可用子命令:

- admin   Software管理员功能，增删改查等
- user   Software普通用户功能，查询等

### Software管理员功能，增删改查等

Software管理员功能，增删改查等

使用方法:

ysadmin cloudapp software admin

可用子命令:

- add   添加软件
- add-users   添加用户
- delete   删除软件
- delete-users   删除用户
- get   查询单个软件
- list   批量查询软件
- modify   修改软件（增量）
- edit 编辑软件初始化脚本

#### 添加软件

添加软件

使用方法:

ysadmin cloudapp software admin add -F req.json [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   request json file |

#### 添加用户

添加用户
支持单个/批量,用逗号隔开 批量: add-users softwareId1,softwareId2 user1,user2
                     单个: add-users softwareId userId

使用方法:

ysadmin cloudapp software admin add-users <softwares> <users>

#### 删除软件

删除软件

使用方法:

ysadmin cloudapp software admin delete <software-id>

#### 删除用户

删除用户
支持单个/批量,用逗号隔开 批量: delete-users <softwares> <users>

使用方法:

ysadmin cloudapp software admin delete-users <softwares> <users>

#### 查询单个软件

查询单个软件

使用方法:

ysadmin cloudapp software admin get <software-id>

#### 批量查询软件

批量查询软件

使用方法:

ysadmin cloudapp software admin list [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| --name |  string |   name |
| --offset |  int |   page offset |
| --platform |  string |   platform \[ WINDOWS \| LINUX \] |
| --size |  int |   page size \(default 1000\) |
| --zone |  string |   zone |

#### 修改软件（增量）

修改软件（增量）

使用方法:

ysadmin cloudapp software admin modify <software-id> -F req.json [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   request json file |

#### 编辑软件初始化脚本

编辑软件初始化脚本

使用方法：

ysadmin cloudapp software admin edit-init-script <software-id>

### Software普通用户功能，查询等

Software普通用户功能，查询等

使用方法:

ysadmin cloudapp software user

可用子命令:

- get   查询单个软件
- list   批量查询软件

#### 查询单个软件

查询单个软件

使用方法:

ysadmin cloudapp software user get <software-id>

#### 批量查询软件

批量查询软件

使用方法:

ysadmin cloudapp software user list [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| --name |  string |   name |
| --offset |  int |   page offset |
| --platform |  string |   platform \[ WINDOWS \| LINUX \] |
| --size |  int |   page size \(default 1000\) |
| --zone |  string |   zone |

