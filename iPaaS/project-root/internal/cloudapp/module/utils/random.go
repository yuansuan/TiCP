package utils

import (
	"math/rand"
	"time"
)

const (
	upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerLetters = "abcdefghijklmnopqrstuvwxyz"
	numbers      = "0123456789"
	// 不要使用'@'和','，agent中有依赖两者作为用户密码分隔符
	// 已知：$经过user-data转换之后会变成^, 故不要用$作为密码特殊字符
	symbolLetters = "~!_."
)

const (
	upper = iota
	lower
	number
	symbol
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomPassword(length int) string {
	nowLen := 0
	res := ""
	for nowLen < length {
		for i := 0; i < 4 && nowLen < length; i++ {
			switch i {
			case upper:
				res += randomChar(upperLetters)
			case lower:
				res += randomChar(lowerLetters)
			case number:
				res += randomChar(numbers)
			case symbol:
				res += randomChar(symbolLetters)
			default:
				res += randomChar(upperLetters)
			}
			nowLen++
		}
	}

	return res
}

func randomChar(str string) string {
	i := rand.Intn(len(str))
	return str[i : i+1]
}
