package widgets

import (
	"sort"
	"strings"

	"github.com/therecipe/qt/widgets"

	"schannel-qt5/config"
)

var (
	// 代理的可用协议
	protocols = []string{
		"http",
		"https",
		"socks5",
	}
)

// schannel-qt5客户端配置控件
type ClientConfigWidget struct {
	widgets.QWidget

	// 配置改变后发送通知
	_ func() `signal:"valueChanged"`

	// client设置
	logFile    *widgets.QLineEdit
	logFileMsg *ColorLabel

	// ssr设置
	nodeConfigPath, ssrConfigPath, binPath          *widgets.QLineEdit
	nodeConfigPathMsg, ssrConfigPathMsg, binPathMsg *ColorLabel

	// 代理设置
	proxy     *widgets.QLineEdit
	proxyType *widgets.QComboBox
	proxyBox  *widgets.QGroupBox
	proxyMsg  *ColorLabel

	// 配置数据和接口
	conf *config.UserConfig
}

// 根据UserConfig创建ClientConfigWidget
func NewClientConfigWidget2(conf *config.UserConfig) *ClientConfigWidget {
	cw := NewClientConfigWidget(nil, 0)
	cw.conf = conf
	cw.InitUI()

	return cw
}

func (cw *ClientConfigWidget) InitUI() {
	// user设置布局
	userBox := widgets.NewQGroupBox2("客户端设置（重启生效）", nil)

	cw.logFile = widgets.NewQLineEdit2(cw.conf.LogFile.String(), nil)
	cw.logFile.SetPlaceholderText("日志文件保存路径")
	cw.logFile.ConnectTextChanged(func(_ string) {
		cw.ValueChanged()
	})
	cw.logFileMsg = NewColorLabelWithColor("路径需要为绝对路径且不能为目录", "red")
	cw.logFileMsg.Hide()

	userLayout := widgets.NewQFormLayout(nil)
	userLayout.AddRow3("日志文件路径:", cw.logFile)
	userLayout.AddRow5(cw.logFileMsg)
	userBox.SetLayout(userLayout)

	// ssr设置布局
	ssrBox := widgets.NewQGroupBox2("ssr设置", nil)
	ssrLayout := widgets.NewQFormLayout(nil)
	cw.ssrConfigPath = widgets.NewQLineEdit2(cw.conf.SSRClientConfigPath.String(), nil)
	cw.ssrConfigPath.SetPlaceholderText("绝对路径")
	cw.ssrConfigPathMsg = NewColorLabelWithColor("路径需要为绝对路径且不能为目录", "red")
	cw.ssrConfigPathMsg.Hide()
	ssrLayout.AddRow3("客户端配置路径：", cw.ssrConfigPath)
	ssrLayout.AddRow5(cw.ssrConfigPathMsg)

	cw.nodeConfigPath = widgets.NewQLineEdit2(cw.conf.SSRNodeConfigPath.String(), nil)
	cw.nodeConfigPath.SetPlaceholderText("绝对路径")
	cw.nodeConfigPathMsg = NewColorLabelWithColor("路径需要为绝对路径且不能为目录", "red")
	cw.nodeConfigPathMsg.Hide()
	ssrLayout.AddRow3("节点配置路径：", cw.nodeConfigPath)
	ssrLayout.AddRow5(cw.nodeConfigPathMsg)

	cw.binPath = widgets.NewQLineEdit2(cw.conf.SSRBin.String(), nil)
	cw.binPath.SetPlaceholderText("绝对路径")
	cw.binPathMsg = NewColorLabelWithColor("路径需要为绝对路径且不能为目录", "red")
	cw.binPathMsg.Hide()
	ssrLayout.AddRow3("程序路径：", cw.binPath)
	ssrLayout.AddRow5(cw.binPathMsg)
	ssrBox.SetLayout(ssrLayout)

	// 对协议列表排序，方便查找
	sort.Strings(protocols)

	// proxy设置，可选
	cw.proxyBox = widgets.NewQGroupBox2("使用代理", nil)
	cw.proxyBox.SetCheckable(true)

	typeLabel := widgets.NewQLabel2("协议类型:", nil, 0)
	cw.proxyType = widgets.NewQComboBox(nil)
	cw.proxyType.AddItems(protocols)
	proxyLabel := widgets.NewQLabel2("代理服务器地址:", nil, 0)
	cw.proxy = widgets.NewQLineEdit(nil)
	cw.proxy.SetPlaceholderText("URL")
	cw.proxyMsg = NewColorLabelWithColor("不是合法的URL", "red")
	cw.proxyMsg.Hide()

	// 根据配置确定是否勾选代理设置
	if cw.conf.Proxy.String() != "" {
		proto, host := cw.splitProtoHost()
		// 显示配置的协议
		cw.proxyType.SetCurrentIndex(sort.SearchStrings(protocols, proto))
		cw.proxy.SetText(host)
		cw.proxyBox.SetChecked(true)
	} else {
		cw.proxyBox.SetChecked(false)
	}

	typeLayout := widgets.NewQHBoxLayout()
	typeLayout.AddWidget(typeLabel, 0, 0)
	typeLayout.AddWidget(cw.proxyType, 0, 0)
	urlLayout := widgets.NewQHBoxLayout()
	urlLayout.AddWidget(proxyLabel, 0, 0)
	urlLayout.AddWidget(cw.proxy, 0, 0)

	proxyLayout := widgets.NewQVBoxLayout()
	proxyLayout.AddLayout(typeLayout, 0)
	proxyLayout.AddLayout(urlLayout, 0)
	proxyLayout.AddWidget(cw.proxyMsg, 0, 0)
	cw.proxyBox.SetLayout(proxyLayout)

	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.AddWidget(userBox, 0, 0)
	mainLayout.AddWidget(ssrBox, 0, 0)
	mainLayout.AddWidget(cw.proxyBox, 0, 0)
	cw.SetLayout(mainLayout)
}

