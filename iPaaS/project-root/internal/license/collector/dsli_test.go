package collector

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
)

const DSLiTest string = `(\w+)\s+maxReleaseNumber:\s*(\d+).*?count:\s*(\d+).*?inuse:\s*(\d+)`

func TestCommandGeneration(t *testing.T) {
	license := &LicenseCollectInfo{
		LicensePath: "localhost 1234",
		LmstatPath:  "/path/to/lmstat",
	}

	path := strings.Split(license.LicensePath, " ")
	host, port := path[0], ""
	if len(path) > 1 {
		port = path[1]
	}
	command := fmt.Sprintf("%v -admin -run \"connect %v %v;getLicenseUsage\"", license.LmstatPath, host, port)

	expectedCommand := "/path/to/lmstat -admin -run \"connect localhost 1234;getLicenseUsage\""
	fmt.Println("Actual command:", command)
	fmt.Println("Expected command:", expectedCommand)
	assert.Equal(t, expectedCommand, command, "生成的命令不符合预期")
}

func TestCustomRegexMatch(t *testing.T) {
	re, err := regexp.Compile(DSLiTest)
	assert.NoError(t, err, "正则表达式编译失败")

	raw := `PAC maxReleaseNumber:  8 maxReleaseDate:2025/9/26 上午7:59:00 expirationDate:2025/9/26 上午7:59:00 type:Token          count:    1 inuse: 0 customerId:200000000159998 pricing structure YLC
		    PAF maxReleaseNumber:  8 maxReleaseDate:2025/9/26 上午7:59:00 expirationDate:2025/9/26 上午7:59:00 type:Token          count:    1 inuse: 0 customerId:200000000159998 pricing structure YLC
		    PAI maxReleaseNumber:  8 maxReleaseDate:2025/9/26 上午7:59:00 expirationDate:2025/9/26 上午7:59:00 type:Token          count:    1 inuse: 0 customerId:200000000159998 pricing structure YLC
		    PAJ maxReleaseNumber:  8 maxReleaseDate:2025/9/26 上午7:59:00 expirationDate:2025/9/26 上午7:59:00 type:Token          count:    6 inuse: 0 customerId:200000000159998 pricing structure YLC
		    PAN maxReleaseNumber:  8 maxReleaseDate:2025/9/26 上午7:59:00 expirationDate:2025/9/26 上午7:59:00 type:Token          count:    1 inuse: 0 customerId:200000000159998 pricing structure YLC
		    PCA maxReleaseNumber:  8 maxReleaseDate:2025/9/26 上午7:59:00 expirationDate:2025/9/26 上午7:59:00 type:Token          count:    1 inuse: 0 customerId:200000000159998 pricing structure YLC`
	matches := re.FindAllStringSubmatch(raw, -1)

	assert.NotEmpty(t, matches, "正则表达式未匹配任何数据")

	expected := [][]string{
		{"PAC", "8", "1", "0"},
		{"PAF", "8", "1", "0"},
		{"PAI", "8", "1", "0"},
		{"PAJ", "8", "6", "0"},
		{"PAN", "8", "1", "0"},
		{"PCA", "8", "1", "0"},
	}
	for i, match := range matches {
		assert.Equal(t, expected[i], match[1:], "第 %d 条记录的匹配结果不一致", i+1)
	}
	for _, match := range matches {
		fmt.Println(match)
	}
}

func TestDSLiCollector_Parse(t *testing.T) {
	license := &LicenseCollectInfo{
		LmstatPath: "/path/to/lmstat",
	}
	collector := NewDSLiCollector(license)

	raw := `PAC maxReleaseNumber:  8 maxReleaseDate:2025/9/26 上午7:59:00 expirationDate:2025/9/26 上午7:59:00 type:Token count:  1 inuse: 0 customerId:200000000159998 pricing structure YLC
		    PAF maxReleaseNumber:  8 maxReleaseDate:2025/9/26 上午7:59:00 expirationDate:2025/9/26 上午7:59:00 type:Token count:  1 inuse: 0 customerId:200000000159998 pricing structure YLC
		    PAI maxReleaseNumber:  8 maxReleaseDate:2025/9/26 上午7:59:00 expirationDate:2025/9/26 上午7:59:00 type:Token count:  1 inuse: 0 customerId:200000000159998 pricing structure YLC
		    PAJ maxReleaseNumber:  8 maxReleaseDate:2025/9/26 上午7:59:00 expirationDate:2025/9/26 上午7:59:00 type:Token count:  6 inuse: 0 customerId:200000000159998 pricing structure YLC
		    PAN maxReleaseNumber:  8 maxReleaseDate:2025/9/26 上午7:59:00 expirationDate:2025/9/26 上午7:59:00 type:Token count:  1 inuse: 0 customerId:200000000159998 pricing structure YLC
		    PCA maxReleaseNumber:  8 maxReleaseDate:2025/9/26 上午7:59:00 expirationDate:2025/9/26 上午7:59:00 type:Token count:  1 inuse: 0 customerId:200000000159998 pricing structure YLC`

	re, err := regexp.Compile(DSLiTest)
	assert.NoError(t, err, "正则表达式编译失败")

	matches := re.FindAllStringSubmatch(raw, -1)
	assert.NotEmpty(t, matches, "正则表达式未匹配任何数据")

	for _, match := range matches {
		parsedData := fmt.Sprintf("%s maxReleaseNumber: %s count: %s inuse: %s", match[1], match[2], match[3], match[4])
		collector.parse(parsedData)
	}

	expectedComponents := map[string]Component{
		"PAC": {Total: 1, Used: 0},
		"PAF": {Total: 1, Used: 0},
		"PAI": {Total: 1, Used: 0},
		"PAJ": {Total: 6, Used: 0},
		"PAN": {Total: 1, Used: 0},
		"PCA": {Total: 1, Used: 0},
	}

	for key, expected := range expectedComponents {
		actual, exists := collector.GetComponents()[key]
		fmt.Println("actual:  ", key, actual)
		assert.True(t, exists, "缺少组件 %s", key)
		assert.Equal(t, expected, actual, "组件 %s 的数据不符合预期", key)
	}
}
