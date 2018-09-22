package config

import (
	"testing"

	"encoding/json"
	"os"
)

func TestConfigPath(t *testing.T) {
	// 设置HOME用于拼接测试
	err := os.Setenv("HOME", "/home/testing")
	if err != nil {
		t.Errorf("set $HOME ERROR: %v\n", err)
	}

	if path, err := ConfigPath(); err != nil {
		t.Errorf("ConfigPath return an ERROR: %v\n", err)
	} else if path != "/home/testing/" + configPath {
		t.Errorf("wrong path: %v", path)
	}
}

func TestMarshalUserConf(t *testing.T) {
	u := new(UserConfig)
	u.SSRNodeConfigPath.Data = "/tmp/testing/t.json"
	u.SSRBin.Data = "/tmp/testing/a.out"
	u.LogFile.Data = "/tmp/a.log"

	data, err := json.MarshalIndent(u, "", "\t")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}

func TestUnmarshalUserConf(t *testing.T) {
	u := new(UserConfig)
	// 需要解析成config的原始数据
	data := `{"ssr_node_config_path":"/tmp/t.json","ssr_bin":"/tmp/a.out","log_file":"/tmp/a.log","ssr_client_config_path":"/tmp/c.json"}`
	if err := json.Unmarshal([]byte(data), u); err != nil {
		t.Error(err)
	}
	t.Log(*u)
}
