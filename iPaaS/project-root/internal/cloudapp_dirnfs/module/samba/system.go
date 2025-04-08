package samba

import (
	"context"
	sysuser "os/user"
	"strings"

	"github.com/pkg/errors"
)

// lookupSysUser 检查系统用户是否存在
func (smb *Samba) lookupSysUser(username string) (bool, error) {
	u, err := sysuser.Lookup(username)
	if err != nil && !strings.Contains(err.Error(), "unknown user") {
		return false, errors.Wrap(err, "lookup")
	}

	smb.log.Infof("lookup system user for %q: %+v", username, u)
	return u != nil, nil
}

// addSysUser 给系统添加一个系统用户
func (smb *Samba) addSysUser(ctx context.Context, username, home string) error {
	ga, err := smb.groupAdd.Run(ctx, username)
	if err != nil {
		smb.log.Errorf("group add fail, error: %s, username: %s, home: %s", err.Error(), username, home)
		return errors.Wrap(err, "groupadd")
	}
	if err = ga.Wait(); err != nil {
		smb.log.Errorf("group add fail, error: %s, username: %s, home: %s", err.Error(), username, home)
		return errors.Wrap(err, "groupadd")
	}

	ua, err := smb.userAdd.Run(ctx, username, home)
	if err != nil {
		smb.log.Errorf("user add fail, error: %s, username: %s, home: %s", err.Error(), username, home)
		return errors.Wrap(err, "useradd")
	}
	if err = ua.Wait(); err != nil {
		smb.log.Errorf("user add fail, error: %s, username: %s, home: %s", err.Error(), username, home)
		return errors.Wrap(err, "useradd")
	}

	return nil
}

// dellSysUser 删除一个系统用户和用户组
func (smb *Samba) delSysUser(ctx context.Context, username string) error {
	run, err := smb.groupDel.Run(ctx, username)
	if err != nil {
		smb.log.Errorf("Delete system group fail, err: %s, username: %s", err.Error(), username)
		return errors.Wrap(err, "groupdel")
	}
	if err = run.Wait(); err != nil {
		smb.log.Errorf("Delete system group fail, err: %s, username: %s", err.Error(), username)
		errors.Wrap(err, "groupdelrun")
	}
	smb.log.Infof("delete sysgroup successfully, username: %s", username)
	run, err = smb.userDel.Run(ctx, username)
	if err != nil {
		smb.log.Errorf("Delete system user fail, err: %s, username: %s", err.Error(), username)
		return errors.Wrap(err, "userdel")
	}
	if err = run.Wait(); err != nil {
		smb.log.Errorf("Delete system user fail, err: %s, username: %s", err.Error(), username)
		return errors.Wrap(err, "userdelrun")
	}
	smb.log.Infof("delete sysuser successfully, username: %s", username)
	return nil
}
