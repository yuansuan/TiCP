package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	openys "github.com/yuansuan/ticp/common/openapi-go"
	"github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/upload"
	cmdutil "github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/cmd/util"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/templates"

	"os"
	"path/filepath"
)

type StorageUploadOptions struct {
	SourcePath   string
	DestPath     string
	IsDir        bool
	Retry        int
	FileCount    int
	CompressType string
	hc           *openys.Client
	clientcmd.IOStreams
}

var uploadExample = templates.Examples(`
    # Upload files
    ysctl storage upload --sourcePath=/tmp/test.txt --destPath=/4TiSsZonTa3/test.txt
	
	# Upload dir
	ysctl storage upload --sourcePath=/tmp/test --destPath=/4TiSsZonTa3/test --isDir=true --retry=10 --fileCount=500 --compressType=GZIP --storage_endpoint=https://root-jn-cloud-test.yuansuan.com:34606 --access_id="fadsfe" --access_secret="d23d" --proxy=""
`)

func NewStorageUploadOptions(ioStreams clientcmd.IOStreams) *StorageUploadOptions {
	return &StorageUploadOptions{
		IOStreams: ioStreams,
	}
}

func NewStorageUpload(f *cmdutil.ApiClient, ioStreams clientcmd.IOStreams) *cobra.Command {
	o := NewStorageUploadOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "upload",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Upload files",
		TraverseChildren:      true,
		Long:                  "Upload files",
		Example:               uploadExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}
	cmd.Flags().StringVar(&o.SourcePath, "sourcePath", "", "Source path")
	cmd.Flags().StringVar(&o.DestPath, "destPath", "", "Dest path")
	cmd.Flags().BoolVar(&o.IsDir, "isDir", false, "Is dir")
	cmd.Flags().IntVar(&o.Retry, "retry", 10, "Retry times")
	cmd.Flags().IntVar(&o.FileCount, "fileSize", 1000, "FileSize")
	cmd.Flags().StringVar(&o.CompressType, "compressType", "GZIP", "Compress type")
	return cmd
}

func (o *StorageUploadOptions) Complete(f *cmdutil.ApiClient, cmd *cobra.Command, args []string) error {
	var err error
	o.hc, _, err = f.StorageClient()
	if err != nil {
		return err
	}
	if o.SourcePath == "" {
		return cmdutil.UsageErrorf(cmd, "Must specify --sourcePath")
	}
	if o.DestPath == "" {
		return cmdutil.UsageErrorf(cmd, "Must specify --destPath")
	}
	if o.Retry < 0 || o.Retry > 10 {
		return cmdutil.UsageErrorf(cmd, "Must specify suitable between 0 and 10 --retry")
	}
	if o.FileCount < 0 || o.FileCount > 1000 {
		return cmdutil.UsageErrorf(cmd, "Must specify suitable between 0 and 1000 --files")
	}
	return nil
}

func (o *StorageUploadOptions) Run(args []string) error {
	fmt.Fprintf(o.Out, "Current config: SourcePath: %s, DestPath: %s, IsDir: %v, Default Worker: %d, Upload Files Limit: %d,Compress Type: %s \n", o.SourcePath, o.DestPath, o.IsDir, 10, o.FileCount, o.CompressType)
	if !o.IsDir {
		err := upload.Upload(o.SourcePath, o.DestPath, o.CompressType, o.hc, nil, nil, nil)
		if err != nil {
			return err
		}
		return nil
	}
	return o.UploadDir(o.SourcePath, o.DestPath)
}

func (o *StorageUploadOptions) UploadDir(sourcePath, destPath string) error {
	_, err := o.hc.Storage.Mkdir(
		o.hc.Storage.Mkdir.Path(destPath),
		o.hc.Storage.Mkdir.IgnoreExist(true),
	)
	if err != nil {
		return err
	}
	err = o.uploadDir(sourcePath, destPath)
	return err
}

func (o *StorageUploadOptions) uploadDir(sourcePath, destPath string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel for files to be uploaded
	files := make(chan string, 1000)

	lsErr := make(chan error)

	// Start a goroutine to list all files in the directory and feed them to the channel
	go func() {
		defer close(files)
		if err := walkDir(ctx, sourcePath, files); err != nil {
			lsErr <- err
		}
	}()

	// Check if there is any error when listing files
	select {
	case err := <-lsErr:
		return err
	case <-time.After(1 * time.Second):
	}

	// Fan-out: create multiple workers
	workerNum := 10
	workerChans := make([]<-chan error, workerNum)
	for i := 0; i < workerNum; i++ {
		workerChans[i] = o.uploadWorker(ctx, files, destPath)
	}

	// Fan-in: merge result from all workers into a single channel
	errChan := merge(ctx, workerChans...)

	// Wait for all workers to finish and return the first error, if any
	for err := range errChan {
		if err != nil {
			cancel()
			return err
		}
	}
	return nil
}

func (o *StorageUploadOptions) uploadWorker(ctx context.Context, files <-chan string, destPath string) <-chan error {
	errChan := make(chan error)
	go func() {
		defer close(errChan)
		for file := range files {
			// Get the relative path of the file
			rel, err := filepath.Rel(o.SourcePath, file)
			if err != nil {
				errChan <- err
				return
			}
			err = o.retryUploadFile(ctx, file, filepath.Join(destPath, rel), o.Retry)
			if err != nil {
				errChan <- err
				return
			}
		}
	}()
	return errChan
}

func walkDir(ctx context.Context, dir string, files chan<- string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files <- path
		}
		return nil
	})
}

func (o *StorageUploadOptions) retryUploadFile(ctx context.Context, file string, destPath string, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = upload.Upload(file, destPath, o.CompressType, o.hc, nil, nil, nil)
		if err == nil {
			return nil
		}
		wait := exponentialBackoff(i)
		fmt.Printf("Upload %s failed, retry after %s\n", file, wait)
		time.Sleep(wait)
	}
	return err
}
