package pyclient

import (
	"net/http"
	"net/url"
	"os/exec"
	"time"

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
	// 使用pkexec在gui程序中请求权限
	cmd := exec.Command("pkexec", "python", p.bin, p.binArg, "-d", "start")
	return cmd.Run()
}

// Restart 重新启动客户端
func (p *PySSRClient) Restart() error {
	cmd := exec.Command("pkexec", "python", p.bin, p.binArg, "-d", "restart")
	return cmd.Run()
}

// Stop 停止客户端
func (p *PySSRClient) Stop() error {
	cmd := exec.Command("pkexec", "python", p.bin, p.binArg, "-d", "stop")
	return cmd.Run()
}

// ConnectionCheck 检查代理是否可用，不可用则返回error
func (p *PySSRClient) ConnectionCheck(timeout time.Duration) error {
	proxyURL, err := url.Parse("socks5://127.0.0.1:1080")
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: timeout,
	}
	client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}

	request, err := http.NewRequest("GET", "https://www.google.com.hk", nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return err
	}

	return nil
}
