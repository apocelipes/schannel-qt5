package pyclient

import (
	"testing"

	"os"
	"time"

	"schannel-qt5/config"
)

var (
	// 测试用空config对象
	dummyClientConfig = &ClientConfig{}
	// 测试用用户配置
	dummyUserConfig = &config.UserConfig{
		Proxy:               config.JSONProxy{},
		LogFile:             config.JSONPath{},
		SSRNodeConfigPath:   config.JSONPath{Data: os.Getenv("SSRNODECONFIG")},
		SSRClientConfigPath: config.JSONPath{},
		SSRBin:              config.JSONPath{Data: os.Getenv("SSRBIN")},
		SSRClientConfig:     dummyClientConfig,
	}
	// 默认的文件路径
	defaultBin  = "~/shadowsocksr/shadowsocks/local.py"
	defaultNode = "~/ssr.json"
)

// checkUserConfig 检查用户配置
// 如果未设置$SSRBIN和$SSRNODECONFIG，则使用默认值替代
func checkUserConfig(t *testing.T) {
	if dummyUserConfig.SSRBin.String() == "" {
		t.Log("not set $SSRBIN, use default\n")
		dummyUserConfig.SSRBin.Data = defaultBin
	}

	if dummyUserConfig.SSRNodeConfigPath.String() == "" {
		t.Log("not set $SSRNODECONFIG, use default\n")
		dummyUserConfig.SSRNodeConfigPath.Data = defaultNode
	}
}

func TestPySSRClient(t *testing.T) {
	// 先对dummyUserConfig初始化
	checkUserConfig(t)
	client := newPySSRClient(dummyUserConfig)
	if client == nil {
		t.Error("newPySSRClient failed: ", dummyUserConfig)
	}

	if err := client.Start(); err != nil {
		t.Errorf("start client failed: %v\n", err)
	}

	if err := client.Restart(); err != nil {
		t.Errorf("restart client failed: %v\n", err)
	}

	if err := client.Stop(); err != nil {
		t.Errorf("stop client failed: %v\n", err)
	}
}

func TestPySSRClientIsRunning(t *testing.T) {
	checkUserConfig(t)
	client := newPySSRClient(dummyUserConfig)
	if client == nil {
		t.Error("newPySSRClient failed: ", dummyUserConfig)
	}

	if err := client.Start(); err != nil {
		t.Errorf("start client failed: %v\n", err)
	}

	// 测试是否正在运行
	if err := client.IsRunning(); err != nil {
		t.Fatalf("test is running error: %v\n", err)
	}

	if err := client.Stop(); err != nil {
		t.Errorf("stop client failed: %v\n", err)
	}

	// 测试是否已经关闭
	if err := client.IsRunning(); err == nil {
		t.Fatalf("test not run error: %v\n", err)
	}
}

func TestPySSRClientConnectionCheck(t *testing.T) {
	checkUserConfig(t)
	client := newPySSRClient(dummyUserConfig)
	if client == nil {
		t.Error("newPySSRClient failed: ", dummyUserConfig)
	}

	// 打开代理
	if err := client.Start(); err != nil {
		t.Errorf("start client failed: %v\n", err)
	}

	if err := client.ConnectionCheck(10 * time.Second); err != nil {
		t.Errorf("connect check failed: %v\n", err)
	}

	// 测试结束关闭代理
	if err := client.Stop(); err != nil {
		t.Errorf("stop client failed: %v\n", err)
	}
}
