package exec

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// Argument 表示一个命令行参数
type Argument interface {
	Value() string
}

// New 创建一个新的命令行程序
func New(bin string, args ...Argument) (*Command, error) {
	max := -1
	for _, arg := range args {
		if p, ok := arg.(pArgument); ok {
			if p.idx > max {
				max = p.idx
			}
		}
	}

	return &Command{bin: bin, args: args, maxi: max + 1}, nil
}

// Command 命令行的包装器实现
type Command struct {
	bin  string
	args []Argument
	maxi int
}

var (
	// ErrInsufficientBindings 传入的绑定参数数量不足
	ErrInsufficientBindings = errors.New("exec: insufficient bindings")
)

// Run 运行一个命令行进程
func (cmd *Command) Run(ctx context.Context, bindings ...interface{}) (*Runner, error) {
	return cmd.RunWithOptions(ctx, bindings)
}

// RunWithOptions 使用额外的参数启动命令行
func (cmd *Command) RunWithOptions(ctx context.Context, bindings []interface{}, opts ...RunOption) (*Runner, error) {
	if len(bindings) < cmd.maxi {
		return nil, ErrInsufficientBindings
	}

	args := make([]string, 0, len(cmd.args))
	for _, arg := range cmd.args {
		if p, ok := arg.(pArgument); ok {
			args = append(args, fmt.Sprintf("%s", bindings[p.idx]))
		} else {
			args = append(args, arg.Value())
		}
	}

	return run(ctx, cmd.bin, args, opts...)
}

// Runner 是一个命令行进程的包装器
type Runner struct {
	cmd    *exec.Cmd
	bin    string
	args   []string
	cancel context.CancelFunc

	stdin  io.WriteCloser
	stdout io.WriteCloser
	stderr io.WriteCloser
}

// Start 启动进程
func (r *Runner) Start() error {
	logging.Default().Infof("executing command: %s %s", r.bin, r.args)
	return r.cmd.Start()
}

// Wait 等待进程结束
func (r *Runner) Wait() error {
	if r.cmd.Process == nil {
		if err := r.Start(); err != nil {
			return err
		}
	}

	defer func() {
		if r.stdin != nil {
			_ = r.stdin.Close()
		}
		if r.stdout != nil {
			_ = r.stdout.Close()
		}
		if r.stderr != nil {
			_ = r.stderr.Close()
		}
	}()
	return r.cmd.Wait()
}

// Input 输入内容到标准输入中
func (r *Runner) Input(arg string) error {
	_, err := r.stdin.Write([]byte(arg))
	return err
}

// RunOption 创建进程的额外参数
type RunOption func(r *Runner) error

// WithStdout 指定标准输出
func WithStdout(stdout io.Writer) RunOption {
	return func(r *Runner) error {
		r.cmd.Stdout = stdout
		return nil
	}
}

// WithStderr 指定标准错误
func WithStderr(stderr io.Writer) RunOption {
	return func(r *Runner) error {
		r.cmd.Stderr = stderr
		return nil
	}
}

// run 创建并启动一个进程包装器
func run(ctx context.Context, bin string, args []string, opts ...RunOption) (*Runner, error) {
	r := &Runner{bin: bin, args: args}
	ctx, r.cancel = context.WithCancel(ctx)
	r.cmd = exec.CommandContext(ctx, bin, args...)
	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}

	var err error
	if r.stdin, err = r.cmd.StdinPipe(); err != nil {
		return nil, err
	}

	return r, nil
}

// sArgument 常规的字符串参数
type sArgument string

// Value 返回参数自身
func (s sArgument) Value() string {
	return string(s)
}

// String 字符串参数的构造方法
func String(s string) Argument {
	return sArgument(s)
}

// pArgument 命令行参数占位符
type pArgument struct {
	idx int
}

// Value 返回占位符对应的参数值
func (p pArgument) Value() string {
	return ""
}

// PlaceHolder 创建一个占位命令参数
func PlaceHolder(i int) Argument {
	return pArgument{idx: i}
}
