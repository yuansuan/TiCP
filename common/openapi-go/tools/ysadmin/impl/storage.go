package impl

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	openapi "github.com/yuansuan/ticp/common/openapi-go"
	"github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/upload"
)

const (
	// 单位为MB
	dataSize         = 200
	testFolder       = "speed-test"
	downloadFileName = "randomfile3G"
	localEndPoint    = "http://localhost:8899"
)

type StorageOptions struct {
	UserId         string
	Zone           string
	Type           string
	Offset         int64
	Limit          int64
	Filter         string
	BasePath       string
	FileName       string
	FileTypes      string
	OperationTypes string
	BeginTime      string
	EndTime        string
	Overwrite      bool
	Size           int64
	Data           []byte
	Downloaded     int64
}

func (o *StorageOptions) validateType() {
	if o.Type != HPCType && o.Type != CloudType {
		fmt.Printf("Unsupportted storage type: %s. Need hpc | cloud\n", o.Type)
		os.Exit(1)
	}
}
func (o *StorageOptions) validateUserID() {
	if o.UserId == "" {
		fmt.Println("UserId is empty, please check your config file or pass it by -U/--user_id")
		os.Exit(1)
	}
}

func (o *StorageOptions) complete() {
	if o.UserId == "" {
		if CurrentCfg.StorageYsID == "" {
			fmt.Println("Warn: config storage_ys_id is empty")
		}
		o.UserId = CurrentCfg.StorageYsID
	}
}

func init() {
	RegisterCmd(NewStorageCommand())
}

// NewStorageCommand 创建存储管理命令
func NewStorageCommand() *cobra.Command {
	o := StorageOptions{}
	cmd := &cobra.Command{
		Use:   "storage",
		Short: "存储管理",
		Long:  "存储管理, 可以操作管理云存储和HPC存储\n存储接口中操作的远程存储路径均需要以/{userID}开头\n例如: /4TpFFZDkFWy/test.txt",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newStorageLsCmd(o),
		newStorageMkdirCmd(o),
		newStorageUploadCmd(o),
		newStorageReadAtCmd(o),
		newStorageDownloadCmd(o),
		newStorageBatchDownloadCmd(o),
		newStorageRmCmd(o),
		newStorageMvCmd(o),
		newStorageUploadSpeedCmd(o),
		newStorageDownloadSpeedCmd(o),
		newStorageBatchDownloadSpeedCmd(o),
		newStorageListQuotaCmd(o),
		newStorageUpdateQuotaCmd(o),
		newStorageQuotaTotalCmd(o),
		newStorageListOperationLogCmd(o),
	)

	return cmd
}

func newStorageLsCmd(o StorageOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "ls \"/{userID}/path\"",
		Short: "列出路径下文件",
		Long:  "列出路径下文件, 路径必须以/{userID}开头",
		Args:  cobra.ExactArgs(1),
		Example: `- 列出某用户根目录下的文件
  - ysadmin storage ls /4TpFFZDkFWy
- 列出某用户根目录下的文件, 过滤掉文件名包含test的文件, 指定存储类型为hpc, 指定区域为az-zhigu, 分页显示, 从第0条开始, 每页显示20条
  - ysadmin storage ls /4TpFFZDkFWy -F test -O 0 -L 20 -T hpc -Z az-zhigu`,
	}

	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "az-zhigu", "区域")
	cmd.Flags().StringVarP(&o.Type, "type", "T", "cloud", "存储类型, hpc | cloud")
	cmd.Flags().Int64VarP(&o.Offset, "offset", "O", 0, "offset")
	cmd.Flags().Int64VarP(&o.Limit, "limit", "L", 1000, "limit")
	cmd.Flags().StringVarP(&o.Filter, "filter", "F", "", "用于过滤的正则表达式")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.validateType()
		o.complete()

		path := args[0]
		c := GetStorageClient(o.Zone, o.Type)
		res, err := c.Storage.AdminLsWithPage(
			c.Storage.AdminLsWithPage.Path(path),
			c.Storage.AdminLsWithPage.PageOffset(o.Offset),
			c.Storage.AdminLsWithPage.PageSize(o.Limit),
			c.Storage.AdminLsWithPage.FilterRegexp(o.Filter),
		)
		PrintResp(res, err, "Ls")

		return nil
	}

	return cmd
}

