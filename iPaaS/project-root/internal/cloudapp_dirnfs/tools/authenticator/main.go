package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_dirnfs/model"
)

var options struct {
	Server   string
	Username string
	Password string
}

func Fatalf(format string, args ...interface{}) {
	fmt.Printf(strings.TrimRight(format, "\n")+"\n", args...)
	os.Exit(127)
}

func init() {
	if options.Server = os.Getenv("REGISTER_CENTER"); len(options.Server) == 0 {
		Fatalf("fatal: invalid environment")
	}
	if options.Username = os.Getenv("SHARE_USERNAME"); len(options.Username) == 0 {
		Fatalf("fatal: invalid environment")
	}
	if options.Password = os.Getenv("SHARE_PASSWORD"); len(options.Password) == 0 {
		Fatalf("fatal: invalid environment")
	}
}

type globalCmd struct {
	cmd *cobra.Command

	subPath string
}

func main() {
	gc := &globalCmd{
		cmd: &cobra.Command{
			Use:  "yuansuan user storage authenticator",
			Long: "yuansuan user storage authenticator",
		},
	}

	gc.cmd.Flags().StringVar(&gc.subPath, "sub-path", "", "user sub path")

	gc.cmd.RunE = gc.run

	if err := gc.cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func (gc *globalCmd) run(_ *cobra.Command, _ []string) error {
	addUserReq := model.AddUserRequest{
		Username: options.Username,
		Password: options.Password,
		SubPath:  gc.subPath,
	}

	hc := resty.New()
	resp, err := hc.R().
		SetBody(addUserReq).
		SetPathParam("userId", options.Username).
		Post(fmt.Sprintf("http://%s/users/{userId}", options.Server))
	if err != nil {
		return fmt.Errorf("fatal: register failed, %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("fatal: authenticate failed")
	}

	return nil
}
