package impl

import (
	"encoding/base64"
	"fmt"
	"github.com/coconutLatte/texteditor"
	"github.com/fatih/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/hardware"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/remoteapp"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/software"

	hardwareAdminAdd "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/hardware/admin/add"
	hardwareAdminList "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/hardware/admin/list"
	hardwareAdminPatch "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/hardware/admin/patch"
	hardwareApiList "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/hardware/api/list"
	remoteAppAdminAdd "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/remoteapp/admin/add"
	remoteAppAdminPatch "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/remoteapp/admin/patch"
	sessionApiStart "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/api/start"
	softwareAdminAdd "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/software/admin/add"
	softwareAdminList "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/software/admin/list"
	softwareAdminPatch "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/software/admin/patch"
	softwareUserList "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/software/api/list"
)

type CloudAppOptions struct {
	BaseOptions
	CloseReason string
	AppName     string
}

type CloudAppBaseOptionsWithReqBody struct {
	ReqJson string
}

func InitCloudAppBaseOptions(cmd *cobra.Command) *CloudAppBaseOptionsWithReqBody {
	o := &CloudAppBaseOptionsWithReqBody{}

	cmd.Flags().StringVarP(&o.ReqJson, "file", "F", "", "request json file")
	return o
}

type CloudAppBaseOptionsByList struct {
	PageOffset int
	PageSize   int
}

func InitCloudAppBaseOptionsByList(cmd *cobra.Command) *CloudAppBaseOptionsByList {
	o := &CloudAppBaseOptionsByList{
		PageOffset: 0,
		PageSize:   1000,
	}

	cmd.Flags().IntVar(&o.PageOffset, "offset", 0, "page offset")
	cmd.Flags().IntVar(&o.PageSize, "size", 1000, "page size")
	return o
}

func (o *CloudAppOptions) AddFlags(cmd *cobra.Command) {
	o.BaseOptions.AddBaseOptions(cmd)
	cmd.Flags().StringVarP(&o.CloseReason, "close_reason", "C", "", "关闭3D云应用会话的原因")
	cmd.Flags().StringVarP(&o.AppName, "app_name", "A", "", "RemoteApp名称")
}

func init() {
	RegisterCmd(NewCloudAppCommand())
}

func NewCloudAppCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloudapp",
		Short: "3D云应用，管理会话/软件/硬件/RemoteApp等资源",
		Long:  "3D云应用，管理会话/软件/硬件/RemoteApp等资源",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newSessionCmd(),
		newSoftwareCmd(),
		newHardwareCmd(),
		newRemoteAppCmd(),
	)

	return cmd
}

func newSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "会话管理，创建/删除/查询等",
		Long:  "会话管理，创建/删除/查询等",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newUserSessionCmd(),
		newAdminSessionCmd(),
	)

	return cmd
}

func newUserSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Session普通用户功能，创建/关闭/删除/查询/ready等",
		Long:  "Session普通用户功能，创建/关闭/删除/查询/ready等",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newUserPostSessionCmd(),
		newUserCloseSessionCmd(),
		newUserCheckSessionReadyCmd(),
		newUserDeleteSessionCmd(),
		newUserGetSessionCmd(),
		newUserListSessionCmd(),
		newStartSessionCmd("user"),
		newStopSessionCmd("user"),
		newRestartSessionCmd("user"),
		newUserRestoreSessionCmd(),
		newExecScriptSessionCmd("user"),
		newMountSessionCmd("user"),
		newUmountSessionCmd("user"),
	)

	return cmd
}

func newUserPostSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "post -F req.json",
		Short: "创建会话",
		Long:  "创建会话",
	}

	o := InitCloudAppBaseOptions(cmd)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		data, err := checkReqJsonFileAndRead(o.ReqJson)
		if err != nil {
			return err
		}

		req := new(session.ApiPostRequest)
		if err = jsoniter.Unmarshal(data, req); err != nil {
			return fmt.Errorf("unmarshan json to *session.ApiPostRequest failed, %w", err)
		}

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Session.User.Start(ensureUserStartSessionOpts(req)...)
		})
	}

	return cmd
}

func ensureUserStartSessionOpts(req *session.ApiPostRequest) []sessionApiStart.Option {
	opts := make([]sessionApiStart.Option, 0)
	if req.HardwareId != nil {
		opts = append(opts, GetYsClient().CloudApp.Session.User.Start.HardwareId(*req.HardwareId))
	}
	if req.SoftwareId != nil {
		opts = append(opts, GetYsClient().CloudApp.Session.User.Start.SoftwareId(*req.SoftwareId))
	}
	if req.MountPaths != nil {
		opts = append(opts, GetYsClient().CloudApp.Session.User.Start.MountPaths(*req.MountPaths))
	}
	if req.ChargeParams != nil {
		opts = append(opts, GetYsClient().CloudApp.Session.User.Start.ChargeParams(*req.ChargeParams))
	}
	if req.PayBy != nil {
		opts = append(opts, GetYsClient().CloudApp.Session.User.Start.PayBy(*req.PayBy))
	}
	return opts
}

func newUserCloseSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close <session-id>",
		Short: "关闭会话（删除可视化机器）",
		Long:  "关闭会话（删除可视化机器）",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		sessionId := args[0]
		if sessionId == "" {
			return fmt.Errorf("session-id cannot be empty")
		}

		return callWithDumpReqResp(&session.ApiCloseRequest{
			SessionId: &sessionId,
		}, func() (interface{}, error) {
			return GetYsClient().CloudApp.Session.User.Close(
				GetYsClient().CloudApp.Session.User.Close.Id(sessionId),
			)
		})
	}

	return cmd
}

func newUserCheckSessionReadyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ready <session-id>",
		Short: "检查会话是否ready",
		Long:  "检查会话是否ready",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		sessionId := args[0]
		if sessionId == "" {
			return fmt.Errorf("session-id cannot be empty")
		}

		return callWithDumpReqResp(&session.ApiReadyRequest{
			SessionId: &sessionId,
		}, func() (interface{}, error) {
			return GetYsClient().CloudApp.Session.User.Ready(
				GetYsClient().CloudApp.Session.User.Ready.Id(sessionId),
			)
		})
	}

	return cmd
}

func newUserDeleteSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <session-id>",
		Short: "删除会话",
		Long:  "删除会话",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		sessionId := args[0]
		if sessionId == "" {
			return fmt.Errorf("session-id cannot be empty")
		}

		return callWithDumpReqResp(&session.ApiDeleteRequest{
			SessionId: &sessionId,
		}, func() (interface{}, error) {
			return GetYsClient().CloudApp.Session.User.Delete(
				GetYsClient().CloudApp.Session.User.Delete.Id(sessionId),
			)
		})
	}

	return cmd
}

func newUserGetSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <session-id>",
		Short: "查询单个会话",
		Long:  "查询单个会话",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		sessionId := args[0]
		if sessionId == "" {
			return fmt.Errorf("session-id cannot be empty")
		}

		return callWithDumpReqResp(&session.ApiGetRequest{
			SessionId: &sessionId,
		}, func() (interface{}, error) {
			return GetYsClient().CloudApp.Session.User.Get(
				GetYsClient().CloudApp.Session.User.Get.Id(sessionId),
			)
		})
	}

	return cmd
}

func newUserListSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "批量查询会话",
		Long:  "批量查询会话",
	}

	o := InitCloudAppBaseOptionsByList(cmd)
	status := new(string)
	cmd.Flags().StringVar(status, "status", "", "\"STARTED\" or \"STARTED,STARTING,CLOSING,CLOSED\"")
	sessionIds := new(string)
	cmd.Flags().StringVar(sessionIds, "session-ids", "", "\"ida\" or \"ida,idb\"")
	zone := new(string)
	cmd.Flags().StringVar(zone, "zone", "", "zone")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := &session.ApiListRequest{
			PageOffset: &o.PageOffset,
			PageSize:   &o.PageSize,
			Status:     status,
			SessionIds: sessionIds,
			Zone:       zone,
		}

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Session.User.List(
				GetYsClient().CloudApp.Session.User.List.PageOffset(*req.PageOffset),
				GetYsClient().CloudApp.Session.User.List.PageSize(*req.PageSize),
				GetYsClient().CloudApp.Session.User.List.Status(*req.Status),
				GetYsClient().CloudApp.Session.User.List.SessionIds(*req.SessionIds),
				GetYsClient().CloudApp.Session.User.List.Zone(*req.Zone),
			)
		})
	}

	return cmd
}

func newStartSessionCmd(authType string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start <session-id>",
		Short: "开机",
		Long:  "开机",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		sessionId := args[0]
		if sessionId == "" {
			return fmt.Errorf("session-id cannot be empty")
		}

		req := &session.PowerOnRequest{
			SessionId: &sessionId,
		}

		switch authType {
		case "user":
			return callWithDumpReqResp(req, func() (interface{}, error) {
				return GetYsClient().CloudApp.Session.User.PowerOn(
					GetYsClient().CloudApp.Session.User.PowerOn.Id(sessionId),
				)
			})
		case "admin":
			return callWithDumpReqResp(req, func() (interface{}, error) {
				return GetYsClient().CloudApp.Session.Admin.PowerOn(
					GetYsClient().CloudApp.Session.Admin.PowerOn.Id(sessionId),
				)
			})
		default:
			return fmt.Errorf("unknown authType: %s", authType)
		}
	}

	return cmd
}

func newStopSessionCmd(authType string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop <session-id>",
		Short: "关机",
		Long:  "关机",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		sessionId := args[0]
		if sessionId == "" {
			return fmt.Errorf("session-id cannot be empty")
		}

		req := &session.PowerOffRequest{
			SessionId: &sessionId,
		}

		switch authType {
		case "user":
			return callWithDumpReqResp(req, func() (interface{}, error) {
				return GetYsClient().CloudApp.Session.User.PowerOff(
					GetYsClient().CloudApp.Session.User.PowerOff.Id(sessionId),
				)
			})
		case "admin":
			return callWithDumpReqResp(req, func() (interface{}, error) {
				return GetYsClient().CloudApp.Session.Admin.PowerOff(
					GetYsClient().CloudApp.Session.Admin.PowerOff.Id(sessionId),
				)
			})
		default:
			return fmt.Errorf("unknown authType: %s", authType)
		}
	}

	return cmd
}

func newRestartSessionCmd(authType string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart <session-id>",
		Short: "重启",
		Long:  "重启",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		sessionId := args[0]
		if sessionId == "" {
			return fmt.Errorf("session-id cannot be empty")
		}

		req := &session.RebootRequest{
			SessionId: &sessionId,
		}

		switch authType {
		case "user":
			return callWithDumpReqResp(req, func() (interface{}, error) {
				return GetYsClient().CloudApp.Session.User.Reboot(
					GetYsClient().CloudApp.Session.User.Reboot.Id(sessionId),
				)
			})
		case "admin":
			return callWithDumpReqResp(req, func() (interface{}, error) {
				return GetYsClient().CloudApp.Session.Admin.Reboot(
					GetYsClient().CloudApp.Session.Admin.Reboot.Id(sessionId),
				)
			})
		default:
			return fmt.Errorf("unknown authType: %s", authType)
		}
	}

	return cmd
}

func newExecScriptSessionCmd(authType string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec-script <session-id> --script-runner <runner> --script-content-encoded <content-base64-encoded> --wait",
		Short: "指定会话执行脚本",
		Long: `exec-script <session-id> --script-runner powershell --script-content-encoded <content-base64-encoded> --wait
or  exec-script <session-id> --script-runner powershell --script-path /path/to/script/path --wait
`,
	}

	var scriptRunner, scriptContentEncoded, scriptPath string
	var wait bool
	cmd.Flags().StringVar(&scriptRunner, "script-runner", "", "powershell is default and only support for now")
	cmd.Flags().StringVar(&scriptContentEncoded, "script-content-encoded", "", "script content encoded by base64")
	cmd.Flags().StringVar(&scriptPath, "script-path", "", "/path/to/script/path")
	cmd.Flags().BoolVar(&wait, "wait", false, "block to wait the script executed")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		sessionId := args[0]
		if sessionId == "" {
			return fmt.Errorf("session-id cannot be empty")
		}

		if scriptRunner == "" {
			scriptRunner = "powershell"
		}

		if scriptContentEncoded == "" && scriptPath == "" {
			return fmt.Errorf("--script-content-encoded or --script-path cannot be empty both")
		}

		if scriptContentEncoded != "" && scriptPath != "" {
			return fmt.Errorf("cannot set --script-content-encoded and --script-path in a time")
		}

		req := &session.ExecScriptRequest{
			SessionId:    &sessionId,
			ScriptRunner: &scriptRunner,
			WaitTillEnd:  &wait,
		}

		if scriptContentEncoded != "" {
			req.ScriptContent = &scriptContentEncoded
		}
		if scriptPath != "" {
			content, err := os.ReadFile(scriptPath)
			if err != nil {
				return fmt.Errorf("read file from [%s] failed, %w", scriptPath, err)
			}

			contentEncoded := base64.StdEncoding.EncodeToString(content)
			req.ScriptContent = &contentEncoded
		}

		switch authType {
		case "admin":
			return callWithDumpReqResp(req, func() (interface{}, error) {
				return GetYsClient().CloudApp.Session.Admin.ExecScript(
					GetYsClient().CloudApp.Session.Admin.ExecScript.SessionId(*req.SessionId),
					GetYsClient().CloudApp.Session.Admin.ExecScript.ScriptRunner(*req.ScriptRunner),
					GetYsClient().CloudApp.Session.Admin.ExecScript.ScriptContent(*req.ScriptContent),
					GetYsClient().CloudApp.Session.Admin.ExecScript.WaitTillEnd(*req.WaitTillEnd),
				)
			})
		case "user":
			return callWithDumpReqResp(req, func() (interface{}, error) {
				return GetYsClient().CloudApp.Session.User.ExecScript(
					GetYsClient().CloudApp.Session.User.ExecScript.SessionId(*req.SessionId),
					GetYsClient().CloudApp.Session.User.ExecScript.ScriptRunner(*req.ScriptRunner),
					GetYsClient().CloudApp.Session.User.ExecScript.ScriptContentEncoded(*req.ScriptContent),
					GetYsClient().CloudApp.Session.User.ExecScript.WaitTillEnd(*req.WaitTillEnd),
				)
			})
		default:
			return fmt.Errorf("unknow authType")
		}
	}

	return cmd
}

func newUserRestoreSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore <old-session-id>",
		Short: "重建",
		Long:  "从原会话的启动盘，重建一个新会话",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		sessionId := args[0]
		if sessionId == "" {
			return fmt.Errorf("old-session-id cannot be empty")
		}

		req := &session.ApiRestoreRequest{
			SessionId: &sessionId,
		}

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Session.User.Restore(
				GetYsClient().CloudApp.Session.User.Restore.SessionId(sessionId),
			)
		})
	}

	return cmd
}

func newMountSessionCmd(authType string) *cobra.Command {
	cmd := &cobra.Command{
		Use: "mount <session-id> --share-directory <sub_dir> --mount-point <mount_point>" + "\n" +
			color.New(color.FgRed).Add(color.Bold).Sprintf("挂载CSP存储需要将share-directory指定为common/<project-id>/<sub-path>，填空或不填会造成挂载CSP根目录的行为，请谨慎操作"),
		Short: "将用户存储挂载至会话中",
		Long: `mount <session-id> --share-directory <sub_dir> --mount-point <mount_point>
example: mount sessionId --share-directory dirA --mount-point X:`,
	}

	var shareDirectory, mountPoint string
	cmd.Flags().StringVar(&shareDirectory, "share-directory", "", "the dir path be mounted")
	cmd.Flags().StringVar(&mountPoint, "mount-point", "", "mount point in session")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		sessionId := args[0]
		if sessionId == "" {
			return fmt.Errorf("session-id cannot be empty")
		}

		if mountPoint == "" {
			return fmt.Errorf("mount point cannot be empty")
		}

		if shareDirectory == "" {
			// confirm if share-directory is empty
			var confirm string
			fmt.Print(color.New(color.FgRed).Add(color.Bold).Sprintf("share-directory is empty, please confirm whether to continue(y/N): "))
			_, err := fmt.Scanln(&confirm)
			if err != nil {
				return fmt.Errorf("scan user input to confirm when share-directory is empty failed, %w", err)
			}

			if strings.ToUpper(confirm) != "Y" {
				return fmt.Errorf("cancel mount")
			}
		}

		req := &session.MountRequest{
			SessionId:      &sessionId,
			ShareDirectory: &shareDirectory,
			MountPoint:     &mountPoint,
		}

		switch authType {
		case "admin":
			return callWithDumpReqResp(req, func() (interface{}, error) {
				return GetYsClient().CloudApp.Session.Admin.Mount(
					GetYsClient().CloudApp.Session.Admin.Mount.SessionId(*req.SessionId),
					GetYsClient().CloudApp.Session.Admin.Mount.MountPoint(*req.MountPoint),
					GetYsClient().CloudApp.Session.Admin.Mount.ShareDirectory(*req.ShareDirectory),
				)
			})
		case "user":
			return callWithDumpReqResp(req, func() (interface{}, error) {
				return GetYsClient().CloudApp.Session.User.Mount(
					GetYsClient().CloudApp.Session.User.Mount.SessionId(*req.SessionId),
					GetYsClient().CloudApp.Session.User.Mount.MountPoint(*req.MountPoint),
					GetYsClient().CloudApp.Session.User.Mount.ShareDirectory(*req.ShareDirectory),
				)
			})
		}

		return nil
	}

	return cmd
}

func newUmountSessionCmd(authType string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "umount <session-id> --mount-point <mount_point>",
		Short: "将用户存储于会话解挂载",
		Long: `umount <session-id> --mount-point <mount_point>
example: umount sessionId --mount-point X:`,
	}

	var mountPoint string
	cmd.Flags().StringVar(&mountPoint, "mount-point", "", "mount point in session")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		sessionId := args[0]
		if sessionId == "" {
			return fmt.Errorf("session-id cannot be empty")
		}

		if mountPoint == "" {
			return fmt.Errorf("mount point cannot be empty")
		}

		req := &session.UmountRequest{
			SessionId:  &sessionId,
			MountPoint: &mountPoint,
		}

		switch authType {
		case "admin":
			return callWithDumpReqResp(req, func() (interface{}, error) {
				return GetYsClient().CloudApp.Session.Admin.Umount(
					GetYsClient().CloudApp.Session.Admin.Umount.SessionId(*req.SessionId),
					GetYsClient().CloudApp.Session.Admin.Umount.MountPoint(*req.MountPoint),
				)
			})
		case "user":
			return callWithDumpReqResp(req, func() (interface{}, error) {
				return GetYsClient().CloudApp.Session.User.Umount(
					GetYsClient().CloudApp.Session.User.Umount.SessionId(*req.SessionId),
					GetYsClient().CloudApp.Session.User.Umount.MountPoint(*req.MountPoint),
				)
			})
		}

		return nil
	}

	return cmd
}

func newAdminSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "Session管理员功能，查询/关闭等",
		Long:  "Session管理员功能，查询/关闭等",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newAdminListSessionCmd(),
		newAdminCloseSessionCmd(),
		newStartSessionCmd("admin"),
		newStopSessionCmd("admin"),
		newRestartSessionCmd("admin"),
		newAdminRestoreSessionCmd(),
		newExecScriptSessionCmd("admin"),
		newMountSessionCmd("admin"),
		newUmountSessionCmd("admin"),
	)

	return cmd
}

func newAdminListSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "批量查询会话",
		Long:  "批量查询会话",
	}

	o := InitCloudAppBaseOptionsByList(cmd)
	status := new(string)
	cmd.Flags().StringVar(status, "status", "", "\"STARTED\" or \"STARTED,STARTING,CLOSING,CLOSED\"")
	sessionIds := new(string)
	cmd.Flags().StringVar(sessionIds, "session-ids", "", "\"ida\" or \"ida,idb\"")
	zone := new(string)
	cmd.Flags().StringVar(zone, "zone", "", "zone")
	userIds := new(string)
	cmd.Flags().StringVar(userIds, "user-ids", "", "\"ida\" or \"ida,idb\"")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := &session.AdminListRequest{
			PageOffset: &o.PageOffset,
			PageSize:   &o.PageSize,
			Status:     status,
			SessionIds: sessionIds,
			Zone:       zone,
			UserIds:    userIds,
		}

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Session.Admin.List(
				GetYsClient().CloudApp.Session.Admin.List.PageOffset(*req.PageOffset),
				GetYsClient().CloudApp.Session.Admin.List.PageSize(*req.PageSize),
				GetYsClient().CloudApp.Session.Admin.List.Status(*req.Status),
				GetYsClient().CloudApp.Session.Admin.List.SessionIds(*req.SessionIds),
				GetYsClient().CloudApp.Session.Admin.List.Zone(*req.Zone),
				GetYsClient().CloudApp.Session.Admin.List.UserIds(*req.UserIds),
			)
		})
	}

	return cmd
}

func newAdminCloseSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close <session-id> <reason>",
		Short: "关闭会话（删除可视化机器）",
		Long:  "关闭会话（删除可视化机器）",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("expect 2 args but got %d", len(args))
		}

		sessionId := args[0]
		if sessionId == "" {
			return fmt.Errorf("session-id cannot be empty")
		}

		reason := args[1]
		if reason == "" {
			return fmt.Errorf("reason cannot be empty")
		}

		return callWithDumpReqResp(&session.AdminCloseRequest{
			SessionId: &sessionId,
			Reason:    &reason,
		}, func() (interface{}, error) {
			return GetYsClient().CloudApp.Session.Admin.Close(
				GetYsClient().CloudApp.Session.Admin.Close.Id(sessionId),
				GetYsClient().CloudApp.Session.Admin.Close.Reason(reason),
			)
		})
	}

	return cmd
}

func newAdminRestoreSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore <user-id> <old-session-id>",
		Short: "重建",
		Long:  "以指定用户原会话的启动盘，重建一个新会话",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("expect 2 args but got %d", len(args))
		}

		userId := args[0]
		if userId == "" {
			return fmt.Errorf("user-id cannot be empty")
		}

		sessionId := args[1]
		if sessionId == "" {
			return fmt.Errorf("old-session-id cannot be empty")
		}

		req := &session.ApiRestoreRequest{
			SessionId: &sessionId,
		}

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Session.Admin.Restore(
				GetYsClient().CloudApp.Session.Admin.Restore.UserId(userId),
				GetYsClient().CloudApp.Session.Admin.Restore.SessionId(sessionId),
			)
		})
	}

	return cmd
}

func newSoftwareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "software",
		Short: "软件管理，创建/删除/更新/查询等",
		Long:  "软件管理，创建/删除/更新/查询等",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newUserSoftwareCmd(),
		newAdminSoftwareCmd(),
	)

	return cmd
}

func newUserSoftwareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Software普通用户功能，查询等",
		Long:  "Software普通用户功能，查询等",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newUserSoftwareGetCmd(),
		newUserSoftwareListCmd(),
	)

	return cmd
}

func newUserSoftwareGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <software-id>",
		Short: "查询单个软件",
		Long:  "查询单个软件",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		softwareId := args[0]
		if softwareId == "" {
			return fmt.Errorf("software-id cannot be empty")
		}

		return callWithDumpReqResp(&software.APIGetRequest{
			SoftwareId: &softwareId,
		}, func() (interface{}, error) {
			return GetYsClient().CloudApp.Software.User.Get(
				GetYsClient().CloudApp.Software.User.Get.Id(softwareId),
			)
		})
	}

	return cmd
}

func newUserSoftwareListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "批量查询软件",
		Long:  "批量查询软件",
	}

	o := InitCloudAppBaseOptionsByList(cmd)
	zone := ""
	cmd.Flags().StringVar(&zone, "zone", "", "zone")
	name := ""
	cmd.Flags().StringVar(&name, "name", "", "name")
	platform := ""
	cmd.Flags().StringVar(&platform, "platform", "", "platform [ WINDOWS | LINUX ]")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := &software.APIListRequest{
			PageOffset: &o.PageOffset,
			PageSize:   &o.PageSize,
		}
		if zone != "" {
			req.Zone = &zone
		}
		if name != "" {
			req.Name = &name
		}
		if platform != "" {
			req.Platform = &platform
		}

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Software.User.List(ensureUserSoftwareListOpts(req)...)
		})
	}

	return cmd
}

func ensureUserSoftwareListOpts(req *software.APIListRequest) []softwareUserList.Option {
	opts := make([]softwareUserList.Option, 0)
	if req.Zone != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.User.List.Zone(*req.Zone))
	}
	if req.Name != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.User.List.Name(*req.Name))
	}
	if req.Platform != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.User.List.Platform(*req.Platform))
	}
	if req.PageSize != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.User.List.PageSize(*req.PageSize))
	}
	if req.PageOffset != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.User.List.PageOffset(*req.PageOffset))
	}

	return opts
}

func newAdminSoftwareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "Software管理员功能，增删改查等",
		Long:  "Software管理员功能，增删改查等",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newAdminAddSoftwareCmd(),
		newAdminAddSoftwareUsersCmd(),
		newAdminDeleteSoftwareCmd(),
		newAdminDeleteSoftwareUsersCmd(),
		newAdminGetSoftwareCmd(),
		newAdminListSoftwareCmd(),
		newAdminModifySoftwareCmd(),
		newAdminEditInitScriptCmd(),
	)

	return cmd
}

func newAdminAddSoftwareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add -F req.json",
		Short: "添加软件",
		Long:  "添加软件",
	}

	o := InitCloudAppBaseOptions(cmd)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		data, err := checkReqJsonFileAndRead(o.ReqJson)
		if err != nil {
			return err
		}

		req := new(software.AdminPostRequest)
		if err = jsoniter.Unmarshal(data, req); err != nil {
			return fmt.Errorf("unmarshan json to *software.AdminPostRequest failed, %w", err)
		}

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Software.Admin.Add(ensureAdminAddSoftwareOpts(req)...)
		})
	}

	return cmd
}

func ensureAdminAddSoftwareOpts(req *software.AdminPostRequest) []softwareAdminAdd.Option {
	opts := make([]softwareAdminAdd.Option, 0)
	if req.Zone != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Add.Zone(*req.Zone))
	}
	if req.Name != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Add.Name(*req.Name))
	}
	if req.Desc != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Add.Desc(*req.Desc))
	}
	if req.Icon != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Add.Icon(*req.Icon))
	}
	if req.Platform != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Add.Platform(*req.Platform))
	}
	if req.ImageId != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Add.ImageId(*req.ImageId))
	}
	if req.InitScript != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Add.InitScript(*req.InitScript))
	}
	if req.GpuDesired != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Add.GpuDesired(*req.GpuDesired))
	}

	return opts
}

func newAdminAddSoftwareUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-users <softwares> <users>",
		Short: "添加用户",
		Long: `添加用户
支持单个/批量,用逗号隔开 批量: add-users softwareId1,softwareId2 user1,user2
                     单个: add-users softwareId userId`,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("expect 2 args but got %d", len(args))
		}

		softwares, users := args[0], args[1]
		if softwares == "" {
			return fmt.Errorf("softwares cannot be empty")
		}

		if users == "" {
			return fmt.Errorf("users cannot be empty")
		}

		req := &software.AdminPostUsersRequest{
			Softwares: strings.Split(softwares, ","),
			Users:     strings.Split(users, ","),
		}

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Software.Admin.AddUsers(
				GetYsClient().CloudApp.Software.Admin.AddUsers.Softwares(req.Softwares),
				GetYsClient().CloudApp.Software.Admin.AddUsers.Users(req.Users),
			)
		})
	}

	return cmd
}

func newAdminDeleteSoftwareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <software-id>",
		Short: "删除软件",
		Long:  "删除软件",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		softwareId := args[0]
		if softwareId == "" {
			return fmt.Errorf("software-id cannot be empty")
		}

		return callWithDumpReqResp(&software.AdminDeleteRequest{
			SoftwareId: &softwareId,
		}, func() (interface{}, error) {
			return GetYsClient().CloudApp.Software.Admin.Delete(
				GetYsClient().CloudApp.Software.Admin.Delete.Id(softwareId),
			)
		})
	}

	return cmd
}

func newAdminDeleteSoftwareUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-users <softwares> <users>",
		Short: "删除用户",
		Long: `删除用户
支持单个/批量,用逗号隔开 批量: delete-users <softwares> <users>`,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("expect 2 args but got %d", len(args))
		}

		softwares, users := args[0], args[1]
		if softwares == "" {
			return fmt.Errorf("softwares cannot be empty")
		}

		if users == "" {
			return fmt.Errorf("users cannot be empty")
		}

		req := &software.AdminDeleteUsersRequest{
			Softwares: strings.Split(softwares, ","),
			Users:     strings.Split(users, ","),
		}

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Software.Admin.DeleteUsers(
				GetYsClient().CloudApp.Software.Admin.DeleteUsers.Softwares(req.Softwares),
				GetYsClient().CloudApp.Software.Admin.DeleteUsers.Users(req.Users),
			)
		})
	}

	return cmd
}

func newAdminGetSoftwareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <software-id>",
		Short: "查询单个软件",
		Long:  "查询单个软件",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		softwareId := args[0]
		if softwareId == "" {
			return fmt.Errorf("software-id cannot be empty")
		}

		return callWithDumpReqResp(&software.AdminGetRequest{
			SoftwareId: &softwareId,
		}, func() (interface{}, error) {
			return GetYsClient().CloudApp.Software.Admin.Get(
				GetYsClient().CloudApp.Software.Admin.Get.Id(softwareId),
			)
		})
	}

	return cmd
}

func newAdminListSoftwareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "批量查询软件",
		Long:  "批量查询软件",
	}

	o := InitCloudAppBaseOptionsByList(cmd)
	userId := ""
	cmd.Flags().StringVar(&userId, "userId", "", "userId")
	zone := ""
	cmd.Flags().StringVar(&zone, "zone", "", "zone")
	name := ""
	cmd.Flags().StringVar(&name, "name", "", "name")
	platform := ""
	cmd.Flags().StringVar(&platform, "platform", "", "platform [ WINDOWS | LINUX ]")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := &software.AdminListRequest{
			PageOffset: &o.PageOffset,
			PageSize:   &o.PageSize,
		}
		if userId != "" {
			req.UserId = &userId
		}
		if zone != "" {
			req.Zone = &zone
		}
		if name != "" {
			req.Name = &name
		}
		if platform != "" {
			req.Platform = &platform
		}

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Software.Admin.List(ensureAdminSoftwareListOpts(req)...)
		})
	}

	return cmd
}

func ensureAdminSoftwareListOpts(req *software.AdminListRequest) []softwareAdminList.Option {
	opts := make([]softwareAdminList.Option, 0)
	if req.UserId != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.List.UserId(*req.UserId))
	}
	if req.Zone != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.List.Zone(*req.Zone))
	}
	if req.Name != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.List.Name(*req.Name))
	}
	if req.Platform != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.List.Platform(*req.Platform))
	}
	if req.PageSize != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.List.PageSize(*req.PageSize))
	}
	if req.PageOffset != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.List.PageOffset(*req.PageOffset))
	}

	return opts
}

func newAdminModifySoftwareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify <software-id> -F req.json",
		Short: "修改软件（增量）",
		Long:  "修改软件（增量）",
	}

	o := InitCloudAppBaseOptions(cmd)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		softwareId := args[0]
		if softwareId == "" {
			return fmt.Errorf("software-id cannot be empty")
		}

		data, err := checkReqJsonFileAndRead(o.ReqJson)
		if err != nil {
			return err
		}

		req := new(software.AdminPatchRequest)
		if err = jsoniter.Unmarshal(data, req); err != nil {
			return fmt.Errorf("unmarshal to *software.AdminPatchRequest struct failed, %w", err)
		}
		req.SoftwareId = &softwareId

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Software.Admin.Patch(ensureAdminPatchSoftwareOpts(req)...)
		})
	}

	return cmd
}

func newAdminEditInitScriptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-init-script <software-id>",
		Short: "编辑软件初始化脚本",
		Long:  "编辑软件初始化脚本",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but go %d", len(args))
		}

		softwareId := args[0]
		if softwareId == "" {
			return fmt.Errorf("software-id cannot be empty")
		}

		getReq := &software.AdminGetRequest{
			SoftwareId: &softwareId,
		}

		var err error
		var resp *software.AdminGetResponse
		err = callWithDumpReqResp(getReq, func() (interface{}, error) {
			resp, err = GetYsClient().CloudApp.Software.Admin.Get(
				GetYsClient().CloudApp.Software.Admin.Get.Id(softwareId),
			)
			return resp, err
		})
		if err != nil {
			return err
		}
		if resp == nil || resp.Data == nil {
			return fmt.Errorf("invalid resp")
		}

		result, err := texteditor.EditorStatic([]byte(resp.Data.InitScript))
		if err != nil {
			return fmt.Errorf("edit init script failed, %w", err)
		}
		r := string(result)

		patchRep := &software.AdminPatchRequest{
			SoftwareId: &softwareId,
			InitScript: &r,
		}
		return callWithDumpReqResp(patchRep, func() (interface{}, error) {
			return GetYsClient().CloudApp.Software.Admin.Patch(
				GetYsClient().CloudApp.Software.Admin.Patch.Id(*patchRep.SoftwareId),
				GetYsClient().CloudApp.Software.Admin.Patch.InitScript(*patchRep.InitScript),
			)
		})
	}

	return cmd
}

