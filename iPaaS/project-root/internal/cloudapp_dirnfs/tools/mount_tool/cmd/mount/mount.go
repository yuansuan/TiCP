package mount

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_dirnfs/tools/mount_tool/cmd/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/hashid"
)

const (
	//MimeType used for register
	MimeType = "application/sharing"
)

// Options Options
type Options struct {
	UserID    string
	ProjectID string
	Password  string
	Server    string
	UserName  string
	Disk      string
}

// NewCmdMount mountctl mount -h will show information
func NewCmdMount() *cobra.Command {
	o := &Options{}
	cmd := &cobra.Command{
		Use:   "mount",
		Short: "mount project directory to windows filesystem",
		Long:  `Linux samba, Windows net use`,
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(o.Validate(cmd))
			util.CheckErr(o.Register())
			util.CheckErr(o.Mount())
		},
	}
	cmd.PersistentFlags().String("p", "", "projectID, not null, like 4pBJ6JWMH6S")
	cmd.PersistentFlags().String("u", "", "userID, not null, like 4Afa4q2V3BS")
	cmd.PersistentFlags().String("s", "", "remote file server url, not null, like 10.0.1.123:8081")
	cmd.PersistentFlags().String("d", "", "disk or dir, not null, like X in windows or /tmp/test in linux")
	return cmd
}

// Validate validate flags
func (o *Options) Validate(cmd *cobra.Command) error {
	log.Println("Validate")
	projectID, err := cmd.Flags().GetString("p")
	if err != nil || len(projectID) == 0 {
		return errors.New("projectID is null")
	}
	log.Printf("projectID: %s", projectID)
	o.ProjectID = projectID

	userID, err := cmd.Flags().GetString("u")
	if err != nil || len(userID) == 0 {
		return errors.New("userID is null")
	}
	log.Printf("userID: %s", userID)
	o.UserID = userID

	server, err := cmd.Flags().GetString("s")
	if err != nil || len(server) == 0 {
		return errors.New("server url is null")
	}
	log.Printf("server url: %s", server)
	o.Server = server

	disk, err := cmd.Flags().GetString("d")
	if err != nil || len(disk) == 0 {
		return errors.New("disk or dir is null")
	}
	log.Printf("disk : %s", disk)
	o.Disk = disk

	log.Println("Validate Successfully")
	return nil
}

// Register register a user
func (o *Options) Register() error {
	log.Println("Encoding options")

	userName, err := hashid.Encode(snowflake.MustParseString(o.ProjectID))
	if err != nil {
		return errors.New("load username")
	}
	o.UserName = userName

	password, err := hashid.Encode(snowflake.MustParseString(o.UserID))
	if err != nil {
		return errors.New("load password")
	}
	o.Password = password

	log.Println("Register")
	resp, err := http.Post(Endpoint(o), MimeType, strings.NewReader(password))
	if err != nil {
		return errors.New("register failed ")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return errors.New("authenticate failed")
	}

	log.Println("Register Successfully")
	return nil
}

// Endpoint  returns the endpoint
func Endpoint(o *Options) string {
	return "http://" + o.Server + "/" + o.UserName
}

// trim trim port
func trim(o *Options) {
	url := o.Server
	log.Printf("before trim %s", url)
	index := strings.Index(url, ":")
	trimmedURL := url[:index]
	o.Server = trimmedURL
	log.Printf("after trim %s", trimmedURL)
}

// Mount to dir
func (o *Options) Mount() error {
	log.Println("start mount")

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		log.Println("current is windows")
		cmd = windowsMount(o)
		if cmd == nil {
			log.Println("encode error")
			return errors.New("encode error")
		}
	} else if runtime.GOOS == "linux" {
		log.Println("current is linux ")
		cmd = linuxMount(o)
		if cmd == nil {
			// 文件夹已存在
			log.Printf("dir already exists")
			return errors.New("dir already exists")
		}
	}

	log.Println("start to execute")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New("Failed to mount " + string(out))
	}
	log.Println("mount complete")
	return nil
}

func windowsMount(o *Options) *exec.Cmd {

	trim(o)
	prefix := "YS"
	disk := o.Disk + ":"

	// 设置编码
	// chcp 65001
	err := exec.Command("chcp", "65001").Run()
	if err != nil {
		return nil
	}

	path := string(os.PathSeparator) + string(os.PathSeparator) + o.Server + string(os.PathSeparator) + prefix + o.UserName
	userName := "/" + "user:" + prefix + o.UserName
	password := o.Password

	// net use X: \\10.0.64.187\YSfa0ff67b22753630966652a6576bb50a /user:YSfa0ff67b22753630966652a6576bb50a db602b590ffbcb584b75afd6c12f63c8
	cmd := exec.Command("net", "use", disk, path, userName, password)
	return cmd
}

func linuxMount(o *Options) *exec.Cmd {
	// 如果文件不存在，则创建
	dir := o.Disk
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		return nil
	}

	trim(o)
	// -o param
	prefix := "YS"
	param := "username=" + prefix + o.UserName + "," + "password=" + o.Password
	path := string(os.PathSeparator) + string(os.PathSeparator) + o.Server + string(os.PathSeparator) + prefix + o.UserName

	// mount -t cifs -o username=YS11e8af3e60e6066c003830697dbdb706,password=09dc7ee850cacfb6999839ca55996c88 //172.16.0.3/YS11e8af3e60e6066c003830697dbdb706 /tmp/test
	cmd := exec.Command("mount", "-t", "cifs", "-o", param, path, dir)
	return cmd
}
