package widgets

import (
	"fmt"

	"github.com/therecipe/qt/widgets"

	"schannel-qt5/parser"
)

// ServicePanel 显示服务信息
type ServicePanel struct {
	widgets.QWidget

	// 服务信息
	serviceName  *widgets.QLabel
	user         *widgets.QLabel
	port         *widgets.QLabel
	password     *widgets.QLabel
	payment      *widgets.QLabel
	expireDate   *widgets.QLabel
	serviceState *widgets.QLabel
}

// NewServicePanel2 根据service生成服务信息面板
func NewServicePanel2(user string, info *parser.SSRInfo) *ServicePanel {
	if info == nil {
		return nil
	}
	panel := NewServicePanel(nil, 0)
	panel.InitUI(user, info)

	return panel
}

// InitUI 初始化UI组件
func (panel *ServicePanel) InitUI(user string, info *parser.SSRInfo) {
	group := widgets.NewQGroupBox2("服务信息：", nil)

	serviceNameLabel := widgets.NewQLabel2("服务名称：", nil, 0)
	panel.serviceName = widgets.NewQLabel(nil, 0)
	userLabel := widgets.NewQLabel2("用户：", nil, 0)
	panel.user = widgets.NewQLabel(nil, 0)
	portLabel := widgets.NewQLabel2("端口：", nil, 0)
	panel.port = widgets.NewQLabel(nil, 0)
	passwordLabel := widgets.NewQLabel2("密码：", nil, 0)
	panel.password = widgets.NewQLabel(nil, 0)
	paymentLabel := widgets.NewQLabel2("费用：", nil, 0)
	panel.payment = widgets.NewQLabel(nil, 0)
	expireLabel := widgets.NewQLabel2("过期时间：", nil, 0)
	panel.expireDate = widgets.NewQLabel(nil, 0)
	serviceStateLabel := widgets.NewQLabel2("服务状态：", nil, 0)
	panel.serviceState = widgets.NewQLabel(nil, 0)

	infoLayout := widgets.NewQGridLayout2()
	infoLayout.AddWidget(serviceNameLabel, 0, 0, 0)
	infoLayout.AddWidget(panel.serviceName, 0, 1, 0)
	infoLayout.AddWidget(userLabel, 1, 0, 0)
	infoLayout.AddWidget(panel.user, 1, 1, 0)
	infoLayout.AddWidget(portLabel, 2, 0, 0)
	infoLayout.AddWidget(panel.port, 2, 1, 0)
	infoLayout.AddWidget(passwordLabel, 3, 0, 0)
	infoLayout.AddWidget(panel.password, 3, 1, 0)
	infoLayout.AddWidget(paymentLabel, 4, 0, 0)
	infoLayout.AddWidget(panel.payment, 4, 1, 0)
	infoLayout.AddWidget(expireLabel, 5, 0, 0)
	infoLayout.AddWidget(panel.expireDate, 5, 1, 0)
	infoLayout.AddWidget(serviceStateLabel, 6, 0, 0)
	infoLayout.AddWidget(panel.serviceState, 6, 1, 0)

	group.SetLayout(infoLayout)
	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.AddWidget(group, 0, 0)
	panel.SetLayout(mainLayout)

	// 向panel填充信息
	panel.UpadteInfo(user, info)
}

// UpdateInfo 更新panel信息
func (panel *ServicePanel) UpadteInfo(user string, info *parser.SSRInfo) {
	panel.user.SetText(user)
	panel.serviceName.SetText(info.Name)
	panel.port.SetText(fmt.Sprint(info.Port))
	panel.password.SetText(info.Passwd)
	panel.payment.SetText(info.Price)
	panel.expireDate.SetText(time2string(info.Expires))
	panel.serviceState.SetText(info.State)
}
