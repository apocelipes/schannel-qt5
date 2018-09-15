package config

import (
	"os"
	"path"
	"strings"
)

// JSONPath unmarshal/marshal时验证路径是否为绝对路径
type JSONPath struct {
	Data string
}

func (jp *JSONPath) UnmarshalJSON(b []byte) error {
	data := string(b)
	// 对于字符串的json值，需要手动去除双引号
	data = strings.TrimSuffix(data, "\"")
	data = strings.TrimPrefix(data, "\"")
	if !strings.HasPrefix(data, "~") && !path.IsAbs(data) {
		return ErrNotAbs
	}

	jp.Data = data
	return nil
}

func (jp JSONPath) String() string {
	return jp.Data
}

func (jp *JSONPath) MarshalJSON() ([]byte, error) {
	if !strings.HasPrefix(jp.Data, "~") && !path.IsAbs(jp.Data) {
		return nil, ErrNotAbs
	}

	return []byte("\"" + jp.Data + "\""), nil
}

// AbsPath 返回对应的绝对路径
func (jp *JSONPath) AbsPath() (string, error) {
	if path.IsAbs(jp.Data) {
		return jp.Data, nil
	}

	if !strings.HasPrefix(jp.Data, "~") {
		return "", ErrNotAbs
	}

	home, exist := os.LookupEnv("HOME")
	if !exist {
		return "", ErrHOME
	}

	return strings.Replace(jp.Data, "~", home, 1), nil
}
