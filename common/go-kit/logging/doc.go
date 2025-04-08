/*
Package logging
1. 支持打印至console/file；
2. 打印至file时支持自动归档。

导入该包时init会生成默认logger，并支持运行时替换默认logger

logger, err := logging.NewLogger(...opts)

_, err := logging.SetDefault(logging.WithLogger(logger))

API均兼容原有
*/
package logging
