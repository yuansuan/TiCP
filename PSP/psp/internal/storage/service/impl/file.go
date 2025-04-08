package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/batchDownload"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/compress/start"
	pbcompressstatus "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/compress/status"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/copyRange"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/create"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/link"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/ls"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/mkdir"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/mv"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/readAt"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/realpath"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/rm"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/stat"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/complete"
	upinit "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/init"
	uploadslice "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/slice"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/cache"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	openApiConfig "github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/storage"
	pbnotice "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/notice"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/project"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/collectionutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/csvutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/lockutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/maputil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type FileServiceImpl struct {
	api                *openapi.OpenAPI
	isCloud            bool
	shareFileRecordDao dao.ShareRecordDao
}

func (srv *FileServiceImpl) UpdateRecordState(ctx *gin.Context, userId int64, recordIds []string) error {
	logger := logging.GetLogger(ctx)

	sids, fids := snowflake.BatchParseString(recordIds)
	if len(fids) > 0 {
		logger.Warnf("record ids [%v] is invalid, cannot parsed to snowflake id", fids)
	}

	return srv.shareFileRecordDao.UpdateRecordState(userId, sids)
}

func (srv *FileServiceImpl) GetShareCount(ctx *gin.Context, userId int64, share int8) (int64, error) {
	return srv.shareFileRecordDao.Count(userId, share)
}

func (srv *FileServiceImpl) CheckUserHomePath(ctx context.Context, userName string, cross bool) error {
	if !cross && !srv.Exist(ctx, userName, false, common.PersonalFolderPath) {
		return srv.CreateDir(ctx, userName, common.PersonalFolderPath, false)
	}
	return nil
}

func (srv *FileServiceImpl) ReadAll(ctx *gin.Context, userId int64) error {
	return srv.shareFileRecordDao.ReadAll(userId)
}

func (srv *FileServiceImpl) CheckOnlyReadDir(ctx *gin.Context, dir string, files []*dto.File) ([]*dto.File, error) {
	workspace := strings.Join([]string{common.Dot, common.WorkspaceFolderPath}, common.Slash)
	if dir != common.Dot && dir != workspace {
		return files, nil
	}

	onlyReadPathList := config.GetConfig().OnlyReadPathList

	switch dir {
	// 用户workspace目录
	case workspace:
		rsp, err := client.GetInstance().Project.GetMemberProjectsByUserId(ctx, &project.GetMemberProjectsByUserIdRequest{
			UserId: snowflake.ID(ginutil.GetUserID(ctx)).String(),
		})

		if err != nil {
			return nil, err
		}

		if rsp.Projects != nil {
			for _, projectObject := range rsp.Projects {
				index := strings.Index(projectObject.LinkPath, common.WorkspaceFolderPath)
				if index != -1 {
					onlyReadPathList = append(onlyReadPathList, projectObject.LinkPath[index:])
				}

			}
		}
		break
	default:
		break
	}

	onlyReadPathMap := maputil.ConvertSliceToMap(onlyReadPathList)

	for _, fileInfo := range files {
		if onlyReadPathMap[fileInfo.Path] {
			fileInfo.OnlyRead = true
		}
	}

	return files, nil
}

func (srv *FileServiceImpl) GetCopyStatus(ctx *gin.Context, key string) (dto.CopyState, error) {
	state, ok := cache.Cache.Get(key)

	if !ok {
		return dto.CopyStateFailure, status.Error(errcode.ErrFileCopyFileNotExist, "")
	}

	return state.(dto.CopyState), nil
}

func (srv *FileServiceImpl) HpcDownload(ctx *gin.Context, req *dto.HpcDownloadRequest) error {

	relRootPath := srv.ConvertLocalPath(req.DestDirPath, ginutil.GetUserName(ctx), false)
	// after username path, Prefix with req.CurrentPath，for example CurrentPath=/xxx, needMkdirPath=/xxx/111/
	needMkdirPaths := req.SrcDirPaths
	// after username path, Prefix with req.CurrentPath，for example CurrentPath=/xxx, needMkdirPath=/xxx/test.yml
	needDownloadFiles := req.SrcFilePaths

	if len(req.SrcDirPaths) > 0 {
		files, err := srv.ListOfRecur(ctx, req.UserName, req.SrcDirPaths, false, true, nil)
		if err != nil {
			return status.Errorf(errcode.ErrFileList, "")
		}

		if len(files) > 0 {
			for _, file := range files {
				if file.IsDir {
					needMkdirPaths = append(needMkdirPaths, file.Path)
				} else {
					needDownloadFiles = append(needDownloadFiles, file.Path)
				}
			}
		}
	}

	// 先把用户选中的文件夹目录创建好
	if len(needMkdirPaths) > 0 {
		needMkdirPaths = collectionutil.RemoveDuplicates(needMkdirPaths)
		for _, dirPath := range needMkdirPaths {
			// mkdir no CurrentPath
			err := os.MkdirAll(path.Join(relRootPath, strings.TrimPrefix(path.Clean(dirPath), path.Clean(req.CurrentPath))), 0777)
			if err != nil {
				return err
			}
		}
	}

	// 下载文件，耗时操作进异步
	go func() {
		if len(needDownloadFiles) > 0 {
			needDownloadFiles = collectionutil.RemoveDuplicates(needDownloadFiles)
			for _, filePath := range needDownloadFiles {
				// 先创建父级文件夹，因为用户有可能不选中文件夹，直接选中文件夹下的文件
				err := os.MkdirAll(path.Join(relRootPath, strings.TrimPrefix(path.Clean(path.Dir(filePath)), path.Clean(req.CurrentPath))), 0777)
				if err != nil {
					continue
				}
				// {relRootPath}/test.yml.downloading
				tmpFilePath := path.Join(relRootPath, fmt.Sprintf("%s%s", strings.TrimPrefix(filePath, req.CurrentPath), consts.TmpFileSuffix))
				finalFilePath := strings.TrimSuffix(tmpFilePath, consts.TmpFileSuffix)

				// 判断是否覆盖
				info, _ := os.Stat(finalFilePath)
				if info != nil && !info.IsDir() {
					if req.Overwrite {
						os.Remove(finalFilePath)
					} else {
						continue
					}
				}

				// 下载文件
				err = srv.linkFile(req.UserName, filePath, tmpFilePath)
				if err != nil {
					os.Remove(tmpFilePath)
				} else {
					os.Rename(tmpFilePath, finalFilePath)
				}
			}
		}
	}()

	return nil
}

func (srv *FileServiceImpl) linkFile(userName string, filePath string, tmpFilePath string) (err error) {
	openapiPath, err := srv.ConvertOpenApiReqPath(userName, filePath, false)
	if err != nil {
		return
	}

	resp, err := storage.DownloadRange(srv.api, openapiPath, 0, 0)
	if err != nil {
		return
	}
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(resp.Data))
	return
}

func NewLocalFileService() (service.FileService, error) {
	api, err := openapi.NewLocalHPCAPI()
	if err != nil {
		return nil, err
	}

	shareDao := dao.NewShareFileRecordDaoImpl()

	fileService := &FileServiceImpl{
		api:                api,
		isCloud:            false,
		shareFileRecordDao: shareDao,
	}

	return fileService, nil
}

func (srv *FileServiceImpl) GetUserHome(username string, cross bool) string {
	var openApiConf *openApiConfig.Settings
	openApiConf = openApiConfig.GetConfig().Local.Settings

	if !cross {
		return filepath.Join("/", openApiConf.UserId, username)
	}

	return filepath.Join("/" + openApiConf.UserId)
}

// ConvertOpenApiReqPath 将用户传入的path转为openapi接受的路径
// 如果传的是result/2022-08
// openapi: /{userName}/result/2022-08
func (srv *FileServiceImpl) ConvertOpenApiReqPath(userName, p string, cross bool) (string, error) {
	var openApiConf *openApiConfig.Settings
	openApiConf = openApiConfig.GetConfig().Local.Settings

	if !cross {
		return path.Join("/", openApiConf.UserId, userName, p), nil
	}

	// /{diskPrefix}/...
	p = path.Clean(path.Join("/", p))

	//var isWhiteListPath bool
	//for _, whitePath := range storageConf.WhitePathList {
	//	if strings.HasPrefix(p, whitePath) {
	//		isWhiteListPath = true
	//		break
	//	}
	//}

	//// 是否为临时目录
	//if strings.HasPrefix(p, fmt.Sprintf("/%v", userName)) || isWhiteListPath {
	return path.Join("/", openApiConf.UserId, p), nil
	//} else {
	//	return "", status.Error(errcode.ErrFileNoPermission, "no permission")
	//}

}

