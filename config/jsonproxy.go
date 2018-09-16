package config

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrNotURL 不是合法的URL
	ErrNotURL = errors.New("not an valid URL")

	matcher = regexp.MustCompile(`^(http(s)?|socks5)://([\w\-]+\.)+[\w\-]+(:\d+)?(/[\w\- ./?%&=]*)?$`)
)

// JSONProxy 验证给定字符串是否是合法的URL
type JSONProxy struct {
	Data string
}

// IsURL 如果p的值是合法的URL，则返回true
func (p *JSONProxy) IsURL() bool {
	return matcher.MatchString(p.Data)
}

func (p JSONProxy) String() string {
	return p.Data
}

func (p *JSONProxy) UnmarshalJSON(b []byte) error {
	data := string(b)
	// 对于字符串的json值，需要手动去除双引号
	data = strings.TrimSuffix(data, "\"")
	data = strings.TrimPrefix(data, "\"")
	p.Data = data

	// 允许""表示不使用proxy
	if !p.IsURL() && p.Data != "" {
		return ErrNotURL
	}

	return nil
}

func (p *JSONProxy) MarshalJSON() ([]byte, error) {
	if !p.IsURL() && p.Data != "" {
		return nil, ErrNotURL
	}

	return []byte("\"" + p.Data + "\""), nil
}
