# 存储管理

存储管理, 可以操作管理云存储和HPC存储
存储接口中操作的远程存储路径均需要以/{userID}开头
例如: /4TpFFZDkFWy/test.txt

使用方法:

ysadmin storage

可用子命令:

- batch-download   批量下载文件
- batch-download-speed   批量下载文件(测速)
- download   下载文件
- download-speed   下载文件(测速)
- list-operation-log   列出操作日志
- list-quota   列出配额
- ls   列出路径下文件
- mkdir   指定路径创建文件夹
- mv   移动文件
- quota-total   配额总量
- readat   读取文件
- rm   删除文件
- update-quota   更新配额
- upload   上传文件
- upload-speed   上传文件(测速)

## 批量下载文件

批量下载文件, batch-download [远程存储目标路径] [本地路径], 目标路径必须以/{userID}开头, 本地路径必须是最终文件路径且需要.zip后缀
例如: /{userID}/path/dir /tmp/filename.zip, 不支持下载单个文件

使用方法:

ysadmin storage batch-download "/{userID}/path" "/localPath" [flags]

示例:

- 批量下载某用户根目录下的testdir目录, 指定存储类型为hpc, 指定区域为az-zhigu, 本地路径为/tmp/testdir.zip, 最后压缩包仅包含testdir目录下的所有文件(-B 不传默认与传"/4TpFFZDkFWy/testdir"相同, 即压缩包中不包含远程路径的一层目录)
  - ysadmin storage batch-download /4TpFFZDkFWy/testdir /tmp/testdir.zip -T hpc -Z az-zhigu [-B /4TpFFZDkFWy/testdir]
- 批量下载某用户根目录下的testdir目录, 指定存储类型为hpc, 指定区域为az-zhigu, 本地路径为/tmp/testdir.zip, 最后压缩包会包含一层testdir目录
  - ysadmin storage batch-download /4TpFFZDkFWy/testdir /tmp/testdir.zip -T hpc -Z az-zhigu -B /4TpFFZDkFWy

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -B, --base-path |  string |   压缩包的起始路径（不包含）, 不传默认与远程路径相同, 即压缩包中不包含远程路径的一层目录 |
| -T, --type |  string |   存储类型, hpc \| cloud \(default "cloud"\) |
| -Z, --zone |  string |   区域 \(default "az-zhigu"\) |

## 批量下载文件(测速)

批量下载文件(测速), 不指定path默认是下载配置文件对应user的某个测速目录下的特定目录, 如果指定了path, 则下载指定的path下的文件。(只测速, 不会实际下载, Size大于文件大小会报错)

使用方法:

ysadmin storage batch-download-speed ["/{userID}/path"] [flags]

示例:

- 从配置用户目录下下载文件(测速), 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage batch-download-speed -T hpc -Z az-zhigu
- 从某用户根目录下的testdir目录下下载文件(测速), 指定存储类型为cloud, 指定区域为az-zhigu
  - ysadmin storage batch-download-speed /4TpFFZDkFWy/testdir -T cloud -Z az-zhigu

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -S, --size |  int |   文件大小 单位为MB \(default 1000\) |
| -T, --type |  string |   存储类型, hpc \| cloud \(default "cloud"\) |
| -Z, --zone |  string |   区域 \(default "az-zhigu"\) |

## 下载文件

下载文件, download [远程存储目标路径] [本地路径], 目标路径必须以/{userID}开头, 本地路径必须是最终文件路径(包含文件名)
例如: /{userID}/path/filename.txt /tmp/filename.txt, 不支持下载目录

使用方法:

ysadmin storage download "/{userID}/path" "/localPath" [flags]

示例:

- 下载某用户根目录下的文件, 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage download /4TpFFZDkFWy/test.txt /tmp/test.txt -T hpc -Z az-zhigu

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -T, --type |  string |   存储类型, hpc \| cloud \(default "cloud"\) |
| -Z, --zone |  string |   区域 \(default "az-zhigu"\) |

## 下载文件(测速)

下载文件(测速), 不指定path默认是下载配置文件对应user的某个测速目录下的特定文件, 如果指定了path, 则下载指定的path下的文件。(只测速, 不会实际下载, Size大于文件大小会报错)

使用方法:

ysadmin storage download-speed ["/{userID}/path"] [flags]

示例:

- 从配置用户目录下下载文件(测速), 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage download-speed -T hpc -Z az-zhigu
- 从某用户根目录下的testdir目录下下载文件, 指定存储类型为cloud, 指定区域为az-zhigu
  - ysadmin storage download-speed /4TpFFZDkFWy/testdir/testfile -T cloud -Z az-zhigu
