package config

import (
	"testing"

	"encoding/json"
	"os"
)

func TestConfigPath(t *testing.T) {
	testData := []*struct {
		// 设置环境变量HOME的值
		home string
		res  string
	}{
		{
			home: "/home/test",
			res:  "/home/test/" + configPath,
		},
		{
			home: "/home/test/",
			res:  "/home/test/" + configPath,
		},
	}

	for _, v := range testData {
		err := os.Setenv("HOME", v.home)
		if err != nil {
			t.Fatalf("无法设置$HOME: %v\n", err)
		}
		res, err := ConfigPath()
		if err != nil {
			t.Fatalf("获取Config Path错误：%v\n", err)
		}
		if v.res != res {
			format := "不正确的Config Path:\n\twant: %s\n\thave: %v\n"
			t.Errorf(format, v.res, res)
		}
	}
}

func TestMarshalUserConf(t *testing.T) {
	u := new(UserConfig)
	u.SSRNodeConfigPath.Data = "/tmp/testing/t.json"
	u.SSRClientConfigPath.Data = "/tmp/testing/client.json"
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
