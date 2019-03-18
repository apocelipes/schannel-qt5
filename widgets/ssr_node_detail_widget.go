package widgets

import (
	"fmt"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"

	"schannel-qt5/parser"
)

// NodeDetailWidget 显示节点的详细信息
type NodeDetailWidget struct {
	widgets.QWidget

	// 节点名称
	name *widgets.QLabel
	// 代理类型
	proxyType *widgets.QLabel
	// IP和端口信息
	ip       *widgets.QLabel
	port     *widgets.QLabel
	password *widgets.QLabel
	// 加密算法
	crypt *widgets.QLabel
	// 混淆算法
	mixin *widgets.QLabel
	// 连接协议
	proto *widgets.QLabel
	// 服务器地理信息
	geo *widgets.QLabel
}

// NewNodeDetailWidgetWithNode 根据参数给出的节点显示其详细信息
func NewNodeDetailWidgetWithNode(node *parser.SSRNode) *NodeDetailWidget {
	n := NewNodeDetailWidget(nil, 0)

	n.InitUI()
	n.SetNodeDetail(node)

	return n
}

// InitUI 初始化UI
func (n *NodeDetailWidget) InitUI() {
	mainLayout := widgets.NewQGridLayout2()
	title := widgets.NewQLabel2("节点信息", nil, 0)
	titleFontSize := float64(title.FontMetrics().AverageCharWidth()) * 1.5
	title.Font().SetPixelSize(int(titleFontSize))
	titleAlign := core.Qt__AlignHCenter | core.Qt__AlignTop
	mainLayout.AddWidget3(title, 0, 0, 1, 2, titleAlign)

	nameLabel := widgets.NewQLabel2("节点名称：", nil, 0)
	n.name = widgets.NewQLabel(nil, 0)
	mainLayout.AddWidget(nameLabel, 1, 0, 0)
	mainLayout.AddWidget(n.name, 1, 1, 0)

	proxyLabel := widgets.NewQLabel2("代理类型：", nil, 0)
	n.proxyType = widgets.NewQLabel(nil, 0)
	mainLayout.AddWidget(proxyLabel, 2, 0, 0)
	mainLayout.AddWidget(n.proxyType, 2, 1, 0)

	ipLabel := widgets.NewQLabel2("IP：", nil, 0)
	n.ip = widgets.NewQLabel(nil, 0)
	mainLayout.AddWidget(ipLabel, 3, 0, 0)
	mainLayout.AddWidget(n.ip, 3, 1, 0)
	portLabel := widgets.NewQLabel2("端口：", nil, 0)
	n.port = widgets.NewQLabel(nil, 0)
	mainLayout.AddWidget(portLabel, 4, 0, 0)
	mainLayout.AddWidget(n.port, 4, 1, 0)
	passwordLabel := widgets.NewQLabel2("密码：", nil, 0)
	n.password = widgets.NewQLabel(nil, 0)
	mainLayout.AddWidget(passwordLabel, 5, 0, 0)
	mainLayout.AddWidget(n.password, 5, 1, 0)

	cryptLabel := widgets.NewQLabel2("加密算法：", nil, 0)
	n.crypt = widgets.NewQLabel(nil, 0)
	mainLayout.AddWidget(cryptLabel, 6, 0, 0)
	mainLayout.AddWidget(n.crypt, 6, 1, 0)

	mixinLabel := widgets.NewQLabel2("混淆算法：", nil, 0)
	n.mixin = widgets.NewQLabel(nil, 0)
	mainLayout.AddWidget(mixinLabel, 7, 0, 0)
	mainLayout.AddWidget(n.mixin, 7, 1, 0)

	protoLabel := widgets.NewQLabel2("连接协议：", nil, 0)
	n.proto = widgets.NewQLabel(nil, 0)
	mainLayout.AddWidget(protoLabel, 8, 0, 0)
	mainLayout.AddWidget(n.proto, 8, 1, 0)

	geoLabel := widgets.NewQLabel2("国家/地区：", nil, 0)
	n.geo = widgets.NewQLabel(nil, 0)
	mainLayout.AddWidget(geoLabel, 9, 0, 0)
	mainLayout.AddWidget(n.geo, 9, 1, 0)

	n.SetLayout(mainLayout)
}

// SetNodeDetail 设置需要显示详细信息的节点
func (n *NodeDetailWidget) SetNodeDetail(node *parser.SSRNode) {
	if node == nil {
		return
	}

	n.name.SetText(node.NodeName)
	n.proxyType.SetText(node.Type)
	n.ip.SetText(node.IP)
	port := fmt.Sprintf("%v", node.Port)
	n.port.SetText(port)
	n.password.SetText(node.Passwd)
	n.crypt.SetText(node.Crypto)
	n.mixin.SetText(node.Minx)
	n.proto.SetText(node.Proto)
	n.geo.SetText(getGeoName(node.IP))
}