func newStorageMkdirCmd(o StorageOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "mkdir \"/{userID}/path\"",
		Short: "指定路径创建文件夹",
		Long:  "指定路径创建文件夹, 路径必须以/{userID}开头",
		Args:  cobra.ExactArgs(1),
		Example: `- 创建某用户根目录下的文件夹, 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage mkdir /4TpFFZDkFWy/test -T hpc -Z az-zhigu`,
	}

	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "az-zhigu", "区域")
	cmd.Flags().StringVarP(&o.Type, "type", "T", "cloud", "存储类型, hpc | cloud")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.validateType()
		o.complete()

		path := args[0]
		c := GetStorageClient(o.Zone, o.Type)
		res, err := c.Storage.AdminMkdir(
			c.Storage.AdminMkdir.Path(path),
		)
		PrintResp(res, err, "Mkdir")

		return nil
	}

	return cmd
}

func newStorageUploadCmd(o StorageOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "upload \"/localPath\" \"/{userID}/path\"",
		Short: "上传文件",
		Long:  "上传文件, upload [本地路径] [远程存储目标路径], 目标路径必须以/{userID}开头, 需要指定最终的文件名/目录名\n例如: /{userID}/path/filename.txt, 支持上传目录",
		Args:  cobra.ExactArgs(2),
		Example: `- 上传本地文件到某用户根目录下, 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage upload /tmp/test.txt /4TpFFZDkFWy/test.txt -T hpc -Z az-zhigu
- 上传本地目录到某用户根目录下, 指定存储类型为cloud, 指定区域为az-zhigu
  - ysadmin storage upload /tmp/testdir /4TpFFZDkFWy/testdir -T cloud -Z az-zhigu`,
	}

	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "az-zhigu", "区域")
	cmd.Flags().StringVarP(&o.Type, "type", "T", "cloud", "存储类型, hpc | cloud")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.validateType()
		o.complete()

		localPath := args[0]
		remotePath := args[1]
		c := GetStorageClient(o.Zone, o.Type)

		stat, err := os.Stat(localPath)
		if err != nil {
			fmt.Printf("StatLocalPathFail, Error: %s, Path: %s\n", err.Error(), localPath)
			return nil
		}
		if stat.IsDir() {
			fmt.Printf("Uploading Dir: %s ...\n", localPath)
			filepath.WalkDir(localPath, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() {
					dest := filepath.Join(remotePath, strings.TrimPrefix(path, localPath))
					err = o.uploadFile(c, path, dest)
					if err != nil {
						return err
					}
				}
				return nil
			})
		} else {
			o.uploadFile(c, localPath, remotePath)
		}

		return nil
	}

	return cmd
}

func newStorageReadAtCmd(o StorageOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "readat \"/{userID}/path\"",
		Short: "读取文件",
		Long:  "读取文件, 路径必须以/{userID}开头",
		Args:  cobra.ExactArgs(1),
		Example: `- 读取某用户根目录下的文件, 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage readat /4TpFFZDkFWy/test.txt -T hpc -Z az-zhigu
- 读取某用户根目录下的文件, 指定存储类型为cloud, 指定区域为az-zhigu, 从第200字节开始读取, 读取1000字节
  - ysadmin storage readat /4TpFFZDkFWy/test.txt -T cloud -Z az-zhigu -O 200 -L 1000`,
	}

	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "az-zhigu", "区域")
	cmd.Flags().StringVarP(&o.Type, "type", "T", "cloud", "存储类型, hpc | cloud")
	cmd.Flags().Int64VarP(&o.Offset, "offset", "O", 0, "offset")
	cmd.Flags().Int64VarP(&o.Limit, "limit", "L", 1000, "limit")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.validateType()
		o.complete()

		path := args[0]
		c := GetStorageClient(o.Zone, o.Type)
		res, err := c.Storage.AdminReadAt(
			c.Storage.AdminReadAt.Path(path),
			c.Storage.AdminReadAt.Offset(o.Offset),
			c.Storage.AdminReadAt.Length(o.Limit),
		)
		if err != nil {
			fmt.Printf("ReadAtFail, Error:\n%s\n", err.Error())
			return nil
		}
		fmt.Printf("ReadAt From %s, Offset: %d, Length: %d\n", path, o.Offset, o.Limit)
		fmt.Println("--------------------Content:--------------------")
		fmt.Println(string(res.Data))

		return nil
	}

	return cmd
}

