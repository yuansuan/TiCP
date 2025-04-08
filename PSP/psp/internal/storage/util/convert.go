package util

import (
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/storage"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/dto"
)

// ConvertFiles ConvertFiles
func ConvertFiles(files []*dto.File) []*storage.File {
	var result []*storage.File

	for _, f := range files {
		result = append(result, ConvertFile(f))
	}

	return result
}

// ConvertFile ConvertFile
func ConvertFile(f *dto.File) *storage.File {
	return &storage.File{
		Name:      f.Name,
		Size:      f.Size,
		MDate:     f.MDate,
		Type:      f.Type,
		IsDir:     f.IsDir,
		IsSymLink: f.IsSymLink,
		Path:      f.Path,
		IsText:    f.IsText,
	}
}

func ConvertFileRsp(f *dto.File) *dto.ListResponse {
	return &dto.ListResponse{
		Name:     f.Name,
		Mode:     f.Mode,
		Size:     f.Size,
		MDate:    f.MDate,
		Type:     f.Type,
		IsDir:    f.IsDir,
		Path:     f.Path,
		IsText:   f.IsText,
		OnlyRead: f.OnlyRead,
	}
}

func ConvertFileRsps(files []*dto.File) []*dto.ListResponse {
	var result []*dto.ListResponse

	for _, f := range files {
		result = append(result, ConvertFileRsp(f))
	}
	return result
}

func ConvertHPCUploadTask(file *dto.File) *dto.HPCUploadTask {
	return &dto.HPCUploadTask{
		FileName:    file.Name,
		SrcPath:     file.Path,
		TotalSize:   file.Size,
		CurrentSize: 0,
		State:       dto.UploadStatePending,
	}
}
