package config

import (
	"testing"

	"encoding/json"
)

// 测试路径json转换
type e struct {
	Path *JSONEmptyPath `json:"path"`
}

func TestJSONEmptyPathMarshalJSON(t *testing.T) {
	testData := []*struct {
		data    *e
		success bool
		res     string
	}{
		{
			data: &e{
				Path: &JSONEmptyPath{JSONPath{Data: ""}},
			},
			success: true,
			res:     `{"path":""}`,
		},
		{
			data: &e{
				Path: &JSONEmptyPath{JSONPath{Data: "/tmp/a"}},
			},
			success: true,
			res:     `{"path":"/tmp/a"}`,
		},
		{
			data: &e{
				Path: &JSONEmptyPath{JSONPath{Data: "tmp/a"}},
			},
			success: false,
			res:     ``,
		},
	}

	for _, v := range testData {
		res, err := json.Marshal(v.data)
		if (err == nil) != v.success {
			format := "marshal error:\n\tdata: %s\n\terror: %v\n\twant :%v\n"
			t.Errorf(format, v.data.Path, err, v.success)
		} else if v.success && string(res) != v.res {
			format := "marshal wrong:\n\twant: %s\n\thave: %s\n"
			t.Errorf(format, v.res, string(res))
		}
	}
}

func TestJSONEmptyPathUnmarshalJSON(t *testing.T) {
	testData := []*struct {
		data    string
		success bool
		// 存放解析结果
		res *e
		// 正确的路径
		resPath string
	}{
		{
			data:    `{"path":""}`,
			success: true,
			res: &e{
				Path: &JSONEmptyPath{JSONPath{}},
			},
			resPath: "",
		},
		{
			data:    `{"path":"/tmp/a"}`,
			success: true,
			res: &e{
				Path: &JSONEmptyPath{JSONPath{}},
			},
			resPath: "/tmp/a",
		},
		{
			data:    `{"path":"tmp/a"}`,
			success: false,
			res: &e{
				Path: &JSONEmptyPath{JSONPath{}},
			},
			resPath: "",
		},
	}

	for _, v := range testData {
		err := json.Unmarshal([]byte(v.data), v.res)
		if (err == nil) != v.success {
			format := "unmarshal error:\n\tdata: %s\n\terror: %v\n\twant :%v\n"
			t.Errorf(format, v.data, err, v.success)
		} else if v.success && v.res.Path.String() != v.resPath {
			format := "unmarshal wrong:\n\twant: %s\n\thave: %s\n"
			t.Errorf(format, v.resPath, v.res.Path)
		}
	}
}

func TestJSONEmptyPathIsEmpty(t *testing.T) {
	testData := []*struct{
		data string
		empty bool
	}{
		{
			data:  "/tmp/a",
			empty: false,
		},
		{
			data:  "",
			empty: true,
		},
	}

	for _,v := range testData {
		ep := &JSONEmptyPath{JSONPath{Data: v.data}}
		if ep.IsEmpty() != v.empty {
			format := "isEmpty wrong:\n\tpath: %s\n\twant: %v\n\thave: %v\n"
			t.Errorf(format, v.data, v.empty, ep.IsEmpty())
		}
	}
}
