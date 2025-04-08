package image

import (
	"errors"
	"fmt"
	"regexp"
)

// Locator 用户定位一个具体的镜像
type Locator interface {
	fmt.Stringer

	// Name 返回需要定位的镜像的名称
	Name() string

	// Tag 返回需要定位的镜像的具体标签
	Tag() string

	// Hash 返回需要定位的镜像的具体哈希值
	Hash() string

	// ShortString 返回字符串形式的短标识
	ShortString() string
}

// LocateOption 定位镜像的额外参数
type LocateOption func(l *locator)

// WithLocateTag 指定特定标签的镜像
func WithLocateTag(tag string) LocateOption {
	return func(l *locator) {
		l.tag = tag
	}
}

// WithLocateHash 指定特定名字的镜像
func WithLocateHash(hash string) LocateOption {
	return func(l *locator) {
		l.hash = hash
	}
}

// Locate 定位一个镜像
func Locate(name string, options ...LocateOption) Locator {
	l := &locator{name: name}
	for _, option := range options {
		option(l)
	}

	return l
}

var (
	// ErrInvalidLocatorString 镜像定位器的格式无效
	ErrInvalidLocatorString = errors.New("locator: invalid form for locator string")

	// _StringLocatorRegexp 字符串形式的定位器正则
	_StringLocatorRegexp = regexp.MustCompile("^([\\w-_]+)(:([\\w-_.]+))?(:([\\w-_.]+)@(\\w+))?$")
)

// FromString 从字符串构建定位器
func FromString(s string) (Locator, error) {
	if matches := _StringLocatorRegexp.FindStringSubmatch(s); matches != nil {
		// aa:bb = [0:"aa:bb", 1:"aa", 2:":bb", 3:"bb", 4:"", 5:"", 6:"" ]
		// aa:bb@cc = [0:"aa:bb@cc", 1:"aa", 2:"", 3:"", 4:":bb@cc", 5:"bb", 6:"cc" ]
		if len(matches[4]) == 0 {
			return Locate(matches[1], WithLocateTag(matches[3])), nil
		}
		return Locate(matches[1], WithLocateTag(matches[5]), WithLocateHash(matches[6])), nil
	}
	return nil, ErrInvalidLocatorString
}

// locator 一个简单的定位器
type locator struct {
	name string
	tag  string
	hash string
}

// String 返回字符串形式标识
func (l *locator) String() string {
	return fmt.Sprintf("%s:%s@%s", l.name, l.tag, l.hash)
}

// ShortString 返回字符串形式的短标识
func (l *locator) ShortString() string {
	return fmt.Sprintf("%s:%s", l.name, l.tag)
}

// Name 返回需要定位的镜像的名称
func (l *locator) Name() string {
	return l.name
}

// Tag 返回需要定位的镜像的具体标签
func (l *locator) Tag() string {
	return l.tag
}

// Hash 返回需要定位的镜像的具体哈希值
func (l *locator) Hash() string {
	return l.hash
}
