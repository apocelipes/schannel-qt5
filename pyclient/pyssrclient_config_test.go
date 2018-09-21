package pyclient

import (
	"strings"
	"testing"
)

func TestClientConfigDefault(t *testing.T) {
	conf := &ClientConfig{}
	if conf.LocalPort() != defaultPort {
		t.Errorf("wrong default port\n")
	}

	if conf.LocalAddr() != defaultAddr {
		t.Errorf("wrong default addr\n")
	}

	if conf.PidFilePath() != defaultPidFile {
		t.Errorf("wrong default pid-file\n")
	}

	if conf.FastOpen() {
		t.Errorf("wrong default fast-open")
	}
}

func TestClientConfigSetLocalPort(t *testing.T) {
	conf := &ClientConfig{}
	trueData := []string{"0", "200", "2000", "1080", "8888"}
	wrongData := []string{"test", "端口", "", "99999"}

	for _, v := range trueData {
		if err := conf.SetLocalPort(v); err != nil {
			t.Errorf("set port failed: %v\n", v)
		} else if conf.LocalPort() != v {
			t.Errorf("set port wrong: have %v; want %v\n", conf.LocalPort(), v)
		}
	}

	for _, v := range wrongData {
		if err := conf.SetLocalPort(v); err == nil {
			t.Errorf("set wrong port but didn't fail: %v\n", v)
		}
	}
}

func TestClientConfigSetFastOpen(t *testing.T) {
	conf := &ClientConfig{}

	conf.SetFastOpen(true)
	if !conf.FastOpen() {
		t.Errorf("set fast-open failed")
	}

	conf.SetFastOpen(false)
	if conf.FastOpen() {
		t.Errorf("set fast-open failed")
	}
}

func TestClientConfigSetLocalAddr(t *testing.T) {
	conf := &ClientConfig{}
	trueData := []string{"0.0.0.0", "127.0.0.1", "10.1.1.2", "210.199.200.201"}
	wrongData := []string{"", "12345", "12.22.33.", "255.256.0.1", "test"}

	for _, v := range trueData {
		if err := conf.SetLocalAddr(v); err != nil {
			t.Errorf("set addr failed: %v\n", v)
		} else if conf.LocalAddr() != v {
			t.Errorf("set addr wrong: have %v; want %v\n", conf.LocalAddr(), v)
		}
	}

	for _, v := range wrongData {
		if err := conf.SetLocalAddr(v); err == nil {
			t.Errorf("set wrong addr but didn't fail: %v\n", v)
		}
	}
}

func TestClientConfigSetPidFilePath(t *testing.T) {
	conf := &ClientConfig{}
	trueData := []string{"/tmp/a.pid", "~/.tmp/a.pid"}
	wrongData := []string{"", "tmp/a.pid", "a.pid"}

	for _, v := range trueData {
		if err := conf.SetPidFilePath(v); err != nil {
			t.Errorf("set pidfile failed: %v\n", v)
		} else if conf.PidFilePath() != v {
			t.Errorf("set pidfile wrong: have %v; want %v\n", conf.PidFilePath(), v)
		}
	}

	for _, v := range wrongData {
		if err := conf.SetPidFilePath(v); err == nil {
			t.Errorf("set wrong pidfile but didn't fail: %v\n", v)
		}
	}
}

func TestClientConfigGenArgs(t *testing.T) {
	testData := []*struct {
		c    *ClientConfig
		args string
	}{
		{
			c: &ClientConfig{
				Addr:       "172.17.0.1",
				Port:       "8888",
				IsFastOpen: false,
				PidFile:    "/tmp/a.pid",
			},
			args: "-b 172.17.0.1 -l 8888 --pid-file /tmp/a.pid",
		},
		{
			c: &ClientConfig{
				Addr:       "",
				Port:       "",
				IsFastOpen: false,
				PidFile:    "",
			},
			args: "-b " + defaultAddr + " -l " + defaultPort + " --pid-file " + defaultPidFile,
		},
		{
			c: &ClientConfig{
				Addr:       "172.17.0.1",
				Port:       "8888",
				IsFastOpen: true,
				PidFile:    "/tmp/a.pid",
			},
			args: "-b 172.17.0.1 -l 8888 --pid-file /tmp/a.pid --fast-open",
		},
	}

	for _, v := range testData {
		args := v.c.GenArgs()
		if strings.Join(args, " ") != v.args {
			t.Errorf("genargs failed:\nArgs: %v\n", args)
		}
	}
}
