# 资金账号管理

资金账号管理, 不同于远算账号, 用于管理用户的账户余额, 账单等

使用方法:

ysadmin account

可用子命令:

- addbalance   充值
- create   创建资金账号
- example   创建资金账号示例文件
- get   获取账户信息
- getbyysid   根据远算账号ID获取账户信息
- listbill   获取账单列表
- reduce   扣费
- refund   退款

## 充值

充值, 单位为0.00001元

使用方法:

ysadmin account addbalance [flags]

示例:

- 充值, 指定账户ID, 充值金额, 交易ID, 备注
  - ysadmin account addbalance -I 5314rXEJwrf -M 10000000 -T 531jv4i44nJ -C "充值100元"

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -C, --comment |  string |   备注 \(必填\) |
| -I, --id |  string |   账户ID \(必填\) |
| -M, --money |  int |   充值金额, 单位为0.00001元 \(必填\) |
| -T, --trade_id |  string |   交易ID \(必填\) |

## 创建资金账号

创建资金账号, 指定一个.json文件作为参数文件, 可使用ysadmin account example生成参考参数文件

使用方法:

ysadmin account create [flags]

示例:

 - 创建资金账号, 参数文件为account.json
  - ysadmin account create -F account.json

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   JSON 文件路径 \(必填\) |

## 创建资金账号示例文件

创建资金账号示例文件

使用方法:

ysadmin account example

示例:

- 创建资金账号示例文件
  - ysadmin account example

## 获取账户信息

获取账户信息

使用方法:

ysadmin account get [flags]

示例:

- 获取账户信息, 指定账户ID
  - ysadmin account get -I 5314rXEJwrf

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   账户ID \(必填\) |

## 根据远算账号ID获取账户信息

根据远算账号ID获取账户信息

使用方法:

ysadmin account getbyysid [flags]

示例:

- 根据远算账号ID获取账户信息, 指定远算账号ID
  - ysadmin account getbyysid -I 5314rXEJwrf

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   远算账号ID \(必填\) |

## 获取账单列表

获取账单列表, 可输出到CSV文件, 可聚合消费账单

使用方法:

ysadmin account listbill [flags]

示例:

- 获取账单列表, 所有的条目
  - ysadmin account listbill --all
- 获取账单列表, 分页获取, 每页10条, 第1页
  - ysadmin account listbill -L 10 -O 1
- 获取账单列表, 分页获取, 每页10条, 第1页, 指定开始时间和结束时间
  - ysadmin account listbill -L 10 -O 1 --start_time "2021-01-01 00:00:00" --end_time "2021-01-31 23:59:59"
- 获取账单列表, 分页获取, 每页10条, 第1页, 指定产品名称, 并输出到CSV文件
  - ysadmin account listbill -L 10 -O 1 -P 3D云应用 --csv_file /tmp/bill.csv
- 获取账单列表, 分页获取, 每页10条, 第1页, 指定账户ID, 并输出到CSV文件, 并聚合消费账单
  - ysadmin account listbill -I 5314rXEJwrf -L 10 -O 1 --csv_file /tmp/bill.csv --merge
- 获取账单列表, 分页获取, 每页10条, 第1页, 指定账户ID, 并输出到CSV文件, 并聚合消费账单, 只展示消费类型的账单
  - ysadmin account listbill -I 5314rXEJwrf -L 10 -O 1 --csv_file /tmp/bill.csv --merge --only_consume

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| --all |  |   所有的条目 |
| --csv_file |  string |   输出CSV格式的数据到指定文件 |
| --end_time |  string |   结束时间 |
| -I, --id |  string |   账户ID |
| -L, --limit |  int |   limit \(default 1000\) |
| --merge |  |   是否需要聚合消费账单，按资源ID展示 |
| -O, --offset |  int |   offset |
| --only_consume |  |   只展示消费类型的账单 |
| -P, --product_name |  string |   产品名称 |
| --start_time |  string |   开始时间 |

## 扣费

扣费, 单位为0.00001元

使用方法:

ysadmin account reduce [flags]

示例:

- 扣费, 指定账户ID, 扣费金额, 交易ID, 备注, 产品名称
  - ysadmin account reduce -I 5314rXEJwrf -M 10000000 -T 531jv4i44nJ -C "扣费100元" -P 3D云应用

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -C, --comment |  string |   备注 \(必填\) |
| -I, --id |  string |   账户ID \(必填\) |
| -M, --money |  int |   扣费金额, 单位为0.00001元 \(必填\) |
| -P, --product_name |  string |   产品名称 |
| -T, --trade_id |  string |   交易ID \(必填\) |

## 退款

退款, 单位为0.00001元

使用方法:

ysadmin account refund [flags]

示例:

- 退款, 指定账户ID, 退款金额, 交易ID, 备注
  - ysadmin account refund -I 5314rXEJwrf -M 10000000 -T 531jv4i44nJ -C "退款100元"

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -C, --comment |  string |   备注 \(必填\) |
| -I, --id |  string |   账户ID \(必填\) |
| -M, --money |  int |   退款金额, 单位为0.00001元 \(必填\) |
| -T, --trade_id |  string |   交易ID \(必填\) |

