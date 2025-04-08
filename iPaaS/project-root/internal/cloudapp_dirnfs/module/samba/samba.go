package samba

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	exec2 "os/exec"
	"time"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_dirnfs/module/jsonmap"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_dirnfs/module/xos/exec"
)

const (
	_SMBNameServerDaemon = "/usr/sbin/nmbd"
	_SMBSambaDaemon      = "/usr/sbin/smbd"

	_SMBControl = "/usr/bin/smbcontrol"
	_SMBPasswd  = "/usr/bin/smbpasswd"
	_SMBStatus  = "/usr/bin/smbstatus"

	_SysGroupAdd = "/usr/sbin/addgroup"
	_SysUserAdd  = "/usr/sbin/adduser"

	_SysGroupDel        = "/usr/sbin/delgroup"
	_SysUserDel         = "/usr/sbin/deluser"
	deleteInterval      = 24 * time.Hour
	startMonitorWaiting = 1 * time.Hour
)

// Samba 表示一个SAMBA共享实例
type Samba struct {
	jm  *jsonmap.Map
	log *logging.Logger

	cfg  string
	nmbd *exec.Command
	smbd *exec.Command

	control    *exec.Command
	passwd     *exec.Command
	delSmbUser *exec.Command

	groupAdd *exec.Command
	userAdd  *exec.Command

	groupDel *exec.Command
	userDel  *exec.Command
}

// Start 启动共享服务
func (smb *Samba) Start(ctx context.Context) error {
	proc, err := smb.nmbd.Run(ctx)
	if err != nil {
		return errors.Wrap(err, "nmbd")
	}

	if err = proc.Wait(); err != nil {
		return errors.Wrap(err, "nmbd")
	}

	proc, err = smb.smbd.RunWithOptions(ctx, []interface{}{smb.cfg},
		exec.WithStdout(os.Stdout), exec.WithStderr(os.Stderr))
	if err != nil {
		return errors.Wrap(err, "smbd")
	}

	go func() {
		time.Sleep(10 * time.Second)
		smb.jm.Visit(func(k string, v interface{}) bool {
			smb.log.Debugf("visiting item %q => %+v", k, v)

			if u, ok := loadUser(v); ok {
				// user already exists in jsonmap, do not persistent again that
				// avoid deadlock when adding the user to system
				if err := smb.AddUser(ctx, u.Username, u.Password, u.Home, false); err != nil {
					smb.log.Infof("reloading user (%q) failed: %s", k, err)
					return true
				}
				smb.log.Infof("reloading user (%+v) succeed", u)
			}

			return false
		})
		smb.log.Infof("Old data load end")
		// 刚启动，smbstatus里被清空，等待一段时间，允许客户端重连
		time.Sleep(startMonitorWaiting)
		smb.log.Infof("start monitor unused user")
		smb.Monitor(ctx)
	}()

	return proc.Wait()
}

// Reload 重载共享存储的所有配置
func (smb *Samba) Reload(ctx context.Context) error {
	proc, err := smb.control.Run(ctx)
	if err != nil {
		return err
	}

	smb.log.Infof("reload the samba services")
	return proc.Wait()
}

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Home     string `json:"home"`
}

func loadUser(v interface{}) (*user, bool) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, false
	}

	u := new(user)
	if err := json.Unmarshal(data, u); err != nil {
		return nil, false
	}
	return u, true
}

// AddUser 添加一个共享用户
func (smb *Samba) AddUser(ctx context.Context, username, password, home string, persistent bool) error {
	exists, err := smb.ensureSysUser(ctx, username, home)
	if err != nil {
		return err
	} else if !exists && persistent {
		smb.log.Debugf("user(%s) not found in system and persistent into disk", username)
		smb.jm.Set(username, &user{Username: username, Password: password, Home: home})
	}

	return smb.addShareUser(ctx, username, password)
}

// DelUser 删除一个用户的持久化数据
func (smb *Samba) DelUser(ctx context.Context, username string) error {
	smb.jm.Del(username)
	defer smb.delSysUser(ctx, username)

	proc, err := smb.delSmbUser.Run(ctx, username)
	if err != nil {
		smb.log.Errorf("Delete smb user fail, err: %s, username: %s", err.Error(), username)
		return err
	}

	if err = proc.Wait(); err != nil {
		smb.log.Errorf("Delete smb user fail, err: %s, username: %s", err.Error(), username)
		return err
	}
	smb.log.Infof("delete smb user successfully, username: %s", username)
	return nil
}

// addShareUser 添加一个共享用户
func (smb *Samba) addShareUser(ctx context.Context, username, password string) error {
	proc, err := smb.passwd.Run(ctx, username)
	if err != nil {
		return err
	}

	smb.log.Infof("add samba user %q with password %q", username, password)
	if err = proc.Start(); err != nil {
		return err
	}

	// input password
	if err = proc.Input(password + "\n"); err != nil {
		return errors.Wrap(err, "input password")
	}

	// password confirmed
	if err = proc.Input(password + "\n"); err != nil {
		return errors.Wrap(err, "confirm password")
	}

	return proc.Wait()
}

