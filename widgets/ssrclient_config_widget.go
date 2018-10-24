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

	addrLabel := widgets.NewQLabel2("本地地址：", nil, 0)
	s.localAddr = widgets.NewQLineEdit(nil)
	s.localAddr.SetPlaceholderText("绑定本地ip地址")
	s.localAddr.SetText(s.conf.LocalAddr())
	s.localAddrMsg = NewColorLabelWithColor("不是合法的ip地址", "red")
	s.localAddrMsg.Hide()

	portLabel := widgets.NewQLabel2("本地端口：", nil, 0)
	s.localPort = widgets.NewQSpinBox(nil)
	// 端口从1024-65535
	s.localPort.SetRange(1024, 65535)
	port, _ := strconv.Atoi(s.conf.LocalPort())
	s.localPort.SetValue(port)
	s.localPortMsg = NewColorLabelWithColor("不是合法的端口值", "red")
	s.localPortMsg.Hide()

	pidFileLabel := widgets.NewQLabel2("pid-file存放路径：", nil, 0)
	s.pidFilePath = widgets.NewQLineEdit(nil)
	s.pidFilePath.SetPlaceholderText("绝对路径")
	s.pidFilePath.SetText(s.conf.PidFilePath())
	s.pidFilePathMsg = NewColorLabelWithColor("不是合法的路径", "red")
	s.pidFilePathMsg.Hide()

	// 检查内核版本
	versionInfo := widgets.NewQLabel(nil, 0)
	s.fastOpen = widgets.NewQCheckBox2("启用fast-open (Linux kernel 3.7+)", nil)
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

	groupLayout := widgets.NewQGridLayout2()
	groupLayout.AddWidget(addrLabel, 0, 0, 0)
	groupLayout.AddWidget(s.localAddr, 0, 1, 0)
	groupLayout.AddWidget(s.localAddrMsg, 1, 0, 0)
	groupLayout.AddWidget(portLabel, 2, 0, 0)
	groupLayout.AddWidget(s.localPort, 2, 1, 0)
	groupLayout.AddWidget(s.localPortMsg, 3, 0, 0)
	groupLayout.AddWidget(pidFileLabel, 4, 0, 0)
	groupLayout.AddWidget(s.pidFilePath, 4, 1, 0)
	groupLayout.AddWidget(s.pidFilePathMsg, 5, 0, 0)
	groupLayout.AddWidget(s.fastOpen, 6, 0, 0)
	groupLayout.AddWidget(versionInfo, 7, 0, 0)

	group.SetLayout(groupLayout)
	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.AddWidget(group, 0, 0)
	s.SetLayout(mainLayout)
}

// UpdateSSRClientConfig 更新config，如果数据不合法则返回error
// 因为传递了引用类型，所以直接修改config对象
func (s *SSRConfigWidget) UpdateSSRClientConfig() error {
	// 记录返回值
	var errRes error

	// 可选时才设置fast-open值
	if s.fastOpen.IsEnabled() {
		s.conf.SetFastOpen(s.fastOpen.IsChecked())
	}

	err := s.conf.SetLocalPort(strconv.Itoa(s.localPort.Value()))
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
