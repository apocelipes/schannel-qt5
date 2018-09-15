package widgets

import (
	"errors"
	"github.com/therecipe/qt/widgets"
	"regexp"
	"schannel-qt5/config"
	"sort"
	"strings"
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
	name, passwd, logFile *widgets.QLineEdit
	nameMsg, passwdMsg, logFileMsg   *ColorLabel
	// ssr设置
	configPath, binPath       *widgets.QLineEdit
	configPathMsg, binPathMsg *ColorLabel
	// 代理设置
	proxy     *widgets.QLineEdit
	proxyType *widgets.QComboBox
	proxyBox  *widgets.QGroupBox
	proxyMsg  *ColorLabel

	// 配置数据和接口
	conf *config.UserConfig
}

// NewConfigWidget2 根据conf生成ConfigWidget
func NewConfigWidget2(conf *config.UserConfig) *ConfigWidget {
	if conf == nil {
		return nil
	}
	widget := NewConfigWidget(nil, 0)
	widget.conf = conf
	widget.InitUI()

	return widget
}

// InitUI 初始化并显示
func (w *ConfigWidget) InitUI() {
	// client设置布局
	clientBox := widgets.NewQGroupBox2("客户端设置（用户名和密码重启生效）", nil)
	nameLabel := widgets.NewQLabel2("用户名:", nil, 0)
	w.name = widgets.NewQLineEdit2(w.conf.UserName, nil)
	w.name.SetPlaceholderText("邮箱")
	w.nameMsg = NewColorLabelWithColor("用户名必须为邮箱", "red")
	w.nameMsg.Hide()

	passwdLabel := widgets.NewQLabel2("密码:     ", nil, 0)
	w.passwd = widgets.NewQLineEdit(nil)
	w.passwd.SetPlaceholderText("密码")
	// 设置密码默认不可见
	w.passwd.SetEchoMode(widgets.QLineEdit__Password)
	w.passwd.SetText(w.conf.Passwd)
	w.passwdMsg = NewColorLabelWithColor("密码不能为空", "red")
	w.passwdMsg.Hide()
	echoCheck := widgets.NewQCheckBox2("密码可见", nil)
	echoCheck.ConnectClicked(func(_ bool) {
		if echoCheck.IsChecked() {
			w.passwd.SetEchoMode(widgets.QLineEdit__Normal)
			return
		}

		w.passwd.SetEchoMode(widgets.QLineEdit__Password)
	})

	logFileLabel := widgets.NewQLabel2("日志文件路径", nil, 0)
	w.logFile = widgets.NewQLineEdit(nil)
	w.logFile.SetPlaceholderText("日志文件保存路径")
	w.logFile.SetText(w.conf.LogFile.String())
	w.logFileMsg = NewColorLabelWithColor("路径需要为绝对路径", "red")
	w.logFileMsg.Hide()

	nameLayout := widgets.NewQHBoxLayout()
	nameLayout.AddWidget(nameLabel, 0, 0)
	nameLayout.AddWidget(w.name, 0, 0)
	passwdLayout := widgets.NewQHBoxLayout()
	passwdLayout.AddWidget(passwdLabel, 0, 0)
	passwdLayout.AddWidget(w.passwd, 0, 0)
	logFileLayout := widgets.NewQHBoxLayout()
	logFileLayout.AddWidget(logFileLabel, 0, 0)
	logFileLayout.AddWidget(w.logFile, 0, 0)
	clientLayout := widgets.NewQVBoxLayout()
	clientLayout.AddLayout(nameLayout, 0)
	clientLayout.AddWidget(w.nameMsg, 0, 0)
	clientLayout.AddLayout(passwdLayout, 0)
	clientLayout.AddWidget(echoCheck, 0, 0)
	clientLayout.AddWidget(w.passwdMsg, 0, 0)
	clientLayout.AddLayout(logFileLayout, 0)
	clientLayout.AddWidget(w.logFileMsg, 0, 0)
	clientBox.SetLayout(clientLayout)

	// ssr设置布局
	ssrBox := widgets.NewQGroupBox2("ssr设置", nil)
	configLabel := widgets.NewQLabel2("ssr配置文件路径:", nil, 0)
	w.configPath = widgets.NewQLineEdit2(w.conf.SSRConfigPath.String(), nil)
	w.configPath.SetPlaceholderText("绝对路径")
	w.configPathMsg = NewColorLabelWithColor("路径需要为绝对路径", "red")
	w.configPathMsg.Hide()

	binLabel := widgets.NewQLabel2("ssr执行文件路径:", nil, 0)
	w.binPath = widgets.NewQLineEdit2(w.conf.SSRBin.String(), nil)
	w.binPath.SetPlaceholderText("绝对路径")
	w.binPathMsg = NewColorLabelWithColor("路径需要为绝对路径", "red")
	w.binPathMsg.Hide()

	configLayout := widgets.NewQHBoxLayout()
	configLayout.AddWidget(configLabel, 0, 0)
	configLayout.AddWidget(w.configPath, 0, 0)
	binLayout := widgets.NewQHBoxLayout()
	binLayout.AddWidget(binLabel, 0, 0)
	binLayout.AddWidget(w.binPath, 0, 0)
	ssrLayout := widgets.NewQVBoxLayout()
	ssrLayout.AddLayout(configLayout, 0)
	ssrLayout.AddWidget(w.configPathMsg, 0, 0)
	ssrLayout.AddLayout(binLayout, 0)
	ssrLayout.AddWidget(w.binPathMsg, 0, 0)
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

	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.AddWidget(clientBox, 0, 0)
	mainLayout.AddWidget(ssrBox, 0, 0)
	mainLayout.AddWidget(w.proxyBox, 0, 0)
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

	err = w.validName()
	if showErrorMsg(w.nameMsg, err) {
		flag = true
	}

	err = w.validPassword()
	if showErrorMsg(w.passwdMsg, err) {
		flag = true
	}

	err = w.validLogFile()
	if showErrorMsg(w.logFileMsg, err) {
		flag = true
	}

	err = w.validConfigPath()
	if showErrorMsg(w.configPathMsg, err) {
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

	if flag {
		return
	}

	w.conf = w.getConfig()
	if err := w.conf.StoreConfig(); err != nil {
		errorMsg := widgets.NewQErrorMessage(nil)
		errorMsg.ShowMessage("保存出错: " + err.Error())
		errorMsg.Exec()
		return
	}

	// 通知其他组件配置发生变化
	w.ConfigChanged(w.conf)
}

// showErrorMsg 控制error label的显示
// err为nil则代表没有错误发生，如果label可见则设为隐藏
// err不为nil时设置label可见
// 设置label可见时返回true，否则返回false（不受label原有状态影响）
func showErrorMsg(label *ColorLabel, err error) bool {
	if err != nil {
		label.Show()
		return true
	}

	label.Hide()
	return false
}

// getConfig 根据设置信息生成新的config对象
func (w *ConfigWidget) getConfig() *config.UserConfig {
	conf := &config.UserConfig{}
	conf.UserName = w.name.Text()
	conf.Passwd = w.passwd.Text()
	conf.LogFile = config.JSONPath{Data: w.logFile.Text()}
	conf.SSRConfigPath = config.JSONPath{Data: w.configPath.Text()}
	conf.SSRBin = config.JSONPath{Data: w.binPath.Text()}
	conf.Proxy = config.JSONProxy{Data: w.GetProxyUrl()}

	return conf
}

// GetProxyURL 返回拼接了type后的URL
func (w *ConfigWidget) GetProxyUrl() string {
	if !w.proxyBox.IsChecked() {
		return ""
	}

	ptype := w.proxyType.CurrentText()
	return ptype + "://" + w.proxy.Text()
}

// validProxy 验证proxy URL是否合法
func (w *ConfigWidget) validProxy() error {
	p := config.JSONProxy{Data: w.GetProxyUrl()}
	if !p.IsURL() && p.String() != "" {
		return config.ErrNotURL
	}

	return nil
}

// validName 验证username是否是合法的邮箱地址
func (w *ConfigWidget) validName() error {
	email := regexp.MustCompile(`^[A-Za-z0-9\\u4e00-\\u9fa5]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
	if !email.MatchString(w.name.Text()) {
		return errors.New("not a valid email address")
	}

	return nil
}

// validPassword 验证password是否合法
func (w *ConfigWidget) validPassword() error {
	if w.passwd.Text() == "" {
		return errors.New("no valid password")
	}

	return nil
}

// validLogFile 验证日志文件保存路径是否在$HOME下或者是绝对路径
func (w *ConfigWidget) validLogFile() error {
	jpath := config.JSONPath{Data: w.logFile.Text()}
	if _, err := jpath.AbsPath(); err != nil {
		return err
	}

	return nil
}

// validConfigPath 验证ssr配置文件路径是否在$HOME下或者是绝对路径
func (w *ConfigWidget) validConfigPath() error {
	jpath := config.JSONPath{Data: w.configPath.Text()}
	if _, err := jpath.AbsPath(); err != nil {
		return err
	}

	return nil
}

// validBinPath 验证ssr可执行文件路径是否在$HOME下或者是绝对路径
func (w *ConfigWidget) validBinPath() error {
	jpath := config.JSONPath{Data: w.binPath.Text()}
	if _, err := jpath.AbsPath(); err != nil {
		return err
	}

	return nil
}
