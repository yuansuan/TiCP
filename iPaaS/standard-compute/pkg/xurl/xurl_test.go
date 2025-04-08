package xurl

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	table := []struct {
		Src string
		Dst string
		Out string
	}{
		{Src: "https://example.com", Dst: "/foo/abr?a=1", Out: "https://example.com/foo/abr?a=1"},
		{Src: "https://example.com", Dst: "foo/abr?a=1", Out: "https://example.com/foo/abr?a=1"},
		{Src: "https://example.com/", Dst: "/foo/abr?a=1", Out: "https://example.com/foo/abr?a=1"},
		{Src: "https://example.com/", Dst: "foo/abr?a=1", Out: "https://example.com/foo/abr?a=1"},
		{Src: "https://example.com/prefix", Dst: "/foo/abr?a=1", Out: "https://example.com/prefix/foo/abr?a=1"},
		{Src: "https://example.com/prefix/", Dst: "/foo/abr?a=1", Out: "https://example.com/prefix/foo/abr?a=1"},
		{Src: "https://example.com", Dst: "https://foobar.com/hello", Out: "https://foobar.com/hello"},
		{Src: "https://example.com", Dst: "//foobar.com/hello", Out: "https://foobar.com/hello"},
	}

	for _, testcase := range table {
		if src, err := url.Parse(testcase.Src); assert.NoError(t, err) {
			if dst, err := url.Parse(testcase.Dst); assert.NoError(t, err) {
				assert.Equal(t, testcase.Out, Merge(dst, src).String())
			}
		}
	}
}
