package sessionaction

import (
	"encoding/base64"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/utils"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/utils/template"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

func Mount(c *gin.Context, opts ...Option) {
	//conf := new(config)
	//for _, opt := range opts {
	//	opt.apply(conf)
	//}
	//
	//logger := conf.logger
	//if logger == nil {
	//	logger = trace.GetLogger(c).Base()
	//}
	//
	//req := new(session.MountRequest)
	//err := bindMountRequest(c, req)
	//if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
	//	logger.Warnf("bind mount request failed, %v", err)
	//	return
	//}
	//
	//err, errResp := validator.ValidateMountRequest(req)
	//if err = response.BadRequestIfError(c, err, errResp); err != nil {
	//	logger.Warnf("validate mount request failed, %v", err)
	//	return
	//}
	//
	//// get session from db
	//sessionDetail, exist, err := dao.GetSessionDetailsBySessionID(c, conf.userId, snowflake.MustParseString(*req.SessionId))
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, fmt.Sprintf("get session [%s] from database failed", *req.SessionId))); err != nil {
	//	logger.Warnf("get session [%s] from database failed, %v", *req.SessionId, err)
	//	return
	//}
	//if !exist {
	//	err = fmt.Errorf("session [%s] not found", *req.SessionId)
	//	_ = response.NotFoundIfError(c, err, response.WrapErrorResp(common.SessionNotFound, fmt.Sprintf("session [%s] not found", *req.SessionId)))
	//	logger.Warn(err)
	//	return
	//}
	//logger = logger.With("session-id", sessionDetail.Session.Id.String())
	//
	//// validate mount point
	//if err = validateMountPoint(c, sessionDetail.Platform, *req.MountPoint); err != nil {
	//	logger.Warnf("invalid mount point, %v", err)
	//	return
	//}
	//
	//// check if ready
	//if !util.IsSessionReadyFromSignalServer(sessionDetail.RoomId, logger) {
	//	err = fmt.Errorf("session [%s] not ready yet", sessionDetail.Session.Id)
	//	_ = response.ForbiddenIfError(c, err, response.WrapErrorResp(common.ForbiddenSessionNotReady, err.Error()))
	//	logger.Warn(err)
	//	return
	//}
	//
	//// call zone list api to get cloud storage endpoint
	//s, err := util.GetState(c)
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "get state failed")); err != nil {
	//	logger.Errorf("get state from gin ctx failed, %v", err)
	//	return
	//}
	//
	//if req.ShareDirectory == nil {
	//	req.ShareDirectory = utils.PString("")
	//}
	//
	//// call create share directory api in cloud storage
	//shareDirectory, err := util.CreateShareDirectory(s, util.CreateShareDirectoryArgs{
	//	Zone:         sessionDetail.Session.Zone.String(),
	//	AssumeUserId: sessionDetail.Session.UserId,
	//	SubPath:      *req.ShareDirectory,
	//})
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "create share directory failed")); err != nil {
	//	logger.Errorf("create share directory failed, %v", err)
	//	return
	//}
	//
	//execMountScriptRequest, err := createExecMountScriptRequest(sessionDetail.Platform, shareDirectory, *req.MountPoint)
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "create exec mount script request failed")); err != nil {
	//	logger.Errorf("create exec mount script request failed, %v", err)
	//	return
	//}
	//
	//// call exec script api to mount share directory
	//execScriptResp, err := util.ExecScript(logger, s, sessionDetail, execMountScriptRequest, trace.GetRequestId(c))
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "call session to mount share directory failed")); err != nil {
	//	logger.Errorf("call session to mount share directory failed, %v", err)
	//	return
	//}
	//
	//// check scriptResp
	//if execScriptResp.Data.ExitCode != 0 {
	//	err = fmt.Errorf("exec mount script occurred error, exitCode is %d", execScriptResp.Data.ExitCode)
	//	_ = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, err.Error()))
	//	logger.Error(err)
	//	return
	//}
	//
	//// update db to mark the new mount point
	//sp := new(cloud.ScriptParams)
	//err = jsoniter.UnmarshalFromString(sessionDetail.Instance.UserParams, sp)
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "unmarshal to scriptParams failed")); err != nil {
	//	logger.Errorf("unmarshal [%s] to scriptParams failed, %v", sessionDetail.Instance.UserParams, err)
	//	return
	//}
	//
	//sp.UpdateByMount(shareDirectory.UserName, shareDirectory.Password, shareDirectory.SharedSrc, *req.MountPoint)
	//
	//userParamsStr, err := jsoniter.MarshalToString(sp)
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "marshal user params to string failed")); err != nil {
	//	logger.Errorf("marshal user params to string failed, %v", err)
	//	return
	//}
	//
	//err = dao.UpdateInstance(c, &models.Instance{
	//	Id:         sessionDetail.Instance.Id,
	//	UserParams: userParamsStr,
	//}, []string{"user_params"})
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "update instance failed")); err != nil {
	//	logger.Errorf("update instance failed, %v", err)
	//	return
	//}

	response.RenderJson(nil, c)
}

func bindMountRequest(c *gin.Context, req *session.MountRequest) error {
	var err error
	if err = c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("should bind uri failed, %w", err)
	}

	if err = c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("should bind json failed, %w", err)
	}

	return nil
}

