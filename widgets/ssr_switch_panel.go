package widgets

import (
	"fmt"
	"sort"
	"time"

	"github.com/therecipe/qt/widgets"

	"schannel-qt5/config"
	"schannel-qt5/parser"
	_ "schannel-qt5/pyclient"
	"schannel-qt5/ssr"
)

// SSRSwitchPanel 显示简略信息，开关ssr
type SSRSwitchPanel struct {
	widgets.QWidget

	// ssr client是否启动
	ssrStat *ColorLabel

	// node缩略信息
	nodeInfo *NodeInfoPanel

	// 链接是否可用
	// 未打开客户端为"未开启"(gray)
	// 链接可用为"OK"(green)
	// 不可用为"error: [error info]"(red)
	connStat *ColorLabel

	// ssr开关
	switchButton *widgets.QPushButton
	// 选择节点对话框按钮
	selectNodeButton *widgets.QPushButton

	// 当前使用的节点
	currentNode *parser.SSRNode
	// ssr client程序和配置文件
	ssrClient ssr.Launcher
	conf      *config.UserConfig
	// 可用节点信息
	nodes []*parser.SSRNode
}

// NewSSRSwitchPanel2 创建ssr开关面板组件
func NewSSRSwitchPanel2(conf *config.UserConfig, nodes []*parser.SSRNode) *SSRSwitchPanel {
	if conf == nil || nodes == nil {
		return nil
	}
	panel := NewSSRSwitchPanel(nil, 0)

	panel.conf = conf
	panel.nodes = make([]*parser.SSRNode, len(nodes))
	copy(panel.nodes, nodes)
	panel.SortNode()

	panel.ssrClient = ssr.NewLauncher("python", panel.conf)
	if panel.ssrClient == nil {
		return nil
	}

	panel.currentNode = &parser.SSRNode{}
	nodePath, err := panel.conf.SSRNodeConfigPath.AbsPath()
	if err != nil {
		return nil
	}
	panel.currentNode.Load(nodePath)

	panel.InitUI()
	return panel
}

// SortNode 将节点排序，方便查找节点信息
// 按照NodeName排序
func (s *SSRSwitchPanel) SortNode() {
	sort.Slice(s.nodes, func(i, j int) bool {
		return s.nodes[i].NodeName < s.nodes[j].NodeName
	})
}

// InitUI 初始化界面
func (s *SSRSwitchPanel) InitUI() {
	group := widgets.NewQGroupBox2("ssr开关", nil)
	componentLayout := widgets.NewQGridLayout2()

	s.nodeInfo = NewNodeInfoPanelWithNode(s.currentNode)
	componentLayout.AddWidget3(s.nodeInfo, 0, 0, 1, 3, 0)

	ssrStatLabel := widgets.NewQLabel2("ssr状态:", nil, 0)
	s.ssrStat = NewColorLabelWithColor("", "")
	s.setSSRStat()
	componentLayout.AddWidget(ssrStatLabel, 1, 0, 0)
	componentLayout.AddWidget3(s.ssrStat, 1, 1, 1, 2, 0)

	s.connStat = NewColorLabelWithColor("", "")
	// 设置自动换行
	s.connStat.AdjustSize()
	s.connStat.SetWordWrap(true)
	s.setConnStat()
	connStatLabel := widgets.NewQLabel2("连接状态:", nil, 0)
	componentLayout.AddWidget(connStatLabel, 2, 0, 0)
	componentLayout.AddWidget3(s.connStat, 2, 1, 1, 2, 0)

	s.selectNodeButton = widgets.NewQPushButton2("选择节点", nil)
	s.selectNodeButton.ConnectClicked(func(_ bool) {
		dialog := NewNodeSelectDialog2(s.currentNode, s.nodes)
		if dialog.Exec() == int(widgets.QDialog__Accepted) {
			s.currentNode = dialog.CurrentNode
			s.nodeInfo.DataRefresh(s.currentNode)
		}
	})
	s.switchButton = widgets.NewQPushButton(nil)
	s.setSwitchLabel()
	s.switchButton.ConnectClicked(func(_ bool) {
		switch s.switchButton.Text() {
		case "打开":
			if err := s.ssrClient.Start(); err != nil {
				errMsg := widgets.NewQErrorMessage(nil)
				errMsg.ShowMessage(fmt.Sprintf("启动客户端错误: %v", err))
				errMsg.Exec()
				return
			}
		case "关闭":
			if err := s.ssrClient.Stop(); err != nil {
				errMsg := widgets.NewQErrorMessage(nil)
				errMsg.ShowMessage(fmt.Sprintf("关闭客户端错误: %v", err))
				errMsg.Exec()
				return
			}
		}

		s.setSSRStat()
		s.setConnStat()
		s.setSwitchLabel()
	})
	componentLayout.AddWidget3(s.selectNodeButton, 3, 1, 1, 1, 0)
	componentLayout.AddWidget3(s.switchButton, 3, 2, 1, 1, 0)

	group.SetLayout(componentLayout)
	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.AddWidget(group, 0, 0)
	s.SetLayout(mainLayout)
}

// setSSRStat 设置ssr客户端是否正在运行的状态信息
func (s *SSRSwitchPanel) setSSRStat() {
	if err := s.ssrClient.IsRunning(); err != nil {
		s.ssrStat.SetColorText("未运行", "red")
		return
	}

	s.ssrStat.SetColorText("正在运行", "green")
}

// setConnStat 设置代理节点是否可用的信息
func (s *SSRSwitchPanel) setConnStat() {
	if err := s.ssrClient.IsRunning(); err != nil {
		s.connStat.SetColorText("未开启客户端", "gray")
		return
	}

	if err := s.ssrClient.ConnectionCheck(5 * time.Second); err != nil {
		errInfo := fmt.Sprintf("error: %v", err)
		s.connStat.SetColorText(errInfo, "red")
		return
	}

	s.connStat.SetColorText("OK", "green")
}

// setSwitchLabel 设置开关按钮的label
func (s *SSRSwitchPanel) setSwitchLabel() {
	if err := s.ssrClient.IsRunning(); err != nil {
		s.switchButton.SetText("打开")
		return
	}

	s.switchButton.SetText("关闭")
}

// DataRefresh 更新config和nodes
func (s *SSRSwitchPanel) DataRefresh(conf *config.UserConfig, nodes []*parser.SSRNode) {
	// 停止旧的客户端运行
	if running := s.ssrClient.IsRunning(); running == nil {
		s.ssrClient.Stop()
		s.switchButton.SetText("打开")
	}
	s.conf = conf
	s.ssrClient = ssr.NewLauncher("python", s.conf)
	if s.ssrClient == nil {
		errMsg := widgets.NewQErrorMessage(nil)
		// TODO 更详细的错误信息
		errMsg.ShowMessage("初始化ssr客户端错误")
		errMsg.Show()
		return
	}

	s.nodes = make([]*parser.SSRNode, len(nodes))
	copy(s.nodes, nodes)
	s.SortNode()

	//TODO 检测当前节点不在节点列表的情况
	s.currentNode = &parser.SSRNode{}
	s.currentNode.Load(s.conf.SSRNodeConfigPath.String())
	s.nodeInfo.DataRefresh(s.currentNode)
	s.setConnStat()
	s.setSSRStat()
}
