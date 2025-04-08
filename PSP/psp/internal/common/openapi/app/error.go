package openapiapp

import (
	"errors"
)

var (
	// ErrAppIDEmpty 应用ID不能为空
	ErrAppIDEmpty = errors.New("the app id is empty")
	// ErrAppNameEmpty 应用名称不能为空
	ErrAppNameEmpty = errors.New("the app name is empty")
	// ErrAppTypeEmpty 应用类型不能为空
	ErrAppTypeEmpty = errors.New("the app type is empty")
	// ErrAppVersionEmpty 应用版本不能为空
	ErrAppVersionEmpty = errors.New("the app version is empty")
)
