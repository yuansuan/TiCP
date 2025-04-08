package xurl

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendValues(t *testing.T) {
	table := []struct {
		Left   url.Values
		Right  url.Values
		Expect url.Values
	}{
		{
			Left:   nil,
			Right:  nil,
			Expect: nil,
		},
		{
			Left:   nil,
			Right:  url.Values{"a": []string{"1"}},
			Expect: url.Values{"a": []string{"1"}},
		},
		{
			Left:   url.Values{"a": []string{"1"}},
			Right:  nil,
			Expect: url.Values{"a": []string{"1"}},
		},
		{
			Left:   url.Values{"a": []string{"1"}},
			Right:  url.Values{"a": []string{"2"}},
			Expect: url.Values{"a": []string{"2"}},
		},
		{
			Left:   url.Values{"a": []string{"1"}},
			Right:  url.Values{"b": []string{"1"}},
			Expect: url.Values{"a": []string{"1"}, "b": []string{"1"}},
		},
	}

	for idx, testcase := range table {
		assert.Equal(t, testcase.Expect, AppendValues(testcase.Left, testcase.Right), "Index: %d", idx)
	}
}

func TestMerge(t *testing.T) {
	table := []struct {
		Ref    string
		Target string
		Expect string
	}{
		{
			Ref:    "https://example.com",
			Target: "/api",
			Expect: "https://example.com/api",
		},
		{
			Ref:    "https://example.com",
			Target: "api",
			Expect: "https://example.com/api",
		},
		{
			Ref:    "https://example.com/api",
			Target: "/echo",
			Expect: "https://example.com/api/echo",
		},
		{
			Ref:    "https://example.com/api",
			Target: "https://example.com/v1/api",
			Expect: "https://example.com/v1/api",
		},
		{
			Ref:    "https://example.com/api",
			Target: "//example.com/v1/api",
			Expect: "https://example.com/v1/api",
		},
	}

	for idx, testcase := range table {
		if ur, err := url.Parse(testcase.Ref); assert.NoError(t, err) {
			if ut, err := url.Parse(testcase.Target); assert.NoError(t, err) {
				if ue, err := url.Parse(testcase.Expect); assert.NoError(t, err) {
					assert.Equal(t, ue.String(), Merge(ut, ur).String(), "Index: %d", idx)
				}
			}
		}
	}
}
