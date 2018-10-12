package ssr

import (
	"time"

	"schannel-qt5/config"
)

// Launcher 用于控制ssr客户端
type Launcher interface {
	// Start 打开客户端
	Start() error
	// Restart 重启客户端
	Restart() error
	// Stop 关闭客户端
	Stop() error
	// IsRunning 客户端是否正在运行
	IsRunning() error
	// ConnectionCheck 检查ssr代理是否可用
	ConnectionCheck(timeout time.Duration) error
}

// LauncherMaker 生成Launcher的工厂函数
type LauncherMaker func(*config.UserConfig) Launcher

var (
	// 保存注册的LauncherMaker
	launchers = make(map[string]LauncherMaker)
	// ssr config注册获取
	configs = make(map[string]config.ClientConfigMaker)
)

// SetLuancherMaker 注册Launcher生成器
func SetLuancherMaker(name string, maker LauncherMaker) {
	if name == "" || maker == nil {
		panic("SetLauncher error: wrong name or LuancherMaker")
	}

	launchers[name] = maker
}

// SetClientConfigMaker 注册ClientConfig生成器
func SetClientConfigMaker(name string, maker config.ClientConfigMaker) {
	if name == "" || maker == nil {
		panic("SetClientConfigMaker error: wrong name or ClientConfigMaker")
	}

	configs[name] = maker
}

// NewLauncher 返回由name指定的Launcher生成器使用config.UserConfig生成的Launcher
func NewLauncher(name string, conf *config.UserConfig) Launcher {
	maker, ok := launchers[name]
	if !ok {
		return nil
	}

	return maker(conf)
}

// NewClientConfig 根据名字返回默认值的ClientConfig
func NewClientConfig(name string) config.ClientConfig {
	maker, ok := configs[name]
	if !ok {
		return nil
	}

	return maker()
}
