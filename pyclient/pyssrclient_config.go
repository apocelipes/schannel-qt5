package pyclient

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"schannel-qt5/config"
	"strconv"
)

// ClientConfig pyssrclient的本地配置
type ClientConfig struct {
	// 本地端口和ip(default: 127.0.0.1:1080)
	Addr string `json:"local_addr,omitempty"`
	Port string `json:"local_port,omitempty"`
	// fast-open 需要linux 3.7+(default: false)
	IsFastOpen bool `json:"fast-open,omitempty"`
	// pidfile存放位置(default: /tmp/ssr_client.pid)
	PidFile string `json:"pid-file,omitempty"`
}

// 实现ClientConfigGetter
func (c *ClientConfig) LocalPort() string {
	if c.Port == "" {
		return "1080"
	}

	return c.Port
}

func (c *ClientConfig) LocalAddr() string {
	if c.Addr == "" {
		return "127.0.0.1"
	}

	return c.Addr
}

func (c *ClientConfig) FastOpen() bool {
	return c.IsFastOpen
}

func (c *ClientConfig) PidFilePath() string {
	if c.PidFile == "" {
		return "/tmp/ssr_pyclient.pid"
	}

	return c.PidFile
}

func (c *ClientConfig) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientConfig) Store(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	f.Write(data)

	return nil
}

// 实现ClientConfigSetter
func (c *ClientConfig) SetLocalPort(port string) error {
	_, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	c.Port = port
	return nil
}

func (c *ClientConfig) SetFastOpen(isFOP bool) {
	c.IsFastOpen = isFOP
}

func (c *ClientConfig) SetPidFilePath(path string) error {
	jpath := config.JSONPath{Data: path}
	if _, err := jpath.AbsPath(); err != nil {
		return err
	}

	c.PidFile = path
	return nil
}

func (c *ClientConfig) SetLocalAddr(addr string) error {
	IP := regexp.MustCompile(`^((25[0-5]|2[0-4]\d|[1]{1}\d{1}\d{1}|[1-9]{1}\d{1}|\d{1})($|(?!\.$)\.)){4}$`)
	if !IP.MatchString(addr) {
		return errors.New("not a valid ip addr")
	}

	c.Addr = addr
	return nil
}

func (c *ClientConfig) GenArgs() []string {
	args := make([]string, 0)
	args = append(args, "-b", c.LocalAddr())
	args = append(args, "-l", c.LocalPort())
	args = append(args, "--pid-file ", c.PidFilePath())
	if c.IsFastOpen {
		args = append(args, "--fast-open")
	}

	return args
}
