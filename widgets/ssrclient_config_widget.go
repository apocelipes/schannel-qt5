package widgets

import (
	"fmt"
	"strconv"

	"github.com/therecipe/qt/widgets"

	"schannel-qt5/config"
)

// 设置ssr client
type SSRConfigWidget struct {
	widgets.QWidget

	// 通知配置值发生变化
	_ func() `signal:"valueChanged"`

	// 本地监听地址
	localAddr    *widgets.QLineEdit
	localAddrMsg *ColorLabel
	localPort    *widgets.QSpinBox
	localPortMsg *ColorLabel

	// pid-file存放路径
	pidFilePath    *widgets.QLineEdit
	pidFilePathMsg *ColorLabel

	// 是否使用fast-open
	fastOpen *widgets.QCheckBox

	conf config.ClientConfig
}

// NewSSRConfigWidget2 创建ssr client config widget
func NewSSRConfigWidget2(conf config.ClientConfig) *SSRConfigWidget {
	widget := NewSSRConfigWidget(nil, 0)
	widget.conf = conf
	widget.InitUI()

	return widget
}

// InitUI 初始化界面
func (s *SSRConfigWidget) InitUI() {
	group := widgets.NewQGroupBox2("ssr client设置", nil)
	groupLayout := widgets.NewQFormLayout(nil)

	s.localAddr = widgets.NewQLineEdit(nil)
	s.localAddr.SetPlaceholderText("绑定本地ip地址")
	s.localAddr.SetText(s.conf.LocalAddr())
	s.localAddr.ConnectTextChanged(func(_ string) {
		s.ValueChanged()
	})
	s.localAddrMsg = NewColorLabelWithColor("不是合法的ip地址", "red")
	s.localAddrMsg.Hide()
	groupLayout.AddRow3("本地地址：", s.localAddr)
	groupLayout.AddRow5(s.localAddrMsg)

	s.localPort = widgets.NewQSpinBox(nil)
	// 端口从1024-65535
	s.localPort.SetRange(1024, 65535)
	port, _ := strconv.Atoi(s.conf.LocalPort())
	s.localPort.SetValue(port)
	s.localPort.ConnectValueChanged(func(_ int) {
		s.ValueChanged()
	})
	s.localPortMsg = NewColorLabelWithColor("不是合法的端口值", "red")
	s.localPortMsg.Hide()
	groupLayout.AddRow3("本地端口：", s.localPort)
	groupLayout.AddRow5(s.localPortMsg)

	s.pidFilePath = widgets.NewQLineEdit(nil)
	s.pidFilePath.SetPlaceholderText("绝对路径")
	s.pidFilePath.SetText(s.conf.PidFilePath())
	s.pidFilePath.ConnectTextChanged(func(_ string) {
		s.ValueChanged()
	})
	s.pidFilePathMsg = NewColorLabelWithColor("不是合法的路径", "red")
	s.pidFilePathMsg.Hide()
	groupLayout.AddRow3("pid-file路径：", s.pidFilePath)
	groupLayout.AddRow5(s.pidFilePathMsg)

	// 检查内核版本
	versionInfo := widgets.NewQLabel(nil, 0)
	s.fastOpen = widgets.NewQCheckBox2("启用fast-open", nil)
	s.fastOpen.SetToolTip("此功能需要Linux kernel版本 >= 3.7.0")
	s.fastOpen.SetEnabled(false)
	if version, err := kernelVersion(); err != nil {
		versionInfo.SetText(err.Error())
	} else if fastOpenAble(version) {
		s.fastOpen.SetEnabled(true)
		s.fastOpen.SetChecked(s.conf.FastOpen())
		versionInfo.SetText(fmt.Sprintf("Linux kernel: %v", version))
	} else {
		versionInfo.SetText(fmt.Sprintf("内核版本不支持fast-open: %v", version))
	}
	s.fastOpen.ConnectClicked(func(_ bool) {
		s.ValueChanged()
	})
	groupLayout.AddRow5(s.fastOpen)
	groupLayout.AddRow5(versionInfo)

	group.SetLayout(groupLayout)
	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.AddWidget(group, 0, 0)
	s.SetLayout(mainLayout)
}

// UpdateSSRClientConfig 更新config，如果数据不合法则返回error
// 因为传递了引用类型，所以直接修改config对象
func (s *SSRConfigWidget) UpdateSSRClientConfig() error {
	// 可选时才设置fast-open值
	if s.fastOpen.IsEnabled() {
		s.conf.SetFastOpen(s.fastOpen.IsChecked())
	}

	// 记录返回值
	var err error
	var errRes error
	err = s.conf.SetLocalPort(strconv.Itoa(s.localPort.Value()))
	if showErrorMsg(s.localPortMsg, err) {
		errRes = err
	}

	err = s.conf.SetLocalAddr(s.localAddr.Text())
	if showErrorMsg(s.localAddrMsg, err) {
		errRes = err
	}

	err = s.conf.SetPidFilePath(s.pidFilePath.Text())
	if showErrorMsg(s.pidFilePathMsg, err) {
		errRes = err
	}

	return errRes
}
