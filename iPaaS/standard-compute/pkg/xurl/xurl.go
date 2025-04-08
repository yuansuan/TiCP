package xurl

import (
	"net/url"
	"path/filepath"
	"strings"
)

// Merge 将 src 中的基础数据合并到 dst 中
func Merge(dst *url.URL, src *url.URL) *url.URL {
	if len(dst.Scheme) == 0 || len(dst.Host) == 0 {
		if len(dst.Host) == 0 {
			dst.Host = src.Host
			dst.Path = filepath.Join(src.Path, dst.Path)
		}
		if len(dst.Scheme) == 0 {
			dst.Scheme = src.Scheme
		}
	}
	return dst
}

// Join 将多个字符串拼接成一个URL
func Join(els ...string) string {
	var sb strings.Builder
	for i, el := range els {
		sb.WriteString(strings.Trim(el, "/"))
		if i != len(els)-1 {
			sb.WriteString("/")
		}
	}
	return "/" + sb.String()
}

// Trim 删除两端多余的 "/"
func Trim(s string) string {
	return strings.Trim(s, "/")
}
