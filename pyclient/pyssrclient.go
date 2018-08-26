package pyclient

import (
	"os"
	"os/exec"

	"schannel-qt5/config"
	"schannel-qt5/ssr"
)

// PySSRClient 调用Python实现的ssr客户端
type PySSRClient struct {
	bin string
	// binArg 运行参数
	binArg string
	config string
}

func init() {
	// 注册为可用的Launcher，name为python
	ssr.SetLuancherMaker("python", ssr.LauncherMaker(NewPySSRClient))
}

// NewPySSRClient 这个函数供ssr.LauncherMaker调用，用于生成ssr.Launcher
func NewPySSRClient(c *config.UserConfig) ssr.Launcher {
	p := new(PySSRClient)
	bin, err := c.SSRBin.AbsPath()
	if err != nil {
		return nil
	}
	p.bin = bin
	p.config, err = c.SSRConfigPath.AbsPath()
	if err != nil {
		return nil
	}

	// -c ssr_config_file
	p.binArg = "-c" + p.config

	return p
}

// Start 启动客户端
func (p *PySSRClient) Start() error {
	// env tells sudo to find commands in $PATH
	cmd := exec.Command("sudo", "env PATH=$PATH", "python", p.bin, p.binArg, "-d", "start")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Restart 重新启动客户端
func (p *PySSRClient) Restart() error {
	cmd := exec.Command("sudo", "env PATH=$PATH", "python", p.bin, p.binArg, "-d", "restart")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Stop 停止客户端
func (p *PySSRClient) Stop() error {
	cmd := exec.Command("sudo", "env PATH=$PATH", "python", p.bin, p.binArg, "-d", "stop")
	cmd.Stderr = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
