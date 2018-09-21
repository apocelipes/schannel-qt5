package pyclient

import (
	"os"
	"schannel-qt5/config"
	"testing"
	"time"
)

var (
	dummyClientConfig = &ClientConfig{}
	dummyUserConfig = &config.UserConfig{
		Proxy:               config.JSONProxy{},
		LogFile:             config.JSONPath{},
		SSRNodeConfigPath:   config.JSONPath{Data: os.Getenv("SSRNODECONFIG")},
		SSRClientConfigPath: config.JSONPath{},
		SSRBin:              config.JSONPath{Data: os.Getenv("SSRBIN")},
		SSRClientConfig:     dummyClientConfig,
	}
	defaultBin = "~/shadowsocksr/shadowsocks/local.py"
	defaultNode = "~/ssr.json"
)

func checkSSRNodeConfig(t *testing.T) {
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
	checkSSRNodeConfig(t)
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

func TestPySSRClientConnectionCheck(t *testing.T) {
	checkSSRNodeConfig(t)
	client := newPySSRClient(dummyUserConfig)
	if client == nil {
		t.Error("newPySSRClient failed: ", dummyUserConfig)
	}

	if err := client.Start(); err != nil {
		t.Errorf("start client failed: %v\n", err)
	}

	if err := client.ConnectionCheck(10 * time.Second); err != nil {
		t.Errorf("connect check failed: %v\n", err)
	}

	if err := client.Stop(); err != nil {
		t.Errorf("stop client failed: %v\n", err)
	}
}