func newStorageDownloadCmd(o StorageOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "download \"/{userID}/path\" \"/localPath\"",
		Short: "下载文件",
		Long:  "下载文件, download [远程存储目标路径] [本地路径], 目标路径必须以/{userID}开头, 本地路径必须是最终文件路径(包含文件名)\n例如: /{userID}/path/filename.txt /tmp/filename.txt, 不支持下载目录",
		Args:  cobra.ExactArgs(2),
		Example: `- 下载某用户根目录下的文件, 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage download /4TpFFZDkFWy/test.txt /tmp/test.txt -T hpc -Z az-zhigu`,
	}

	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "az-zhigu", "区域")
	cmd.Flags().StringVarP(&o.Type, "type", "T", "cloud", "存储类型, hpc | cloud")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.validateType()
		o.complete()

		remotePath := args[0]
		localPath := args[1]

		c := GetStorageClient(o.Zone, o.Type)
		fmt.Printf("Downloading: %s   ===>>>  %s  ...\n", remotePath, localPath)
		res, err := c.Storage.AdminDownload(
			c.Storage.AdminDownload.Path(remotePath),
		)
		if err != nil {
			fmt.Printf("Download Fail, Error: %s", err.Error())
			return nil
		}

		err = os.WriteFile(localPath, res.Data, 0755)
		if err != nil {
			fmt.Printf("WriteLocalPathFail, Error: %s, Path: %s", err.Error(), localPath)
		} else {
			fmt.Printf("DownloadOK: %s   ===>>>  %s\n", remotePath, localPath)
		}

		return nil
	}

	return cmd
}

func newStorageBatchDownloadCmd(o StorageOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "batch-download \"/{userID}/path\" \"/localPath\"",
		Short: "批量下载文件",
		Long:  "批量下载文件, batch-download [远程存储目标路径] [本地路径], 目标路径必须以/{userID}开头, 本地路径必须是最终文件路径且需要.zip后缀\n例如: /{userID}/path/dir /tmp/filename.zip, 不支持下载单个文件",
		Args:  cobra.ExactArgs(2),
		Example: `- 批量下载某用户根目录下的testdir目录, 指定存储类型为hpc, 指定区域为az-zhigu, 本地路径为/tmp/testdir.zip, 最后压缩包仅包含testdir目录下的所有文件(-B 不传默认与传"/4TpFFZDkFWy/testdir"相同, 即压缩包中不包含远程路径的一层目录)
  - ysadmin storage batch-download /4TpFFZDkFWy/testdir /tmp/testdir.zip -T hpc -Z az-zhigu [-B /4TpFFZDkFWy/testdir]
- 批量下载某用户根目录下的testdir目录, 指定存储类型为hpc, 指定区域为az-zhigu, 本地路径为/tmp/testdir.zip, 最后压缩包会包含一层testdir目录
  - ysadmin storage batch-download /4TpFFZDkFWy/testdir /tmp/testdir.zip -T hpc -Z az-zhigu -B /4TpFFZDkFWy`,
	}

	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "az-zhigu", "区域")
	cmd.Flags().StringVarP(&o.Type, "type", "T", "cloud", "存储类型, hpc | cloud")
	cmd.Flags().StringVarP(&o.BasePath, "base-path", "B", "", "压缩包的起始路径（不包含）, 不传默认与远程路径相同, 即压缩包中不包含远程路径的一层目录")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.validateType()
		o.complete()

		remotePath := args[0]
		localPath := args[1]
		basePath := o.BasePath
		if basePath == "" {
			basePath = remotePath
		}

		c := GetStorageClient(o.Zone, o.Type)
		fmt.Printf("BatchDownloading: %s   ===>>>  %s  ...\n", remotePath, localPath)
		_, err := c.Storage.AdminBatchDownload(
			c.Storage.AdminBatchDownload.Paths(remotePath),
			c.Storage.AdminBatchDownload.WithResolver(o.getDownloadResolver(localPath)),
			c.Storage.AdminBatchDownload.FileName(filepath.Base(localPath)),
			c.Storage.AdminBatchDownload.IsCompress(true),
			c.Storage.AdminBatchDownload.BasePath(basePath),
		)
		if err != nil {
			fmt.Printf("BatchDownload Fail, Error: %s", err.Error())
			return nil
		}
		fmt.Printf("BatchDownload OK: %s   ===>>>  %s\n", remotePath, localPath)

		return nil
	}

	return cmd

}

