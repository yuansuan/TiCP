package impl

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	RegisterCmd(NewEnvCommand())
}

// NewEnvCommand 创建环境管理命令
func NewEnvCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "envs",
		Short: "🎄 环境管理",
		Long:  "🎄 环境管理, 用于显示和切换API的环境, 环境以一组config配置文件的形式存在, 默认命令会显示有哪些环境以及当前环境(仅名称), 通过子命令可以切换环境",
		Args:  cobra.NoArgs,
		Example: `- 显示所有环境
  - ysadmin envs`,
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		// 列出所有环境（仅名称）
		current := Cfg.CurrentEnvironment
		keys := make([]string, 0, len(Cfg.Environments))
		for k := range Cfg.Environments {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		fmt.Println("Environments:")
		for _, env := range keys {
			if env == current {
				fmt.Printf(" * %s  [%s]\n", color.New(color.FgGreen).Add(color.Bold).SprintFunc()(env), Cfg.Environments[env].Endpoint)
			} else {
				fmt.Printf("   %s  [%s]\n", env, Cfg.Environments[env].Endpoint)
			}
		}
	}

	cmd.AddCommand(
		newSwitchEnvCmd(),
		newShowEnvCmd(),
	)
	return cmd
}

func newSwitchEnvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch [environment]",
		Short: "切换环境",
		Long:  "切换环境, 通过指定环境名称, 切换到对应环境的配置",
		Args:  cobra.ExactArgs(1),
		Example: `- 切换到test环境
  - ysadmin envs switch test`,
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		// 切换到指定环境
		newEnv := args[0]
		if err := SwitchEnv(newEnv); err != nil {
			fmt.Printf("切换环境失败: %v\n", err)
			return
		}
	}

	return cmd
}

func newShowEnvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [environment]",
		Short: "显示环境信息",
		Long:  "显示环境信息, 默认显示当前环境信息。 如果指定了具体环境, 则显示指定环境的信息",
		Args:  cobra.MaximumNArgs(1),
		Example: `- 显示当前环境信息
  - ysadmin envs show
- 显示指定环境信息
  - ysadmin envs show test`,
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		// 列出指定环境的详细信息
		env := Cfg.CurrentEnvironment
		if len(args) > 0 {
			env = args[0]
		}

		if _, exist := Cfg.Environments[env]; !exist {
			fmt.Println("环境不存在: ", env)
			return
		}

		fmt.Printf("Environment: %s\n", color.New(color.FgGreen).Add(color.Bold).SprintFunc()(env))
		printCompute(Cfg.Environments[env])
		printStorage(Cfg.Environments[env])
		printIam(Cfg.Environments[env])

	}

	return cmd
}

// SwitchEnv switches the current environment
func SwitchEnv(env string) error {
	if _, exist := Cfg.Environments[env]; !exist {
		return fmt.Errorf("环境不存在: %s", env)
	}

	Cfg.CurrentEnvironment = env
	fmt.Println("Switched to environment:", color.New(color.FgGreen).Add(color.Bold).SprintFunc()(env))

	if err := SaveConfig(); err != nil {
		return err
	}

	return nil
}

func printCompute(cfg EnvironmentConfig) {
	fmt.Printf(" * %s\n", color.New(color.FgBlue).Add(color.Bold).SprintFunc()("Compute:"))
	fmt.Printf("   Endpoint: %s\n", cfg.Endpoint)
	fmt.Printf("   YsID: %s\n", cfg.ComputeYsID)
	fmt.Printf("   AccessKey: %s\n", cfg.ComputeAccessKeyID)
	fmt.Printf("   SecretKey: %s\n", cfg.ComputeAccessKeySecret)
}

func printStorage(cfg EnvironmentConfig) {
	fmt.Printf(" * %s\n", color.New(color.FgBlue).Add(color.Bold).SprintFunc()("Storage:"))
	fmt.Printf("   YsID: %s\n", cfg.StorageYsID)
	fmt.Printf("   AccessKey: %s\n", cfg.StorageAccessKeyID)
	fmt.Printf("   SecretKey: %s\n", cfg.StorageAccessKeySecret)
}

func printIam(cfg EnvironmentConfig) {
	fmt.Printf(" * %s\n", color.New(color.FgBlue).Add(color.Bold).SprintFunc()("IAM:"))
	fmt.Printf("   Endpoint: %s\n", cfg.IamAdminEndpoint)
	fmt.Printf("   AdminAccessKey: %s\n", cfg.IamAdminAccessKeyID)
	fmt.Printf("   AdminSecretKey: %s\n", cfg.IamAdminAccessKeySecret)
}
