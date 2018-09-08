package ssr

import (
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
	// ConnectionCheck 检查ssr代理是否可用
	ConnectionCheck() bool
}

// LauncherMaker 生成Launcher的工厂函数
type LauncherMaker func(*config.UserConfig) Launcher

// 保存注册的LauncherMaker
var launchers = make(map[string]LauncherMaker)

// SetLuancherMaker 注册Launcher生成器
func SetLuancherMaker(name string, l LauncherMaker) {
	if name == "" || l == nil {
		panic("SetLauncher error: wrong name or LuancherMaker")
	}

	launchers[name] = l
}

// NewLauncher 返回由name指定的Launcher生成器使用config.UserConfig生成的Launcher
func NewLauncher(name string, conf *config.UserConfig) Launcher {
	l, ok := launchers[name]
	if !ok {
		return nil
	}

	return l(conf)
}