func newStorageRmCmd(o StorageOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "rm \"/{userID}/path\"",
		Short: "删除文件",
		Long:  "删除文件, 路径必须以/{userID}开头",
		Args:  cobra.ExactArgs(1),
		Example: `- 删除某用户根目录下的文件, 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage rm /4TpFFZDkFWy/test.txt -T hpc -Z az-zhigu
- 删除某用户根目录下的文件夹, 指定存储类型为cloud, 指定区域为az-zhigu
  - ysadmin storage rm /4TpFFZDkFWy/testdir -T cloud -Z az-zhigu`,
	}

	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "az-zhigu", "区域")
	cmd.Flags().StringVarP(&o.Type, "type", "T", "cloud", "存储类型, hpc | cloud")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.validateType()
		o.complete()

		path := args[0]
		c := GetStorageClient(o.Zone, o.Type)
		res, err := c.Storage.AdminRm(
			c.Storage.AdminRm.Path(path),
		)
		PrintResp(res, err, "Delete")

		return nil
	}

	return cmd
}

func newStorageMvCmd(o StorageOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "mv \"/{userID}/srcPath\" \"/{userID}/destPath\"",
		Short: "移动文件",
		Long:  "移动文件, src和dest路径必须以/{userID}开头, 不存在的上级目录会被创建",
		Args:  cobra.ExactArgs(2),
		Example: `- 移动某用户根目录下的文件, 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage mv /4TpFFZDkFWy/test.txt /4TpFFZDkFWy/test2.txt -T hpc -Z az-zhigu
- 移动某用户根目录下的文件夹, 指定存储类型为cloud, 指定区域为az-zhigu
  - ysadmin storage mv /4TpFFZDkFWy/testdir /4TpFFZDkFWy/testdir2 -T cloud -Z az-zhigu`,
	}

	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "az-zhigu", "区域")
	cmd.Flags().StringVarP(&o.Type, "type", "T", "cloud", "存储类型, hpc | cloud")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.validateType()
		o.complete()

		src := args[0]
		dest := args[1]
		c := GetStorageClient(o.Zone, o.Type)
		res, err := c.Storage.AdminMv(
			c.Storage.AdminMv.Src(src),
			c.Storage.AdminMv.Dest(dest),
		)
		PrintResp(res, err, "Move")

		return nil
	}

	return cmd
}

