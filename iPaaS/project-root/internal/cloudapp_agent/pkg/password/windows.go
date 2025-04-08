//go:build windows

package password

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"text/template"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_agent/pkg/environment"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_agent/pkg/log"
)

// TODO check if change password success
const resetPasswordScript = `
$password = ConvertTo-SecureString '{{ .LoginPassword }}' -AsPlainText -Force; 
$userAccount = GetQuota-LocalUser -Name '{{ .LoginUsername }}'; 
$userAccount | Set-LocalUser -Password $password -PasswordNeverExpires $true;
`

type args struct {
	LoginUsername string
	LoginPassword string
}

func Reset(customEnv *environment.CustomEnv) error {
	var err error
	if err = resetAdministratorPassword(customEnv); err != nil {
		return fmt.Errorf("reset adminstrator password failed, %w", err)
	}

	if err = resetRemoteAppsPassword(customEnv); err != nil {
		return fmt.Errorf("reset remoteApps password failed, %w", err)
	}

	return nil
}

func resetAdministratorPassword(customEnv *environment.CustomEnv) error {
	loginPass := customEnv.Get(loginPasswordEnvKey)
	if loginPass == "" {
		log.Warnf("%s is empty from custom env, determined not to reset password.", loginPasswordEnvKey)
		return nil
	}
	log.Infof("get %s from custom env success, value: %s", loginPasswordEnvKey, loginPass)

	return resetPassword("Administrator", loginPass)
}

// 示例配置
// REMOTE_APP_USER_PASSES=starccm@pass1,cmd@pass2
func resetRemoteAppsPassword(customEnv *environment.CustomEnv) error {
	remoteAppUserPassesV := customEnv.Get(remoteAppUserPassesEnvKey)

	for _, remoteAppEnv := range strings.Split(remoteAppUserPassesV, ",") {
		userPass := strings.Split(remoteAppEnv, userPassSplitCharacter)
		if len(userPass) != 2 {
			log.Warnf("parse %s to username@password failed", remoteAppEnv)
			continue
		}

		username, password := userPass[0], userPass[1]

		if err := resetPassword(username, password); err != nil {
			err = fmt.Errorf("reset password failed username [%s] password [%s], %w", username, password, err)
			log.Error(err)
			return err
		}
		log.Infof("reset username [%s] to password [%s] success", username, password)
	}

	return nil
}

func resetPassword(username, password string) error {
	arg := args{
		LoginUsername: username,
		LoginPassword: password,
	}

	tpl, err := template.New("reset_password").Parse(resetPasswordScript)
	if err != nil {
		err = fmt.Errorf("new template failed, %w", err)
		log.Error(err)
		return err
	}

	buf := &bytes.Buffer{}
	if err = tpl.Execute(buf, &arg); err != nil {
		err = fmt.Errorf("template execute failed, %w", err)
		log.Error(err)
		return err
	}

	cmd := exec.Command("powershell", "-Command", buf.String())

	// 目前认为不报错则修改密码成功
	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("execute reset password cmd failed, %w", err)
		log.Error(err)
		return err
	}

	return nil
}
