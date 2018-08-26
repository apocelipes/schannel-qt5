package config

import (
	"testing"

	"encoding/json"
	"os"
)

func TestAbsPath(t *testing.T) {
	p1 := JSONPath{"/testing/abs/path"}

	data1, err := p1.AbsPath()
	if data1 != p1.string {
		t.Error("wrong on abs path: " + data1)
	} else if err != nil {
		t.Error(err)
	}

	err = os.Setenv("HOME", "/home/testing")
	if err != nil {
		t.Errorf("set $HOME ERROR: %v\n", err)
	}

	p2 := JSONPath{"~/testing/path"}

	data2, err := p2.AbsPath()
	if data2 != "/home/testing/testing/path" {
		t.Error("wrong on home path: " + data2)
	} else if err != nil {
		t.Error(err)
	}

	p3 := JSONPath{"testing/path"}
	_, err = p3.AbsPath()
	if err == nil {
		t.Error("ErrNotAbs didn't work")
	}
}

// testing type
type p struct {
	Port int64    `json:"port"`
	Path JSONPath `json:"path"`
}

func TestUnmarshalJson(t *testing.T) {
	data := "{\"port\":12345,\"path\":\"~/.test/\"}"

	j := new(p)
	err := json.Unmarshal([]byte(data), j)
	if err != nil {
		t.Error(err)
	}

	if j.Port != 12345 || j.Path.string != "~/.test/" {
		t.Error("unmarshal error")
	}
}

func TestMarshalJson(t *testing.T) {
	j := new(p)
	j.Port = 12345
	j.Path.string = "~/.test/"
	data, err := json.Marshal(j)
	if err != nil {
		t.Error(err)
	}

	if string(data) != "{\"port\":12345,\"path\":\"~/.test/\"}" {
		t.Error("marshal error")
	}
}