- 从某用户根目录下的testdir目录下下载文件, 指定存储类型为cloud, 指定区域为az-zhigu, 指定读取文件大小为200MB
  - ysadmin storage download-speed /4TpFFZDkFWy/testdir/testfile -T cloud -Z az-zhigu -S 200

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -S, --size |  int |   文件大小 单位为MB \(default 1000\) |
| -T, --type |  string |   存储类型, hpc \| cloud \(default "cloud"\) |
| -Z, --zone |  string |   区域 \(default "az-zhigu"\) |

## 列出操作日志

列出操作日志

使用方法:

ysadmin storage list-operation-log [flags]

示例:

- 列出操作日志, 指定存储类型为hpc, 指定区域为az-zhigu, 分页显示, 从第0条开始, 每页显示20条
  - ysadmin storage list-operation-log -T hpc -Z az-zhigu -O 0 -S 20
- 列出操作日志, 指定存储类型为hpc, 指定区域为az-zhigu, 指定文件名为test.txt, 指定操作类型为UPLOAD, 指定文件类型为FILE, 指定开始时间为2021-01-01 00:00:00, 指定结束时间为2021-01-02 00:00:00
  - ysadmin storage list-operation-log -T hpc -Z az-zhigu --file-name "test.txt" --operation-type UPLOAD --file-type FILE -b "2021-01-01 00:00:00" -e "2021-01-02 00:00:00"

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -b, --begin-time |  string |   开始时间, 格式: YYYY-MM-DD HH:mm:ss |
| -e, --end-time |  string |   结束时间, 格式: YYYY-MM-DD HH:mm:ss |
| -f, --file-name |  string |   文件名 |
| -t, --file-type |  string |   文件类型, 可选值: FILE-文件, DIRECTORY-目录 |
| -L, --limit |  int |   limit \(default 1000\) |
| -O, --offset |  int |   offset |
| -o, --operation-type |  string |   操作类型, 可选值: UPLOAD-上传, DOWNLOAD-下载, DELETE-删除, MOVE-移动, MKDIR-添加文件夹, COPY-拷贝, COPY\_RANGE-指定范围拷贝,COMPRESS-压缩, CREATE-创建, LINK-链接, READ\_AT-读, WRITE\_AT-写 |
| -T, --type |  string |   存储类型, hpc \| cloud \(default "cloud"\) |
| -U, --user_id |  string |   用户ID, 不填默认使用配置文件的storage\_ys\_id |
| -Z, --zone |  string |   区域 \(default "az-zhigu"\) |

## 列出配额

列出配额, 列出分区存储下各目录(用户)配额

使用方法:

ysadmin storage list-quota [flags]

示例:

- 列出分区存储下各目录(用户)配额, 指定存储类型为hpc, 指定区域为az-zhigu, 分页显示, 从第0条开始, 每页显示20条
  - ysadmin storage list-quota -T hpc -Z az-zhigu -O 0 -S 20

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -O, --offset |  int |   offset |
| -S, --size |  int |   size \(default 1000\) |
| -T, --type |  string |   存储类型, hpc \| cloud \(default "cloud"\) |
| -Z, --zone |  string |   区域 \(default "az-zhigu"\) |

## 列出路径下文件

列出路径下文件, 路径必须以/{userID}开头

使用方法:

ysadmin storage ls "/{userID}/path" [flags]

示例:

- 列出某用户根目录下的文件
  - ysadmin storage ls /4TpFFZDkFWy
- 列出某用户根目录下的文件, 过滤掉文件名包含test的文件, 指定存储类型为hpc, 指定区域为az-zhigu, 分页显示, 从第0条开始, 每页显示20条
  - ysadmin storage ls /4TpFFZDkFWy -F test -O 0 -L 20 -T hpc -Z az-zhigu

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --filter |  string |   用于过滤的正则表达式 |
| -L, --limit |  int |   limit \(default 1000\) |
| -O, --offset |  int |   offset |
| -T, --type |  string |   存储类型, hpc \| cloud \(default "cloud"\) |
| -Z, --zone |  string |   区域 \(default "az-zhigu"\) |

## 指定路径创建文件夹

指定路径创建文件夹, 路径必须以/{userID}开头

使用方法:

ysadmin storage mkdir "/{userID}/path" [flags]

示例:

- 创建某用户根目录下的文件夹, 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage mkdir /4TpFFZDkFWy/test -T hpc -Z az-zhigu

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -T, --type |  string |   存储类型, hpc \| cloud \(default "cloud"\) |
| -Z, --zone |  string |   区域 \(default "az-zhigu"\) |

## 移动文件

移动文件, src和dest路径必须以/{userID}开头, 不存在的上级目录会被创建

使用方法:

ysadmin storage mv "/{userID}/srcPath" "/{userID}/destPath" [flags]

示例:

