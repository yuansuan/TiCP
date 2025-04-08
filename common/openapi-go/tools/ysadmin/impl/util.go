package impl

import (
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
)

func helpRun(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func callWithDumpReqResp(req interface{}, f func() (interface{}, error)) error {
	// log req
	fmt.Println(structToString(req))
	fmt.Println()
	resp, err := f()
	if err != nil {
		return err
	}
	// log resp
	fmt.Println(structToString(resp))
	fmt.Println()
	return nil
}

func structToString(v interface{}) string {
	s, err := jsoniter.MarshalIndent(v, "", "  ")
	if err != nil {
		return spew.Sdump(v)
	}
	return string(s)
}

func checkReqJsonFileAndRead(file string) ([]byte, error) {
	if file == "" {
		return nil, fmt.Errorf("request json file cannot be empty")
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read request json file [%s] failed, %w", file, err)
	}
	return data, nil
}

func checkFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func StringToUnixTimestamp(str string) (int64, error) {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, str)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}
