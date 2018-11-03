package widgets

import (
	"github.com/therecipe/qt/widgets"

	"schannel-qt5/parser"
)

// SummarizedWidget 综合服务信息显示，包括用户信息，服务信息
type SummarizedWidget struct {
	widgets.QWidget

	// 修改服务信息
	_ func(*parser.SSRInfo) `signal:"serviceInfoChanged,auto"`
	// 修改支付状态
	_ func(parser.PaymentState) `signal:"paymentChanged,auto"`
	// 收到数据变动
	_ func() `signal:"dataRefresh,auto"`

	// 用户数据接口
	dataBridge UserDataBridge

	// 服务信息面板
	servicePanel *ServicePanel
	// ssr开关面板
	switchPanel  *SSRSwitchPanel
	// 使用量信息
	usedPanel    *UsedPanel
	// 是否需要付款
	invoicePanel *InvoicePanel

	// 用户名-email
	user string
}
