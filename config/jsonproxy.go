package config

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrNotURL 不是合法的URL
	ErrNotURL = errors.New("not an valid URL")

	matcher = regexp.MustCompile(`(http(s)?|socks5)://([\w\-]+\.)+[\w\-]+(:\d+)?(/[\w\- ./?%&=]*)?`)
)

// JSONProxy 验证给定字符串是否是合法的URL
type JSONProxy struct {
	string
}

// IsURL 如果p的值是合法的URL，则返回true
func (p *JSONProxy) IsURL() bool {
	return matcher.MatchString(p.string)
}

func (p JSONProxy) String() string {
	return p.string
}

func (p *JSONProxy) UnmarshalJSON(b []byte) error {
	data := string(b)
	// 对于字符串的json值，需要手动去除双引号
	data = strings.TrimSuffix(data, "\"")
	data = strings.TrimPrefix(data, "\"")
	p.string = data

	if !p.IsURL() {
		return ErrNotURL
	}

	return nil
}

func (p *JSONProxy) MarshalJSON() ([]byte, error) {
	if !p.IsURL() {
		return nil, ErrNotURL
	}

	return []byte("\"" + p.string + "\""), nil
}