// splitProtoHost 分割返回协议和主机名
func (w *ClientConfigWidget) splitProtoHost() (proto, host string) {
	data := strings.Split(w.conf.Proxy.String(), "://")
	proto = data[0]
	host = data[1]

	return
}

// UpdateClientConfig 更新UserClient，传递的为UserConfig的引用，可以直接修改
func (cw *ClientConfigWidget) UpdateClientConfig() error {
	var err error

	err = cw.validLogFile()
	showErrorMsg(cw.logFileMsg, err)

	err = cw.validSSRConfigPath()
	showErrorMsg(cw.ssrConfigPathMsg, err)

	err = cw.validNodeConfigPath()
	showErrorMsg(cw.nodeConfigPathMsg, err)

	err = cw.validBinPath()
	showErrorMsg(cw.binPathMsg, err)

	err = cw.validProxy()
	showErrorMsg(cw.proxyMsg, err)

	return err
}

// GetProxyURL 返回拼接了type后的URL
func (cw *ClientConfigWidget) GetProxyUrl() string {
	if !cw.proxyBox.IsChecked() {
		return ""
	}

	pType := cw.proxyType.CurrentText()
	return pType + "://" + cw.proxy.Text()
}

// validProxy 验证proxy URL是否合法
func (cw *ClientConfigWidget) validProxy() error {
	url := cw.GetProxyUrl()
	p := config.JSONProxy{Data: url}
	if !p.IsURL() && p.String() != "" {
		return config.ErrNotURL
	}

	return nil
}

// validLogFile 验证日志文件保存路径是否在$HOME下或者是绝对路径
func (cw *ClientConfigWidget) validLogFile() error {
	text := cw.logFile.Text()
	return checkEmptyPath(text)
}

// validNodeConfigPath 验证ssr配置文件路径是否在$HOME下或者是绝对路径
func (cw *ClientConfigWidget) validNodeConfigPath() error {
	text := cw.nodeConfigPath.Text()
	return checkPath(text)
}

// validBinPath 验证ssr可执行文件路径是否在$HOME下或者是绝对路径
func (cw *ClientConfigWidget) validBinPath() error {
	text := cw.binPath.Text()
	return checkPath(text)
}

// validSSRConfigPath 验证ssr可执行文件路径是否在$HOME下或者是绝对路径
func (cw *ClientConfigWidget) validSSRConfigPath() error {
	text := cw.ssrConfigPath.Text()
	return checkPath(text)
}