// PreUpload 文件上传前置接口，获取upload_id
func (srv *FileServiceImpl) PreUpload(ctx context.Context, req *dto.PreUploadRequest) (*dto.PreUploadResponse, error) {
	logger := logging.GetLogger(ctx)

	if err := validatePath(req.Path); err != nil {
		return nil, status.Error(errcode.ErrFileNameValidate, "")
	}

	openApiPath, err := srv.ConvertOpenApiReqPath(req.UserName, req.Path, req.Cross)
	if err != nil {
		return nil, err
	}

	uploadInitRsp, err := storage.UploadInit(srv.api, upinit.Request{
		Path: openApiPath,
		Size: req.FileSize,
	})

	if err != nil || uploadInitRsp == nil || uploadInitRsp.UploadID == "" {
		logger.Errorf("failed to get upload_id, err: %v", err)
		return nil, err
	}

	return &dto.PreUploadResponse{UploadId: uploadInitRsp.UploadID}, err
}

// Upload 文件切片上传
func (srv *FileServiceImpl) Upload(ctx *gin.Context, req *dto.UploadRequest, slice []*multipart.FileHeader) error {

	logger := logging.GetLogger(ctx)

	// 文件上传结束
	if req.Finish {
		// 结束标识为true调用UploadComplete
		openApiPath, err := srv.ConvertOpenApiReqPath(req.UserName, req.Path, req.Cross)
		if err != nil {
			return err
		}

		_, err = storage.UploadComplete(srv.api, complete.Request{
			Path:     openApiPath,
			UploadID: req.UploadID,
		})
		if err != nil {
			logger.Errorf("failed to upload slice, err: %v", err)
			return status.Errorf(errcode.ErrFileUploadSlice, "failed to upload complete")
		}

		tracelog.Info(ctx, fmt.Sprintf("upload file success, path[%s]", openApiPath))
	} else {
		// 获取切片文件
		sliceReader, err := slice[0].Open()

		if err != nil {
			logger.Errorf("failed to read slice, err: %v", err)
			return status.Errorf(errcode.ErrFileFailToOpenSlice, "failed to read slice")

		}

		fileBytes, err := io.ReadAll(sliceReader)

		if err != nil {
			logger.Errorf("failed to read the whole slice, err: %v", err)
			return status.Errorf(errcode.ErrFileSliceIncompleted, "failed to read the whole slice")
		}

		// 上传文件分片
		_, err = storage.UploadSlice(srv.api, uploadslice.Request{
			UploadID: req.UploadID,
			Offset:   req.Offset,
			Length:   req.SliceSize,
			Slice:    fileBytes,
		})

		if err != nil {
			logger.Errorf("failed to upload slice, err: ", err)
			return status.Errorf(errcode.ErrFileUploadSlice, "failed to upload slice")
		}
	}
	return nil
}

func (srv *FileServiceImpl) BatchDownloadPre(ctx *gin.Context, req *dto.BatchDownloadPreRequest) (*dto.BatchDownloadPreResponse, error) {
	logger := logging.GetLogger(ctx)

	if !CheckDownloadPermission(context.Background(), getRequestIP(ctx), req.UserName) {
		logger.Error("no download permission")
		return nil, status.Errorf(errcode.ErrFileDownloadPermission, "no download permission")
	}

	token := uuid.New().String()
	cache.Cache.Set(token, dto.DownloadCache{
		FilePaths:  req.FilePaths,
		FileName:   req.FileName,
		Cross:      req.Cross,
		IsCompress: req.IsCompress,
		IsCloud:    req.IsCloud,
		UserName:   req.UserName,
	}, gocache.DefaultExpiration)

	return &dto.BatchDownloadPreResponse{Token: token}, nil
}

// Compress
//
//	@GET	/Compress
func (srv *FileServiceImpl) Compress(ctx *gin.Context, req *dto.CompressRequest) (*dto.CompressResponse, error) {
	logger := logging.GetLogger(ctx)
	if len(req.SrcPaths) == 0 || req.DstPath == "" {
		logger.Error("compress start error")
		return nil, fmt.Errorf("compress path is null")
	}

	userName := ginutil.GetUserName(ctx)

	basePath, err := srv.ConvertOpenApiReqPath(userName, req.BasePath, false)
	if err != nil {
		return nil, fmt.Errorf("convert path is error")
	}
	dstPath, err := srv.ConvertOpenApiReqPath(userName, req.DstPath, false)
	if err != nil {
		return nil, fmt.Errorf("convert path is error")
	}

	srcPaths := make([]string, 0, len(req.SrcPaths))
	for _, path := range req.SrcPaths {
		srcPath, err := srv.ConvertOpenApiReqPath(userName, path, false)
		if err != nil {
			return nil, fmt.Errorf("convert path is error")
		}
		srcPaths = append(srcPaths, srcPath)
	}

	data, err := storage.CompressStart(srv.api, start.Request{
		Paths:      srcPaths,
		TargetPath: dstPath,
		BasePath:   basePath,
	})
	if err != nil {
		logger.Error("compress start error")
		return nil, fmt.Errorf("compress path is null")
	}

	userId := ginutil.GetUserID(ctx)
	var compressTask dto.CompressTask
	compressTask.CompressId = data.CompressID
	compressTask.Status = int8(dto.CompressStateCompressing)
	key := fmt.Sprintf("%s:%s:%d", common.StorageModule, common.CompressPreKey, userId)
	cache.Cache.Set(key, compressTask, dto.CompressDuration)

	logger.Infof("storage compress cache set task key:%v", key)

	go func() {
		ticker := time.NewTimer(3 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				dataStatus, err := storage.CompressStatus(srv.api, pbcompressstatus.Request{
					CompressID: data.CompressID,
				})
				if err != nil {
					logger.Infof("storage CompressStatus err:%v", err)
					continue
				}
				if err == nil {
					if dataStatus.Status == (dto.CompressStateCompressing.String()) {
						compressTask.Status = int8(dto.CompressStateCompressing)
					} else if dataStatus.Status == (dto.CompressStateFailure.String()) {
						compressTask.Status = int8(dto.CompressStateFailure)
					} else if dataStatus.Status == (dto.CompressStateSuccess.String()) {
						compressTask.Status = int8(dto.CompressStateSuccess)
						cache.Cache.Set(key, compressTask, dto.CompressDuration)
						return
					}
					cache.Cache.Set(key, compressTask, dto.CompressDuration)
				}
			}
		}
	}()

	return &dto.CompressResponse{CompressID: data.CompressID, TargetPath: data.FileName, IsCloud: req.IsCloud}, nil
}

func (srv *FileServiceImpl) CompressTasks(ctx *gin.Context) ([]*dto.CompressTask, error) {
	logger := logging.GetLogger(ctx)
	userId := ginutil.GetUserID(ctx)

	key := fmt.Sprintf("%s:%s:%d", common.StorageModule, common.CompressPreKey, userId)
	task, ok := cache.Cache.Get(key)
	if !ok {
		logger.Infof("cache get task empty, key:%v", key)
		return []*dto.CompressTask{}, nil
	}

	userCompressTask, ok := task.(dto.CompressTask)
	if ok {
		if int8(dto.CompressStateSuccess) == userCompressTask.Status || int8(dto.CompressStateFailure) == userCompressTask.Status {
			cache.Cache.Delete(key)
		}
		return []*dto.CompressTask{&userCompressTask}, nil
	}

	return []*dto.CompressTask{}, nil
}

func (srv *FileServiceImpl) CompressSubmit(ctx *gin.Context) (bool, error) {
	logger := logging.GetLogger(ctx)

	userId := ginutil.GetUserID(ctx)
	key := fmt.Sprintf("%s:%s:%d", common.StorageModule, common.CompressPreKey, userId)

	_, ok := cache.Cache.Get(key)
	if ok {
		logger.Errorf("already exist task !")
		return false, nil
	}

	return true, nil
}

// BatchDownload
//
//	@GET	/batch-download
func (srv *FileServiceImpl) BatchDownload(ctx *gin.Context, userName, token string, isCloud bool) error {
	downloadParam, ok := cache.Cache.Get(token)
	if !ok {
		return fmt.Errorf("by token get download param failed")
	}
	downloadParamDto, ok := downloadParam.(dto.DownloadCache)
	if !ok {
		return fmt.Errorf("convert download param to download param dto failed")
	}
	if downloadParamDto.IsCloud != isCloud {
		return fmt.Errorf("is cloud not match")
	}
	if downloadParamDto.UserName != userName {
		return fmt.Errorf("user name not match")
	}

	ctx.Writer.Header()["Content-Type"] = []string{"application/zip"}
	ctx.Writer.Header()["Content-Disposition"] = []string{"attachment; filename=" + url.QueryEscape(downloadParamDto.FileName)}

	openApiPathList := make([]string, 0, len(downloadParamDto.FilePaths))
	for _, filePath := range downloadParamDto.FilePaths {
		openApiPath, err := srv.ConvertOpenApiReqPath(userName, filePath, downloadParamDto.Cross)
		if err != nil {
			return status.Error(errcode.ErrFileBatchDownloadPermission, "no batch download permission")
		}
		openApiPathList = append(openApiPathList, openApiPath)
	}

	_, err := storage.BatchDownload(srv.api, batchDownload.Request{
		Paths:      openApiPathList,
		FileName:   downloadParamDto.FileName,
		BasePath:   srv.GetUserHome(userName, downloadParamDto.Cross),
		IsCompress: downloadParamDto.IsCompress,
	}, BuildSingleResolver(ctx))
	if err != nil {
		return err
	}

	return nil
}

