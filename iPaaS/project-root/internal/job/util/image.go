package util

import (
	"bytes"

	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/lmittmann/ppm"
	"golang.org/x/image/bmp"
)

// ImageType 图片类型
type ImageType int

const (
	// JPG jpg/jpeg
	JPG ImageType = iota
	// PNG png
	PNG
	// GIF gif
	GIF
	// PPM ppm
	PPM
	// BMP bmp
	BMP
)

// ConvertImageTo 转换图片格式到
func ConvertImageTo(r io.Reader, imagePath string, target ImageType) ([]byte, error) {
	// 获取文件后缀
	fileExt := strings.ToLower(filepath.Ext(imagePath))
	var img image.Image
	var err error
	var fDecode func(r io.Reader) (image.Image, error)

	// 解析图片
	switch fileExt {
	case ".jpg", ".jpeg":
		fDecode = jpeg.Decode
	case ".png":
		fDecode = png.Decode
	case ".gif":
		fDecode = gif.Decode
	case ".ppm":
		fDecode = ppm.Decode
	case ".bmp":
		fDecode = bmp.Decode
	default:
		return []byte{}, errors.New("unsupported")
	}

	img, err = fDecode(r)
	if err != nil {
		return []byte{}, err
	}

	var res string
	buf := bytes.NewBufferString(res)

	var fEncode func(w io.Writer, m image.Image) error

	// 转为目标格式
	switch target {
	case JPG:
		fEncode = func(w io.Writer, m image.Image) error {
			return jpeg.Encode(w, m, nil)
		}
	case PNG:
		fEncode = png.Encode
	case GIF:
		fEncode = func(w io.Writer, m image.Image) error {
			return gif.Encode(w, m, nil)
		}
	case PPM:
		fEncode = ppm.Encode
	case BMP:
		fEncode = bmp.Encode
	default:
		return []byte{}, errors.New("unsupported")
	}

	err = fEncode(buf, img)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil

}

// IsImageFile 判断是否是图片文件
func IsImageFile(isDir bool, size int64, name string) (bool, error) {
	if isDir || size == 0 {
		return false, nil
	}

	matched, err := regexp.MatchString(`(.png|.jpeg|.ppm|.jpg|.bmp)$`, name)
	if err != nil {
		return false, errors.Wrapf(err, "regexp.MatchString error")
	}

	return matched, nil
}

// ExtractNameFromFileName 从文件名中提取名称 如：test_1.png -> test
func ExtractNameFromFileName(fileName string) string {
	matches := regexp.MustCompile(`(?m)(.*?)_.*`).FindAllStringSubmatch(fileName, -1)
	for _, names := range matches {
		if len(names) > 1 {
			return names[1]
		}
	}
	return ""
}
