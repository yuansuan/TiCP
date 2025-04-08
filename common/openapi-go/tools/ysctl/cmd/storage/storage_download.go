package storage

import (
	"context"
	"math"
	"math/rand"
	"os"
	"path"
	"strings"
	"sync"

	openys "github.com/yuansuan/ticp/common/openapi-go"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/spf13/cobra"

	"fmt"
	"time"
)

type StorageDownloadOptions struct {
	Path       string
	OutputPath string
	Retry      int
	FileCount  int
	hc         *openys.Client
	clientcmd.IOStreams
}

var downloadExample = templates.Examples(`
    # Download file
    ysctl storage download --path=/4TiSsZonTa3/1.txt --outputPath=1.txt

	# Download dir
	ysctl storage download --path=/4TiSsZonTa3 --outputPath=/test --fileCount=500 --retry=10 --storage_endpoint=https://root-jn-cloud-test.yuansuan.com:34606 --access_id="fadsfe" --access_secret="d23d" --proxy=""
`)

func NewStorageDownloadOptions(ioStreams clientcmd.IOStreams) *StorageDownloadOptions {
	return &StorageDownloadOptions{
		IOStreams: ioStreams,
	}
}

func NewStorageDownload(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewStorageDownloadOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "download",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Download file",
		TraverseChildren:      true,
		Long:                  "Download file",
		Example:               downloadExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
		},
		SuggestFor: []string{},
	}
	cmd.Flags().StringVar(&o.Path, "path", "", "Path")
	cmd.Flags().StringVar(&o.OutputPath, "outputPath", "", "Output path")
	cmd.Flags().IntVar(&o.Retry, "retry", 10, "Retry times")
	cmd.Flags().IntVar(&o.FileCount, "fileSize", 1000, "File size")
	return cmd
}

func (o *StorageDownloadOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, _, err = f.StorageClient()
	if err != nil {
		return err
	}
	if o.Path == "" {
		return cmdutil.UsageErrorf(cmd, "Must specify --path")

	}
	if o.OutputPath == "" {
		return cmdutil.UsageErrorf(cmd, "Must specify --outputPath")
	}
	if o.FileCount < 0 || o.FileCount > 1000 {
		return cmdutil.UsageErrorf(cmd, "Must specify suitable between 0 and 1000 --fileSize")
	}
	if o.Retry < 0 || o.Retry > 10 {
		return cmdutil.UsageErrorf(cmd, "Must specify suitable between 0 and 10 --retry")
	}
	return nil
}

func (o *StorageDownloadOptions) createDirWithRetry(dst string, maxRetries int) error {
	dir := path.Dir(dst)
	for i := 1; i <= maxRetries; i++ {
		if _, err := os.Stat(dir); err == nil {
			return nil
		}
		err := os.MkdirAll(dir, 0755)
		if err == nil {
			return nil
		}
		wait := exponentialBackoff(i)
		fmt.Printf("CreateDirFail, Path: %s, Error: %v. Retrying in %v...\n", dir, err, wait)
		time.Sleep(wait)
	}
	return fmt.Errorf("failed to create directory %s after %d attempts", dir, maxRetries)
}

func merge(ctx context.Context, chans ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error)

	output := func(c <-chan error) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}

	wg.Add(len(chans))
	for _, c := range chans {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func exponentialBackoff(i int) time.Duration {
	wait := time.Duration(int64(math.Pow(2, float64(i)))) * time.Second
	wait += time.Duration(rand.Intn(1000)) * time.Millisecond // add some jitter
	return wait
}

func trimPrefix(path string) string {
	return strings.TrimPrefix(path, "/")
}

func trim(path string) string {
	if !strings.HasPrefix(path, "/") {
		return path
	}
	// // If the path starts with a slash, remove it
	path = strings.TrimPrefix(path, "/")

	// If the path contains more than one component, return the full path
	components := strings.Split(path, "/")
	if len(components) > 1 {
		return "/" + path
	}

	// If the path contains only one component, return the component without a leading slash
	return components[0]
}

func (o *StorageDownloadOptions) lsPath(ctx context.Context, src string, q chan<- *v20230530.FileInfo) error {
	pageSize := o.FileCount
	var dirs []string
	for i := 0; ; i += pageSize {
		resp, err := o.hc.API.Storage.LsWithPage(
			o.hc.API.Storage.LsWithPage.Path(src),
			o.hc.API.Storage.LsWithPage.PageOffset(int64(i)),
			o.hc.API.Storage.LsWithPage.PageSize(int64(pageSize)),
		)
		if err != nil {
			fmt.Printf("LsWithPageFail, Path: %s, Offset: %d, PageSize: %d", src, i, pageSize)
			return err
		}
		for _, file := range resp.Data.Files {
			file.Name = path.Join(src, file.Name)
			if !file.IsDir {
				select {
				case q <- file:
				case <-ctx.Done():
					break
				}
			} else {
				dirs = append(dirs, file.Name)
			}
		}
		if len(resp.Data.Files) < pageSize {
			break
		}
	}

	for _, d := range dirs {
		if err := o.lsPath(ctx, d, q); err != nil {
			return err
		}
	}
	return nil
}
