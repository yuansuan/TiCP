package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/oshelp"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/storage/client"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/util"
)

type Downloader interface {
	// Download
	// src: http://storage.domain/dirA/fileA
	// dst: /compute_prefix/dirA/fileA
	Download(ctx context.Context, src string, dst string) error
}

type localDownloader struct {
	cli *client.Client
}

func (d *localDownloader) Download(ctx context.Context, src string, dst string) error {
	log.Debugf("local download, src: %s, dst: %s", src, dst)

	srcEndpoint, srcPath, err := util.ParseRawStorageUrl(src)
	if err != nil {
		err = fmt.Errorf("parse raw storage url failed, %w", err)
		log.Error(err)
		return err
	}

	// get real path
	realPath, err := d.cli.RealPath(srcEndpoint, srcPath)
	if err != nil {
		err = fmt.Errorf("get real path failed, %w", err)
		log.Error(err)
		return err
	}

	if realPath == dst {
		log.Debugf("download file src equal dst [%s], skip to download", dst)
		return nil
	}

	srcFile, err := os.Open(realPath)
	if err != nil {
		err = fmt.Errorf("open src file failed, %w", err)
		log.Error(err)
		return err
	}
	defer srcFile.Close()

	// create dir before file
	dstDir, _ := filepath.Split(dst)
	if err = os.MkdirAll(dstDir, 0755); err != nil {
		err = fmt.Errorf("mkdirall %s failed, %w", dstDir, err)
		log.Error(err)
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		err = fmt.Errorf("create dst file failed, %w", err)
		log.Error(err)
		return err
	}
	defer dstFile.Close()

	opts := make([]oshelp.Option, 0)
	username := config.GetConfig().BackendProvider.SchedulerCommon.SubmitSysUser
	if username != "" {
		opts = append(opts, oshelp.WithChown(username))
	}

	if err = oshelp.CopyToFile(ctx, dstFile, srcFile, opts...); err != nil {
		err = fmt.Errorf("io copy buffer failed, %w", err)
		log.Error(err)
		return err
	}

	return nil
}

type remoteDownloader struct {
	cli *client.Client
}

func (d *remoteDownloader) Download(ctx context.Context, src string, dst string) error {
	log.Debugf("remote download, src: %s, dst: %s", src, dst)

	// call download api
	srcEndpoint, srcPath, err := util.ParseRawStorageUrl(src)
	if err != nil {
		err = fmt.Errorf("parse raw storage url failed, %w", err)
		log.Error(err)
		return err
	}

	dstDir, _ := filepath.Split(dst)
	if err = os.MkdirAll(dstDir, 0755); err != nil {
		err = fmt.Errorf("mkdirall %s failed, %w", dstDir, err)
		log.Error(err)
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		err = fmt.Errorf("create dst file failed, %w", err)
		log.Error(err)
		return err
	}
	defer dstFile.Close()

	if err = d.cli.DownloadByStream(ctx, srcEndpoint, srcPath, dstFile); err != nil {
		err = fmt.Errorf("call storage download failed, %w", err)
		log.Error(err)
		return err
	}

	return nil
}