func newStorageUploadSpeedCmd(o StorageOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "upload-speed [\"/{userID}/path\"]",
		Short: "上传文件(测速)",
		Long:  "上传文件(测速), 不指定path默认是上传到配置文件对应user的某个测速目录下, 如果指定了path, 则上传到指定的path下",
		Args:  cobra.MaximumNArgs(1),
		Example: `- 上传文件到配置用户目录下(测速), 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage upload-speed -T hpc -Z az-zhigu
- 上传文件到某用户根目录下的testdir目录下(测速), 指定存储类型为cloud, 指定区域为az-zhigu
  - ysadmin storage upload-speed /4TpFFZDkFWy/testdir -T cloud -Z az-zhigu
- 上传文件到某用户根目录下的testdir目录下(测速), 指定存储类型为cloud, 指定区域为az-zhigu, 指定文件大小为2000MB
  - ysadmin storage upload-speed /4TpFFZDkFWy/testdir -T cloud -Z az-zhigu -S 2000`,
	}

	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "az-zhigu", "区域")
	cmd.Flags().StringVarP(&o.Type, "type", "T", "cloud", "存储类型, hpc | cloud")
	cmd.Flags().Int64VarP(&o.Size, "size", "S", 1000, "文件大小 单位为MB")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.validateType()
		o.complete()
		o.validateUserID()

		c := GetStorageClient(o.Zone, o.Type)

		path := ""
		if len(args) > 0 {
			path = args[0]
		}

		dataSize := int64(dataSize)
		totalDataSize := int64(0)
		fileName := "test-upload-file-" + time.Now().Format("20060102150405")
		fmt.Printf("Uploading File ...,dest: %v \n", filepath.Join(path, fileName))

		var dest string
		for totalDataSize < o.Size {
			if totalDataSize+dataSize > o.Size {
				dataSize = o.Size - totalDataSize
			}
			totalDataSize += dataSize
			dest = fmt.Sprintf("/%s/%s/%s/%s", o.UserId, testFolder, fileName, uuid.New().String())
			if path != "" {
				dest = filepath.Join(path, fileName, uuid.New().String())
			}
			data := make([]byte, dataSize*1024*1024)
			rand.Read(data)
			startTime := time.Now()
			err := upload.UploadData(data, dest, c)
			if err == nil {
				fmt.Printf("Upload File Size: %dMB, Speed: %f MB/s\n", dataSize, float64(dataSize)/time.Since(startTime).Seconds())
			} else {
				fmt.Println("Upload Fail, Error: ", err.Error())
				return nil
			}
		}
		fmt.Printf("Upload File Successfully, dest: %v \n", dest)

		return nil
	}

	return cmd
}

func newStorageDownloadSpeedCmd(o StorageOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "download-speed [\"/{userID}/path\"]",
		Short: "下载文件(测速)",
		Long:  "下载文件(测速), 不指定path默认是下载配置文件对应user的某个测速目录下的特定文件, 如果指定了path, 则下载指定的path下的文件。(只测速, 不会实际下载, Size大于文件大小会报错)",
		Args:  cobra.MaximumNArgs(1),
		Example: `- 从配置用户目录下下载文件(测速), 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage download-speed -T hpc -Z az-zhigu
- 从某用户根目录下的testdir目录下下载文件, 指定存储类型为cloud, 指定区域为az-zhigu
  - ysadmin storage download-speed /4TpFFZDkFWy/testdir/testfile -T cloud -Z az-zhigu
- 从某用户根目录下的testdir目录下下载文件, 指定存储类型为cloud, 指定区域为az-zhigu, 指定读取文件大小为200MB
  - ysadmin storage download-speed /4TpFFZDkFWy/testdir/testfile -T cloud -Z az-zhigu -S 200`,
	}

	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "az-zhigu", "区域")
	cmd.Flags().StringVarP(&o.Type, "type", "T", "cloud", "存储类型, hpc | cloud")
	cmd.Flags().Int64VarP(&o.Size, "size", "S", 1000, "文件大小 单位为MB")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.validateType()
		o.complete()
		o.validateUserID()

		c := GetStorageClient(o.Zone, o.Type)

		fileName := downloadFileName
		path := fmt.Sprintf("/%s/%s/%s", o.UserId, testFolder, fileName)
		if len(args) > 0 && args[0] != "" {
			path = args[0]
		}

		o.Data = make([]byte, 0, o.Size*1024*1024)
		go o.monitorDownloadSpeed()

		fmt.Printf("Downloading File ...,src: %s\n", path)
		_, err := c.Storage.AdminDownload(
			c.Storage.AdminDownload.Path(path),
			c.Storage.AdminDownload.Range(0, o.Size*1024*1024),
			c.Storage.AdminDownload.WithResolver(o.getResolver()),
		)
		if err != nil {
			fmt.Printf("Download Fail, Error: %s", err.Error())
			return nil
		}
		fmt.Printf("Download File Successfully\n")

		return nil
	}

	return cmd
}

