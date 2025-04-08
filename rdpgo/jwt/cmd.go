package jwt

import (
	"errors"
	"fmt"
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"

	"github.com/yuansuan/ticp/rdpgo/guacamole"
)

func NewEncodeCmd() *cobra.Command {
	var argsFile string
	var arg guacamole.ConnectArgsInToken

	cmd := &cobra.Command{
		Use:   "encode",
		Short: "use --args-file or other flags to set args",
		Long:  "use --args-file or other flags to set args",

		RunE: func(_ *cobra.Command, _ []string) error {
			var err error
			if argsFile == "" {
				return errors.New("--args-file cannot be empty")
			}

			content, err := ioutil.ReadFile(argsFile)
			if err != nil {
				return err
			}

			if err = jsoniter.Unmarshal(content, &arg); err != nil {
				return err
			}

			rawData, err := jsoniter.MarshalToString(&arg)
			if err != nil {
				return err
			}

			token, err := Encode(rawData)
			if err != nil {
				return fmt.Errorf("encode failed, %w", err)
			}

			fmt.Println(token)
			return nil
		},
	}

	initEncodeFlag(cmd, &argsFile)

	return cmd
}

func initEncodeFlag(cmd *cobra.Command, argsFile *string) {
	cmd.Flags().StringVar(argsFile, "args-file", "", "args file in json format")
}

func NewDecodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "decode <token>",
		Short: "decode jwt token",
		Long:  "decode jwt token",
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("should only have 1 argument which is token")
			}

			decodedData, err := Decode(args[0])
			if err != nil {
				return fmt.Errorf("decode token failed, %w", err)
			}

			fmt.Println(decodedData)

			return nil
		},
	}
}
