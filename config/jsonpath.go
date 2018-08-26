package config

import (
	"os"
	"path"
	"strings"
)

// JSONPath unmarshal/marshal时验证路径是否为绝对路径
type JSONPath struct {
	string
}

func (jp *JSONPath) UnmarshalJSON(b []byte) error {
	data := string(b)
	data = strings.TrimSuffix(data, "\"")
	data = strings.TrimPrefix(data, "\"")
	if !strings.HasPrefix(data, "~") && !path.IsAbs(data) {
		return ErrNotAbs
	}

	jp.string = data
	return nil
}

func (jp JSONPath) String() string {
	return jp.string
}

func (jp *JSONPath) MarshalJSON() ([]byte, error) {
	if !strings.HasPrefix(jp.string, "~") && !path.IsAbs(jp.string) {
		return nil, ErrNotAbs
	}

	return []byte("\"" + jp.string + "\""), nil
}

// AbsPath 返回对应的绝对路径
func (jp *JSONPath) AbsPath() (string, error) {
	if path.IsAbs(jp.string) {
		return jp.string, nil
	}

	if !strings.HasPrefix(jp.string, "~") {
		return "", ErrNotAbs
	}

	home, exist := os.LookupEnv("HOME")
	if !exist {
		return "", ErrHOME
	}

	return strings.Replace(jp.string, "~", home, 1), nil
}
