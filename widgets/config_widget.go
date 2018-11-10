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

// ConfigWidget 显示和设置本客户端的配置
type ConfigWidget struct {
	widgets.QWidget

	// 通知conf已经更新
	_ func(*config.UserConfig) `signal:"configChanged"`

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

	// ssr client设置
	ssrClientConfigWidget *SSRConfigWidget

	// 配置数据和接口
	conf *config.UserConfig
}

// NewConfigWidget2 根据conf生成ConfigWidget
func NewConfigWidget2(conf *config.UserConfig) *ConfigWidget {
	if conf == nil || conf.SSRClientConfig == nil {
		return nil
	}
	widget := NewConfigWidget(nil, 0)
	widget.conf = conf
	widget.InitUI()

	return widget
}

// InitUI 初始化并显示
func (w *ConfigWidget) InitUI() {
	// user设置布局
	userBox := widgets.NewQGroupBox2("客户端设置（重启生效）", nil)

	logFileLabel := widgets.NewQLabel2("日志文件路径:", nil, 0)
	w.logFile = widgets.NewQLineEdit2(w.conf.LogFile.String(), nil)
	w.logFile.SetPlaceholderText("日志文件保存路径")
	w.logFileMsg = NewColorLabelWithColor("路径需要为绝对路径且不能为目录", "red")
	w.logFileMsg.Hide()

	logFileLayout := widgets.NewQHBoxLayout()
	logFileLayout.AddWidget(logFileLabel, 1, 0)
	logFileLayout.AddWidget(w.logFile, 1, 0)

	userLayout := widgets.NewQVBoxLayout()
	userLayout.AddLayout(logFileLayout, 1)
	userLayout.AddWidget(w.logFileMsg, 1, 0)
	userBox.SetLayout(userLayout)

	// ssr设置布局
	ssrBox := widgets.NewQGroupBox2("ssr设置", nil)
	ssrLayout := widgets.NewQFormLayout(nil)
	w.ssrConfigPath = widgets.NewQLineEdit2(w.conf.SSRClientConfigPath.String(), nil)
	w.ssrConfigPath.SetPlaceholderText("绝对路径")
	w.ssrConfigPathMsg = NewColorLabelWithColor("路径需要为绝对路径且不能为目录", "red")
	w.ssrConfigPathMsg.Hide()
	ssrLayout.AddRow3("客户端配置路径：", w.ssrConfigPath)
	ssrLayout.AddRow5(w.ssrConfigPathMsg)

	w.nodeConfigPath = widgets.NewQLineEdit2(w.conf.SSRNodeConfigPath.String(), nil)
	w.nodeConfigPath.SetPlaceholderText("绝对路径")
	w.nodeConfigPathMsg = NewColorLabelWithColor("路径需要为绝对路径且不能为目录", "red")
	w.nodeConfigPathMsg.Hide()
	ssrLayout.AddRow3("节点配置路径：", w.nodeConfigPath)
	ssrLayout.AddRow5(w.nodeConfigPathMsg)

	w.binPath = widgets.NewQLineEdit2(w.conf.SSRBin.String(), nil)
	w.binPath.SetPlaceholderText("绝对路径")
	w.binPathMsg = NewColorLabelWithColor("路径需要为绝对路径且不能为目录", "red")
	w.binPathMsg.Hide()
	ssrLayout.AddRow3("程序路径：", w.binPath)
	ssrLayout.AddRow5(w.binPathMsg)
	ssrBox.SetLayout(ssrLayout)

	// 对协议列表排序，方便查找
	sort.Strings(protocols)

	// proxy设置，可选
	w.proxyBox = widgets.NewQGroupBox2("使用代理", nil)
	w.proxyBox.SetCheckable(true)

	typeLabel := widgets.NewQLabel2("协议类型:", nil, 0)
	w.proxyType = widgets.NewQComboBox(nil)
	w.proxyType.AddItems(protocols)
	proxyLabel := widgets.NewQLabel2("代理服务器地址:", nil, 0)
	w.proxy = widgets.NewQLineEdit(nil)
	w.proxy.SetPlaceholderText("URL")
	w.proxyMsg = NewColorLabelWithColor("不是合法的URL", "red")
	w.proxyMsg.Hide()

	// 根据配置确定是否勾选代理设置
	if w.conf.Proxy.String() != "" {
		proto, host := w.splitProtoHost()
		// 显示配置的协议
		w.proxyType.SetCurrentIndex(sort.SearchStrings(protocols, proto))
		w.proxy.SetText(host)
		w.proxyBox.SetChecked(true)
	} else {
		w.proxyBox.SetChecked(false)
	}

	w.ssrClientConfigWidget = NewSSRConfigWidget2(w.conf.SSRClientConfig)

	typeLayout := widgets.NewQHBoxLayout()
	typeLayout.AddWidget(typeLabel, 0, 0)
	typeLayout.AddWidget(w.proxyType, 0, 0)
	urlLayout := widgets.NewQHBoxLayout()
	urlLayout.AddWidget(proxyLabel, 0, 0)
	urlLayout.AddWidget(w.proxy, 0, 0)

	proxyLayout := widgets.NewQVBoxLayout()
	proxyLayout.AddLayout(typeLayout, 0)
	proxyLayout.AddLayout(urlLayout, 0)
	proxyLayout.AddWidget(w.proxyMsg, 0, 0)
	w.proxyBox.SetLayout(proxyLayout)

	saveButton := widgets.NewQPushButton2("保存", nil)
	saveButton.ConnectClicked(w.saveConfig)

	leftLayout := widgets.NewQVBoxLayout()
	leftLayout.AddWidget(userBox, 0, 0)
	leftLayout.AddWidget(ssrBox, 0, 0)
	leftLayout.AddWidget(w.proxyBox, 0, 0)
	rightLayout := widgets.NewQVBoxLayout()
	rightLayout.AddWidget(w.ssrClientConfigWidget, 0, 0)
	// 防止Grid被过度拉伸
	rightLayout.AddStretch(0)
	topLayout := widgets.NewQHBoxLayout()
	topLayout.AddLayout(leftLayout, 0)
	topLayout.AddLayout(rightLayout, 0)
	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.AddLayout(topLayout, 0)
	mainLayout.AddStretch(0)
	mainLayout.AddWidget(saveButton, 0, 0)
	w.SetLayout(mainLayout)
}

