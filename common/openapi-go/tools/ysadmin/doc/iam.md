# 管理IAM相关资源

管理IAM相关资源

使用方法:

ysadmin iam

可用子命令:

- add   添加IAM相关资源
- delete   删除IAM相关资源
- get   获取IAM相关资源
- list   列出IAM相关资源
- update   更新IAM相关资源

## 添加IAM相关资源

添加IAM相关资源

使用方法:

ysadmin iam add

可用子命令:

- policy   为用户添加IAM策略
- role   为用户添加IAM角色
- rolepolicyrelation   为用户添加IAM角色策略关联
- secret   为用户添加IAM密钥
- user   添加IAM用户

### 为用户添加IAM策略

为用户添加IAM策略

使用方法:

ysadmin iam add policy [flags]

示例:

- 为用户添加IAM策略, 参数文件为param.json
  - ysadmin iam add policy -F param.json
  - param.json内容如下:
    \{
        "PolicyName": "4T4\_VIPBoxPolicy123",
        "RoleName": "4T4\_VIPBoxRole123",
        "Version": "1.0",
        "Effect": "allow",
        "Resources": \[
            "\*"
        \],
        "Actions": \[
            "\*"
        \],
        "UserId": "4TiSxuPtJEm"
    \}


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |

### 为用户添加IAM角色

为用户添加IAM角色

使用方法:

ysadmin iam add role [flags]

示例:

- 为用户添加IAM角色, 参数文件为param.json
  - ysadmin iam add role -F param.json
  - param.json内容如下:
    \{
        "Description": "world peace",
        "RoleName": "4T4\_VIPBoxRole123",
        "TrustPolicy": \{
            "Actions": null,
            "Effect": "allow",
            "Principals": \[
                "4T4ZZvA2tVc"
            \],
            "Resources": null
        \},
        "UserId": "4TiSxuPtJEm"
    \}


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |

### 为用户添加IAM角色策略关联

为用户添加IAM角色策略关联

使用方法:

ysadmin iam add rolepolicyrelation [flags]

示例:

- 为用户添加IAM角色策略关联, 参数文件为param.json
  - ysadmin iam add rolepolicyrelation -F param.json
  - param.json内容如下:
    \{
        "PolicyName": "4T4\_VIPBoxPolicy123",
        "RoleName": "4T4\_VIPBoxRole123",
        "UserId": "4TiSxuPtJEm"
    \}


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |

### 为用户添加IAM密钥

为用户添加IAM密钥

使用方法:

ysadmin iam add secret [flags]

示例:

- 为用户添加IAM密钥
  - ysadmin iam add secret -I 4TiSxuPtJEm -T YS\_admin


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   用户id \(必填\) |
| -T, --tag |  string |   AK tag, YS\_ 开头的为远算云产品账号，不能随意使用 \(必填\) |

### 添加IAM用户

添加IAM用户

使用方法:

ysadmin iam add user [flags]

示例:

- 添加IAM用户, 参数文件为param.json
  - ysadmin iam add user -F param.json
  - param.json内容如下:
    \{
        "Phone": "13800138000",
        "Password": "MyPassword@1234"
    \}


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |

## 删除IAM相关资源

删除IAM相关资源

使用方法:

ysadmin iam delete

可用子命令:

- policy   删除某用户IAM策略
- role   删除某用户IAM角色
- rolepolicyrelation   删除某用户IAM角色策略关联
- secret   删除IAM密钥

### 删除某用户IAM策略

删除某用户IAM策略

使用方法:

ysadmin iam delete policy [flags]

示例:

- 删除某用户IAM策略, 参数文件为param.json
  - ysadmin iam delete policy -F param.json
  - param.json内容如下:
    \{
        "PolicyName": "4T4\_VIPBoxPolicy123",
        "UserId": "4TiSxuPtJEm"
    \}


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |

### 删除某用户IAM角色

删除某用户IAM角色

使用方法:

ysadmin iam delete role [flags]

示例:

- 删除某用户IAM角色, 参数文件为param.json
  - ysadmin iam delete role -F param.json
  - param.json内容如下:
    \{
        "RoleName": "4T4\_VIPBoxRole123",
        "UserId": "4TiSxuPtJEm"
    \}


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |

### 删除某用户IAM角色策略关联

删除某用户IAM角色策略关联

使用方法:

ysadmin iam delete rolepolicyrelation [flags]

示例:

- 删除某用户IAM角色策略关联, 参数文件为param.json
  - ysadmin iam delete rolepolicyrelation -F param.json
  - param.json内容如下:
    \{
        "PolicyName": "4T4\_VIPBoxPolicy123",
        "RoleName": "4T4\_VIPBoxRole123",
        "UserId": "4TiSxuPtJEm"
    \}


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |

### 删除IAM密钥

删除IAM密钥

使用方法:

ysadmin iam delete secret [flags]

示例:

- 删除IAM密钥
  - ysadmin iam delete secret -I 6I02NW57S0IJADN08OXK


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   AK id \(必填\) |

## 获取IAM相关资源

获取IAM相关资源

使用方法:

ysadmin iam get

可用子命令:

- policy   获取某用户某个具体IAM策略信息
- role   获取某用户某个具体IAM角色信息
- secret   获取IAM密钥
- user   获取IAM用户