func newStorageBatchDownloadSpeedCmd(o StorageOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "batch-download-speed [\"/{userID}/path\"]",
		Short: "批量下载文件(测速)",
		Long:  "批量下载文件(测速), 不指定path默认是下载配置文件对应user的某个测速目录下的特定目录, 如果指定了path, 则下载指定的path下的文件。(只测速, 不会实际下载, Size大于文件大小会报错)",
		Args:  cobra.MaximumNArgs(1),
		Example: `- 从配置用户目录下下载文件(测速), 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage batch-download-speed -T hpc -Z az-zhigu
- 从某用户根目录下的testdir目录下下载文件(测速), 指定存储类型为cloud, 指定区域为az-zhigu
  - ysadmin storage batch-download-speed /4TpFFZDkFWy/testdir -T cloud -Z az-zhigu`,
	}

	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "az-zhigu", "区域")
	cmd.Flags().StringVarP(&o.Type, "type", "T", "cloud", "存储类型, hpc | cloud")
	cmd.Flags().Int64VarP(&o.Size, "size", "S", 1000, "文件大小 单位为MB")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.validateType()
		o.complete()
		o.validateUserID()

		c := GetStorageClient(o.Zone, o.Type)

		path := fmt.Sprintf("/%s/%s", o.UserId, testFolder)
		if len(args) > 0 {
			path = args[0]
		}

		o.Data = make([]byte, 0, o.Size*1024*1024)
		go o.monitorDownloadSpeed()

		fmt.Printf("Downloading File ...,src: %s\n", path)
		_, err := c.Storage.AdminBatchDownload(
			c.Storage.AdminBatchDownload.Paths(path),
			c.Storage.AdminBatchDownload.WithResolver(o.getResolver()),
			c.Storage.AdminBatchDownload.FileName("test.zip"),
		)
		if err != nil {
			fmt.Printf("Batch download Fail, Error: %s", err.Error())
			return nil
		}
		fmt.Printf("Batch download File Successfully\n")

		return nil
	}

	return cmd
}

func newStorageListQuotaCmd(o StorageOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list-quota",
		Short: "列出配额",
		Long:  "列出配额, 列出分区存储下各目录(用户)配额",
		Args:  cobra.ExactArgs(0),
		Example: `- 列出分区存储下各目录(用户)配额, 指定存储类型为hpc, 指定区域为az-zhigu, 分页显示, 从第0条开始, 每页显示20条
  - ysadmin storage list-quota -T hpc -Z az-zhigu -O 0 -S 20`,
	}

	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "az-zhigu", "区域")
	cmd.Flags().StringVarP(&o.Type, "type", "T", "cloud", "存储类型, hpc | cloud")
	cmd.Flags().Int64VarP(&o.Offset, "offset", "O", 0, "offset")
	cmd.Flags().Int64VarP(&o.Size, "size", "S", 1000, "size")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.validateType()
		o.complete()

		c := GetStorageClient(o.Zone, o.Type)

		res, err := c.StorageQuota.ListStorageQuotaAdmin(
			c.StorageQuota.ListStorageQuotaAdmin.PageOffset(int(o.Offset)),
			c.StorageQuota.ListStorageQuotaAdmin.PageSize(int(o.Size)),
		)
		PrintResp(res, err, "ListQuota")

		return nil
	}

	return cmd
}

