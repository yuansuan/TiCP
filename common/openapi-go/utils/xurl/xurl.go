package xurl

import (
	"net/url"
	"path"
)

// AppendValues appends the given url.Values together, either value may be nil
func AppendValues(left, right url.Values) url.Values {
	if left == nil {
		return right
	}

	if right == nil {
		return left
	}

	for k, v := range right {
		left[k] = v
	}

	return left
}

// Merge if the scheme or host value in dst is empty, then merge the
// scheme and host from src into dst, and return dst finally
func Merge(dst *url.URL, src *url.URL) *url.URL {
	if len(dst.Scheme) == 0 || len(dst.Host) == 0 {
		if len(dst.Host) == 0 {
			dst.Host = src.Host
			dst.Path = path.Join(src.Path, dst.Path)
		}
		if len(dst.Scheme) == 0 {
			dst.Scheme = src.Scheme
		}
	}
	return dst
}