func BuildSingleResolver(c *gin.Context) xhttp.ResponseResolver {
	return func(resp *http.Response) error {

		if resp.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("download file error: %v", resp.StatusCode))
		}

		_, err := io.Copy(c.Writer, resp.Body)
		if err != nil {
			panic(err)
		}

		return nil
	}
}

func getRequestIP(c *gin.Context) string {
	reqIP := c.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}

// GetDownloadName GetDownloadName
func GetDownloadName(ctx context.Context, userName string, token string) (downloadName string, err error) {
	c, ok := cache.Cache.Get(token)
	if !ok {
		return "", status.Errorf(errcode.ErrFileCacheFail, "cannot get cache with token(%v)", token)
	}
	d, ok := c.(dto.DownloadCache)
	if !ok {
		return "", status.Errorf(errcode.ErrFileCacheFail, "cache of token(%v) should be downloadCache", token)
	}
	if d.UserName != userName {
		return "", status.Errorf(errcode.ErrFileNoPermission, "different userName(%v) with userName(%v) in cache", userName, d.UserName)
	}
	return d.FileName, nil
}

// CheckDownloadPermission check download permission by ip address and username
func CheckDownloadPermission(ctx context.Context, ip string, userName string) bool {
	// todo 白名单配置
	//logger := logging.GetLogger(ctx)
	//
	//reply, _ := client.GetInstance().SysConfig.Config.GetDownloadConfig(ctx, &empty.Empty{})
	//if reply != nil && reply.Download.Enable {
	//	reply, err := client.GetInstance().SysConfig.Config.GetWhitelist(ctx, &sysconfig.GetWhitelistReq{})
	//	if err != nil {
	//		logger.Errorf("Get whitelist from db error %v", err)
	//		return false
	//	}
	//
	//	for _, whitelist := range reply.Whitelist {
	//		logger.Debugf("[CheckDownloadPermission] client_ip=%s, client_username=%s,"+
	//			" whitelist_ip=%s, whitelist_user=%s", ip, userName, whitelist.IpAddress, whitelist.Username)
	//		if (ip == whitelist.IpAddress && (whitelist.Username == "" || userName == whitelist.Username)) ||
	//			(userName == whitelist.Username && (whitelist.IpAddress == "" || ip == whitelist.IpAddress)) {
	//			return true
	//		}
	//	}
	//	return false
	//}

	return true
}

func validatePath(paths ...string) error {
	for _, path := range paths {

		if err := filePathLength(path); err != nil {
			return err
		}

		if err := checkInvalidChars(path, invalidCharsForPath); err != nil {
			return err
		}

	}

	return nil
}

func checkInvalidChars(path string, invalidChars string) error {
	for _, c := range invalidChars {

		if strings.Contains(path, string(c)) {
			return status.Errorf(errcode.ErrFilePathInvalidChar, "invalid character (%v) in %v", string(c), path)
		}
	}
	return nil
}

// GetFileInfoByStat 调用openapi获取单文件信息
func (srv *FileServiceImpl) GetFileInfoByStat(ctx context.Context, userName, path string, cross bool) (file *dto.File, err error) {

	openApiPath, err := srv.ConvertOpenApiReqPath(userName, path, cross)
	if err != nil {
		return nil, err
	}

	data, err := storage.Stat(srv.api, stat.Request{
		Path: openApiPath,
	})

	if err != nil {
		return nil, err
	}

	if data == nil || data.File == nil {
		return nil, err
	}

	var fileType string
	if data.File.IsDir {
		fileType = consts.FileTypeFolder
	} else {
		fileType = filepath.Ext(data.File.Name)
	}

	return &dto.File{
		Name:      data.File.Name,
		Size:      data.File.Size,
		MDate:     data.File.ModTime.Unix(),
		Type:      fileType,
		IsDir:     data.File.IsDir,
		IsSymLink: false,
		Path:      path,
		IsText:    false,
	}, nil
}

func filePathLength(path string) error {
	if len(path) > pathMaxLen {
		return status.Errorf(errcode.ErrFilePathTooLong, "Length of the file_path exceeds max allowed %v", pathMaxLen)
	}
	if len(filepath.Base(path)) > filenameMaxLen {
		return status.Errorf(errcode.ErrFilePathTooLong, "Length of the file_name exceeds max allowed %v", filenameMaxLen)
	}

	return nil
}

const (
	// pathMaxLen defines the max length of the file path
	pathMaxLen = 4096
	/*
		BTRFS   255 bytes
		exFAT   255 UTF-16 characters
		ext2    255 bytes
		ext3    255 bytes
		ext3cow 255 bytes
		ext4    255 bytes
		FAT32   8.3 (255 UCS-2 code units with VFAT LFNs)
		NTFS    255 characters
		XFS     255 bytes
	*/
	filenameMaxLen          = 255
	invalidCharsForPath     = "'\"\\,;`"
	invalidCharsForFileName = "'\"/\\,;`"
)

// Exist is the service that runs ls command to check path existence
func (srv *FileServiceImpl) Exist(ctx context.Context, userName string, cross bool, paths ...string) bool {
	for _, path := range paths {
		fileInfo, err := srv.GetFileInfoByStat(ctx, userName, path, cross)
		if err != nil || fileInfo == nil {
			logging.GetLogger(ctx).Infof("file(%v) not exist", path)
			return false
		}
	}

	return true
}

// CreateDir is a service to make one new directory with mkdir command
func (srv *FileServiceImpl) CreateDir(ctx context.Context, userName, path string, cross bool) error {
	if err := validatePath(path); err != nil {
		return status.Error(errcode.ErrFileNameValidate, "")
	}
	openApiPath, err := srv.ConvertOpenApiReqPath(userName, path, cross)
	if err != nil {
		return err
	}

	// 调用运算云的openapi
	if rsp, err := storage.Mkdir(srv.api, mkdir.Request{
		Path: openApiPath,
	}); err != nil {
		if rsp.ErrorCode == common.PaasErrorCodePathExists {
			return status.Error(errcode.ErrFileDirAlreadyExist, "")
		}
		return status.Error(errcode.ErrFileFailMkdir, err.Error())
	}
	tracelog.Info(ctx, fmt.Sprintf("create dir success, dir path[%s]", openApiPath))

	return nil
}

// Remove removes the files by rm command
func (srv *FileServiceImpl) Remove(ctx context.Context, userName, path string, cross bool) error {
	openApiPath, err := srv.ConvertOpenApiReqPath(userName, path, cross)
	if err != nil {
		return err
	}
	if _, err := storage.Rm(srv.api, rm.Request{
		Path: openApiPath,
	}); err != nil {
		return status.Error(errcode.ErrFileFailRemove, err.Error())
	}
	tracelog.Info(ctx, fmt.Sprintf("remove file success, file path[%s]", openApiPath))

	return nil
}

// Rename Rename
func (srv *FileServiceImpl) Rename(ctx context.Context, userName, path string, newPath string, overWrite, cross bool) error {
	if err := validatePath(newPath); err != nil {
		return status.Error(errcode.ErrFileNameValidate, "")
	}

	// check if there is a file that has the new name in the same dir
	if srv.Exist(ctx, userName, cross, newPath) {
		if !overWrite {
			return status.Errorf(errcode.ErrFileNameWasTaken, "%v exists", newPath)
		}
		err := srv.Remove(ctx, userName, newPath, cross)
		if err != nil {
			return err
		}
	}

	// 调用openapi的路径
	openApiSrcPath, err := srv.ConvertOpenApiReqPath(userName, path, cross)
	if err != nil {
		return err
	}
	openApiDestPath, err := srv.ConvertOpenApiReqPath(userName, newPath, cross)
	if err != nil {
		return err
	}

	if _, err := storage.Mv(srv.api, mv.Request{
		SrcPath:  openApiSrcPath,
		DestPath: openApiDestPath,
	}); err != nil {
		return fmt.Errorf("failed to rename [%v] to [%v], err: %v", openApiSrcPath, openApiDestPath, err)
	}
	tracelog.Info(ctx, fmt.Sprintf("rename file success, srcPath[%s], destPath[%s]", openApiSrcPath, openApiDestPath))

	return nil
}

