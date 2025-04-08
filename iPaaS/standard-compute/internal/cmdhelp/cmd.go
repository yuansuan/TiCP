package cmdhelp

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
)

type CmdConfig struct {
	env      []string
	dir      string
	username string
	flags    map[string]string
}

type Option interface {
	apply(c *CmdConfig)
}

type optionFunc func(c *CmdConfig)

func (f optionFunc) apply(c *CmdConfig) {
	f(c)
}

func WithCmdEnv(env []string) Option {
	return optionFunc(func(c *CmdConfig) {
		c.env = env
	})
}

func WithCmdDir(dir string) Option {
	return optionFunc(func(c *CmdConfig) {
		c.dir = dir
	})
}

func WithCmdUser(username string) Option {
	return optionFunc(func(c *CmdConfig) {
		c.username = username
	})
}

func ExecShellCmd(ctx context.Context, cmdStr string, opts ...Option) (stdOut string, stdErr string, err error) {
	logger := log.GetJobTraceLogger(ctx)

	cmdConfig := &CmdConfig{}
	for _, opt := range opts {
		opt.apply(cmdConfig)
	}

	var cmd *exec.Cmd
	if cmdConfig.username != "" {
		// 使用 cmd.SysProcAttr.Credential 会带来商业软件的报错，例如starCCM会有报错仍是文件权限的问题，暂时无法解决，故此处暂时使用sudo -u
		// 所以若私有云场景，标准计算需要启动于有sudo权限的用户，并且可以任意切换至提交作业用户
		cmd = exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("sudo -u %s %s", cmdConfig.username, cmdStr))

		logger.With("exec-user", cmdConfig.username)
		logger.Infof("assigned user: %s", cmdConfig.username)
	} else {
		cmd = exec.CommandContext(ctx, "bash", "-c", cmdStr)
		logger.Info("not assigned user")
	}

	cmd.Env = append(os.Environ(), cmdConfig.env...)
	cmd.Dir = cmdConfig.dir

	cmdStrForLog := os.Expand(cmd.String(), func(key string) string {
		for _, env := range cmd.Env {
			kv := strings.SplitN(env, "=", 2)
			if len(kv) == 2 && kv[0] == key {
				return kv[1]
			}
		}

		return ""
	})
	logger.Infof("command: %s", cmdStrForLog)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return stdout.String(), stderr.String(), err
}

// parseFlags 入参类似是 -e a.log -o b.log -l select=xxx
// 只包含了flags部分
func parseFlags(s string) map[string]string {
	// 使用 strings.Fields 将输入字符串拆分成切片
	args := strings.Fields(s)

	// 创建一个 map 用于存储参数
	argMap := make(map[string]string)

	// 使用循环来解析命令行参数
	for i := 0; i < len(args); i++ {
		arg := args[i]

		// 如果参数以 "-" 开头，表示是一个标志参数
		if strings.HasPrefix(arg, "-") {
			// 如果下一个参数不以 "-" 开头，则将其作为参数值，否则设置为空字符串
			var value string
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				value = args[i+1]
				i++ // 跳过参数值
			} else {
				value = ""
			}

			// 存储参数和值到 map 中
			argMap[arg] = value
		}
	}

	return argMap
}

func EnsureSubmitCommand(rawCmd string, additionalFlags map[string]string, submitEndWith string) (string, error) {
	if !strings.HasSuffix(rawCmd, submitEndWith) {
		return "", fmt.Errorf("invalid raw submit command [%s], not end with [%s]", rawCmd, submitEndWith)
	}

	fields := strings.Fields(strings.TrimSuffix(rawCmd, submitEndWith))
	if len(fields) <= 1 {
		return "", fmt.Errorf("invalid raw submit command, %s", rawCmd)
	}

	// 保留下执行命令头部
	binName := fields[0]

	// 解析配置文件的flags到map中
	rawFlagsMap := parseFlags(strings.Join(fields[1:], " "))

	// 接口中的覆盖配置文件中的
	for k, v := range additionalFlags {
		rawFlagsMap[k] = v
	}

	// 还原提交作业命令
	flags := ""
	for k, v := range rawFlagsMap {
		flags = fmt.Sprintf("%s %s %s", flags, k, v)
	}

	return fmt.Sprintf("%s %s %s", binName, flags, submitEndWith), nil
}