func newStorageUpdateQuotaCmd(o StorageOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "update-quota",
		Short: "更新配额",
		Long:  "更新配额, 更新某用户配额",
		Args:  cobra.ExactArgs(0),
		Example: `- 更新某用户配额, 指定存储类型为hpc, 指定区域为az-zhigu, 配额为1000GB
  - ysadmin storage update-quota -U 4TpFFZDkFWy -T hpc -Z az-zhigu -L 1000`,
	}

	cmd.Flags().StringVarP(&o.UserId, "user_id", "U", "", "用户ID, 不填默认使用配置文件的storage_ys_id")
	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "az-zhigu", "区域")
	cmd.Flags().StringVarP(&o.Type, "type", "T", "cloud", "存储类型, hpc | cloud")
	cmd.Flags().Int64VarP(&o.Limit, "limit", "L", 1000, "limit")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.validateType()
		o.complete()
		o.validateUserID()

		c := GetStorageClient(o.Zone, o.Type)

		res, err := c.StorageQuota.PutStorageQuotaAdmin(
			c.StorageQuota.PutStorageQuotaAdmin.UserID(o.UserId),
			c.StorageQuota.PutStorageQuotaAdmin.StorageLimit(float64(o.Limit)),
		)
		PrintResp(res, err, "UpdateQuota")
		if err == nil {
			fmt.Printf("update user %s quota limit to %dGB Successfully\n", o.UserId, o.Limit)
		}

		return nil
	}

	return cmd
}

func newStorageQuotaTotalCmd(o StorageOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "quota-total",
		Short: "配额总量",
		Long:  "配额总量, 列出分区存储下配额总量",
		Args:  cobra.ExactArgs(0),
		Example: `- 列出分区存储下配额总量, 指定存储类型为hpc, 指定区域为az-zhigu
  - ysadmin storage quota-total -T hpc -Z az-zhigu`,
	}

	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "az-zhigu", "区域")
	cmd.Flags().StringVarP(&o.Type, "type", "T", "cloud", "存储类型, hpc | cloud")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.validateType()
		o.complete()

		c := GetStorageClient(o.Zone, o.Type)

		res, err := c.StorageQuota.GetStorageQuotaTotalAdmin()
		PrintResp(res, err, "QuotaTotal")

		return nil
	}

	return cmd
}

func newStorageListOperationLogCmd(o StorageOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list-operation-log",
		Short: "列出操作日志",
		Long:  "列出操作日志",
		Args:  cobra.ExactArgs(0),
		Example: `- 列出操作日志, 指定存储类型为hpc, 指定区域为az-zhigu, 分页显示, 从第0条开始, 每页显示20条
  - ysadmin storage list-operation-log -T hpc -Z az-zhigu -O 0 -S 20
- 列出操作日志, 指定存储类型为hpc, 指定区域为az-zhigu, 指定文件名为test.txt, 指定操作类型为UPLOAD, 指定文件类型为FILE, 指定开始时间为2021-01-01 00:00:00, 指定结束时间为2021-01-02 00:00:00
  - ysadmin storage list-operation-log -T hpc -Z az-zhigu --file-name "test.txt" --operation-type UPLOAD --file-type FILE -b "2021-01-01 00:00:00" -e "2021-01-02 00:00:00"`,
	}

	cmd.Flags().StringVarP(&o.UserId, "user_id", "U", "", "用户ID, 不填默认使用配置文件的storage_ys_id")
	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "az-zhigu", "区域")
	cmd.Flags().StringVarP(&o.Type, "type", "T", "cloud", "存储类型, hpc | cloud")
	cmd.Flags().StringVarP(&o.FileName, "file-name", "f", "", "文件名")
	cmd.Flags().StringVarP(&o.FileTypes, "file-type", "t", "", "文件类型, 可选值: FILE-文件, DIRECTORY-目录")
	cmd.Flags().StringVarP(&o.OperationTypes, "operation-type", "o", "", "操作类型, 可选值: UPLOAD-上传, DOWNLOAD-下载, DELETE-删除, MOVE-移动, MKDIR-添加文件夹, COPY-拷贝, COPY_RANGE-指定范围拷贝,COMPRESS-压缩, CREATE-创建, LINK-链接, READ_AT-读, WRITE_AT-写")
	cmd.Flags().StringVarP(&o.BeginTime, "begin-time", "b", "", "开始时间, 格式: YYYY-MM-DD HH:mm:ss")
	cmd.Flags().StringVarP(&o.EndTime, "end-time", "e", "", "结束时间, 格式: YYYY-MM-DD HH:mm:ss")
	cmd.Flags().Int64VarP(&o.Offset, "offset", "O", 0, "offset")
	cmd.Flags().Int64VarP(&o.Limit, "limit", "L", 1000, "limit")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.validateType()
		o.complete()
		o.validateUserID()

		c := GetStorageClient(o.Zone, o.Type)

		var beginTimestamp int64
		var endTimestamp int64
		var err error
		if o.BeginTime != "" {
			beginTimestamp, err = StringToUnixTimestamp(o.BeginTime)
			if err != nil {
				fmt.Printf("BeginTimeError,format: YYYY-MM-DD HH:mm:ss, Error: %s\n", err.Error())
				return nil
			}
		}

		if o.EndTime != "" {
			endTimestamp, err = StringToUnixTimestamp(o.EndTime)
			if err != nil {
				fmt.Printf("EndTimeError,format: YYYY-MM-DD HH:mm:ss, Error: %s\n", err.Error())
				return nil
			}
		}
		res, err := c.StorageOperationLog.ListOperationLogAdmin(
			c.StorageOperationLog.ListOperationLogAdmin.PageOffset(o.Offset),
			c.StorageOperationLog.ListOperationLogAdmin.PageSize(o.Limit),
			c.StorageOperationLog.ListOperationLogAdmin.FileName(o.FileName),
			c.StorageOperationLog.ListOperationLogAdmin.OperationTypes(o.OperationTypes),
			c.StorageOperationLog.ListOperationLogAdmin.FileTypes(o.FileTypes),
			c.StorageOperationLog.ListOperationLogAdmin.UserID(o.UserId),
			c.StorageOperationLog.ListOperationLogAdmin.BeginTime(beginTimestamp),
			c.StorageOperationLog.ListOperationLogAdmin.EndTime(endTimestamp),
		)
		PrintResp(res, err, "ListOperationLog")

		return nil
	}

	return cmd
}