// Move Move
// srcpaths should be existed files, dstpath should be existed directory
func (srv *FileServiceImpl) Move(ctx *gin.Context, userName string, cross, overwrite bool, destDir string, srcPaths ...string) error {

	// make sure source paths and destination path exist, and user has os permission
	exist := srv.Exist(ctx, userName, cross, append(srcPaths, destDir)...)

	if !exist {
		return status.Error(errcode.ErrFileNotExist, "file not exist")
	}

	// if overwrite is not set, check if there is a new path that was taken.
	for _, path := range srcPaths {
		// make sure destination path is not same as source path
		if strings.Compare(destDir, filepath.Join("/", filepath.Dir(path))) == 0 {
			logging.GetLogger(ctx).Infof("dstPath(%v) is same as srcPath(%v)", destDir, path)
			return status.Errorf(errcode.ErrFileSamePath, "dstPath(%v) is same as srcPath(%v)", destDir, path)
		}
		// mv: cannot move folder to a subdirectory of itself

		if strings.HasPrefix(destDir, path) {
			return status.Errorf(errcode.ErrFileMoveItself, "")
		}

		newPath := filepath.Join(destDir, filepath.Base(path))
		if srv.Exist(ctx, userName, cross, newPath) {
			if !overwrite {
				return status.Errorf(errcode.ErrFileAlreadyExist, "file(%v) already exist", newPath)
			}

			// 不能覆盖此文件的父目录
			if strings.HasPrefix(filepath.Join("/", path), newPath) {
				return status.Errorf(errcode.ErrrFileOverwriteParent, "can't overwrite parent file")
			}

			err := srv.Remove(ctx, userName, newPath, cross)
			if err != nil {
				return err
			}
		}

		// 调用运算云的openapi
		openApiSrcPath, err := srv.ConvertOpenApiReqPath(userName, path, cross)
		if err != nil {
			return err
		}

		openApiDestPath, err := srv.ConvertOpenApiReqPath(userName, newPath, cross)
		if err != nil {
			return err
		}

		if _, err := storage.Mv(srv.api, mv.Request{
			SrcPath:  openApiSrcPath,
			DestPath: openApiDestPath,
		}); err != nil {
			return fmt.Errorf("failed to move [%v] to [%v], err: %v", srcPaths, destDir, err)
		}

		tracelog.Info(ctx, fmt.Sprintf("move file success, srcPath[%s], destPath[%s]", openApiSrcPath, openApiDestPath))
	}

	return nil
}

func (srv *FileServiceImpl) Mv(ctx context.Context, userName string, cross, overwrite bool, srcpath, dstpath string) error {

	// make sure source paths and destination path exist, and user has os permission
	exist := srv.Exist(ctx, userName, cross, srcpath)

	if !exist {
		return status.Error(errcode.ErrFileNotExist, "file not exist")
	}

	// if overwrite is not set, check if there is a new path that was taken.
	if srv.Exist(ctx, userName, cross, dstpath) {
		// check if newPath is the same file as the file
		if !overwrite {
			return status.Errorf(errcode.ErrFileAlreadyExist, "file(%v) already exist", dstpath)
		}

		err := srv.Remove(ctx, userName, dstpath, cross)
		if err != nil {
			return err
		}
	}

	openApiSrcPath, err := srv.ConvertOpenApiReqPath(userName, srcpath, cross)
	if err != nil {
		return err
	}

	openApiDestPath, err := srv.ConvertOpenApiReqPath(userName, dstpath, cross)
	if err != nil {
		return err
	}

	// 调用运算云的openapi
	if _, err := storage.Mv(srv.api, mv.Request{
		SrcPath:  openApiSrcPath,
		DestPath: openApiDestPath,
	}); err != nil {
		return fmt.Errorf("failed to move [%v] to [%v], err: %v", srcpath, dstpath, err)
	}
	tracelog.Info(ctx, fmt.Sprintf("move file success, srcPath[%s], destPath[%s]", openApiSrcPath, openApiDestPath))

	return nil
}

func (srv *FileServiceImpl) HardLink(ctx context.Context, req *dto.LinkRequest) error {

	srcFilePaths := req.SrcFilePaths
	srcDirPaths := req.SrcDirPaths
	currentPath := req.CurrentPath
	filterPaths := req.FilterPaths
	dstPath := req.DstPath
	userName := req.UserName
	cross := req.Cross
	// make sure source paths and destination path exist, and user has os permission
	exist := srv.Exist(ctx, userName, cross, srcFilePaths...)

	if !exist {
		return status.Error(errcode.ErrFileNotExist, "file not exist")
	}

	// 所有需要上传的文件和文件夹
	// 前置准备，递归查询所有需要处理的文件和文件夹，文件夹先生成
	filePaths, _, err := srv.PreCopy(ctx, userName, currentPath, dstPath, srcFilePaths, srcDirPaths, cross)
	if err != nil {
		return err
	}

	// 过滤一部分文件
	filePaths = srv.filterFilePath(filePaths, currentPath, filterPaths)

	for _, path := range filePaths {
		// make sure destination path is not same as source path
		if strings.Compare(req.DstPath, filepath.Dir(path)) == 0 {
			logging.GetLogger(ctx).Infof("dstPath(%v) is same as openApiSrcPath(%v)", dstPath, path)
			err = status.Errorf(errcode.ErrFileSamePath, "dstPath(%v) is same as openApiSrcPath(%v)", dstPath, path)
			break
		}

		newPath := filepath.Join(dstPath, strings.TrimPrefix(path, req.CurrentPath))
		if srv.Exist(ctx, userName, cross, newPath) {
			if !req.Overwrite {
				continue
			}

			err = srv.Remove(ctx, userName, newPath, cross)
			if err != nil {
				return err
			}
		}

		openApiSrcPath, err := srv.ConvertOpenApiReqPath(userName, path, cross)
		if err != nil {
			return err
		}

		openApiDestPath, err := srv.ConvertOpenApiReqPath(userName, newPath, cross)
		if err != nil {
			return err
		}

		_, err = storage.Link(srv.api, link.Request{
			SrcPath:  openApiSrcPath,
			DestPath: openApiDestPath,
		})
		if err != nil {
			return err
		}

		tracelog.Info(ctx, fmt.Sprintf("link file success, srcPath[%s], destPath[%s]", openApiSrcPath, openApiDestPath))
	}

	return nil
}

func (srv *FileServiceImpl) Link(ctx context.Context, req *dto.LinkRequest) (err error) {

	srcFilePaths := req.SrcFilePaths
	dstPath := req.DstPath
	userName := req.UserName
	cross := req.Cross
	// make sure source paths and destination path exist, and user has os permission
	exist := srv.Exist(ctx, userName, cross, srcFilePaths...)

	if !exist {
		return status.Error(errcode.ErrFileNotExist, "file not exist")
	}

	for _, path := range srcFilePaths {
		// make sure destination path is not same as source path
		if strings.Compare(req.DstPath, filepath.Dir(path)) == 0 {
			logging.GetLogger(ctx).Infof("dstPath(%v) is same as openApiSrcPath(%v)", dstPath, path)
			err = status.Errorf(errcode.ErrFileSamePath, "dstPath(%v) is same as openApiSrcPath(%v)", dstPath, path)
			break
		}

		newPath := filepath.Join(dstPath, filepath.Base(path))
		if srv.Exist(ctx, userName, cross, newPath) {
			if !req.Overwrite {
				continue
			}

			err = srv.Remove(ctx, userName, newPath, cross)
			if err != nil {
				return err
			}
		}

		openApiSrcPath, err := srv.ConvertOpenApiReqPath(userName, path, cross)
		if err != nil {
			return err
		}

		openApiDestPath, err := srv.ConvertOpenApiReqPath(userName, newPath, cross)
		if err != nil {
			return err
		}

		_, err = storage.Link(srv.api, link.Request{
			SrcPath:  openApiSrcPath,
			DestPath: openApiDestPath,
		})
		if err != nil {
			return err
		}

		tracelog.Info(ctx, fmt.Sprintf("hard link file success, srcPath[%s], destPath[%s]", openApiSrcPath, openApiDestPath))
	}

	return nil
}