### 获取某用户某个具体IAM策略信息

获取某用户某个具体IAM策略信息

使用方法:

ysadmin iam get policy [flags]

示例:

- 获取某用户某个具体IAM策略信息, 参数文件为param.json
  - ysadmin iam get policy -F param.json
  - param.json内容如下:
    \{
        "PolicyName": "4T4\_VIPBoxPolicy123",
        "UserId": "4TiSxuPtJEm"
    \}


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |

### 获取某用户某个具体IAM角色信息

获取某用户某个具体IAM角色信息

使用方法:

ysadmin iam get role [flags]

示例:

- 获取某用户某个具体IAM角色信息, 参数文件为param.json
  - ysadmin iam get role -F param.json
  - param.json内容如下:
    \{
        "RoleName": "4T4\_VIPBoxRole123",
        "UserId": "4TiSxuPtJEm"
    \}


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |

### 获取IAM密钥

获取IAM密钥

使用方法:

ysadmin iam get secret [flags]

示例:

- 获取IAM密钥
  - ysadmin iam get secret -I 6I02NW57S0IJADN08OXK


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   AK id \(必填\) |

### 获取IAM用户

获取IAM用户

使用方法:

ysadmin iam get user [flags]

示例:

- 获取IAM用户
  - ysadmin iam get user -I 4TiSxuPtJEm


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   用户id \(必填\) |

## 列出IAM相关资源

列出IAM相关资源

使用方法:

ysadmin iam list

可用子命令:

- policy   列出某用户IAM策略
- role   列出某用户IAM角色列表
- secret   列出某用户IAM密钥
- secrets   列出所有IAM密钥
- users   列出所有IAM用户

### 列出某用户IAM策略

列出某用户IAM策略

使用方法:

ysadmin iam list policy [flags]

示例:

- 列出某用户IAM策略
  - ysadmin iam list policy -I 4TiSxuPtJEm


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   用户id \(必填\) |

### 列出某用户IAM角色列表

列出某用户IAM角色列表,

使用方法:

ysadmin iam list role [flags]

示例:

- 列出某用户IAM角色列表
  - ysadmin iam list role -I 4TiSxuPtJEm


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   用户id \(必填\) |

### 列出某用户IAM密钥

列出某用户IAM密钥

使用方法:

ysadmin iam list secret [flags]

示例:

- 列出某用户IAM密钥
  - ysadmin iam list secret -I 4TiSxuPtJEm


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   用户id \(必填\) |

### 列出所有IAM密钥

列出所有IAM密钥

使用方法:

ysadmin iam list secrets [flags]

示例:

- 列出所有IAM密钥
  - ysadmin iam list secrets
- 列出所有IAM密钥, 带分页参数
  - ysadmin iam list secrets -O 0 -L 10


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -L, --limit |  int |   limit \(default 1000\) |
| -O, --offset |  int |   offset |

### 列出所有IAM用户

列出所有IAM用户

使用方法:

ysadmin iam list users [flags]

示例:

- 列出所有IAM用户
  - ysadmin iam list users
- 列出所有IAM用户, 带分页参数
  - ysadmin iam list users -O 0 -L 10
- 列出所有IAM用户, 带用户id
  - ysadmin iam list users -I 4TiSxuPtJEm


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   用户id |
| -L, --limit |  int |   limit \(default 1000\) |
| -O, --offset |  int |   offset |

## 更新IAM相关资源

更新IAM相关资源

使用方法:

ysadmin iam update

可用子命令:

- policy   更新某用户IAM策略
- role   更新某用户IAM角色
- user   更新IAM用户

### 更新某用户IAM策略

更新某用户IAM策略

使用方法:

ysadmin iam update policy [flags]

示例:

- 更新某用户IAM策略, 参数文件为param.json
  - ysadmin iam update policy -F param.json
  - param.json内容如下:
    \{
        "PolicyName": "4T4\_VIPBoxPolicy123",
        "RoleName": "4T4\_VIPBoxRole123",
        "Version": "1.0",
        "Effect": "allow",
        "Resources": \[
            "\*"
        \],
        "Actions": \[
            "\*"
        \],
        "UserId": "4TiSxuPtJEm"
    \}


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |

### 更新某用户IAM角色

更新某用户IAM角色

使用方法:

ysadmin iam update role [flags]

示例:

- 更新某用户IAM角色, 参数文件为param.json
  - ysadmin iam update role -F param.json
  - param.json内容如下:
    \{
        "Description": "world peace",
        "RoleName": "4T4\_VIPBoxRole123",
        "TrustPolicy": \{
            "Actions": null,
            "Effect": "allow",
            "Principals": \[
                "4T4ZZvA2tVc"
            \],
            "Resources": null
        \},
        "UserId": "4TiSxuPtJEm"
    \}


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |

### 更新IAM用户

更新IAM用户

使用方法:

ysadmin iam update user [flags]

示例:

- 更新IAM用户, 参数文件为param.json
  - ysadmin iam update user -F param.json
  - param.json内容如下:
    \{
        "UserId": "4TiSxxxxxx",
        "Name": "test",
    \}


可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   json文件 \(必填\) |

