# 🎄 环境管理

🎄 环境管理, 用于显示和切换API的环境, 环境以一组config配置文件的形式存在, 默认命令会显示有哪些环境以及当前环境(仅名称), 通过子命令可以切换环境

使用方法:

ysadmin envs

示例:

- 显示所有环境
  - ysadmin envs

可用子命令:

- show   显示环境信息
- switch   切换环境

## 显示环境信息

显示环境信息, 默认显示当前环境信息。 如果指定了具体环境, 则显示指定环境的信息

使用方法:

ysadmin envs show [environment]

示例:

- 显示当前环境信息
  - ysadmin envs show
- 显示指定环境信息
  - ysadmin envs show test

## 切换环境

切换环境, 通过指定环境名称, 切换到对应环境的配置

使用方法:

ysadmin envs switch [environment]

示例:

- 切换到test环境
  - ysadmin envs switch test

