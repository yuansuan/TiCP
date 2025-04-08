# 📖 生成某个命令的文档

📖 生成某个命令的文档, 通过解析cobra的结构体, 生成markdown格式的文档, 包含描述、用法、示例、子命令、flag等信息, 支持递归生成子命令的文档, 支持输出到文件

使用方法:

ysadmin doc [flags]

示例:

ysadmin doc job
ysadmin doc job submit

可用子命令:

- sub   用于生成子命令文档测试的命令

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -F, --file |  string |   指定输出文件，默认输出到标准输出 |
| -h, --help |  |   help for doc |
| -S, --sub |  |   是否递归生成子命令的文档 |

## 用于生成子命令文档测试的命令

用于生成子命令文档测试的命令

使用方法:

ysadmin doc sub [flags]

示例:

ysadmin doc sub

可用子命令:

- sub   用于生成子命令的子命令文档测试的命令

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -a, --aaa |  string |   用于测试的flagA |
| --bbb |  int |   用于测试的flagB |
| --ccc |  |   用于测试的flagC |

### 用于生成子命令的子命令文档测试的命令

用于生成子命令的子命令文档测试的命令

使用方法:

ysadmin doc sub sub [flags]

示例:

ysadmin doc sub

可用Flags:

| 命令参数 | 类型 | 说明 |
| ---: | :---: | :--- |
| -a, --aaa |  string |   用于测试的flagA |
| --bbb |  int |   用于测试的flagB |
| --ccc |  |   用于测试的flagC |

