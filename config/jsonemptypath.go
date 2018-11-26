package config

import "strings"

// 允许设置路径为空的JSONPath
type JSONEmptyPath struct {
	JSONPath
}

func (ep *JSONEmptyPath) UnmarshalJSON(b []byte) error {
	if string(b) == `""` {
		ep.Data = string(b)
		ep.Data = strings.TrimSuffix(ep.Data, `"`)
		ep.Data = strings.TrimPrefix(ep.Data, `"`)
		return nil
	}

	return ep.JSONPath.UnmarshalJSON(b)
}

func (ep *JSONEmptyPath) MarshalJSON() ([]byte, error) {
	if ep.Data == "" {
		return []byte(`""`), nil
	}

	return ep.JSONPath.MarshalJSON()
}

// 路径数据是否为空
func (ep *JSONEmptyPath) IsEmpty() bool {
	return ep.Data == ""
}
