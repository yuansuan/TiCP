package unmount

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_dirnfs/tools/mount_tool/cmd/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/hashid"
)

// Options options
type Options struct {
	ProjectID string
	Server    string
	UserName  string
	Disk      string
}

// NewCmdUnMount cspctl unmount -h will show information
func NewCmdUnMount() *cobra.Command {
	o := &Options{}
	cmd := &cobra.Command{
		Use:   "unmount",
		Short: "unmount project directory from windows filesystem",
		Long:  `Linux samba, Windows net use`,
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(o.Validate(cmd))
			util.CheckErr(o.UnRegister())
			util.CheckErr(o.UnMount())
		},
	}
	cmd.PersistentFlags().String("p", "", "projectID, not null, like 4pBJ6JWMH6S")
	cmd.PersistentFlags().String("s", "", "remote file server url, not null, like 10.0.1.123:8081")
	cmd.PersistentFlags().String("d", "", "disk, not null, like X")
	return cmd
}

// Validate valid flags
func (o *Options) Validate(cmd *cobra.Command) error {
	log.Println("Validate")
	projectID, err := cmd.Flags().GetString("p")
	if err != nil || len(projectID) == 0 {
		return errors.New("projectID is null")
	}
	log.Printf("projectID: %s", projectID)
	o.ProjectID = projectID

	server, err := cmd.Flags().GetString("s")
	if err != nil || len(server) == 0 {
		return errors.New("server url is null")
	}
	log.Printf("server url: %s", server)
	o.Server = server

	disk, err := cmd.Flags().GetString("d")
	if err != nil || len(disk) == 0 {
		return errors.New("disk is null")
	}
	o.Disk = disk

	log.Println("Validate Successfully")
	return nil
}

// UnRegister unregister
func (o *Options) UnRegister() error {
	log.Println("Encoding options")
	userName, err := hashid.Encode(snowflake.MustParseString(o.ProjectID))
	if err != nil {
		return errors.New("load username")
	}
	o.UserName = userName

	log.Println("UnRegister")
	req, _ := http.NewRequest(http.MethodDelete, Endpoint(o), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("unregister failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println(resp.StatusCode)
		return errors.New("authenticate failed")
	}
	log.Println("Register Successfully")
	return nil
}

// Endpoint endpoint
func Endpoint(o *Options) string {
	return "http://" + o.Server + "/" + o.UserName
}

// UnMount unmount from server
func (o *Options) UnMount() error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = windowsUnMount(o)
		if cmd == nil {
			log.Println("windows encode error")
			return errors.New("encode error")
		}
	} else if runtime.GOOS == "linux" {
		cmd = linuxUnMount(o)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New("Failed to unmount " + string(out))
	}

	// remove dir after umount
	if err = os.RemoveAll(o.Disk); err != nil {
		return errors.New("Failed to remove " + o.Disk)
	}

	log.Println("unmount complete")
	return nil
}

func windowsUnMount(o *Options) *exec.Cmd {
	log.Println("start unmount")

	disk := o.Disk + ":"
	command := "/delete"

	// 设置编码
	// chcp 65001
	exec.Command("chcp", "65001").Run()

	// net use
	cmd := exec.Command("net", "use", disk, command)
	return cmd
}

func linuxUnMount(o *Options) *exec.Cmd {
	log.Println("start unmount")
	dir := o.Disk
	// umount /xx/xx
	cmd := exec.Command("umount", dir)
	return cmd
}