func (srv *FileServiceImpl) SymLink(ctx context.Context, req *dto.SymLinkRequest) error {
	srcPath := req.SrcPath
	dstPath := req.DstPath
	userName := req.UserName
	cross := req.Cross
	// make sure source paths and destination path exist, and user has os permission
	exist := srv.Exist(ctx, userName, cross, srcPath)

	if !exist {
		return status.Error(errcode.ErrFileNotExist, "file not exist")
	}

	if srv.Exist(ctx, userName, cross, dstPath) {
		if !req.Overwrite {
			return status.Error(errcode.ErrFileAlreadyExist, "")
		}

		err := srv.Remove(ctx, userName, dstPath, cross)
		if err != nil {
			return err
		}
	}

	openApiSrcPath, err := srv.ConvertOpenApiReqPath(userName, srcPath, cross)
	if err != nil {
		return err
	}

	openApiDestPath, err := srv.ConvertOpenApiReqPath(userName, dstPath, cross)
	if err != nil {
		return err
	}

	_, err = storage.Link(srv.api, link.Request{
		SrcPath:  openApiSrcPath,
		DestPath: openApiDestPath,
	})
	if err != nil {
		return err
	}
	tracelog.Info(ctx, fmt.Sprintf("symlink file success, srcPath[%s], destPath[%s]", openApiSrcPath, openApiDestPath))

	return nil
}

func (srv *FileServiceImpl) PreCopy(ctx context.Context, userName, currentPath, dstPath string, srcFilePaths, srcDirPaths []string, cross bool) ([]string, []string, error) {

	filePaths := srcFilePaths
	dirPaths := srcDirPaths
	// 递归文件夹
	if len(srcDirPaths) > 0 {
		files, err := srv.ListOfRecur(ctx, userName, srcDirPaths, cross, true, nil)
		if err != nil {
			return nil, nil, status.Errorf(errcode.ErrFileList, "")
		}
		if len(files) > 0 {
			// 查询到的文件夹和srcDir加入needMkdirPaths;查询到的文件和srcFile加入needUploadFiles
			for _, file := range files {
				if file.IsDir {
					dirPaths = append(dirPaths, file.Path)
				} else {
					filePaths = append(filePaths, file.Path)
				}
			}
		}
	}

	// mkdir所有文件夹
	if len(dirPaths) > 0 {
		dirPaths = collectionutil.RemoveDuplicates(dirPaths)
		dirPaths = collectionutil.RemoveString(dirPaths, currentPath)
		for _, dirPath := range dirPaths {
			err := srv.CreateDir(ctx, userName, path.Join(dstPath, strings.TrimPrefix(path.Clean(dirPath), path.Clean(currentPath))), cross)
			if err != nil {
				if status.Code(err) == errcode.ErrFileDirAlreadyExist {
					continue
				}
				return nil, nil, err
			}
		}
	}

	return filePaths, dirPaths, nil
}

func (srv *FileServiceImpl) Copy(ctx *gin.Context, req *dto.CopyRequest) (string, error) {
	srcFilePaths := req.SrcFilePaths
	srcDirPaths := req.SrcDirPaths
	currentPath := req.CurrentPath
	dstPath := req.DstPath
	userName := req.UserName
	cross := req.Cross
	// make sure source paths and destination path exist, and user has os permission
	exist := srv.Exist(ctx, userName, cross, srcFilePaths...)

	if !exist {
		return "", status.Error(errcode.ErrFileNotExist, "file not exist")
	}

	// 前置准备，递归查询所有需要处理的文件和文件夹，文件夹先生成
	filePaths, _, err := srv.PreCopy(ctx, userName, currentPath, dstPath, srcFilePaths, srcDirPaths, cross)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%s:%s", common.CopyFilePreKey, uuid.New())
	cache.Cache.Set(key, dto.CopyStateCopying, gocache.NoExpiration)

	go func() {
		var err error
		// if overwrite is not set, check if there is a new path that was taken.
		for _, path := range filePaths {
			// make sure destination path is not same as source path
			if strings.Compare(dstPath, filepath.Dir(path)) == 0 {
				logging.GetLogger(ctx).Infof("dstPath(%v) is same as srcPath(%v)", dstPath, path)
				err = status.Errorf(errcode.ErrFileSamePath, "dstPath(%v) is same as srcPath(%v)", dstPath, path)
				break
			}
			newPath := filepath.Join(dstPath, strings.TrimPrefix(path, req.CurrentPath))
			if srv.Exist(ctx, userName, cross, newPath) {
				if !req.Overwrite {
					continue
				}
			}

			err = srv.CopyRange(ctx, userName, path, newPath, cross, req.Overwrite)
			if err != nil {
				break
			}
		}

		var content string
		msg := &pbnotice.WebsocketMessage{
			UserId: snowflake.ID(ginutil.GetUserID(ctx)).String(),
			Type:   common.ShareFileEventType,
		}

		allSrcPaths := collectionutil.MergeSlice(req.SrcFilePaths, req.SrcDirPaths)
		if err != nil {
			cache.Cache.Set(key, dto.CopyStateFailure, time.Duration(10)*time.Second)
			content = fmt.Sprintf("%v发送的文件%v保存失败", strings.Split(allSrcPaths[0], "/")[0], allSrcPaths)
		} else {
			cache.Cache.Set(key, dto.CopyStateSuccess, time.Duration(10)*time.Second)
			content = fmt.Sprintf("%v发送的文件%v保存完成", strings.Split(allSrcPaths[0], "/")[0], allSrcPaths)
		}
		msg.Content = content
		_, _ = client.GetInstance().Notice.SendWebsocketMessage(ctx, msg)
	}()

	return key, nil
}

// Get is a service to get the file properties
func (srv *FileServiceImpl) Get(ctx context.Context, userName, path string, cross bool) (file *dto.File, err error) {

	// To form the structure of file
	file, err = srv.GetFileInfoByStat(ctx, userName, path, cross)
	if err != nil {
		return nil, status.Errorf(errcode.ErrFileGetInfo, "failed to get user(%v) info of file(%v), err: %v", userName, path, err)
	}

	return file, nil
}

// Read reads the file content by executing cat "filepath"
func (srv *FileServiceImpl) Read(ctx context.Context, userName, path string, offset int64, len int64, cross bool) ([]byte, error) {
	content := make([]byte, len)
	resolver := func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("download file error: %v", resp.StatusCode))
		}
		var err error
		content, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return nil
	}
	// 调用运算云的openapi
	openApiPath, err := srv.ConvertOpenApiReqPath(userName, path, cross)
	if err != nil {
		return nil, err
	}

	_, err = storage.ReadAt(srv.api, readAt.Request{
		Path:   openApiPath,
		Offset: offset,
		Length: len,
	}, resolver)

	if err != nil {
		return nil, status.Errorf(errcode.ErrFileDownload, "read file error: %v", err)
	}

	return content, nil
}

func (srv *FileServiceImpl) List(ctx context.Context, userName, dir string, page *xtype.Page, cross bool, showHideFile bool, filterRegexps []string) (files []*dto.File, err error) {

	var fileType string
	var scroll bool
	if page == nil {
		page = &xtype.Page{
			Index: 0,
			Size:  consts.DefPageSize,
		}
		scroll = true
	}

	filterRegex := ""
	if !showHideFile {
		filterRegex = consts.DefaultFilterHideFileRegex

		filterHideFileRegex := config.GetConfig().FilterHideFileRegex
		if filterHideFileRegex != "" {
			filterRegex = filterHideFileRegex
		}
		filterRegexps = append(filterRegexps, filterRegex)
	}

	data, err := srv.OnePageFileList(userName, dir, filterRegexps, page, scroll, cross)
	if err != nil {
		return nil, err
	}

	for _, info := range data.Data.Files {
		if info.IsDir {
			fileType = consts.FileTypeFolder
		} else {
			fileType = filepath.Ext(info.Name)
		}

		fileInfo := dto.File{
			Name:  info.Name,
			Size:  info.Size,
			MDate: info.ModTime.Unix(),
			Type:  fileType,
			IsDir: info.IsDir,
			Path:  path.Join(dir, info.Name),
		}

		files = append(files, &fileInfo)
	}

	return
}

// OnePageFileList 每次查询一页量的文件列表 如果scroll:true 则会查询下一页(如果还存在下一页的话)
func (srv *FileServiceImpl) OnePageFileList(userName, dir string, filterRegexps []string, page *xtype.Page, scroll, cross bool) (*ls.Response, error) {
	openApiPath, err := srv.ConvertOpenApiReqPath(userName, dir, cross)
	if err != nil {
		return nil, err
	}
	data, err := storage.Ls(srv.api, ls.Request{
		Path:             openApiPath,
		FilterRegexpList: filterRegexps,
		PageOffset:       page.Index,
		PageSize:         page.Size,
	})
	if err != nil {
		return nil, err
	}

	if scroll && len(data.Data.Files) == consts.DefPageSize {
		page.Index = page.Index + consts.DefPageSize
		nextData, err := srv.OnePageFileList(userName, dir, filterRegexps, page, true, cross)
		if err != nil {
			return nil, status.Error(errcode.ErrFilePaas, "")
		}
		if len(data.Data.Files) > 0 {
			data.Data.Files = append(data.Data.Files, nextData.Data.Files...)
		}
	}

	return data, nil
}

