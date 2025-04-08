package upload

import (
	"github.com/pkg/errors"
	client "github.com/yuansuan/ticp/common/openapi-go"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	"golang.org/x/sync/errgroup"
	"os"
)

const (
	BLOCK_SIZE   = 1024 * 1024 * 10
	MaxGoroutine = 8
)

func Upload(filePath, destPath, compressType string, c *client.Client,
	initFunc func(int64, string, *client.Client) (string, int64, error),
	uploadChunkFunc func(string, int64, []byte, string, *client.Client) error,
	completeFunc func(string, string, *client.Client) error) error {

	if initFunc == nil {
		initFunc = InitUpload
	}
	if uploadChunkFunc == nil {
		uploadChunkFunc = UploadChunk
	}

	if completeFunc == nil {
		completeFunc = CompleteUpload
	}

	// 初始化上传 获得上传id
	file, err := os.Stat(filePath)
	if err != nil {
		return errors.New("file not exist: " + err.Error())
	}
	//fileName := file.Name()
	//fileType := "application/octet-stream"
	uploadId, fileSize, err := initFunc(file.Size(), destPath, c)
	if err != nil {
		return errors.New("init upload failed: " + err.Error())
	}

	// 分片上传文件
	if err = uploadFileConcurrently(filePath, compressType, fileSize, uploadChunkFunc, uploadId, c); err != nil {
		return err
	}

	// 完成上传
	if err = completeFunc(uploadId, destPath, c); err != nil {
		return errors.New("complete upload failed: " + err.Error())
	}

	return nil
}

func UploadData(data []byte, destPath string, c *client.Client) error {

	// 初始化上传 获得上传id
	uploadId, fileSize, err := InitUpload(int64(len(data)), destPath, c)
	if err != nil {
		return errors.New("init upload failed: " + err.Error())
	}

	// 分片上传文件
	var eg errgroup.Group
	var chunkIndex int64
	var offset int64

	sem := make(chan struct{}, MaxGoroutine)

	for {
		sem <- struct{}{}
		if offset >= fileSize {
			break
		}
		chunkSize := int64(BLOCK_SIZE)
		if fileSize-offset < BLOCK_SIZE {
			chunkSize = fileSize - offset
		}
		currentOffset := offset
		currentChunkIndex := chunkIndex
		offset += chunkSize
		chunkIndex++

		chunkData := data[currentOffset : currentOffset+chunkSize]
		eg.Go(func() error {
			if err = UploadChunk(uploadId, currentChunkIndex, chunkData, common.GZIP, c); err != nil {
				<-sem
				return errors.New("upload chunk failed: " + err.Error())
			}
			<-sem
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	// 完成上传
	if err = CompleteUpload(uploadId, destPath, c); err != nil {
		return errors.New("complete upload failed: " + err.Error())
	}
	return nil
}

func uploadFileConcurrently(filePath, compressType string, fileSize int64, uploadChunkFunc func(string, int64, []byte, string, *client.Client) error, uploadId string, c *client.Client) error {
	var eg errgroup.Group
	var chunkIndex int64
	var offset int64

	f, err := os.Open(filePath)
	if err != nil {
		return errors.New("open file failed: " + err.Error())
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)

	sem := make(chan struct{}, MaxGoroutine)

	for {
		sem <- struct{}{}
		if offset >= fileSize {
			break
		}
		chunkSize := int64(BLOCK_SIZE)
		if fileSize-offset < BLOCK_SIZE {
			chunkSize = fileSize - offset
		}
		currentOffset := offset
		currentChunkIndex := chunkIndex
		offset += chunkSize
		chunkIndex++

		chunkData := make([]byte, chunkSize)
		if _, err := f.ReadAt(chunkData, currentOffset); err != nil {
			<-sem
			return errors.New("read file failed: " + err.Error())
		}

		eg.Go(func() error {
			if err = uploadChunkFunc(uploadId, currentChunkIndex, chunkData, compressType, c); err != nil {
				<-sem
				return errors.New("upload chunk failed: " + err.Error())
			}
			<-sem
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil

}

func CompleteUpload(id, path string, c *client.Client) error {
	_, err := c.Storage.UploadComplete(
		c.Storage.UploadComplete.UploadID(id),
		c.Storage.UploadComplete.Path(path),
	)
	if err != nil {
		return err
	}
	return nil
}

func UploadChunk(id string, index int64, data []byte, compressorType string, c *client.Client) error {
	if compressorType == "" {
		compressorType = common.NONE
	}
	_, err := c.Storage.UploadSlice(
		c.Storage.UploadSlice.UploadID(id),
		c.Storage.UploadSlice.Slice(data),
		c.Storage.UploadSlice.Offset(index*BLOCK_SIZE),
		c.Storage.UploadSlice.Length(int64(len(data))),
		c.Storage.UploadSlice.Compressor(compressorType),
	)
	if err != nil {
		return err
	}

	return nil

}

func InitUpload(fileSize int64, destPath string, c *client.Client) (string, int64, error) {
	res, err := c.Storage.UploadInit(
		c.Storage.UploadInit.Path(destPath),
		c.Storage.UploadInit.Size(fileSize),
	)
	if err != nil {
		return "", 0, err
	}
	return res.Data.UploadID, fileSize, nil
}
