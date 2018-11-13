package pyclient

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"

	"schannel-qt5/config"
	"schannel-qt5/ssr"
)

var (
	defaultPidFile = "/tmp/ssr_pyclient.pid"
	defaultPort    = "1080"
	defaultAddr    = "127.0.0.1"
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

func init() {
	// 注册到config生成器
	ssr.SetClientConfigMaker("python", config.ClientConfigMaker(newClientConfig))
}

// newClientConfig 生成config对象
func newClientConfig() config.ClientConfig {
	return &ClientConfig{}
}

// 实现ClientConfigGetter
func (c *ClientConfig) LocalPort() string {
	if c.Port == "" {
		return defaultPort
	}

	return c.Port
}

func (c *ClientConfig) LocalAddr() string {
	if c.Addr == "" {
		return defaultAddr
	}

	return c.Addr
}

func (c *ClientConfig) FastOpen() bool {
	return c.IsFastOpen
}

func (c *ClientConfig) PidFilePath() string {
	if c.PidFile == "" {
		return defaultPidFile
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
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer f.Close()

	// 格式化成易于阅读的形式
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	f.Write(data)

	return nil
}

// 实现ClientConfigSetter
// SetLocalPort 设置本地端口，端口不能大于65535且不能为0
func (c *ClientConfig) SetLocalPort(port string) error {
	i, err := strconv.Atoi(port)
	if err != nil {
		return err
	} else if i > 65535 || i == 0 {
		return errors.New("port over range")
	}

	c.Port = port
	return nil
}

func (c *ClientConfig) SetFastOpen(isFOP bool) {
	c.IsFastOpen = isFOP
}

// SetPidFilePath 设置pidfile存放路径，需要为绝对路径
func (c *ClientConfig) SetPidFilePath(path string) error {
	jpath := config.JSONPath{Data: path}
	if _, err := jpath.AbsPath(); err != nil {
		return err
	}

	c.PidFile = path
	return nil
}

// SetLocalAddr 检查并设置要bind的本地ip
// 暂时只支持IPv4
func (c *ClientConfig) SetLocalAddr(addr string) error {
	IP := regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
	if !IP.MatchString(addr) {
		return errors.New("not a valid ip addr")
	}

	c.Addr = addr
	return nil
}

// GenArgs 根据config对象生成命令行参数选项
func (c *ClientConfig) GenArgs() []string {
	args := make([]string, 0)
	args = append(args, "-b", c.LocalAddr())
	args = append(args, "-l", c.LocalPort())
	args = append(args, "--pid-file", c.PidFilePath())
	if c.IsFastOpen {
		args = append(args, "--fast-open")
	}

	return args
}
