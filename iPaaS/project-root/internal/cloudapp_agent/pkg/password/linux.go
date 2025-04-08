//go:build linux

package password

import (
	"bytes"
	"fmt"
	"os/exec"
	"text/template"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_agent/pkg/environment"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_agent/pkg/log"
)

const (
	defaultLoginUser = "ecpuser"
)

const resetPasswordScript = `#!/bin/bash
passwd {{ .Username }} <<EOF
{{ .Password }}
{{ .Password }}
EOF`

type args struct {
	Username string
	Password string
}

func Reset(customEnv *environment.CustomEnv) error {
	var err error
	if err = resetRootPassword(customEnv); err != nil {
		return fmt.Errorf("reset root password failed, %w", err)
	}

	if err = resetEcpUserPassword(customEnv); err != nil {
		return fmt.Errorf("reset ecpuser password failed, %w", err)
	}

	return nil
}

func resetRootPassword(customEnv *environment.CustomEnv) error {
	loginPass := customEnv.Get(loginPasswordEnvKey)
	if loginPass == "" {
		log.Warnf("%s is empty from custom env, determined not to reset password.", loginPasswordEnvKey)
		return nil
	}

	log.Infof("get %s from custom env success, value: %s", loginPasswordEnvKey, loginPass)

	return resetPassword("root", loginPass)
}

func resetEcpUserPassword(customEnv *environment.CustomEnv) error {
	loginPass := customEnv.Get(loginPasswordEnvKey)
	if loginPass == "" {
		log.Warnf("%s is empty from custom env, determined not to reset password.", loginPasswordEnvKey)
		return nil
	}

	log.Infof("get %s from custom env success, value: %s", loginPasswordEnvKey, loginPass)

	return resetPassword(defaultLoginUser, loginPass)
}

func resetPassword(username, password string) error {
	arg := args{
		Username: username,
		Password: password,
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

	cmd := exec.Command("bash", "-c", buf.String())

	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("execute reset password cmd failed, %w", err)
		log.Error(err)
		return err
	}

	return nil
}