func (srv *FileServiceImpl) ListOfRecur(ctx context.Context, name string, folderPaths []string, cross bool, showHideFile bool, filterRegexp []string) ([]*dto.File, error) {
	var files []*dto.File
	// 遍历每个文件夹路径
	for _, folderPath := range folderPaths {
		childFiles, err := srv.getFiles(ctx, name, folderPath, cross, showHideFile, filterRegexp)
		if err != nil {
			return nil, err
		}
		files = append(files, childFiles...)
	}

	return files, nil
}

func (srv *FileServiceImpl) getFiles(ctx context.Context, name string, folderPath string, cross, showHideFile bool, filterRegexp []string) ([]*dto.File, error) {
	var files []*dto.File

	// 读取文件夹下的所有文件和文件夹
	files, err := srv.List(ctx, name, folderPath, nil, cross, showHideFile, filterRegexp)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range files {
		// 如果是文件夹，则递归获取其下的文件和文件夹信息
		if fileInfo.IsDir {
			subFiles, err := srv.getFiles(ctx, name, filepath.Join(folderPath, fileInfo.Name), cross, showHideFile, filterRegexp)
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		}
	}

	return files, nil
}

// Write writes bytes to existing files.
func (srv *FileServiceImpl) Write(ctx *gin.Context, userName string, path string, fileSize int64, offset int64, sliceData []byte) error {

	if fileSize < 0 || offset < 0 || offset+int64(len(sliceData)) > fileSize {
		return status.Error(errcode.ErrFileArgs, "file_size < 0 or offset < 0 or offset+len(slice_data) > file_size")
	}

	if err := validatePath(path); err != nil {
		return err
	}

	// todo 等远算云提供write接口

	return nil
}

// Realpath 获取文件真实路径
func (srv *FileServiceImpl) Realpath(ctx context.Context, relativePath string) (string, error) {
	response, err := storage.Realpath(srv.api, &realpath.Request{
		RelativePath: relativePath,
	})
	if err != nil {
		return "", err
	}

	return response.Data.RealPath, nil
}

func (srv *FileServiceImpl) executeUploadTask(ctx context.Context, taskKey string, tasks map[string]*dto.HPCUploadTask, userName string, cross bool) {
	logger := logging.Default()

	for uploadID, task := range tasks {
		if task.State != dto.UploadStatePending {
			continue
		}

		task.Update(taskKey, dto.UploadStateUploading, "")
		// 并发分片上传
		success, err := srv.uploadSliceConcurrent(taskKey, uploadID, userName, task)
		if err != nil {
			task.Update(taskKey, dto.UploadStateFailure, "")
			task.ErrMsg = errcode.GetErrCodeMessage(err, errcode.ErrFileUploadHPCFile)
			logger.Infof("[%s] task file [%s] upload fail: src:[%s] dest:[%v], uploadID:[%s], err:[%v]", userName, task.FileName, task.SrcPath, task.DestPath, uploadID, err)
			continue
		}

		if success {
			err = srv.uploadComplete(userName, task.DestPath, uploadID, cross)
			if err != nil {
				logger.Infof("[%s] task file [%s] upload fail: src:[%s] dest:[%v], uploadID:[%s], err:[%v]", userName, task.FileName, task.SrcPath, task.DestPath, err)
				task.Update(taskKey, dto.UploadStateFailure, errcode.GetErrCodeMessage(err, errcode.ErrFileUploadHPCFile))
				continue
			}
			logger.Infof("[%s] task file [%s] upload success: src:[%s] dest:[%v], uploadID:[%s]", userName, task.FileName, task.SrcPath, task.DestPath, uploadID)
			tracelog.Info(ctx, fmt.Sprintf("upload hpc file success, srcPath[%s], destPath[%s]", task.SrcPath, task.DestPath))

			task.Update(taskKey, dto.UploadStateSuccess, "")
		}
	}
}

func (srv *FileServiceImpl) QueryUploadHPCFileTask(ctx context.Context, srcFilePaths, srcDirPaths []string, userName string, cross bool) (*dto.QueryHPCFileTaskResponse, error) {
	// 所有需要上传的文件和文件夹
	fileTasks := make(map[string]*dto.HPCUploadTask, 0)
	dirTasks := srcDirPaths

	if len(srcFilePaths) > 0 {
		for _, filePath := range srcFilePaths {
			file, err := srv.Get(ctx, userName, filePath, false)
			if err != nil {
				return nil, status.Errorf(errcode.ErrFileNotExist, "")
			}
			fileTasks[file.Path] = util.ConvertHPCUploadTask(file)
		}
	}

	if len(srcDirPaths) > 0 {
		files, err := srv.ListOfRecur(ctx, userName, srcDirPaths, false, true, nil)
		if err != nil {
			return nil, status.Errorf(errcode.ErrFileList, "")
		}
		if len(files) > 0 {
			// 查询到的文件夹和srcDir加入needMkdirPaths;查询到的文件和srcFile加入needUploadFiles
			for _, file := range files {
				if file.IsDir {
					dirTasks = append(dirTasks, file.Path)
				} else {
					fileTasks[file.Path] = util.ConvertHPCUploadTask(file)
				}
			}
		}
	}

	return &dto.QueryHPCFileTaskResponse{
		DirTasks:  dirTasks,
		FileTasks: fileTasks,
	}, nil
}

func (srv *FileServiceImpl) GetUploadHPCFileTask(ctx context.Context, taskKey string) ([]*dto.HPCUploadTaskResponse, error) {
	list := make([]*dto.HPCUploadTaskResponse, 0)

	tasks, ok := cache.Cache.Get(taskKey)
	if !ok || tasks == nil {
		return list, nil
	}

	// 考虑到有另一个协程正在处理tasks，如果其中元素正在被删除会导致遍历报错，这里选择拷贝一份再遍历
	// json.Marshal使用了反射，当另一个协程在修改tasks元素时，会导致数据竞争问题,所以要加锁
	taskLock(taskKey)
	jsonStr, err := json.Marshal(tasks)
	taskUnLock(taskKey)
	if err != nil {
		return nil, err
	}

	copyTasks := make(map[string]*dto.HPCUploadTask)
	err = json.Unmarshal(jsonStr, &copyTasks)
	if err != nil {
		return nil, err
	}

	for uploadID, task := range copyTasks {
		list = append(list, &dto.HPCUploadTaskResponse{
			UploadID:    uploadID,
			FileName:    task.FileName,
			SrcPath:     task.SrcPath,
			DestPath:    task.DestPath,
			TotalSize:   task.TotalSize,
			CurrentSize: task.CurrentSize,
			State:       task.State,
			ErrMsg:      task.ErrMsg,
		})
	}

	// 优先根据状态排序，状态一致根据文件名排序
	if len(list) > 0 {
		sort.Slice(list, func(i, j int) bool {
			if list[i].State != list[j].State {
				return list[i].State < list[j].State
			}
			return list[i].FileName < list[j].FileName
		})
	}

	return list, nil
}

func (srv *FileServiceImpl) CancelUploadHPCFileTask(ctx *gin.Context, taskKey string, uploadID string) error {
	tasks, ok := cache.Cache.Get(taskKey)
	if !ok {
		return status.Error(errcode.ErrFileUploadHPCFileNotExist, "")
	}
	task := (tasks.(map[string]*dto.HPCUploadTask))[uploadID]

	if task != nil {
		logging.Default().Debugf("[%s] cancel hpc upload task,from [%s] to [%s], uploadID:[%s]", ginutil.GetUserName(ctx), task.SrcPath, task.DestPath, uploadID)
		task.Update(taskKey, dto.UploadStateCancel, "")
		task.CancelTask(taskKey)
	}
	return nil
}

func (srv *FileServiceImpl) ResumeUploadHPCFileTask(ctx *gin.Context, taskKey string, uploadID string) error {
	tasks, ok := cache.Cache.Get(taskKey)
	if !ok {
		return status.Error(errcode.ErrFileUploadHPCFileNotExist, "")
	}
	task := (tasks.(map[string]*dto.HPCUploadTask))[uploadID]
	if task != nil {
		logging.Default().Debugf("[%s]resume hpc upload task,from [%s] to [%s], uploadID:[%s]", ginutil.GetUserName(ctx), task.SrcPath, task.DestPath, uploadID)
		task.Update(taskKey, dto.UploadStatePending, "")
		task.ErrMsg = ""
		task.CurrentSize = 0
		cache.Cache.Set(taskKey, tasks, gocache.NoExpiration)
	}

	return nil
}

