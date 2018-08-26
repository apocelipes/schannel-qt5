package widgets

import (
  "github.com/therecipe/qt/widgets"

  "schannel-qt5/parser"
)

// ServiceInfoWidget 显示账户和服务信息
type ServiceInfoWidget struct {
  widgets.QWidget

  // 修改服务信息
  _ func(*parser.SSRInfo) `signal:"serviceInfoChanged,auto"`
  // 修改支付状态
  _ func(parser.PaymentState) `signal:"paymentChanged,auto"`
  // 收到数据变动
  _ func() `signal:"dataRefresh,auto"`

  // 用户数据接口
  dataBridge UserDataBridge

  // 服务信息
  user *widgets.QLabel
  name *widgets.QLabel
  payment *widgets.QLabel
  expireDate *widgets.QLabel
  serviceState *widgets.QLabel
  // ssr简略信息
  port *widgets.QLabel
  passwd *widgets.QLabel
  // 使用量信息
  // TODO: 自定义组件，继承QProcessorBar
  // 是否需要付款
  // TODO: 自定义组件，根据invoice来判断是否要付款，付款红字提示
}
