package config

// ssr client配置接口
type ClientConfig interface {
	// LocalPort 获取本地监听端口
	LocalPort() string
	// LocalAddr 获取本地监听地址
	LocalAddr() string
	// FastOpen 是否使用fast-open
	FastOpen() bool
	// PidFilePath pidfile存放路径
	PidFilePath() string
	// Load 从配置文件解析config
	Load(path string) error
	// Store 保存到配置文件
	Store(path string) error
}
