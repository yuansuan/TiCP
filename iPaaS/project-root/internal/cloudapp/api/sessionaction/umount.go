package sessionaction

import (
	"encoding/base64"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/utils"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/utils/template"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

func Umount(c *gin.Context, opts ...Option) {
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
	//req := new(session.UmountRequest)
	//err := bindUmountRequest(c, req)
	//if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
	//	logger.Warnf("bind umount request failed, %v", err)
	//	return
	//}
	//
	//err, errResp := validator.ValidateUmountRequest(req)
	//if err = response.BadRequestIfError(c, err, errResp); err != nil {
	//	logger.Warnf("validate umount request failed, %v", err)
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
	//// check if mount point exist in database record
	//scriptParams := new(cloud.ScriptParams)
	//err = jsoniter.Unmarshal([]byte(sessionDetail.Instance.UserParams), scriptParams)
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "unmarshal json failed")); err != nil {
	//	logger.Warnf("unmarshal %s to *cloud.ScriptParams failed, %v", sessionDetail.Instance.UserParams, err)
	//	return
	//}
	//
	//if !mountPointExistInDBRecord(scriptParams, *req.MountPoint) {
	//	err = fmt.Errorf("mount point not found")
	//	_ = response.ForbiddenIfError(c, err, response.WrapErrorResp(common.MountPointNotFound, err.Error()))
	//	logger.Warnf("mount point [%s] not exist in database record [%s]", *req.MountPoint, sessionDetail.Instance.UserParams)
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
	//s, err := util.GetState(c)
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "get state failed")); err != nil {
	//	logger.Errorf("get state from gin ctx failed, %v", err)
	//	return
	//}
	//
	//// call exec script api to umount share directory
	//execUmountScriptRequest, err := createExecUmountScriptRequest(sessionDetail.Platform, *req.MountPoint)
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "create exec umount script request failed")); err != nil {
	//	logger.Errorf("create exec umount script request failed, %v", err)
	//	return
	//}
	//
	//execScriptResp, err := util.ExecScript(logger, s, sessionDetail, execUmountScriptRequest, trace.GetRequestId(c))
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "call session to umount share directory failed")); err != nil {
	//	logger.Errorf("call session to umount share directory failed, %v", err)
	//	return
	//}
	//
	//// check scriptResp
	//if execScriptResp.Data.ExitCode != 0 {
	//	err = fmt.Errorf("exec umount script occurred error, exitCode is %d", execScriptResp.Data.ExitCode)
	//	_ = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, err.Error()))
	//	logger.Error(err)
	//	return
	//}
	//
	//// update database to remove mount point
	//sp := new(cloud.ScriptParams)
	//err = jsoniter.UnmarshalFromString(sessionDetail.Instance.UserParams, sp)
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "unmarshal to scriptParams failed")); err != nil {
	//	logger.Errorf("unmarshal [%s] to scriptParams failed, %v", sessionDetail.Instance.UserParams, err)
	//	return
	//}
	//
	//sp.UpdateByUMount(*req.MountPoint)
	//
	//userParamsStr, err := jsoniter.MarshalToString(sp)
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "marshal userParams to string failed")); err != nil {
	//	logger.Errorf("marshal userParams to string failed, %v", err)
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

	// share directory should not be deleted from cloud storage as the share directory may be used by other sessions

	response.RenderJson(nil, c)
}

func mountPointExistInDBRecord(scriptParams *cloud.ScriptParams, mountPoint string) bool {
	return scriptParams.ShareMountPaths.IsMountPointExist(mountPoint)
}

func bindUmountRequest(c *gin.Context, req *session.UmountRequest) error {
	var err error
	if err = c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("should bind uri failed, %w", err)
	}

	if err = c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("should bind json failed, %w", err)
	}

	return nil
}

type umountScriptTplArgs struct {
	MountPoint string
}

const windowsUmountScriptTemplate = `
$agentEnvPath = "C:\Windows\ys\agent_env"
function Umount {
	param(
		[string]$mountPoint
	)

	try {
		net use $mountPoint /delete
	} catch {
		Write-Error "Error umount share: $_"
		exit 1
	}
}

function removeElementInValue {
	param(
		[string]$originValue,
		[int]$removedIndex
	)

	$resArray = @()
	$originArray = $originValue.Split(',')
	for ($i = 0; $i -lt $originArray.Length; $i++) {
		if ($i -ne $removedIndex) {
			$resArray += $originArray[$i]
		}
	}

	$resultValue = $resArray -join ','

	return $resultValue
}

function indexOfMountPoint {
	param(
		[string[]]$shareMountPaths,
		[string]$mountPoint
	)

	$foundIndex = -1
	for ($i = 0; $i -lt $shareMountPaths.Length; $i++) {
		if ($mountPoint -eq $shareMountPaths[$i].Split('=')[1])
		{
			return $i
		}
	}

	return $foundIndex
}


function RemoveRecordInAgentEnv {
	param(
		[string]$mountPoint
	)

	$agentEnvContent = GetQuota-Content -Path $agentEnvPath

	# first loop, find the SHARE_MOUNT_PATHS value, check if the mountPoint record exist
	for ($i = 0; $i -lt $agentEnvContent.Length; $i++) {
		if ($agentEnvContent[$i] -notmatch '^(.*?)=(.*)$') {
			continue
		}

		$key = $matches[1]
		$value = $matches[2]
		if ($key -eq "SHARE_MOUNT_PATHS") {
			$removeIndex = indexOfMountPoint -shareMountPaths $value.Split(',') -mountPoint $mountPoint
			if ($removeIndex -eq -1) {
				Write-Host "$mountPoint not exist in agent_env, no need to remove"
				return
			}
			break
		}
	}

	# second loop, delete the index value which should be removed in SHARE_USERNAME SHARE_PASSWORD SHARE_MOUNT_PATHS
	for ($i = 0; $i -lt $agentEnvContent.Length; $i++) {
		if ($agentEnvContent[$i] -notmatch '^(.*?)=(.*)$') {
			continue
		}

		$key = $matches[1]
		$value = $matches[2]
		if ($key -eq "SHARE_USERNAME" -or $key -eq "SHARE_PASSWORD" -or $key -eq "SHARE_MOUNT_PATHS") {
			$newValue = removeElementInValue -originValue $value -removedIndex $removeIndex
			$agentEnvContent[$i] = "$key=" + $newValue
		}
	}

	Set-Content -Path $agentEnvPath -Value $agentEnvContent
}

Umount -mountPoint {{.MountPoint}}
RemoveRecordInAgentEnv -mountPoint {{.MountPoint}}
`

func createExecUmountScriptRequest(platform models.Platform, mountPoint string) (*session.ExecScriptRequest, error) {
	req := &session.ExecScriptRequest{
		WaitTillEnd: utils.PBool(true),
	}

	scriptTpl := ""
	switch platform {
	case models.Windows:
		scriptTpl = windowsUmountScriptTemplate
		req.ScriptRunner = utils.PString("powershell")
	case models.Linux:
		return nil, fmt.Errorf("not support linux session yet")
	default:
		return nil, fmt.Errorf("unknown platform [%s]", platform)
	}

	res, err := template.Render(scriptTpl, &umountScriptTplArgs{
		MountPoint: mountPoint,
	})
	if err != nil {
		return nil, fmt.Errorf("render umount script template failed, %w", err)
	}
	req.ScriptContent = utils.PString(base64.StdEncoding.EncodeToString([]byte(res)))

	return req, nil
}