func (srv *FileServiceImpl) AbortAllTask(ctx *gin.Context, taskKey string) error {
	tasks, ok := cache.Cache.Get(taskKey)
	if !ok {
		return status.Error(errcode.ErrFileUploadHPCFileNotExist, "")
	}
	logging.Default().Debugf("[%s] Abort Task:[%s]", ginutil.GetUserName(ctx), taskKey)

	for _, task := range tasks.(map[string]*dto.HPCUploadTask) {
		task.Update(taskKey, dto.UploadStateCancel, "")
		task.Cancel()
	}

	return nil
}

func (srv *FileServiceImpl) UploadHPCFile(ctx context.Context, req *dto.UploadHPCRequest, taskReq *dto.QueryHPCFileTaskResponse) (string, error) {
	dirTasks := taskReq.DirTasks
	fileTasks := taskReq.FileTasks

	if len(fileTasks) <= 0 {
		return "", status.Errorf(errcode.ErrFileUploadHPCFileEmpty, "")
	}

	// mkdir所有文件夹
	if len(dirTasks) > 0 {
		dirTasks = collectionutil.RemoveDuplicates(dirTasks)
		for _, dirPath := range dirTasks {
			currentDirPath := path.Join(req.DestDirPath, strings.TrimPrefix(path.Clean(dirPath), path.Clean(req.CurrentPath)))
			if !req.Cross {
				currentDirPath = path.Join(strings.TrimPrefix(path.Clean(dirPath), path.Clean(req.CurrentPath)))
			}
			err := srv.CreateDir(ctx, req.UserName, currentDirPath, req.Cross)
			if err != nil {
				if status.Code(err) == errcode.ErrFileDirAlreadyExist {
					continue
				}
				return "", err
			}
		}
	}

	tasks := make(map[string]*dto.HPCUploadTask, 0)

	// uploadInit所有needUploadFiles，生成缓存
	for filePath, fileInfo := range fileTasks {
		destFilePath := path.Join(req.DestDirPath, strings.TrimPrefix(path.Clean(filePath), path.Clean(req.CurrentPath)))
		if !req.Cross {
			destFilePath = path.Join(strings.TrimPrefix(path.Clean(filePath), path.Clean(req.CurrentPath)))
		}

		// 是否覆盖
		exist := srv.Exist(ctx, req.UserName, req.Cross, destFilePath)
		if exist {
			if req.Overwrite {
				err := srv.Remove(ctx, req.UserName, destFilePath, req.Cross)
				if err != nil {
					return "", status.Error(errcode.ErrFileFailCover, "")
				}
			} else {
				return "", status.Error(errcode.ErrFileAlreadyExist, "")
			}
		}
		// 获取uploadID
		rsp, err := srv.PreUpload(ctx, &dto.PreUploadRequest{
			Path:     destFilePath,
			FileSize: fileInfo.TotalSize,
			UserName: req.UserName,
			Cross:    req.Cross,
		})
		if err != nil {
			return "", err
		}
		// 用于中断任务
		cancelCtx, cancel := context.WithCancel(context.Background())
		fileInfo.CancelCtx = cancelCtx
		fileInfo.Cancel = cancel
		fileInfo.DestPath = destFilePath
		tasks[rsp.UploadId] = fileInfo
	}

	if len(tasks) <= 0 {
		return "", status.Errorf(errcode.ErrFileUploadHPCFileEmpty, "")
	}

	// 放入缓存,为了不影响同用户提交的历史任务，加上uuid
	cacheKey := fmt.Sprintf("%s:%s:%s", common.HpcUploadTaskPreKey, req.UserName, uuid.New().String())

	cache.Cache.Set(cacheKey, tasks, gocache.NoExpiration)

	// 耗时操作，上传文件
	go func() {
		for {
			// 执行上传任务,捞取pending任务
			// 这里从cache取值是防止cache过期被删除但原tasks仍被引用导致循环无法中断
			cacheTasks, ok := cache.Cache.Get(cacheKey)
			if !ok {
				return
			}

			tasks = cacheTasks.(map[string]*dto.HPCUploadTask)
			srv.executeUploadTask(ctx, cacheKey, tasks, req.UserName, req.Cross)
			// 等待10s后移除成功与取消的任务(如果任务完成太快直接删除，用户可能无法查询到任务列表)；失败任务等待人工恢复
			// 也可以当做定时任务每10s执行一轮
			time.Sleep(time.Duration(10) * time.Second)

			removeEndTask(cacheKey, tasks)

			// 任务列表为空，删除缓存，退出循环
			if len(tasks) == 0 {
				cache.Cache.Delete(cacheKey)
				break
			}

			_, ttl, _ := cache.Cache.GetWithExpiration(cacheKey)
			// 如果没有过期时间，给个过期时间
			if time.Time.IsZero(ttl) {
				// 如果用户在executeUploadTask到这之间(10s内)执行了恢复上传(会将缓存重新设置为不过期)，则会出现意料外的未开始执行任务却给了过期时间,所以需要判断下是否有待上传任务
				var notSetExpire bool
				for _, task := range tasks {
					if task.State == dto.UploadStatePending {
						notSetExpire = true
						break
					}
				}
				if !notSetExpire {
					cache.Cache.Set(cacheKey, tasks, time.Duration(config.GetConfig().HpcUploadConfig.WaitResumeTime)*time.Minute)
				}
			}
		}
	}()

	return cacheKey, nil
}

func removeEndTask(taskKey string, tasks map[string]*dto.HPCUploadTask) {
	deleteKeys := make([]string, 0)
	for key, task := range tasks {
		if task.State == dto.UploadStateSuccess || task.State == dto.UploadStateCancel {
			deleteKeys = append(deleteKeys, key)
		}
	}

	taskLock(taskKey)
	for _, key := range deleteKeys {
		delete(tasks, key)
	}
	taskUnLock(taskKey)
}

