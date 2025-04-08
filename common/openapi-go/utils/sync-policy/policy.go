package syncpolicy

import (
	"path/filepath"
	"regexp"
	"strings"
)

const (
	minCompressSize = 8 * 1024 // 8k, do not compress if data less than it
)

const (
	defaultCompress = "gzip"
)

type RuleMethod string

const (
	// 文件名扩展 大小写不敏感
	RuleMethodExt = "ext"
	// 是否有前缀 大小写不敏感
	RuleMethodPrefix = "prefix"
	// 正则匹配
	RuleMethodRegexp = "regexp"
	// 直接对比
	RuleMethodDirect = "direct"
)

type Rules struct {
	Method   RuleMethod `yaml:"method"`
	Express  string     `yaml:"express"`
	Compress string     `yaml:"compress"`
}

type RuleFn func(filename string) (bool, string)

type Policy struct {
	rules       []RuleFn
	originRules []*Rules
}

func NewPolicy(rules []*Rules) (*Policy, error) {
	r, err := compileRules(rules)
	if err != nil {
		return nil, err
	}

	return &Policy{
		rules:       r,
		originRules: rules,
	}, nil
}

func (p *Policy) Do(filename string, size int64) string {
	if size < minCompressSize {
		return ""
	}
	return p.do(filename)
}

func (p *Policy) do(filename string) string {
	for _, f := range p.rules {
		b, m := f(filename)
		if b {
			return m
		}
	}
	return defaultCompress
}

func compileRules(r []*Rules) ([]RuleFn, error) {
	var ret []RuleFn
	for _, v := range r {
		f, err := compileRule(v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, f)
	}
	return ret, nil
}

func compileRule(r *Rules) (RuleFn, error) {
	switch r.Method {
	case RuleMethodExt:
		express := strings.ToLower(r.Express)
		return func(filename string) (bool, string) {
			ext := strings.ToLower(filepath.Ext(filename))
			return ext == express, string(r.Method)
		}, nil
	case RuleMethodPrefix:
		express := strings.ToLower(r.Express)
		return func(filename string) (bool, string) {
			filename = strings.ToLower(filename)
			return strings.HasPrefix(filename, express), string(r.Method)
		}, nil
	case RuleMethodRegexp:
		reg, err := regexp.Compile(r.Express)
		if err != nil {
			return nil, err
		}
		return func(filename string) (bool, string) {
			return reg.MatchString(filename), string(r.Method)
		}, nil
	case RuleMethodDirect:
		return func(filename string) (bool, string) {
			return filename == r.Express, string(r.Method)
		}, nil
	}

	return nil, nil
}