// splitProtoHost 分割返回协议和主机名
func (w *ConfigWidget) splitProtoHost() (proto, host string) {
	data := strings.Split(w.conf.Proxy.String(), "://")
	proto = data[0]
	host = data[1]

	return
}

// saveConfig 验证并保存配置
func (w *ConfigWidget) saveConfig(_ bool) {
	var err error
	// flag为true时代表验证不通过
	var flag bool

	// 保存时不可修改设置信息
	w.SetEnabled(false)
	defer w.SetEnabled(true)

	err = w.validLogFile()
	if showErrorMsg(w.logFileMsg, err) {
		flag = true
	}

	err = w.validSSRConfigPath()
	if showErrorMsg(w.ssrConfigPathMsg, err) {
		flag = true
	}

	err = w.validNodeConfigPath()
	if showErrorMsg(w.nodeConfigPathMsg, err) {
		flag = true
	}

	err = w.validBinPath()
	if showErrorMsg(w.binPathMsg, err) {
		flag = true
	}

	err = w.validProxy()
	if showErrorMsg(w.proxyMsg, err) {
		flag = true
	}

	// 更新ssr client config
	err = w.ssrClientConfigWidget.UpdateSSRClientConfig()
	if err != nil {
		flag = true
	}

	if flag {
		return
	}

	conf := w.getConfig()
	// ssr client conf被直接更新
	conf.SSRClientConfig = w.conf.SSRClientConfig

	w.conf = conf
	if err := w.conf.StoreConfig(); err != nil {
		errorMsg := widgets.NewQErrorMessage(nil)
		errorMsg.ShowMessage("保存出错: " + err.Error())
		errorMsg.Exec()
		return
	}

	// 通知其他组件配置发生变化
	w.ConfigChanged(w.conf)
}

// getConfig 根据设置信息生成新的config对象
func (w *ConfigWidget) getConfig() *config.UserConfig {
	conf := &config.UserConfig{}
	conf.LogFile = config.JSONPath{Data: w.logFile.Text()}
	conf.SSRNodeConfigPath = config.JSONPath{Data: w.nodeConfigPath.Text()}
	conf.SSRBin = config.JSONPath{Data: w.binPath.Text()}
	conf.Proxy = config.JSONProxy{Data: w.GetProxyUrl()}
	conf.SSRClientConfigPath = config.JSONPath{Data: w.ssrConfigPath.Text()}

	return conf
}

// GetProxyURL 返回拼接了type后的URL
func (w *ConfigWidget) GetProxyUrl() string {
	if !w.proxyBox.IsChecked() {
		return ""
	}

	pType := w.proxyType.CurrentText()
	return pType + "://" + w.proxy.Text()
}

// validProxy 验证proxy URL是否合法
func (w *ConfigWidget) validProxy() error {
	url := w.GetProxyUrl()
	p := config.JSONProxy{Data: url}
	if !p.IsURL() && p.String() != "" {
		return config.ErrNotURL
	}

	return nil
}

// validLogFile 验证日志文件保存路径是否在$HOME下或者是绝对路径
func (w *ConfigWidget) validLogFile() error {
	text := w.logFile.Text()
	return checkPath(text)
}

// validNodeConfigPath 验证ssr配置文件路径是否在$HOME下或者是绝对路径
func (w *ConfigWidget) validNodeConfigPath() error {
	text := w.nodeConfigPath.Text()
	return checkPath(text)
}

// validBinPath 验证ssr可执行文件路径是否在$HOME下或者是绝对路径
func (w *ConfigWidget) validBinPath() error {
	text := w.binPath.Text()
	return checkPath(text)
}

// validSSRConfigPath 验证ssr可执行文件路径是否在$HOME下或者是绝对路径
func (w *ConfigWidget) validSSRConfigPath() error {
	text := w.ssrConfigPath.Text()
	return checkPath(text)
}
