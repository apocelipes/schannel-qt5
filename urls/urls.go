package urls

const (
	// RootPath 网站主URL
	RootPath        = `https://sgchannel.me/`
	// AuthPath 登录页面的URL
	AccountPath        = RootPath + `clientarea.php`
	// ServiceListPath 服务列表URL
	ServiceListPath = AccountPath + `?action=services`
	// LoginPath 登录验证的URL
	LoginPath       = RootPath + `dologin.php`
	// InvoicePath 账单列表的URL
	InvoicePath     = AccountPath + `?action=invoices`
	// 测试代理的URL
	ProxyTestPath   = `https://golang.org`
)
