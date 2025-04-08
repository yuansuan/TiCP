package rpc

import (
	"context"
	"encoding/json"
	"path/filepath"

	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/cache"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/storage"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

func (srv *GRPCService) Exist(ctx context.Context, req *storage.ExistReq) (*storage.ExistResp, error) {
	if !req.Cross && strutil.IsEmpty(req.UserName) {
		return nil, status.Error(errcode.ErrFileUnknownUser, "username is required!")
	}
	resp := &storage.ExistResp{}
	isExistArray := make([]bool, 0)

	fileService, err := srv.GetFileService()
	if err != nil {
		return resp, err
	}

	// Loop the paths to check
	for _, p := range req.Paths {
		path := filepath.Clean(p)

		// Check the path existence
		isExistArray = append(isExistArray, fileService.Exist(ctx, req.UserName, req.Cross, path))
	}

	resp.IsExist = isExistArray

	return resp, nil
}

func (srv *GRPCService) CreateDir(ctx context.Context, req *storage.CreateDirReq) (*storage.CreateDirResp, error) {
	if !req.Cross && strutil.IsEmpty(req.UserName) {
		return nil, status.Error(errcode.ErrFileUnknownUser, "username is required!")
	}
	resp := &storage.CreateDirResp{}

	path := filepath.Clean(req.Path)

	fileService, err := srv.GetFileService()
	if err != nil {
		return resp, err
	}

	// Check the path existence
	if fileService.Exist(ctx, req.UserName, req.Cross, path) {
		err := status.Errorf(errcode.ErrFileAlreadyExist, "file(%v) already exists", path)

		return resp, err
	}

	// Create one new directory with mkdir
	err = fileService.CreateDir(ctx, req.UserName, path, req.Cross)

	return resp, err
}

func (srv *GRPCService) Mv(ctx context.Context, req *storage.MvReq) (resp *storage.MvResp, err error) {
	resp = &storage.MvResp{}

	fileService, err := srv.GetFileService()
	if err != nil {
		return resp, err
	}

	// Run this service to move this file to other directory
	err = fileService.Mv(ctx, "", true, req.Overwrite, req.Srcpath, req.Dstpath)
	if err != nil {
		return resp, err
	}

	return
}

func (srv *GRPCService) Rm(ctx context.Context, req *storage.RmReq) (resp *storage.RmResp, err error) {
	if !req.Cross && strutil.IsEmpty(req.UserName) {
		return nil, status.Error(errcode.ErrFileUnknownUser, "username is required!")
	}

	resp = &storage.RmResp{}

	fileService, err := srv.GetFileService()
	if err != nil {
		return
	}

	for _, path := range req.Paths {

		// Check the path existence
		if !fileService.Exist(ctx, req.UserName, req.Cross, path) {
			return resp, status.Errorf(errcode.ErrFileNotExist, "file(%v) not exist", path)
		}

		// Remove the file
		err = fileService.Remove(ctx, req.UserName, path, req.Cross)
		if err != nil {
			return
		}
	}

	return
}

func (srv *GRPCService) GetFileService() (service.FileService, error) {
	return srv.LocalFileService, nil
}

func (srv *GRPCService) Realpath(ctx context.Context, req *storage.RealpathReq) (*storage.RealpathResp, error) {
	fileService, err := srv.GetFileService()
	if err != nil {
		return nil, err
	}

	realpath, err := fileService.Realpath(ctx, req.RelativePath)
	if err != nil {
		return nil, err
	}

	return &storage.RealpathResp{
		Realpath: realpath,
	}, nil
}

// List lists all files under the path.
func (srv *GRPCService) List(ctx context.Context, req *storage.ListReq) (*storage.ListResp, error) {
	if !req.Cross && strutil.IsEmpty(req.UserName) {
		return nil, status.Error(errcode.ErrFileUnknownUser, "username is required!")
	}

	resp := &storage.ListResp{}

	path := filepath.Clean(req.Path)

	fileService, err := srv.GetFileService()
	if err != nil {
		return resp, err
	}
	// Get the file list under this path

	var page *xtype.Page
	if req.Page != nil {
		page = &xtype.Page{
			Index: req.Page.Index,
			Size:  req.Page.Size,
		}
	}

	files, err := fileService.List(ctx, req.UserName, path, page, req.Cross, req.ShowHideFile, req.FilterRegexpList)
	if err != nil {
		return resp, nil
	}

	resp.Files = util.ConvertFiles(files)

	return resp, nil
}

func (srv *GRPCService) SubmitUploadHpcFileTask(ctx context.Context, req *storage.SubmitUploadHpcFileTaskReq) (resp *storage.SubmitUploadHpcFileTaskResp, err error) {
	return &storage.SubmitUploadHpcFileTaskResp{
		TaskKey: "",
	}, nil
}

func (srv *GRPCService) GetUploadHpcFileTaskStatus(ctx context.Context, req *storage.GetUploadHpcFileTaskStatusReq) (resp *storage.GetUploadHpcFileTaskStatusResp, err error) {

	tasks, ok := cache.Cache.Get(req.TaskKey)
	if !ok || tasks == nil {
		return nil, status.Error(errcode.ErrFileUploadHPCFileNotExist, "")
	}

	// 考虑到有另一个协程正在处理tasks，如果其中元素正在被删除会导致遍历报错，这里选择拷贝一份再遍历
	jsonStr, err := json.Marshal(tasks)
	if err != nil {
		return nil, err
	}
	copyTasks := make(map[string]*dto.HPCUploadTask)
	err = json.Unmarshal(jsonStr, &copyTasks)
	if err != nil {
		return nil, err
	}

	var totalSize, currentSize int64
	var uploadStatus storage.UploadTaskStatusEnum

	// 优先级 失败 > 上传中 > 成功
	for _, task := range copyTasks {
		totalSize += task.TotalSize
		currentSize += task.CurrentSize
		if task.State == dto.UploadStateFailure && storage.UploadTaskStatusEnum_Failure > uploadStatus {
			uploadStatus = storage.UploadTaskStatusEnum_Failure
		} else if (task.State == dto.UploadStatePending || task.State == dto.UploadStateUploading) && storage.UploadTaskStatusEnum_Uploading > uploadStatus {
			uploadStatus = storage.UploadTaskStatusEnum_Uploading
		}
	}

	return &storage.GetUploadHpcFileTaskStatusResp{
		Status:      uploadStatus,
		TotalSize:   totalSize,
		CurrentSize: currentSize,
	}, nil
}

func (srv *GRPCService) InitUserHome(ctx context.Context, req *storage.InitUserHomeReq) (resp *storage.InitUserHomeResp, err error) {
	err2 := srv.LocalFileService.CheckUserHomePath(ctx, req.UserName, false)
	if err2 != nil {
		err = err2
	}
	return
}

func (srv *GRPCService) Link(ctx context.Context, req *storage.LinkReq) (resp *storage.LinkResp, err error) {
	if !req.Cross && strutil.IsEmpty(req.UserName) {
		return nil, status.Error(errcode.ErrFileUnknownUser, "username is required!")
	}

	fileService, err := srv.GetFileService()
	if err != nil {
		return resp, err
	}

	err = fileService.Link(ctx, &dto.LinkRequest{
		DstPath:      req.DstPath,
		Overwrite:    req.Overwrite,
		Cross:        req.Cross,
		IsCloud:      req.IsCloud,
		UserName:     req.UserName,
		SrcFilePaths: req.FilePaths,
	})

	if err != nil {
		return nil, err
	}

	return &storage.LinkResp{}, nil
}

func (srv *GRPCService) SymLink(ctx context.Context, req *storage.SymLinkReq) (resp *storage.SymLinkResp, err error) {
	if !req.Cross && strutil.IsEmpty(req.UserName) {
		return nil, status.Error(errcode.ErrFileUnknownUser, "username is required!")
	}
	fileService, err := srv.GetFileService()
	if err != nil {
		return resp, err
	}

	err = fileService.SymLink(ctx, &dto.SymLinkRequest{
		DstPath:   req.DstPath,
		Overwrite: req.Overwrite,
		Cross:     req.Cross,
		IsCloud:   req.IsCloud,
		UserName:  req.UserName,
		SrcPath:   req.SrcPath,
	})

	if err != nil {
		return nil, err
	}

	return &storage.SymLinkResp{}, nil
}

func (srv *GRPCService) HardLink(ctx context.Context, req *storage.HardLinkReq) (resp *storage.HardLinkResp, err error) {
	if !req.Cross && strutil.IsEmpty(req.UserName) {
		return nil, status.Error(errcode.ErrFileUnknownUser, "username is required!")
	}
	fileService, err := srv.GetFileService()
	if err != nil {
		return resp, err
	}

	err = fileService.HardLink(ctx, &dto.LinkRequest{
		DstPath:      req.DstPath,
		Overwrite:    req.Overwrite,
		Cross:        req.Cross,
		IsCloud:      req.IsCloud,
		UserName:     req.UserName,
		SrcDirPaths:  req.SrcDirPaths,
		SrcFilePaths: req.SrcFilePaths,
		CurrentPath:  req.CurrentPath,
		FilterPaths:  req.FilterPaths,
	})

	if err != nil {
		return nil, err
	}

	return &storage.HardLinkResp{}, nil
}

func (srv *GRPCService) Read(req *storage.ReadReq, res storage.Storage_ReadServer) error {
	if !req.Cross && strutil.IsEmpty(req.UserName) {
		return status.Error(errcode.ErrFileUnknownUser, "username is required!")
	}

	ctx := context.Background()

	fileService, err := srv.GetFileService()
	if err != nil {
		return err
	}

	fileInfo, err := fileService.GetFileInfoByStat(ctx, req.UserName, req.Path, req.Cross)
	if err != nil {
		return status.Error(errcode.ErrFileNotExist, "")
	}

	if fileInfo.IsDir {
		return status.Error(errcode.ErrFileFailReadFolder, "")
	}

	chunkSize := int64(64 * 1024)
	totalSize := fileInfo.Size

	// 分片流式传输
	for offset := int64(0); offset < totalSize; offset += chunkSize {
		lenth := chunkSize
		if offset+chunkSize > totalSize {
			lenth = totalSize - offset
		}

		content, err := fileService.Read(ctx, req.UserName, req.Path, offset, lenth, req.Cross)
		if err != nil {
			return err
		}

		if err := res.Send(&storage.ReadResp{
			Chunk: content,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (srv *GRPCService) GetLocalRootPathConfig(ctx context.Context, req *storage.GetLocalRootPathConfigReq) (resp *storage.GetLocalRootPathConfigResp, err error) {
	return &storage.GetLocalRootPathConfigResp{
		LocalRootPath: config.GetConfig().LocalRootPath,
	}, nil
}
