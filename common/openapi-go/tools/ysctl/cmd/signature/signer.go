package signature

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	openys "github.com/yuansuan/ticp/common/openapi-go"
	"github.com/yuansuan/ticp/common/openapi-go/credential"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
	"github.com/yuansuan/ticp/common/openapi-go/utils/signer"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type GenSignerOptions struct {
	hc *openys.Client
	clientcmd.IOStreams
	file      string
	keyId     string
	secret    string
	urlParams string
	timestamp int64
	body      string
}

var genSignerExample = templates.Examples(`
    # gen Signature
    ysctl signer gen --key_id abcd --key_secret efg --url_params='FileSize=101&Offset=53' --timestamp=1718710952 --http_body='{"size":"2259426912"}'  --http_body_file=file1
`)

func NewGenSignerOptions(ioStreams clientcmd.IOStreams) *GenSignerOptions {
	return &GenSignerOptions{
		IOStreams: ioStreams,
	}
}

func NewGenSigner(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewGenSignerOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "genSignature",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "gen http signature",
		TraverseChildren:      true,
		Long:                  "gen http signature",
		Example:               genSignerExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}
	cmd.Flags().StringVar(&o.keyId, "key_id", "", "access key id")
	cmd.Flags().StringVar(&o.secret, "key_secret", "", "access key secret")
	cmd.Flags().Int64Var(&o.timestamp, "timestamp", 0, "timestamp")
	cmd.Flags().StringVar(&o.urlParams, "url_params", "", "params in url, format like: FileSize=101&Offset=53&FilePath=/abc")
	cmd.Flags().StringVar(&o.body, "http_body", "", "body in http")
	cmd.Flags().StringVar(&o.file, "http_body_file", "", "http body data in file")
	return cmd
}

func (o *GenSignerOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	return nil
}

type ReaderAndCloser struct {
	*strings.Reader
}

func (r *ReaderAndCloser) Close() error {
	return nil
}

func (o *GenSignerOptions) Run(args []string) error {
	if o.keyId == "" {
		return errors.New("empty key_id")
	}
	if o.secret == "" {
		return errors.New("empty secret")
	}
	if o.timestamp == 0 {
		return errors.New("empty timestamp")
	}
	if err := o.printSig("AppKey"); err != nil {
		return err
	}
	return o.printSig("AccessKeyId")
}

func (o *GenSignerOptions) printSig(keyIdName string) error {
	rawUrl := fmt.Sprintf("http://127.0.0.1/?%s=%s&Timestamp=%d", keyIdName, o.keyId, o.timestamp)
	sigGenerator, _ := signer.NewSigner(credential.NewCredential(o.keyId, o.secret))

	var bodyReader io.ReadCloser
	var err error
	bodyReader = &ReaderAndCloser{
		Reader: strings.NewReader(o.body),
	}

	if o.body == "" && o.file != "" {
		if bodyReader, err = getFileReader(o.file); err != nil {
			return err
		}
	}

	if o.urlParams != "" {
		rawUrl = fmt.Sprintf("%s&%s", rawUrl, o.urlParams)
	}
	reqUrl, err := url.Parse(rawUrl)
	if err != nil {
		fmt.Printf("parser url fail: %s, raw url is: %s, 请检查url_params\n", err.Error(), rawUrl)
		return errors.New("parser url fail")
	}
	req := &http.Request{
		Body:   bodyReader,
		Header: map[string][]string{"Content-Type": []string{"application/json"}},
		URL:    reqUrl,
	}
	sig, err := sigGenerator.SignHttp(req)
	if err != nil {
		fmt.Printf("gen error: %s\n", err.Error())
	} else {
		fmt.Printf("http request url: %s\n", rawUrl)
		fmt.Printf("http body: %s\n", o.body)
		fmt.Printf("http body file: %s\n", o.file)
		fmt.Printf("signature raw str(without secret): %s\n", sig.SourceStr)
		fmt.Printf("signature is: %s \n", sig.Signature)
	}
	return nil
}

func getFileReader(path string) (io.ReadCloser, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("open file fail: %s", err.Error()))
	}
	return f, nil
}
