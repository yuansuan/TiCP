// Copyright (C) 2019 LambdaCal Inc.

package util

import (
	"bytes"
	"io/ioutil"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Utf8ToGbk encode utf8 strings to gbk
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// ReplaceWindowsLineBreak replace the Windows line break character '\r'
func ReplaceWindowsLineBreak(input string) string {
	return strings.ReplaceAll(input, "\r", "")
}

// ReplaceLinuxLineBreak replace the Windows line break character '\n'
func ReplaceLinuxLineBreak(input string) string {
	return strings.ReplaceAll(input, "\n", "")
}

// SpaceTrimmedSplit splits a string with a separator
// It trims white space from all elements in the result array
func SpaceTrimmedSplit(s string, sep string) []string {
	subs := strings.Split(s, sep)
	var lines []string
	for i := 0; i < len(subs); i++ {
		line := strings.TrimSpace(subs[i])
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}
