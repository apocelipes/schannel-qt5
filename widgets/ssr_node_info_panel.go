package widgets

import (
	"github.com/therecipe/qt/widgets"

	"schannel-qt5/parser"
)

// NodeInfoPanel ssr节点简略信息面板
// 显示: 名字，ip，国家(根据节点名称得出)
type NodeInfoPanel struct {
	widgets.QGroupBox

	// 处理数据更新，重新计算摘要信息
	_ func(node *parser.SSRNode) `signal:"dataRefresh,auto"`

	// 节点名字
	nodeNameLabel *widgets.QLabel
	// ip
	ipLabel *widgets.QLabel
	// 国家信息
	geoInfoLabel *widgets.QLabel

	node *parser.SSRNode
}

// NewNodeInfoPanelWithNode 根据node参数生成简略信息面板
func NewNodeInfoPanelWithNode(node *parser.SSRNode) *NodeInfoPanel {
	panel := NewNodeInfoPanel(nil)
	panel.node = node
	panel.InitUI()

	return panel
}

// InitUI 初始化UI
func (n *NodeInfoPanel) InitUI() {
	n.SetTitle("ssr摘要:")

	geoInfo := getGeoName(n.node.IP)
	n.geoInfoLabel = widgets.NewQLabel2(geoInfo, nil, 0)
	n.ipLabel = widgets.NewQLabel2(n.node.IP, nil, 0)
	n.nodeNameLabel = widgets.NewQLabel2(n.node.NodeName, nil, 0)

	mainLayout := widgets.NewQGridLayout2()
	name := widgets.NewQLabel2("节点名称:", nil, 0)
	mainLayout.AddWidget(name, 0, 0, 0)
	mainLayout.AddWidget(n.nodeNameLabel, 0, 1, 0)
	ip := widgets.NewQLabel2("IP:", nil, 0)
	mainLayout.AddWidget(ip, 1, 0, 0)
	mainLayout.AddWidget(n.ipLabel, 1, 1, 0)
	geo := widgets.NewQLabel2("地区/国家:", nil, 0)
	mainLayout.AddWidget(geo, 2, 0, 0)
	mainLayout.AddWidget(n.geoInfoLabel, 2, 1, 0)

	n.SetLayout(mainLayout)
}

// dataRefresh 刷新节点摘要信息
func (n *NodeInfoPanel) dataRefresh(node *parser.SSRNode) {
	n.node = node
	geoInfo := getGeoName(n.node.IP)
	n.nodeNameLabel.SetText(n.node.NodeName)
	n.ipLabel.SetText(n.node.IP)
	n.geoInfoLabel.SetText(geoInfo)
}
