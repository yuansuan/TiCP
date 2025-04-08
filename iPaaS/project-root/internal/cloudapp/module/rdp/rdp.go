package rdp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"text/template"

	sprig "github.com/go-task/slim-sprig"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/rdpgo/guacamole"
	"github.com/yuansuan/ticp/rdpgo/jwt"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
)

// example: remoteAppArgsTemplate = "c:\\tmp\\chrome\\{{ randNumeric 5 }}"
// render result be like: remoteAppArgsTemplate = "c:\\tmp\\chrome\\29384"
func renderRemoteAppArgsTemplate(remoteAppArgsTemp string) (string, error) {
	tpl := template.Must(template.New("_").Funcs(sprig.TxtFuncMap()).Parse(remoteAppArgsTemp))

	buf := bytes.Buffer{}
	if err := tpl.Execute(&buf, nil); err != nil {
		return "", fmt.Errorf("template execute failed, %w", err)
	}

	return buf.String(), nil
}

func GenerateRemoteAppURLBase64(sessionDetail *models.SessionWithDetail, remoteAppUserName *models.RemoteAppUserPass, remoteApp *models.RemoteApp, c cloud.Aggregator) (string, error) {
	url, err := generateURL(sessionDetail, c, remoteAppUserName.Username, remoteAppUserName.Password, remoteApp)
	if err != nil {
		return "", fmt.Errorf("generate url failed, %w", err)
	}

	return base64.StdEncoding.EncodeToString([]byte(url)), nil
}

func generateURL(sessionDetail *models.SessionWithDetail, c cloud.Aggregator, username, password string, remoteApp *models.RemoteApp) (string, error) {
	privateIp, err := c.ParseRawInstanceIP(sessionDetail.Instance.Zone, sessionDetail.Instance.InstanceData)
	if err != nil {
		return "", errors.Wrap(err, "remoteapp.GenerateRemoteAppURL")
	}

	zoneOpts, err := c.GetZoneOpts(sessionDetail.Instance.Zone)
	if err != nil {
		return "", fmt.Errorf("get zone opts failed, zone: %s, err: %w", sessionDetail.Instance.Zone, err)
	}

	tokenArgs := guacamole.ConnectArgsInToken{
		GuacadAddr:    zoneOpts.GuacdAddress,
		AssetProtocol: "rdp",
		AssetHost:     privateIp,
		AssetPort:     "3389",
		AssetUser:     username,
		AssetPassword: password,
		DisableGfx:    "false",
	}

	if remoteApp != nil {
		remoteAppArgs, err := renderRemoteAppArgsTemplate(remoteApp.Args)
		if err != nil {
			return "", fmt.Errorf("render remote app args template failed, %w", err)
		}

		tokenArgs.AssetRemoteApp = remoteAppName(remoteApp.Name)
		tokenArgs.AssetRemoteAppDir = remoteApp.Dir
		tokenArgs.AssetRemoteAppArgs = remoteAppArgs
	}

	tokenArgsJSON, err := jsoniter.MarshalToString(&tokenArgs)
	if err != nil {
		return "", fmt.Errorf("marshal json to string failed, %w", err)
	}

	token, err := jwt.Encode(tokenArgsJSON)
	if err != nil {
		return "", fmt.Errorf("encode token failed, raw string: %s, err: %w", tokenArgsJSON, err)
	}

	return accessURL(zoneOpts.AccessOrigin, token), nil
}

func accessURL(gateway, token string) string {
	return fmt.Sprintf("%s/?token=%s", gateway, token)
}

func remoteAppName(name string) string {
	if name == "" {
		return ""
	}
	return fmt.Sprintf("||%s", name)
}

func GetDefaultUsernameByPlatform(platform models.Platform) string {
	switch platform {
	case models.Windows:
		return "administrator"
	case models.Linux:
		return "ecpuser"
	default:
		return ""
	}
}
