package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"

	"github.com/yuansuan/ticp/common/project-root-api/rdpgo/v1/execscript"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_agent/pkg/log"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_agent/pkg/script"
)

const (
	cloudAppExecScriptArchiveDirName = "cloudapp-exec-script-archives"
)

func init() {
	registerEndpoint(endPoint{
		Method:       http.MethodPost,
		RelativePath: "/execScript",
		Handler:      execScript,
	})
}

func execScript(c *gin.Context) {
	req := new(execscript.Request)
	var err error
	if err = c.ShouldBindJSON(req); err != nil {
		log.Errorf("bind json failed, %v", err)
		c.JSON(http.StatusBadRequest, &execscript.ResponseData{
			ExitCode: 0,
			Stdout:   "",
			Stderr:   err.Error(),
		})
		return
	}

	scriptRunner, err := script.NewRunner(req.ScriptRunner)
	if err != nil {
		log.Errorf("new script runner failed, %v", err)
		c.JSON(http.StatusBadRequest, &execscript.ResponseData{
			ExitCode: 0,
			Stdout:   "",
			Stderr:   err.Error(),
		})
		return
	}

	scriptContent, err := base64.StdEncoding.DecodeString(req.ScriptContentEncoded)
	if err != nil {
		log.Errorf("decode scriptContent by base64 failed, %v, content: %s", err, req.ScriptContentEncoded)
		c.JSON(http.StatusBadRequest, &execscript.ResponseData{
			ExitCode: 0,
			Stdout:   "",
			Stderr:   err.Error(),
		})
		return
	}
	log.Infof("scriptRunner: %s, scriptContent: %s, WaitTillEnd: %v", req.ScriptRunner, scriptContent, req.WaitTillEnd)

	var archiveDir string

	archiveDir, err = os.UserHomeDir()
	if err != nil {
		archiveDir = os.TempDir()
	}

	if err = os.MkdirAll(filepath.Join(archiveDir, cloudAppExecScriptArchiveDirName), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, &execscript.ResponseData{
			ExitCode: 0,
			Stdout:   "",
			Stderr:   err.Error(),
		})
		log.Errorf("mkdir %s failed, %v", filepath.Join(archiveDir, cloudAppExecScriptArchiveDirName), err)
		return
	}

	requestId := trace.GetRequestId(c)

	scriptFile := filepath.Join(archiveDir, cloudAppExecScriptArchiveDirName, fmt.Sprintf("%s.ps1", requestId))
	if err = os.WriteFile(scriptFile, scriptContent, 0666); err != nil {
		log.Errorf("write script content to file %s failed, %v", scriptFile, err)
		c.JSON(http.StatusInternalServerError, &execscript.ResponseData{
			ExitCode: 0,
			Stdout:   "",
			Stderr:   err.Error(),
		})
		return
	}

	saveExecResultToFile := func(res *script.ExecResult) {
		log.Infof("going to save exec result log to [%s], requestId [%s]", archiveDir, requestId)

		stdoutFile := filepath.Join(archiveDir, cloudAppExecScriptArchiveDirName, fmt.Sprintf("%s.stdout.log", requestId))
		stderrFile := filepath.Join(archiveDir, cloudAppExecScriptArchiveDirName, fmt.Sprintf("%s.stderr.log", requestId))

		if err = os.WriteFile(stdoutFile, []byte(res.Stdout), 0666); err != nil {
			log.Warn("write stdout to %s failed, %v", stdoutFile, err)
		}

		if err = os.WriteFile(stderrFile, []byte(res.Stderr), 0666); err != nil {
			log.Warn("write stderr to %s failed, %v", stderrFile, err)
		}
	}

	log.Infof("going to exec scriptFile: %s", scriptFile)
	execResult, err := scriptRunner.Exec(scriptFile, !req.WaitTillEnd, saveExecResultToFile)
	if err != nil {
		log.Errorf("exec script failed, %v", err)
		if execResult != nil {
			log.Errorf("exitCode: %d", execResult.ExitCode)
			log.Errorf("stdout: %s", execResult.Stdout)
			log.Errorf("stderr: %s", execResult.Stderr)

			c.JSON(http.StatusOK, &execscript.ResponseData{
				ExitCode: execResult.ExitCode,
				Stdout:   execResult.Stdout,
				Stderr:   execResult.Stderr,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, &execscript.ResponseData{
			ExitCode: 0,
			Stdout:   "",
			Stderr:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &execscript.ResponseData{
		ExitCode: execResult.ExitCode,
		Stdout:   execResult.Stdout,
		Stderr:   execResult.Stderr,
	})
}
