package osutil

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"

	"github.com/pkg/errors"
)

// CmdContext ...
type CmdContext struct {
	context.Context

	uid    uint32
	gid    uint32
	env    map[string]string
	chroot string
}

// WithUser ...
func (ctx *CmdContext) WithUser(uid, gid uint32) {
	ctx.uid = uid
	ctx.gid = gid
}

// WithChroot ...
func (ctx *CmdContext) WithChroot(chroot string) {
	ctx.chroot = chroot
}

// WithEnv ...
func (ctx *CmdContext) WithEnv(key, value string) {
	if ctx.env == nil {
		ctx.env = map[string]string{}
	}
	ctx.env[key] = value
}

// Env ...
func (ctx *CmdContext) Env() []string {
	ret := []string{}

	for _, e := range os.Environ() {
		ret = append(ret, e)
	}

	for k, v := range ctx.env {
		ret = append(ret, fmt.Sprintf("%v=%v", k, v))
	}

	return ret
}

// CmdHelper ...
type CmdHelper struct{}

// CommandHelper ...
var CommandHelper = CmdHelper{}

// NewCtx ...
func (h *CmdHelper) NewCtx(ctx context.Context, username string) (*CmdContext, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return nil, errors.Wrap(err, "cmd_helper new ctx:")
	}

	uid, _ := strconv.Atoi(u.Uid)
	gid, _ := strconv.Atoi(u.Gid)

	return &CmdContext{
		Context: ctx,
		uid:     uint32(uid),
		gid:     uint32(gid),
	}, nil
}

// NewCtxWithUID ...
func (h *CmdHelper) NewCtxWithUID(ctx context.Context, uid int64) (*CmdContext, error) {
	u, err := user.LookupId(fmt.Sprintf("%v", uid))
	if err != nil {
		return nil, errors.Wrap(err, "cmd_helper new ctx:")
	}

	gid, _ := strconv.Atoi(u.Gid)

	return &CmdContext{
		Context: ctx,
		uid:     uint32(uid),
		gid:     uint32(gid),
	}, nil
}

// NewCtxWithOperator ...
func (h *CmdHelper) NewCtxWithOperator(ctx context.Context, uid, gid int64) *CmdContext {
	return &CmdContext{
		Context: ctx,
		uid:     uint32(uid),
		gid:     uint32(gid),
	}
}

// NewCtxWithCurrent ...
func (h *CmdHelper) NewCtxWithCurrent(ctx context.Context) *CmdContext {
	return &CmdContext{
		Context: ctx,
		uid:     uint32(os.Getuid()),
		gid:     uint32(os.Getgid()),
	}
}

// Execf ...
func (h *CmdHelper) Execf(ctx *CmdContext, command string, args ...string) (stdout []byte, stderr []byte, err error) {
	cmd := exec.CommandContext(ctx, command, args...)
	var b, d bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &d

	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: ctx.uid, Gid: ctx.gid}

	cmd.Env = ctx.Env()
	cmd.SysProcAttr.Chroot = ctx.chroot

	err = cmd.Run()
	return b.Bytes(), d.Bytes(), errors.Wrap(err, "cmd_helper execf:")
}

// WithHomeEnv ...
func (h *CmdHelper) WithHomeEnv(ctx *CmdContext) (*CmdContext, error) {
	u, err := user.LookupId(fmt.Sprintf("%v", ctx.uid))
	if err != nil {
		return nil, err
	}

	ctx.WithEnv("HOME", u.HomeDir)
	ctx.WithEnv("USER", u.Username)
	return ctx, nil
}

// Bash ...
func (h *CmdHelper) Bash(ctx *CmdContext, cmd string) (string, string, error) {
	ctx, err := h.WithHomeEnv(ctx)
	if err != nil {
		return "", "", err
	}

	ctx.WithEnv("BASH_ENV", "~/.bashrc")
	stdout, stderr, err := h.Execf(ctx, "bash", "-c", cmd)
	return string(stdout), string(stderr), err
}

// Bashf ...
func (h *CmdHelper) Bashf(ctx *CmdContext, format string, a ...interface{}) (string, string, error) {
	return h.Bash(ctx, fmt.Sprintf(format, a...))
}

// BashWithCurrent ...
func (h *CmdHelper) BashWithCurrent(ctx context.Context, cmd string) (string, string, error) {
	cmdCtx := h.NewCtxWithCurrent(ctx)

	cmdCtx, err := h.WithHomeEnv(cmdCtx)
	if err != nil {
		return "", "", err
	}

	cmdCtx.WithEnv("BASH_ENV", "~/.bashrc")

	stdout, stderr, err := h.Execf(cmdCtx, "bash", "-c", cmd)
	return string(stdout), string(stderr), err
}

// BashfWithCurrent ...
func (h *CmdHelper) BashfWithCurrent(ctx context.Context, format string, a ...interface{}) (string, string, error) {
	return h.BashWithCurrent(ctx, fmt.Sprintf(format, a...))
}

// ------------------------- deprecated -----------------------

// Exec Exec
// Deprecated method
func Exec(command string, args ...string) (stdout []byte, stderr []byte, err error) {
	cmd := exec.Command(command, args...)
	var b, d bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &d

	err = cmd.Run()
	return b.Bytes(), d.Bytes(), err
}

// Bash bash
// Deprecated method
func Bash(command string) (stdout []byte, stderr []byte, err error) {
	return Exec("bash", "-c", command)
}

// Sh Sh
// Deprecated method
func Sh(command string) (stdout []byte, stderr []byte, err error) {
	return Exec("sh", "-c", command)
}
