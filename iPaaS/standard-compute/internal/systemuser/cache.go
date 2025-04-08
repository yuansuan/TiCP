package systemuser

import (
	"fmt"
	"os/user"
	"strconv"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
)

var users = map[string]*User{}

type User struct {
	Uid int
	Gid int
}

func Init() error {
	commonCfg := config.GetConfig().BackendProvider.SchedulerCommon
	username := commonCfg.SubmitSysUser
	if username != "" {
		var uid, gid int
		if commonCfg.SubmitSysUserUid != 0 && commonCfg.SubmitSysUserGid != 0 {
			uid, gid = commonCfg.SubmitSysUserUid, commonCfg.SubmitSysUserGid
		} else {
			ui, err := user.Lookup(username)
			if err != nil {
				return fmt.Errorf("lookup user [%s] failed, %w", username, err)
			}

			uid, err = strconv.Atoi(ui.Uid)
			if err != nil {
				return fmt.Errorf("parse uid [%s] to int failed, %w", ui.Uid, err)
			}

			gid, err = strconv.Atoi(ui.Gid)
			if err != nil {
				return fmt.Errorf("parse gid [%s] to int failed, %w", ui.Gid, err)
			}
		}

		users[username] = &User{
			Uid: uid,
			Gid: gid,
		}
	}

	return nil
}

func Get(username string) (*User, error) {
	u, exist := users[username]
	if !exist {
		return nil, fmt.Errorf("cannot find user [%s] in cache", username)
	}

	return u, nil
}
