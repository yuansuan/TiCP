package util

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// GoVersionCompare returns a bool comparing two strings lexicographically.
// The result will be false if a < b, and true if a >= b
func GoVersionCompare(a, b string) (bool, error) {
	pointNum := max(strings.Count(a, "."), strings.Count(b, "."))
	// 字符串转换为可比较字符
	aNum, err := versionStrToInt(a, pointNum)
	if err != nil {
		return false, err
	}

	bNum, err := versionStrToInt(b, pointNum)
	if err != nil {
		return false, err
	}

	return aNum >= bNum, nil
}

func versionStrToInt(version string, pointNum int) (int, error) {
	currentNum := ""
	var versionNum float64 = 0

	//过滤掉非数字及点
	for _, v := range version {
		vStr := string(v)
		if strings.Contains("0123456789.", vStr) {
			currentNum += vStr
		}
	}

	intSlice := strings.Split(currentNum, ".")
	for k, v := range intSlice {
		intNum, err := strconv.Atoi(v)
		if err != nil {
			return 0, errors.New("check the go version your input")
		}
		versionNum += float64(intNum) * math.Pow(100, float64(pointNum-k))
	}

	return int(versionNum), nil
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