// ensureSysUser 添加一个系统用户，如果用户已存在则不添加
func (smb *Samba) ensureSysUser(ctx context.Context, username, home string) (bool, error) {
	if exists, err := smb.lookupSysUser(username); err != nil {
		return false, err
	} else if !exists {
		// 对于相同的用户的多次重复注册只有第一次的HOME会生效，这个是由
		// 添加系统用户时检测用户是否存在实现的
		return false, smb.addSysUser(ctx, username, home)
	}
	// @TODO update password when the user already exists
	return true, nil
}

func (smb *Samba) ConnExisted(ctx context.Context, username string) (bool, error) {
	cmd := exec2.Command("sh", "-c", fmt.Sprintf("%s | grep %s", _SMBStatus, username))
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec2.ExitError); ok && exitError.ExitCode() == 1 {
			return false, nil
		}
		smb.log.Errorf("Exec Cmd fail, Error: %s, username: %s", err.Error(), username)
		return false, err
	} else {
		return true, nil
	}
}

func (smb *Samba) Monitor(ctx context.Context) {
	oldNoExisted := map[string]time.Time{}
	tick := time.Tick(1 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			var toDeleted []string
			smb.jm.Visit(func(k string, v interface{}) bool {
				if existed, err := smb.ConnExisted(ctx, k); err != nil {
					smb.log.Errorf("Run conn existed cmd fail, error: %s, username: %s", err.Error(), k)
				} else if !existed {
					if oldTime, ok := oldNoExisted[k]; !ok {
						oldNoExisted[k] = time.Now()
						smb.log.Infof("Find username connection not existed, will delete after %v, username: %s", deleteInterval, k)
					} else if time.Now().Sub(oldTime) > deleteInterval {
						//超过 DeleteInterval时间不存在连接的用户将被删除
						toDeleted = append(toDeleted, k)
						smb.log.Infof("Find username connection not existed, will delete immediately, username: %s", k)
					}
				}
				return false
			})
			if len(toDeleted) == 0 {
				smb.log.Infof("not unused sysgroup and sysuser to delete")
			}
			for _, k := range toDeleted {
				if err := smb.DelUser(ctx, k); err != nil {
					smb.log.Errorf("Delete user fail, error: %s, username: %s", err.Error(), k)
				}
				delete(oldNoExisted, k)
			}
		}
	}
}

// New 创建并重新加载配置文件
func New(samba, users string) (*Samba, error) {
	// all = smbd + nmbd
	control, err := exec.New(_SMBControl, exec.String("all"), exec.String("reload-config"))
	if err != nil {
		return nil, err
	}

	// -s => use stdin for password prompt, -a => add user
	passwd, err := exec.New(_SMBPasswd, exec.String("-s"), exec.String("-a"), exec.PlaceHolder(0))
	if err != nil {
		return nil, err
	}

	delSmbUser, err := exec.New(_SMBPasswd, exec.String("-x"), exec.PlaceHolder(0))
	if err != nil {
		return nil, err
	}

	// NetBIOS name server to provide NetBIOS over IP naming services to clients
	nmbd, err := exec.New(_SMBNameServerDaemon)
	if err != nil {
		return nil, err
	}

	// server to provide SMB/CIFS services to clients
	smbd, err := exec.New(_SMBSambaDaemon,
		exec.String("-S"),                 // Log to stdout
		exec.String("-F"),                 // Run daemon in foreground
		exec.String("--no-process-group"), // Don't create a new process group
		exec.String("-s"),                 // Use alternate configuration file
		exec.PlaceHolder(0),
	)
	if err != nil {
		return nil, err
	}

	// -S => create a system group
	groupAdd, err := exec.New(_SysGroupAdd, exec.String("-S"), exec.PlaceHolder(0))
	if err != nil {
		return nil, err
	}

	groupDel, err := exec.New(_SysGroupDel, exec.PlaceHolder(0))
	if err != nil {
		return nil, err
	}

	// -S => create a system users
	// -D => don't assign a password
	userAdd, err := exec.New(_SysUserAdd,
		exec.String("-h"),
		exec.PlaceHolder(1),
		exec.String("-G"),
		exec.String("root"),
		exec.String("-SD"),
		exec.PlaceHolder(0))
	if err != nil {
		return nil, err
	}

	userDel, err := exec.New(_SysUserDel, exec.PlaceHolder(0))
	if err != nil {
		return nil, err
	}

	// loading users from json file
	m, err := jsonmap.New(users)
	if err != nil {
		return nil, err
	}

	return &Samba{
		jm:  m,
		log: logging.Default(),

		cfg:  samba,
		nmbd: nmbd,
		smbd: smbd,

		control:    control,
		passwd:     passwd,
		delSmbUser: delSmbUser,

		groupAdd: groupAdd,
		userAdd:  userAdd,

		groupDel: groupDel,
		userDel:  userDel,
	}, nil
}
