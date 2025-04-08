# license server配置管理

license server配置管理, 可以配置远算自有license和外部license

使用方法:

ysadmin license

可用子命令:

- add   添加license server配置
- delete   删除license server配置
- example   输出配置文件示例
- get   获取license server配置
- list   列出license server配置
- put   修改license server配置

## 添加license server配置

添加license server配置

使用方法:

ysadmin license add

可用子命令:

- licinfo   添加license info
- licmanager   添加license manager
- moduleconfig   添加module config

### 添加license info

添加license info

使用方法:

ysadmin license add licinfo [flags]

示例:

- 添加license info, 参数文件为lic_info.json
  - ysadmin license add licinfo -F lic_info.json

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |

### 添加license manager

添加license manager

使用方法:

ysadmin license add licmanager [flags]

示例:

- 添加license manager, 参数文件为lic_manager.json
  - ysadmin license add licmanager -F lic_manager.json

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |

### 添加module config

添加module config

使用方法:

ysadmin license add moduleconfig [flags]

示例:

- 添加module config, 参数文件为module_config.json
  - ysadmin license add moduleconfig -F module_config.json

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |

## 删除license server配置

删除license server配置

使用方法:

ysadmin license delete

可用子命令:

- licinfo   删除license info
- licmanager   删除license manager
- moduleconfig   删除module config

### 删除license info

删除license info

使用方法:

ysadmin license delete licinfo [flags]

示例:

- 删除ID为52Zzc3ycfEU的license info
  - ysadmin license delete licinfo -I 52Zzc3ycfEU

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   license info id \(必填\) |

### 删除license manager

删除license manager

使用方法:

ysadmin license delete licmanager [flags]

示例:

- 删除ID为52Zzc3ycfEU的license manager
  - ysadmin license delete licmanager -I 52Zzc3ycfEU

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   license manager id \(必填\) |

### 删除module config

删除module config

使用方法:

ysadmin license delete moduleconfig [flags]

示例:

- 删除ID为52Zzc3ycfEU的module config
  - ysadmin license delete moduleconfig -I 52Zzc3ycfEU

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   module config id \(必填\) |

## 输出配置文件示例

输出配置文件示例

使用方法:

ysadmin license example

可用子命令:

- licinfo   输出license info配置文件示例
- licmanager   输出license manager配置文件示例
- moduleconfig   输出module config配置文件示例

### 输出license info配置文件示例

输出license info配置文件示例

使用方法:

ysadmin license example licinfo

示例:

- 输出license info配置文件示例
  - ysadmin license example licinfo

### 输出license manager配置文件示例

输出license manager配置文件示例

使用方法:

ysadmin license example licmanager

示例:

- 输出license manager配置文件示例
  - ysadmin license example licmanager

### 输出module config配置文件示例

输出module config配置文件示例

使用方法:

ysadmin license example moduleconfig

示例:

- 输出module config配置文件示例
  - ysadmin license example moduleconfig

## 获取license server配置

获取license server配置

使用方法:

ysadmin license get

可用子命令:

- licmanager   获取license manager

### 获取license manager

获取license manager

使用方法:

ysadmin license get licmanager [flags]

示例:

- 获取ID为52Zzc3ycfEU的license manager
  - ysadmin license get licmanager -I 52Zzc3ycfEU

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   license manager id \(必填\) |

## 列出license server配置

列出license server配置

使用方法:

ysadmin license list

可用子命令:

- licmanager   列出license manager

### 列出license manager

列出license manager

使用方法:

ysadmin license list licmanager

示例:

- 列出所有license manager
  - ysadmin license list licmanager

## 修改license server配置

修改license server配置

使用方法:

ysadmin license put

可用子命令:

- licinfo   修改license info
- licmanager   修改license manager
- moduleconfig   修改module config

### 修改license info

修改license info

使用方法:

ysadmin license put licinfo [flags]

示例:

- 修改ID为52Zzc3ycfEU的license info, 参数文件为lic_info.json
  - ysadmin license put licinfo -I 52Zzc3ycfEU -F lic_info.json

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |
| -I, --id |  string |   license info id \(必填\) |

### 修改license manager

修改license manager

使用方法:

ysadmin license put licmanager [flags]

示例:

- 修改ID为52Zzc3ycfEU的license manager, 参数文件为lic_manager.json
  - ysadmin license put licmanager -I 52Zzc3ycfEU -F lic_manager.json

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |
| -I, --id |  string |   license manager id \(必填\) |

### 修改module config

修改module config

使用方法:

ysadmin license put moduleconfig [flags]

示例:

- 修改ID为52Zzc3ycfEU的module config, 参数文件为module_config.json
  - ysadmin license put moduleconfig -I 52Zzc3ycfEU -F module_config.json

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |
| -I, --id |  string |   module config id \(必填\) |

