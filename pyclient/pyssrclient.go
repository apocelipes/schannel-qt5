package pyclient

import (
	"net/http"
	"net/url"
	"os/exec"
	"time"

	"schannel-qt5/config"
	"schannel-qt5/ssr"
	"schannel-qt5/urls"
)

// PySSRClient 调用Python实现的ssr客户端
type PySSRClient struct {
	bin string
	// binArg 运行参数
	binArgs []string
	// 配置
	conf config.ClientConfig
}

func init() {
	// 注册为可用的Launcher，name为python
	ssr.SetLuancherMaker("python", ssr.LauncherMaker(newPySSRClient))
}

// newPySSRClient 这个函数供ssr.LauncherMaker调用，用于生成ssr.Launcher
func newPySSRClient(c *config.UserConfig) ssr.Launcher {
	p := new(PySSRClient)
	bin, err := c.SSRBin.AbsPath()
	if err != nil {
		return nil
	}
	p.bin = bin

	nodeConfigFile, err := c.SSRNodeConfigPath.AbsPath()
	if err != nil {
		return nil
	}

	p.conf = c.SSRClientConfig

	// -c ssr_node_config_file
	p.binArgs = []string{"python", p.bin}
	p.binArgs = append(p.binArgs, "-c", nodeConfigFile)
	p.binArgs = append(p.binArgs, p.conf.(*ClientConfig).GenArgs()...)

	return p
}

// Start 启动客户端
func (p *PySSRClient) Start() error {
	// 使用pkexec在gui程序中请求权限
	args := make([]string, len(p.binArgs))
	copy(args, p.binArgs)
	args = append(args, "-d", "start")
	cmd := exec.Command("pkexec", args...)
	return cmd.Run()
}

// Restart 重新启动客户端
func (p *PySSRClient) Restart() error {
	args := make([]string, len(p.binArgs))
	copy(args, p.binArgs)
	args = append(args, "-d", "restart")
	cmd := exec.Command("pkexec", args...)
	return cmd.Run()
}

// Stop 停止客户端
func (p *PySSRClient) Stop() error {
	args := make([]string, len(p.binArgs))
	copy(args, p.binArgs)
	args = append(args, "-d", "stop")
	cmd := exec.Command("pkexec", args...)
	return cmd.Run()
}

// ConnectionCheck 检查代理是否可用，不可用则返回error
func (p *PySSRClient) ConnectionCheck(timeout time.Duration) error {
	proxyURL, err := url.Parse("socks5://" + p.conf.LocalAddr() + ":" + p.conf.LocalPort())
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: timeout,
	}
	client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}

	request, err := http.NewRequest("GET", urls.RootPath, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