func ensureAdminPatchSoftwareOpts(req *software.AdminPatchRequest) []softwareAdminPatch.Option {
	opts := make([]softwareAdminPatch.Option, 0)
	if req.SoftwareId != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Patch.Id(*req.SoftwareId))
	}
	if req.Zone != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Patch.Zone(*req.Zone))
	}
	if req.Name != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Patch.Name(*req.Name))
	}
	if req.Desc != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Patch.Desc(*req.Desc))
	}
	if req.Icon != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Patch.Icon(*req.Icon))
	}
	if req.Platform != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Patch.Platform(*req.Platform))
	}
	if req.ImageId != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Patch.ImageId(*req.ImageId))
	}
	if req.InitScript != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Patch.InitScript(*req.InitScript))
	}
	if req.GpuDesired != nil {
		opts = append(opts, GetYsClient().CloudApp.Software.Admin.Patch.GpuDesired(*req.GpuDesired))
	}

	return opts
}

func newHardwareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hardware",
		Short: "硬件管理，创建/删除/更新/查询等",
		Long:  "硬件管理，创建/删除/更新/查询等",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newUserHardwareCmd(),
		newAdminHardwareCmd(),
	)

	return cmd
}

func newUserHardwareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Hardware普通用户功能，查询等",
		Long:  "Hardware普通用户功能，查询等",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newUserHardwareGetCmd(),
		newUserHardwareListCmd(),
	)

	return cmd
}

func newUserHardwareGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <hardware-id>",
		Short: "查询单个硬件",
		Long:  "查询单个硬件",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		hardwareId := args[0]
		if hardwareId == "" {
			return fmt.Errorf("hardware-id cannot be empty")
		}

		return callWithDumpReqResp(&hardware.APIGetRequest{
			HardwareId: &hardwareId,
		}, func() (interface{}, error) {
			return GetYsClient().CloudApp.Hardware.User.Get(
				GetYsClient().CloudApp.Hardware.User.Get.Id(hardwareId),
			)
		})
	}

	return cmd
}

func newUserHardwareListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "批量查询硬件",
		Long:  "批量查询硬件",
	}

	o := InitCloudAppBaseOptionsByList(cmd)
	zone := ""
	cmd.Flags().StringVar(&zone, "zone", "", "zone")
	name := ""
	cmd.Flags().StringVar(&name, "name", "", "name")
	cpuStr := ""
	cmd.Flags().StringVar(&cpuStr, "cpu", "", "cpu")
	memStr := ""
	cmd.Flags().StringVar(&memStr, "mem", "", "memory [MB]")
	gpuStr := ""
	cmd.Flags().StringVar(&gpuStr, "gpu", "", "gpu count")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := &hardware.APIListRequest{
			PageOffset: &o.PageOffset,
			PageSize:   &o.PageSize,
		}
		if zone != "" {
			req.Zone = &zone
		}
		if name != "" {
			req.Name = &name
		}
		if cpuStr != "" {
			cpu, err := strconv.Atoi(cpuStr)
			if err != nil {
				return fmt.Errorf("parse cpu to int failed, %w", err)
			}

			req.Cpu = &cpu
		}
		if memStr != "" {
			mem, err := strconv.Atoi(memStr)
			if err != nil {
				return fmt.Errorf("parse mem to int failed, %w", err)
			}

			req.Mem = &mem
		}
		if gpuStr != "" {
			gpu, err := strconv.Atoi(gpuStr)
			if err != nil {
				return fmt.Errorf("parse gpu to int failed, %w", err)
			}

			req.Gpu = &gpu
		}

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Hardware.User.List(ensureUserHardwareListOpts(req)...)
		})
	}

	return cmd
}

func ensureUserHardwareListOpts(req *hardware.APIListRequest) []hardwareApiList.Option {
	opts := make([]hardwareApiList.Option, 0)
	if req.PageOffset != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.User.List.PageOffset(*req.PageOffset))
	}
	if req.PageSize != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.User.List.PageSize(*req.PageSize))
	}
	if req.Zone != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.User.List.Zone(*req.Zone))
	}
	if req.Name != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.User.List.Name(*req.Name))
	}
	if req.Cpu != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.User.List.Cpu(*req.Cpu))
	}
	if req.Mem != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.User.List.Mem(*req.Mem))
	}
	if req.Gpu != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.User.List.Gpu(*req.Gpu))
	}

	return opts
}

func newAdminHardwareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "Hardware管理员功能，增删改查等",
		Long:  "Hardware管理员功能，增删改查等",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newAdminHardwareAddCmd(),
		newAdminHardwareAddUsersCmd(),
		newAdminHardwareDeleteCmd(),
		newAdminHardwareDeleteUsersCmd(),
		newAdminHardwareGetCmd(),
		newAdminHardwareListCmd(),
		newAdminHardwareModifyCmd(),
	)

	return cmd
}

func newAdminHardwareAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add -F req.json",
		Short: "创建硬件",
		Long:  "创建硬件",
	}

	o := InitCloudAppBaseOptions(cmd)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		data, err := checkReqJsonFileAndRead(o.ReqJson)
		if err != nil {
			return err
		}

		req := new(hardware.AdminPostRequest)
		if err = jsoniter.Unmarshal(data, req); err != nil {
			return fmt.Errorf("unmarshan json to *hardware.AdminPostRequest failed, %w", err)
		}

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Hardware.Admin.Add(ensureAdminAddHardwareOpts(req)...)
		})
	}

	return cmd
}

func ensureAdminAddHardwareOpts(req *hardware.AdminPostRequest) []hardwareAdminAdd.Option {
	opts := make([]hardwareAdminAdd.Option, 0)
	if req.Zone != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Add.Zone(*req.Zone))
	}
	if req.Name != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Add.Name(*req.Name))
	}
	if req.Desc != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Add.Desc(*req.Desc))
	}
	if req.InstanceType != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Add.InstanceType(*req.InstanceType))
	}
	if req.InstanceFamily != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Add.InstanceFamily(*req.InstanceFamily))
	}
	if req.Network != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Add.Network(*req.Network))
	}
	if req.Cpu != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Add.Cpu(*req.Cpu))
	}
	if req.CpuModel != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Add.CpuModel(*req.CpuModel))
	}
	if req.Mem != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Add.Mem(*req.Mem))
	}
	if req.Gpu != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Add.Gpu(*req.Gpu))
	}
	if req.GpuModel != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Add.GpuModel(*req.GpuModel))
	}

	return opts
}

func newAdminHardwareAddUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-users <hardwares> <users>",
		Short: "添加用户",
		Long: `添加用户
支持单个/批量,用逗号隔开 批量: add-users hardwareId1,hardwareId2 user1,user2
                     单个: add-users hardwareId userId`,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("expect 2 args but got %d", len(args))
		}

		hardwares, users := args[0], args[1]
		if hardwares == "" {
			return fmt.Errorf("hardwares cannot be empty")
		}
		if users == "" {
			return fmt.Errorf("users cannot be empty")
		}

		req := &hardware.AdminPostUsersRequest{
			Hardwares: strings.Split(hardwares, ","),
			Users:     strings.Split(users, ","),
		}
		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Hardware.Admin.AddUsers(
				GetYsClient().CloudApp.Hardware.Admin.AddUsers.Hardwares(req.Hardwares),
				GetYsClient().CloudApp.Hardware.Admin.AddUsers.Users(req.Users),
			)
		})
	}

	return cmd
}

func newAdminHardwareDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <hardware-id>",
		Short: "删除硬件",
		Long:  "删除硬件",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 2 args but got %d", len(args))
		}

		hardwareId := args[0]
		if hardwareId == "" {
			return fmt.Errorf("hardware-id cannot be empty")
		}

		return callWithDumpReqResp(&hardware.AdminDeleteRequest{
			HardwareId: &hardwareId,
		}, func() (interface{}, error) {
			return GetYsClient().CloudApp.Hardware.Admin.Delete(
				GetYsClient().CloudApp.Hardware.Admin.Delete.Id(hardwareId),
			)
		})
	}

	return cmd
}

func newAdminHardwareDeleteUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-users <hardwares> <users>",
		Short: "删除用户",
		Long: `删除用户
支持单个/批量,用逗号隔开 批量: delete-users hardwareId1,hardwareId2 user1,user2
                     单个: delete-users hardwareId userId`,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("expect 2 args but got %d", len(args))
		}

		hardwares, users := args[0], args[1]
		if hardwares == "" {
			return fmt.Errorf("hardwares cannot be empty")
		}
		if users == "" {
			return fmt.Errorf("users cannot be empty")
		}

		req := &hardware.AdminDeleteUsersRequest{
			Hardwares: strings.Split(hardwares, ","),
			Users:     strings.Split(users, ","),
		}
		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Hardware.Admin.DeleteUsers(
				GetYsClient().CloudApp.Hardware.Admin.DeleteUsers.Hardwares(req.Hardwares),
				GetYsClient().CloudApp.Hardware.Admin.DeleteUsers.Users(req.Users),
			)
		})
	}

	return cmd
}

func newAdminHardwareGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <hardware-id>",
		Short: "查询单个硬件",
		Long:  "查询单个硬件",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		hardwareId := args[0]
		if hardwareId == "" {
			return fmt.Errorf("hardware-id cannot be empty")
		}

		return callWithDumpReqResp(&hardware.AdminGetRequest{
			HardwareId: &hardwareId,
		}, func() (interface{}, error) {
			return GetYsClient().CloudApp.Hardware.Admin.Get(
				GetYsClient().CloudApp.Hardware.Admin.Get.Id(hardwareId),
			)
		})
	}

	return cmd
}

func newAdminHardwareListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "批量查询硬件",
		Long:  "批量查询硬件",
	}

	o := InitCloudAppBaseOptionsByList(cmd)
	userId := ""
	cmd.Flags().StringVar(&userId, "userId", "", "userId")
	zone := ""
	cmd.Flags().StringVar(&zone, "zone", "", "zone")
	name := ""
	cmd.Flags().StringVar(&name, "name", "", "name")
	cpuStr := ""
	cmd.Flags().StringVar(&cpuStr, "cpu", "", "cpu")
	memStr := ""
	cmd.Flags().StringVar(&memStr, "mem", "", "memory [MB]")
	gpuStr := ""
	cmd.Flags().StringVar(&gpuStr, "gpu", "", "gpu count")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := &hardware.AdminListRequest{
			PageOffset: &o.PageOffset,
			PageSize:   &o.PageSize,
		}
		if userId != "" {
			req.UserId = &userId
		}
		if zone != "" {
			req.Zone = &zone
		}
		if name != "" {
			req.Name = &name
		}
		if cpuStr != "" {
			cpu, err := strconv.Atoi(cpuStr)
			if err != nil {
				return fmt.Errorf("parse cpu to int failed, %w", err)
			}

			req.Cpu = &cpu
		}
		if memStr != "" {
			mem, err := strconv.Atoi(memStr)
			if err != nil {
				return fmt.Errorf("parse mem to int failed, %w", err)
			}

			req.Mem = &mem
		}
		if gpuStr != "" {
			gpu, err := strconv.Atoi(gpuStr)
			if err != nil {
				return fmt.Errorf("parse gpu to int failed, %w", err)
			}

			req.Gpu = &gpu
		}

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Hardware.Admin.List(ensureAdminHardwareListOpts(req)...)
		})
	}

	return cmd
}

func ensureAdminHardwareListOpts(req *hardware.AdminListRequest) []hardwareAdminList.Option {
	opts := make([]hardwareAdminList.Option, 0)
	if req.UserId != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.List.UserId(*req.UserId))
	}
	if req.PageOffset != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.List.PageOffset(*req.PageOffset))
	}
	if req.PageSize != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.List.PageSize(*req.PageSize))
	}
	if req.Zone != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.List.Zone(*req.Zone))
	}
	if req.Name != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.List.Name(*req.Name))
	}
	if req.Cpu != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.List.Cpu(*req.Cpu))
	}
	if req.Mem != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.List.Mem(*req.Mem))
	}
	if req.Gpu != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.List.Gpu(*req.Gpu))
	}

	return opts
}

func newAdminHardwareModifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify <hardware-id> -F req.json",
		Short: "修改硬件（增量）",
		Long:  "修改硬件（增量）",
	}

	o := InitCloudAppBaseOptions(cmd)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		hardwareId := args[0]
		if hardwareId == "" {
			return fmt.Errorf("hardware-id cannot be empty")
		}

		data, err := checkReqJsonFileAndRead(o.ReqJson)
		if err != nil {
			return err
		}

		req := new(hardware.AdminPatchRequest)
		if err = jsoniter.Unmarshal(data, req); err != nil {
			return fmt.Errorf("unmarshal to *hardware.AdminPatchRequest failed, %w", err)
		}
		req.HardwareId = &hardwareId

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.Hardware.Admin.Patch(ensureAdminHardwareModifyOpts(req)...)
		})
	}

	return cmd
}