func taskLock(taskKey string) {
	for {
		// 锁只是保险用，如果redis不可用，不阻塞业务
		successFlag, err := lockutil.TryLock(taskKey)
		if err != nil {
			successFlag = true
		}
		if successFlag {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func taskUnLock(taskKey string) {
	lockutil.UnLock(taskKey)
}

func (srv *FileServiceImpl) uploadSliceConcurrent(taskKey, uploadID, userName string, task *dto.HPCUploadTask) (bool, error) {
	logger := logging.Default()
	var eg errgroup.Group
	var offset int64
	var err error

	hpcConfig := config.GetConfig().HpcUploadConfig
	// 控制上传并发数
	sem := make(chan struct{}, hpcConfig.ConcurrencyLimit)

	// hpc文件目录转为本地文件目录 /path/to/file => /homePath/{userName}/path/to/file
	localFilePath := srv.ConvertLocalPath(task.SrcPath, userName, false)
	blockSize := int64(1024 * 1024 * hpcConfig.BlockSize)

	file, err := os.Open(localFilePath)
	if err != nil {
		return false, status.Error(errcode.ErrFileNotExist, "")
	}
	defer file.Close()
	// 用来控制
	var onceErr atomic.Bool

	for {
		if onceErr.Load() || offset >= task.TotalSize {
			break
		}

		sem <- struct{}{}
		select {
		// 用户取消了任务
		case <-task.CancelCtx.Done():
			task.Update(taskKey, dto.UploadStateCancel, "")
			<-sem
			return false, nil
		default:
			// 计算offset和分片大小
			chunkSize := blockSize
			if task.TotalSize-offset < blockSize {
				chunkSize = task.TotalSize - offset
			}
			currentOffset := offset
			offset += chunkSize

			// 并发分片上传
			eg.Go(
				func() error {
					logger.Debugf("[%s] task file [%s] begin upload: src:[%s] dest:[%v], offset:[%v], chunkSize:[%v], uploadID:[%s]", userName, task.FileName, task.SrcPath, task.DestPath, currentOffset, chunkSize, uploadID)

					err = retry.Do(
						func() error {
							return srv.uploadChunk(currentOffset, chunkSize, file, uploadID)
						},
						// 重试次数
						retry.Attempts(uint(hpcConfig.RetryCount)),
						// 设置重试间隔类型为固定间隔
						retry.DelayType(retry.FixedDelay),
						// 重试间隔
						retry.Delay(time.Duration(hpcConfig.RetryDelay)*time.Second),
					)

					if err != nil {
						onceErr.Store(true)
						return err
					}
					atomic.AddInt64(&task.CurrentSize, chunkSize)
					<-sem
					return nil
				})

		}
	}

	// 阻塞直到此文件所有协程结束
	err = eg.Wait()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (srv *FileServiceImpl) uploadChunk(offset int64, chunkSize int64, file *os.File, uploadID string) error {
	buffer := make([]byte, chunkSize)
	_, err := file.ReadAt(buffer, offset)
	if err != nil {
		return err
	}
	_, err = storage.UploadSlice(srv.api, uploadslice.Request{
		UploadID: uploadID,
		Offset:   offset,
		Length:   chunkSize,
		Slice:    buffer,
	})
	return err
}

func (srv *FileServiceImpl) uploadComplete(userName, path, uploadID string, cross bool) error {
	openApiPath, err := srv.ConvertOpenApiReqPath(userName, path, cross)
	if err != nil {
		return err
	}

	_, err = storage.UploadComplete(srv.api, complete.Request{
		Path:     openApiPath,
		UploadID: uploadID,
	})
	return err
}

// ConvertLocalPath hpc文件路径转化为本机存储路径
func (srv *FileServiceImpl) ConvertLocalPath(hpcPath, userName string, cross bool) string {
	rootPath := config.GetConfig().LocalRootPath

	if !cross {
		return path.Join(rootPath, userName, hpcPath)
	}

	return path.Join(rootPath, hpcPath)
}

func (srv *FileServiceImpl) CopyRange(ctx context.Context, userName, path, newPath string, cross, overwrite bool) error {

	fileInfo, err := srv.GetFileInfoByStat(ctx, userName, path, cross)
	if err != nil {
		return err
	}

	srcPath, err := srv.ConvertOpenApiReqPath(userName, path, cross)
	if err != nil {
		return err
	}
	destPath, err := srv.ConvertOpenApiReqPath(userName, newPath, cross)
	if err != nil {
		return err
	}

	totalSize := fileInfo.Size

	_, err = storage.Create(srv.api, create.Request{
		Path:      destPath,
		Size:      totalSize,
		Overwrite: overwrite,
	})

	if err != nil {
		return err
	}
	hpcConfig := config.GetConfig().HpcUploadConfig

	var offset int64
	blockSize := int64(1024 * 1024 * hpcConfig.BlockSize)

	for {
		if offset >= totalSize {
			break
		}

		chunkSize := blockSize
		if totalSize-offset < blockSize {
			chunkSize = totalSize - offset
		}

		_, err = storage.CopyRange(srv.api, copyRange.Request{
			SrcPath:    srcPath,
			DestPath:   destPath,
			SrcOffset:  offset,
			DestOffset: offset,
			Length:     chunkSize,
		})

		offset += chunkSize
		if err != nil {
			return err
		}

	}
	tracelog.Info(ctx, fmt.Sprintf("copy file success, srcPath[%s], destPath[%s]", srcPath, destPath))

	return nil
}

func (srv *FileServiceImpl) GetRecordList(ctx *gin.Context, userID int64, req dto.GetRecordListRequest) (*dto.ShareRecordListResponse, error) {
	recordList, total, err := srv.shareFileRecordDao.GetFileUserList(userID, req.Page, req.Filter)

	if err != nil {
		return nil, err
	}

	for _, record := range recordList {
		shareType := "分享"
		if record.ShareType == 1 {
			shareType = "发送"
		}
		id, _ := strconv.Atoi(record.Id)
		record.Id = snowflake.ID(id).String()
		record.Content = fmt.Sprintf("您收到一份来自[%s]%s的文件[%s]", record.Owner, shareType, filepath.Base(record.FilePath))
	}
	return &dto.ShareRecordListResponse{
		Page: &xtype.PageResp{
			Index: req.Page.Index,
			Size:  req.Page.Size,
			Total: total,
		},
		List: recordList,
	}, nil
}

func (srv *FileServiceImpl) GenerateShareLink(ctx *gin.Context, req dto.GenerateShareRequest) error {
	userName := ginutil.GetUserName(ctx)

	if !srv.Exist(ctx, userName, false, req.ShareFilePath) {
		return status.Error(errcode.ErrFileNotExist, "")
	}

	share := model.ShareFileRecord{
		FilePath: req.ShareFilePath,
		Owner:    userName,
		Type:     int8(req.ShareType),
	}

	recordID, err := srv.shareFileRecordDao.Add(share)
	if err != nil {
		return err
	}

	for _, userID := range req.ShareUserList {
		shareFileUser := model.ShareFileUser{
			ShareRecordId: recordID,
			UserId:        snowflake.MustParseString(userID),
		}
		if err := srv.shareFileRecordDao.AddShareFileUser(shareFileUser); err != nil {
			return err
		}
	}

	//for _, userID := range req.ShareUserList {
	//	content := fmt.Sprintf("您收到一份来自[%s]的分享文件:%s(%s)", userName, filepath.Base(req.ShareFilePath), share.Code)
	//	msg := &pbnotice.WebsocketMessage{
	//		UserId:  userID,
	//		Type:    common.ShareFileEventType,
	//		Content: content,
	//	}
	//
	//	if _, err = client.GetInstance().Notice.SendWebsocketMessage(ctx, msg); err != nil {
	//		return err
	//	}
	//}

	return nil
}

func (srv *FileServiceImpl) GetShareFile(ctx *gin.Context, id snowflake.ID) (*dto.ShareFileInfo, error) {
	shareInfo, ok, err := srv.shareFileRecordDao.Get(id)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, status.Error(errcode.ErrFileShareNotExist, "")
	}

	fileInfo, err := srv.GetFileInfoByStat(ctx, shareInfo.Owner, shareInfo.FilePath, false)
	if err != nil && status.Code(err) == errcode.ErrFileNotExist {
		err = status.Error(errcode.ErrFileShareNotExist, "")
		return nil, err
	}

	return &dto.ShareFileInfo{
		Name:      fileInfo.Name,
		Size:      fileInfo.Size,
		MDate:     fileInfo.MDate,
		Type:      fileInfo.Type,
		IsDir:     fileInfo.IsDir,
		Path:      filepath.Join(shareInfo.Owner, fileInfo.Path),
		ShareType: shareInfo.Type,
	}, nil
}

func (srv *FileServiceImpl) CheckSharePath(ctx context.Context, userName string, cross bool) error {
	if cross {
		return nil
	}

	exist := srv.Exist(ctx, userName, false, common.PublicFolderPath)
	if config.GetConfig().PublicFolderEnable {
		// 如果公共目录开启且目录不存在，给用户创建软连接
		if !exist {
			return srv.SymLink(ctx, &dto.SymLinkRequest{
				Overwrite: true,
				Cross:     true,
				UserName:  userName,
				SrcPath:   common.PublicFolderPath,
				DstPath:   filepath.Join(userName, common.PublicFolderPath),
			})
		}
	} else {
		// 如果未开启且目录存在，则删除用户共享软连接目录
		if exist {
			return srv.Remove(ctx, userName, common.PublicFolderPath, false)
		}
	}

	return nil
}

func (srv *FileServiceImpl) filterFilePath(originFilePaths []string, currentPath string, filterPaths []string) []string {
	if filterPaths == nil || len(filterPaths) == 0 {
		return originFilePaths
	}

	for i, filterPath := range filterPaths {
		filterPaths[i] = filepath.Join(currentPath, filterPath)
	}

	originFilePaths = collectionutil.RemoveStrings(originFilePaths, filterPaths)

	return originFilePaths
}

const (
	KB int64 = 1024
	MB       = 1024 * 1024
	GB       = 1024 * 1024 * 1024

	KbStr = "KB"
	MbStr = "MB"
	GbStr = "GB"
)

// GetStorageSize 获取存储空间大小
func GetStorageSize(storageSize int64) string {
	if storageSize == 0 {
		return "0B"
	}

	if storageSize > GB {
		return fmt.Sprintf("%.2f", float64(storageSize)/float64(GB)) + GbStr
	}
	if storageSize > MB {
		return fmt.Sprintf("%.2f", float64(storageSize)/float64(MB)) + MbStr
	}
	if storageSize > KB || storageSize >= 0 {
		return fmt.Sprintf("%.2f", float64(storageSize)/float64(KB)) + KbStr
	}

	return ""
}

func buildFileLogHeader() []csvutil.CsvHeaderEntity {
	return []csvutil.CsvHeaderEntity{
		{
			Name:   "文件名称",
			Column: "FileName",
		},
		{
			Name:   "文件路径",
			Column: "FilePath",
		},
		{
			Name:   "文件类型",
			Column: "FileType",
			Converter: func(i interface{}) string {
				return i.(dto.FileTypeEnum).String()
			},
		},
		{
			Name:   "操作类型",
			Column: "OperateType",
			Converter: func(i interface{}) string {
				return i.(dto.OpTypeEnum).String()
			},
		},
		{
			Name:   "文件大小",
			Column: "StorageSize",
			Converter: func(i interface{}) string {
				return GetStorageSize(i.(int64))
			},
		},
		{
			Name:   "操作时间",
			Column: "OperateTime",
			Converter: func(i interface{}) string {
				return csvutil.CSVFormatTime(i.(time.Time), common.DatetimeFormat, common.Bar)
			},
		},
	}
}