- 移动某用户根目录下的文件, 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage mv /4TpFFZDkFWy/test.txt /4TpFFZDkFWy/test2.txt -T hpc -Z az-zhigu
- 移动某用户根目录下的文件夹, 指定存储类型为cloud, 指定区域为az-zhigu
  - ysadmin storage mv /4TpFFZDkFWy/testdir /4TpFFZDkFWy/testdir2 -T cloud -Z az-zhigu

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -T, --type |  string |   存储类型, hpc \| cloud \(default "cloud"\) |
| -Z, --zone |  string |   区域 \(default "az-zhigu"\) |

## 配额总量

配额总量, 列出分区存储下配额总量

使用方法:

ysadmin storage quota-total [flags]

示例:

- 列出分区存储下配额总量, 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage quota-total -T hpc -Z az-zhigu

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -T, --type |  string |   存储类型, hpc \| cloud \(default "cloud"\) |
| -Z, --zone |  string |   区域 \(default "az-zhigu"\) |

## 读取文件

读取文件, 路径必须以/{userID}开头

使用方法:

ysadmin storage readat "/{userID}/path" [flags]

示例:

- 读取某用户根目录下的文件, 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage readat /4TpFFZDkFWy/test.txt -T hpc -Z az-zhigu
- 读取某用户根目录下的文件, 指定存储类型为cloud, 指定区域为az-zhigu, 从第200字节开始读取, 读取1000字节
  - ysadmin storage readat /4TpFFZDkFWy/test.txt -T cloud -Z az-zhigu -O 200 -L 1000

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -L, --limit |  int |   limit \(default 1000\) |
| -O, --offset |  int |   offset |
| -T, --type |  string |   存储类型, hpc \| cloud \(default "cloud"\) |
| -Z, --zone |  string |   区域 \(default "az-zhigu"\) |

## 删除文件

删除文件, 路径必须以/{userID}开头

使用方法:

ysadmin storage rm "/{userID}/path" [flags]

示例:

- 删除某用户根目录下的文件, 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage rm /4TpFFZDkFWy/test.txt -T hpc -Z az-zhigu
- 删除某用户根目录下的文件夹, 指定存储类型为cloud, 指定区域为az-zhigu
  - ysadmin storage rm /4TpFFZDkFWy/testdir -T cloud -Z az-zhigu

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -T, --type |  string |   存储类型, hpc \| cloud \(default "cloud"\) |
| -Z, --zone |  string |   区域 \(default "az-zhigu"\) |

## 更新配额

更新配额, 更新某用户配额

使用方法:

ysadmin storage update-quota [flags]

示例:

- 更新某用户配额, 指定存储类型为hpc, 指定区域为az-zhigu, 配额为1000GB
  - ysadmin storage update-quota -U 4TpFFZDkFWy -T hpc -Z az-zhigu -L 1000

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -L, --limit |  int |   limit \(default 1000\) |
| -T, --type |  string |   存储类型, hpc \| cloud \(default "cloud"\) |
| -U, --user_id |  string |   用户ID, 不填默认使用配置文件的storage\_ys\_id |
| -Z, --zone |  string |   区域 \(default "az-zhigu"\) |

## 上传文件

上传文件, upload [本地路径] [远程存储目标路径], 目标路径必须以/{userID}开头, 需要指定最终的文件名/目录名
例如: /{userID}/path/filename.txt, 支持上传目录

使用方法:

ysadmin storage upload "/localPath" "/{userID}/path" [flags]

示例:

- 上传本地文件到某用户根目录下, 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage upload /tmp/test.txt /4TpFFZDkFWy/test.txt -T hpc -Z az-zhigu
- 上传本地目录到某用户根目录下, 指定存储类型为cloud, 指定区域为az-zhigu
  - ysadmin storage upload /tmp/testdir /4TpFFZDkFWy/testdir -T cloud -Z az-zhigu

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -T, --type |  string |   存储类型, hpc \| cloud \(default "cloud"\) |
| -Z, --zone |  string |   区域 \(default "az-zhigu"\) |

## 上传文件(测速)

上传文件(测速), 不指定path默认是上传到配置文件对应user的某个测速目录下, 如果指定了path, 则上传到指定的path下

使用方法:

ysadmin storage upload-speed ["/{userID}/path"] [flags]

示例:

- 上传文件到配置用户目录下(测速), 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage upload-speed -T hpc -Z az-zhigu
- 上传文件到某用户根目录下的testdir目录下(测速), 指定存储类型为cloud, 指定区域为az-zhigu
  - ysadmin storage upload-speed /4TpFFZDkFWy/testdir -T cloud -Z az-zhigu
- 上传文件到某用户根目录下的testdir目录下(测速), 指定存储类型为cloud, 指定区域为az-zhigu, 指定文件大小为2000MB
  - ysadmin storage upload-speed /4TpFFZDkFWy/testdir -T cloud -Z az-zhigu -S 2000

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -S, --size |  int |   文件大小 单位为MB \(default 1000\) |
| -T, --type |  string |   存储类型, hpc \| cloud \(default "cloud"\) |
| -Z, --zone |  string |   区域 \(default "az-zhigu"\) |

