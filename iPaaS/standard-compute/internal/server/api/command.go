package api

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/command"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/cmdhelp"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/response"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/errorcode"
)

// PostCommand white list
func PostCommand(c *gin.Context) {
	logger := trace.GetLogger(c)

	req, err := bindPostCommandReq(c)
	if err = response.BadRequestIfError(c, err, errorcode.InvalidArgument); err != nil { // TODO 整理错误码，将错误吗整改至project-root-api库中
		logger.Errorf("bind post command request failed, %v", err)
		return
	}

	if !validateCommand(req.Command) {
		err = fmt.Errorf("validate command failed, %w", err)
		_ = response.BadRequestIfError(c, err, errorcode.InvalidArgument)
		logger.Error(err)
		return
	}

	timeout, err := time.ParseDuration(fmt.Sprintf("%ds", req.Timeout))
	if err != nil {
		err = fmt.Errorf("invalid timeout %d", req.Timeout)
		_ = response.BadRequestIfError(c, err, errorcode.InvalidArgument)
		logger.Error(err)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	wd, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("get workdir failed, %w", err)
		_ = response.InternalErrorIfError(c, err, errorcode.InternalServerError)
		logger.Error(err)
		return
	}

	exitCode := 0
	stdout, stderr, err := cmdhelp.ExecShellCmd(timeoutCtx, req.Command, cmdhelp.WithCmdDir(wd))
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			exitCode = exitErr.ExitCode()
		}
	}

	isTimeout := false
	if timeoutCtx.Err() != nil {
		isTimeout = true
	}

	response.OK(c, &command.SystemPostResponseData{
		Stdout:    stdout,
		Stderr:    stderr,
		ExitCode:  exitCode,
		IsTimeout: isTimeout,
	})
}

func bindPostCommandReq(c *gin.Context) (*command.SystemPostRequest, error) {
	req := new(command.SystemPostRequest)
	if err := c.Bind(req); err != nil {
		return nil, fmt.Errorf("bind failed, %w", err)
	}

	return req, nil
}

// TODO fulfill this validate rules
func validateCommand(command string) bool {
	if command == "" {
		return false
	}

	return true
}
