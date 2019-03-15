package config

import (
	"testing"

	"encoding/json"
	"os"
)

func TestAbsPath(t *testing.T) {
	// 测试绝对路径
	abs := JSONPath{"/testing/abs/path"}
	data1, err := abs.AbsPath()
	if data1 != abs.Data {
		t.Error("wrong on abs path: " + data1)
	} else if err != nil {
		t.Error(err)
	}

	// 测试~开头的HOME下路径
	os.Clearenv()
	err = os.Setenv("HOME", "/home/testing")
	if err != nil {
		t.Errorf("set $HOME ERROR: %v\n", err)
	}

	underHome := JSONPath{"~/testing/path"}
	data2, err := underHome.AbsPath()
	if data2 != "/home/testing/testing/path" {
		t.Error("wrong on home path: " + data2)
	} else if err != nil {
		t.Error(err)
	}

	// 测试非绝对路径
	notAbs := JSONPath{"testing/path"}
	_, err = notAbs.AbsPath()
	if err == nil {
		t.Error("ErrNotAbs didn't work")
	}

	empty := JSONPath{""}
	_, err = empty.AbsPath()
	if err == nil {
		t.Error("empty string passed")
	}
}

// testing type
type p struct {
	Port int64    `json:"port"`
	Path JSONPath `json:"path"`
}

func TestUnmarshalJson(t *testing.T) {
	// 测试json数据
	data := "{\"port\":12345,\"path\":\"~/.test/\"}"

	j := new(p)
	err := json.Unmarshal([]byte(data), j)
	if err != nil {
		t.Error(err)
	}

	if j.Port != 12345 || j.Path.Data != "~/.test/" {
		t.Error("unmarshal error")
	}
}

func TestMarshalJson(t *testing.T) {
	j := &p{
		Port: 12345,
		Path: JSONPath{Data: "~/.test/"},
	}
	data, err := json.Marshal(j)
	if err != nil {
		t.Error(err)
	}

	if string(data) != "{\"port\":12345,\"path\":\"~/.test/\"}" {
		t.Error("marshal error")
	}
}