func ensureAdminHardwareModifyOpts(req *hardware.AdminPatchRequest) []hardwareAdminPatch.Option {
	opts := make([]hardwareAdminPatch.Option, 0)
	if req.HardwareId != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Patch.Id(*req.HardwareId))
	}
	if req.Zone != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Patch.Zone(*req.Zone))
	}
	if req.Name != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Patch.Name(*req.Name))
	}
	if req.Desc != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Patch.Desc(*req.Desc))
	}
	if req.InstanceType != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Patch.InstanceType(*req.InstanceType))
	}
	if req.InstanceFamily != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Patch.InstanceFamily(*req.InstanceFamily))
	}
	if req.Network != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Patch.Network(*req.Network))
	}
	if req.Cpu != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Patch.Cpu(*req.Cpu))
	}
	if req.CpuModel != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Patch.CpuModel(*req.CpuModel))
	}
	if req.Mem != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Patch.Mem(*req.Mem))
	}
	if req.Gpu != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Patch.Gpu(*req.Gpu))
	}
	if req.GpuModel != nil {
		opts = append(opts, GetYsClient().CloudApp.Hardware.Admin.Patch.GpuModel(*req.GpuModel))
	}

	return opts
}

func newRemoteAppCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remoteapp",
		Short: "remoteapp相关功能",
		Long:  "remoteapp相关功能",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newUserRemoteAppCmd(),
		newAdminRemoteAppCmd(),
	)

	return cmd
}

func newUserRemoteAppCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "普通用户相关功能，查询等",
		Long:  "普通用户相关功能，查询等",
		RunE:  helpRun,
	}

	cmd.AddCommand(newUserRemoteAppGetCmd())

	return cmd
}

func newUserRemoteAppGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <session-id> <remoteapp-name>",
		Short: "查询对应会话的远程应用信息",
		Long:  "查询对应会话的远程应用信息",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("expect 2 args but got %d", len(args))
		}

		sessionId, remoteAppName := args[0], args[1]
		if sessionId == "" {
			return fmt.Errorf("session-id cannot be empty")
		}

		if remoteAppName == "" {
			return fmt.Errorf("remoteapp-name cannot be empty")
		}

		return callWithDumpReqResp(&remoteapp.ApiGetRequest{
			SessionId:     &sessionId,
			RemoteAppName: &remoteAppName,
		}, func() (interface{}, error) {
			return GetYsClient().CloudApp.RemoteApp.User.Get(
				GetYsClient().CloudApp.RemoteApp.User.Get.SessionId(sessionId),
				GetYsClient().CloudApp.RemoteApp.User.Get.RemoteAppName(remoteAppName),
			)
		})
	}

	return cmd
}

func newAdminRemoteAppCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "管理员相关功能，增删改等",
		Long:  "管理员相关功能，增删改等",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newAdminRemoteAppAddCmd(),
		newAdminRemoteAppDeleteCmd(),
		newAdminRemoteAppModifyCmd(),
	)

	return cmd
}

func newAdminRemoteAppAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add -F req.json",
		Short: "创建远程应用",
		Long:  "创建远程应用",
	}

	o := InitCloudAppBaseOptions(cmd)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		data, err := checkReqJsonFileAndRead(o.ReqJson)
		if err != nil {
			return err
		}

		req := new(remoteapp.AdminPostRequest)
		if err = jsoniter.Unmarshal(data, req); err != nil {
			return fmt.Errorf("unmarshal to *remoteapp.AdminPostRequest failed, %w", err)
		}

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.RemoteApp.Admin.Add(ensureAdminRemoteAppAddOpts(req)...)
		})
	}

	return cmd
}

func ensureAdminRemoteAppAddOpts(req *remoteapp.AdminPostRequest) []remoteAppAdminAdd.Option {
	opts := make([]remoteAppAdminAdd.Option, 0)
	if req.SoftwareId != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Add.SoftwareId(*req.SoftwareId))
	}
	if req.Desc != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Add.Desc(*req.Desc))
	}
	if req.Name != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Add.Name(*req.Name))
	}
	if req.Dir != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Add.Desc(*req.Dir))
	}
	if req.Args != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Add.Args(*req.Args))
	}
	if req.Logo != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Add.Logo(*req.Logo))
	}
	if req.DisableGfx != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Add.DisableGfx(*req.DisableGfx))
	}
	if req.LoginUser != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Add.LoginUser(*req.LoginUser))
	}

	return opts
}

func newAdminRemoteAppDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <remoteapp-id>",
		Short: "创建远程应用",
		Long:  "创建远程应用",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		remoteAppId := args[0]
		if remoteAppId == "" {
			return fmt.Errorf("remoteapp-id cannot be empty")
		}

		return callWithDumpReqResp(&remoteapp.AdminDeleteRequest{
			RemoteAppId: &remoteAppId,
		}, func() (interface{}, error) {
			return GetYsClient().CloudApp.RemoteApp.Admin.Delete(
				GetYsClient().CloudApp.RemoteApp.Admin.Delete.Id(remoteAppId),
			)
		})
	}

	return cmd
}

func newAdminRemoteAppModifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify <remoteapp-id> -F req.json",
		Short: "修改远程应用（增量）",
		Long:  "修改远程应用（增量）",
	}

	o := InitCloudAppBaseOptions(cmd)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expect 1 args but got %d", len(args))
		}

		remoteAppId := args[0]
		if remoteAppId == "" {
			return fmt.Errorf("remoteapp-id cannot be empty")
		}

		data, err := checkReqJsonFileAndRead(o.ReqJson)
		if err != nil {
			return err
		}

		req := new(remoteapp.AdminPatchRequest)
		if err = jsoniter.Unmarshal(data, req); err != nil {
			return fmt.Errorf("unmarshal to *remoteapp.AdminPostRequest failed, %w", err)
		}
		req.RemoteAppId = &remoteAppId

		return callWithDumpReqResp(req, func() (interface{}, error) {
			return GetYsClient().CloudApp.RemoteApp.Admin.Patch(ensureAdminRemoteAppModifyOpts(req)...)
		})
	}

	return cmd
}

func ensureAdminRemoteAppModifyOpts(req *remoteapp.AdminPatchRequest) []remoteAppAdminPatch.Option {
	opts := make([]remoteAppAdminPatch.Option, 0)
	if req.RemoteAppId != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Patch.Id(*req.RemoteAppId))
	}
	if req.SoftwareId != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Patch.SoftwareId(*req.SoftwareId))
	}
	if req.Desc != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Patch.Desc(*req.Desc))
	}
	if req.Name != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Patch.Name(*req.Name))
	}
	if req.Dir != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Patch.Dir(*req.Dir))
	}
	if req.Args != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Patch.Args(*req.Args))
	}
	if req.Logo != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Patch.Logo(*req.Logo))
	}
	if req.DisableGfx != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Patch.DisableGfx(*req.DisableGfx))
	}
	if req.LoginUser != nil {
		opts = append(opts, GetYsClient().CloudApp.RemoteApp.Admin.Patch.LoginUser(*req.LoginUser))
	}

	return opts
}