func validateMountPoint(c *gin.Context, platform models.Platform, mountPoint string) error {
	var err error
	switch platform {
	//case models.Windows:
	//	if !util.StringInSlice(mountPoint, util.WindowsMountPathPermit) {
	//		err = fmt.Errorf("invalid mount point %s, should be in %v", mountPoint, util.WindowsMountPathPermit)
	//		_ = response.BadRequestIfError(c, err, response.WrapErrorResp(common.InvalidArgumentMountPoint, err.Error()))
	//		return err
	//	}
	case models.Linux:
		if !util.IsAbsPath(mountPoint) {
			err = fmt.Errorf("invalid mount point %s, should be absolute", mountPoint)
			_ = response.BadRequestIfError(c, err, response.WrapErrorResp(common.InvalidArgumentMountPoint, err.Error()))
			return err
		}

	default:
		err = fmt.Errorf("unknown platform %s", platform)
		_ = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, err.Error()))
		return err
	}

	return nil
}

type mountScriptTplArgs struct {
	NewLine    string
	MountPoint string
	MountSrc   string
	MountHost  string
	MountUser  string
	MountPass  string
}

// add to host first
// do mount
const windowsMountScriptTemplate = `
$hostsPath = "C:\Windows\System32\drivers\etc\hosts"
$agentEnvPath = "C:\Windows\ys\agent_env"

function AddHost {
	param(
		[string]$shareHostIp,
		[string]$shareServerHost
	)
	try {
		$newHostsEntry = "{{.NewLine}}$shareHostIp    $shareServerHost"
		Write-Host "trying to add a host entry $newHostsEntry"
		$hostsContent = GetQuota-Content -Path $hostsPath -Raw
		if ($hostsContent -notmatch $newHostsEntry){
			Add-Content -Path $hostsPath -Value $newHostsEntry
		}
	}
	catch {
		Write-Error "Error Add Host: $newHostsEntry, $_"
		exit 1
	}
}

function MountShare {
	param(
		[string]$mountSrc,
		[string]$mountPoint,
		[string]$mountHost,
		[string]$mountUser,
		[string]$mountPass
	)
	$shareServerHost = $mountSrc

	AddHost -shareHostIp $mountHost -shareServerHost $shareServerHost

	$command = "net use '$mountPoint' '\\$shareServerHost\$mountSrc' /USER:'$mountUser' '$mountPass' 2>&1"

	try {
		Write-Host "trying to mount share directory, cmd: $command"
		$output = Invoke-Expression $command
		if ($LASTEXITCODE -eq 0) {
			Write-Host "mount from $mountSrc to $mountPoint success!"
		} else {
			Write-Error "mount from $mountSrc to $mountPoint failed!, error output in below"
			Write-Error "$output"
			exit 2
		}
	}
	catch {
		Write-Error "Error mounting share: $_"
		exit 3
	}
}

function RenameVolume {
	param(
		[string]$volumeName,
		[string]$mountPoint
	)

	if ($volumeName -eq "") {
		Write-Host "skip to rename volume"
		return
	}

	$app = New-Object -ComObject shell.application
	$app.NameSpace("$mountPoint").self.name = "$volumeName"
}

function appendToLine {
	param(
		[string]$line,
		[string]$originValue,
		[string]$appendValue
	)

	if ($originValue -ne "") {
		$line += ","
	}
	$line += $appendValue

	$line
}

function SaveRecordToAgentEnv {
	param(
		[string]$mountUser,
		[string]$mountPass,
		[string]$mountSrc,
		[string]$mountPoint
	)

	$agentEnvContent = GetQuota-Content -Path $agentEnvPath

	for ($i = 0; $i -lt $agentEnvContent.Length; $i++) {
		if ($agentEnvContent[$i] -match '^(.*?)=(.*)$') {
			$key = $matches[1]
			$value = $matches[2]

			if ($key -eq "SHARE_USERNAME") {
				$agentEnvContent[$i] = appendToLine -line $agentEnvContent[$i] -originValue $value -appendValue $mountUser
			}

			if ($key -eq "SHARE_PASSWORD") {
				$agentEnvContent[$i] = appendToLine -line $agentEnvContent[$i] -originValue $value -appendValue $mountPass
			}

			if ($key -eq "SHARE_MOUNT_PATHS") {
				$agentEnvContent[$i] = appendToLine -line $agentEnvContent[$i] -originValue $value -appendValue "$mountSrc=$mountPoint"
			}
		}
	}

	Set-Content -Path $agentEnvPath -Value $agentEnvContent
}

MountShare -mountPoint {{.MountPoint}} -mountSrc {{.MountSrc}} -mountHost {{.MountHost}} -mountUser {{.MountUser}} -mountPass {{.MountPass}}
RenameVolume -volumeName {{.MountSrc}} -mountPoint {{.MountPoint}}
SaveRecordToAgentEnv -mountUser {{.MountUser}} -mountPass {{.MountPass}} -mountSrc {{.MountSrc}} -mountPoint {{.MountPoint}}
`

func createExecMountScriptRequest(platform models.Platform, sharedDirectory *schema.SharedDirectory, mountPoint string) (*session.ExecScriptRequest, error) {
	req := &session.ExecScriptRequest{
		WaitTillEnd: utils.PBool(true),
	}

	scriptTpl := ""
	switch platform {
	case models.Windows:
		scriptTpl = windowsMountScriptTemplate
		req.ScriptRunner = utils.PString("powershell")
	case models.Linux:
		return nil, fmt.Errorf("not support linux session yet")
	default:
		return nil, fmt.Errorf("unknown platform [%s]", platform)
	}

	res, err := template.Render(scriptTpl, &mountScriptTplArgs{
		NewLine:    "`n",
		MountPoint: mountPoint,
		MountSrc:   sharedDirectory.SharedSrc,
		MountHost:  sharedDirectory.SharedHost,
		MountUser:  sharedDirectory.UserName,
		MountPass:  sharedDirectory.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("render umount script template failed, %w", err)
	}
	req.ScriptContent = utils.PString(base64.StdEncoding.EncodeToString([]byte(res)))

	return req, nil
}
