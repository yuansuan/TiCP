package util

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ppmData = "P6\n#comment\n2 2\n255\n\x00\x00\x00\xff\xff\xff\x00\x00\x00\xff\xff\xff"
)

func TestConvertImageTo(t *testing.T) {

	t.Run("ppm to png", func(t *testing.T) {
		img := bytes.NewReader([]byte(ppmData))

		res, err := ConvertImageTo(img, "test.ppm", PNG)
		if !assert.NoError(t, err) {
			return
		}

		assert.NotEmpty(t, res)

		// // 直接写入文件
		// err = os.WriteFile("./test.png", res, 0644)
		// if !assert.NoError(t, err) {
		// 	return
		// }

		// // 转换成base64
		// data := base64.StdEncoding.EncodeToString(res)
		// data = "data:image/png;base64," + data

		// err = os.WriteFile("./test.txt", []byte(data), 0644)
		// if !assert.NoError(t, err) {
		// 	return
		// }
	})

}

func TestIsImageFile(t *testing.T) {
	testCases := []struct {
		fileInfo struct {
			Name  string
			IsDir bool
			Size  int64
		}
		expectedIsImage bool
	}{
		{
			fileInfo: struct {
				Name  string
				IsDir bool
				Size  int64
			}{
				Name:  "image.png",
				IsDir: false,
				Size:  1024,
			},
			expectedIsImage: true,
		},
		{
			fileInfo: struct {
				Name  string
				IsDir bool
				Size  int64
			}{
				Name:  "document.txt",
				IsDir: false,
				Size:  512,
			},
			expectedIsImage: false,
		},
		{
			fileInfo: struct {
				Name  string
				IsDir bool
				Size  int64
			}{
				Name:  "dir",
				IsDir: true,
				Size:  10,
			},
			expectedIsImage: false,
		},
		{
			fileInfo: struct {
				Name  string
				IsDir bool
				Size  int64
			}{
				Name:  "empty",
				IsDir: false,
				Size:  0,
			},
			expectedIsImage: false,
		},
	}

	for _, tc := range testCases {
		isImage, err := IsImageFile(tc.fileInfo.IsDir, tc.fileInfo.Size, tc.fileInfo.Name)
		assert.NoError(t, err)
		assert.Equal(t, tc.expectedIsImage, isImage)
	}
}

func TestExtractNameFromFileName(t *testing.T) {
	testCases := []struct {
		fileName     string
		expectedName string
	}{
		{
			fileName:     "example_name.png",
			expectedName: "example",
		},
		{
			fileName:     "filewithoutunderscore.png",
			expectedName: "",
		},
		{
			fileName:     "this_is_a_test_name.png",
			expectedName: "this",
		},
	}

	for _, tc := range testCases {
		result := ExtractNameFromFileName(tc.fileName)
		assert.Equal(t, tc.expectedName, result)
	}
}
