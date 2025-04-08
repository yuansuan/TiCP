# 商品服务，配置商品/特价，查看订单等

商品服务，配置商品/特价，查看订单等

使用方法:

ysadmin merchandise

可用子命令:

- add   添加商品
- delete   删除商品
- get   查询单个商品
- list   批量查询商品
- order   查询订单等
- patch   修改商品
- publish   上架商品
- specialprice   添加/修改/删除/查询特价
- unpublish   下架商品

## 添加商品

添加商品

使用方法:

ysadmin merchandise add -F req.json [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   request json file |

## 删除商品

删除商品

使用方法:

ysadmin merchandise delete -I <merchandise_id> [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   merchandise id |

## 查询单个商品

查询单个商品

使用方法:

ysadmin merchandise get -I <merchandise_id> [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   merchandise id |

## 批量查询商品

批量查询商品

使用方法:

ysadmin merchandise list [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| --charge-type |  string |   \[PrePaid \| PostPaid\] |
| --offset |  int |   page offset |
| --out-resource-id |  string |   out resource id \[appId \| softwareId \| hardwareId\] |
| --publish-state |  string |   \[Up \| Down\] |
| --size |  int |   page size \(default 1000\) |
| --ysproduct |  string |   \[CloudCompute \| CloudApp\] |

## 查询订单等

查询订单等

使用方法:

ysadmin merchandise order

可用子命令:

- list   批量查询订单

### 批量查询订单

批量查询订单

使用方法:

ysadmin merchandise order list [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| --account-id |  string |   account id |
| --charge-type |  string |   \[PrePaid \| PostPaid\] |
| --merchandise-id |  string |   merchandise id |
| --offset |  int |   page offset |
| --size |  int |   page size \(default 1000\) |

## 修改商品

修改商品

使用方法:

ysadmin merchandise patch -I <merchandise_id> -F req.json [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   request json file |
| -I, --id |  string |   merchandise id |

## 上架商品

上架商品

使用方法:

ysadmin merchandise publish -I <merchandise_id> [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   merchandise id |

## 添加/修改/删除/查询特价

添加/修改/删除/查询特价

使用方法:

ysadmin merchandise specialprice

可用子命令:

- add   添加特价
- delete   删除特价
- list   批量查询特价
- put   修改特价

### 添加特价

添加特价

使用方法:

ysadmin merchandise specialprice add --merchandise-id <merchandise_id> --account-id <account_id> <unit-price> [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| --account-id |  string |   account id |
| --merchandise-id |  string |   merchandise id |

### 删除特价

删除特价

使用方法:

ysadmin merchandise specialprice delete --merchandise-id <merchandise_id> --account-id <account_id> [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| --account-id |  string |   account id |
| --merchandise-id |  string |   merchandise id |

### 批量查询特价

批量查询特价

使用方法:

ysadmin merchandise specialprice list [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| --account-id |  string |   account id |
| --merchandise-id |  string |   merchandise id |
| --offset |  int |   page offset |
| --size |  int |   page size \(default 1000\) |

### 修改特价

修改特价

使用方法:

ysadmin merchandise specialprice put --merchandise-id <merchandise_id> --account-id <account_id> <unit-price> [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| --account-id |  string |   account id |
| --merchandise-id |  string |   merchandise id |

## 下架商品

下架商品

使用方法:

ysadmin merchandise unpublish -I <merchandise_id> [flags]

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -I, --id |  string |   merchandise id |

