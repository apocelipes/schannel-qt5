package config

// ClientConfigGetter 获取ClientConfig
type ClientConfigGetter interface {
	// LocalPort 获取本地监听端口
	LocalPort() string
	// LocalAddr 获取本地监听地址
	LocalAddr() string
	// FastOpen 是否使用fast-open
	FastOpen() bool
	// PidFilePath pidfile存放路径
	PidFilePath() string
}

// ClientConfigSetter 设置ClientConfig
type ClientConfigSetter interface {
	// SetLocalPort 设置本地监听端口
	SetLocalPort(port string) error
	// SetLocalAddr 设置本地监听地址
	SetLocalAddr(addr string) error
	// SetFastOpen 设置是否使用fast-open
	SetFastOpen(isFOP bool)
	// SetPidFilePath 设置pidfile存放路径
	SetPidFilePath(path string) error
}

// ssr client配置接口
type ClientConfig interface {
	ClientConfigGetter
	ClientConfigSetter
	// Load 从配置文件解析config
	Load(path string) error
	// Store 保存到配置文件
	Store(path string) error
}

// ClientConfigMaker 产生新的config对象，所有选项使用初始默认值
type ClientConfigMaker func() ClientConfig