// ------------------

func (o *StorageOptions) uploadFile(c *openapi.Client, src, dest string) error {
	file, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("Stat File %s fail, Error: %s\n", src, err.Error())
	}
	fmt.Printf("Uploading file(size: %d): %s  ===>>>  %s...\n", file.Size(), src, dest)
	err = upload.Upload(src, dest, c, nil, nil, nil)
	if err == nil {
		fmt.Printf("Upload File Successfully, %s  ===>>>  %s\n", src, dest)
	} else {
		fmt.Println("Upload Fail, Error: ", err.Error())
	}
	return nil
}

func (o *StorageOptions) monitorDownloadSpeed() {
	interval := 1
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	var lastSize int
	for {
		select {
		case <-ticker.C:
			currentSize := len(o.Data)
			sizeIncrease := currentSize - lastSize
			speed := float64(sizeIncrease) / float64(interval)
			fmt.Printf("Downloaded: %dMB, Speed: %f MB/s\n", currentSize/1024/1024, speed/1024/1024)

			lastSize = currentSize

			if int64(currentSize) >= o.Size*1024*1024 {
				break
			}
		}
	}
}

func (o *StorageOptions) getResolver() func(resp *http.Response) error {

	return func(resp *http.Response) error {

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
			body, _ := io.ReadAll(resp.Body)
			defer func() { _ = resp.Body.Close() }()
			return errors.Errorf("http: %v, body: %v", resp.Status, string(body))
		}

		buf := make([]byte, 1024)
		var err error
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				o.Data = append(o.Data, buf[:n]...)
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
		}

		defer func() { _ = resp.Body.Close() }()

		return err
	}
}

func (o *StorageOptions) getDownloadResolver(filepath string) func(resp *http.Response) error {

	return func(resp *http.Response) error {

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
			body, _ := io.ReadAll(resp.Body)
			defer func() { _ = resp.Body.Close() }()
			return errors.Errorf("http: %v, body: %v", resp.Status, string(body))
		}

		fd, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		defer func() { _ = fd.Close() }()

		_, err = io.Copy(fd, resp.Body)
		defer func() { _ = resp.Body.Close() }()
		if err != nil {
			return err
		}

		return err
	}
}
