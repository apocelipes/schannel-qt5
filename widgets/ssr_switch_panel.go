package widgets

import (
	"fmt"
	"log"
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

	// node缩略信息
	nodeInfo *NodeInfoPanel

	// 链接是否可用
	// 未打开客户端为"未开启"(gray)
	// 链接可用为"OK"(green)
	// 不可用为"error: [error info]"(red)
	connStat *ColorLabel

	// ssr开关
	switchButton *SwitchButton
	// 选择节点对话框按钮
	selectNodeButton *widgets.QPushButton

	// 当前使用的节点
	currentNode *parser.SSRNode
	// ssr client程序和配置文件
	ssrClient ssr.Launcher
	conf      *config.UserConfig
	// 可用节点信息
	nodes []*parser.SSRNode

	logger *log.Logger
}

// NewSSRSwitchPanel2 创建ssr开关面板组件
func NewSSRSwitchPanel2(conf *config.UserConfig, nodes []*parser.SSRNode, logger *log.Logger) *SSRSwitchPanel {
	if conf == nil || nodes == nil {
		return nil
	}
	panel := NewSSRSwitchPanel(nil, 0)

	panel.conf = conf
	panel.nodes = make([]*parser.SSRNode, len(nodes))
	copy(panel.nodes, nodes)
	panel.SortNode()
	panel.logger = logger

	panel.ssrClient = ssr.NewLauncher("python", panel.conf)
	if panel.ssrClient == nil {
		panel.logger.Println("ssr client created failed")
		return nil
	}

	panel.currentNode = &parser.SSRNode{}
	nodePath, err := panel.conf.SSRNodeConfigPath.AbsPath()
	if err != nil {
		panel.logger.Println("node load failed: ", err)
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
	componentLayout := widgets.NewQGridLayout2()

	s.nodeInfo = NewNodeInfoPanelWithNode(s.currentNode)
	componentLayout.AddWidget3(s.nodeInfo, 0, 0, 1, 3, 0)

	s.connStat = NewColorLabelWithColor("", "")
	// 设置自动换行
	s.connStat.AdjustSize()
	s.connStat.SetWordWrap(true)
	s.setConnStat()
	connStatLabel := widgets.NewQLabel2("连接状态:", nil, 0)
	componentLayout.AddWidget(connStatLabel, 1, 0, 0)
	componentLayout.AddWidget3(s.connStat, 1, 1, 1, 2, 0)

	s.switchButton = NewSwitchButton2(s.ssrClient.IsRunning() == nil)
	s.switchButton.ConnectClicked(func(checked bool) {
		info := ""
		switch checked {
		case true:
			if err := s.ssrClient.Start(); err != nil {
				errInfo := fmt.Sprintf("启动客户端错误: %v", err)
				showErrorDialog(errInfo, s)
				return
			}

			info = "已打开"
		case false:
			if err := s.ssrClient.Stop(); err != nil {
				errInfo := fmt.Sprintf("关闭客户端错误: %v", err)
				showErrorDialog(errInfo, s)
				return
			}

			info = "已关闭"
		}

		ShowNotification("SSR客户端", info, "", -1)
		s.setConnStat()
	})
	switchLabel := widgets.NewQLabel2("ssr开关：", nil, 0)
	componentLayout.AddWidget(switchLabel, 2, 0, 0)
	componentLayout.AddWidget3(s.switchButton, 2, 1, 1, 2, 0)

	s.selectNodeButton = widgets.NewQPushButton2("选择节点", nil)
	s.selectNodeButton.ConnectClicked(func(_ bool) {
		dialog := NewNodeSelectDialog2(s.currentNode, s.nodes)
		shade := NewShadeWidget2(s.NativeParentWidget())
		if dialog.Exec() == int(widgets.QDialog__Accepted) {
			s.currentNode = dialog.CurrentNode
			s.nodeInfo.DataRefresh(s.currentNode)
		}
		shade.Close()
		// goqt无法自动释放QWidget
		// 且此处不适合DeleteOnClose，所以需要手动调用DestroyNodeSelectDialog
		dialog.DestroyNodeSelectDialog()
	})
	componentLayout.AddWidget3(s.selectNodeButton, 3, 2, 1, 1, 0)

	s.SetLayout(componentLayout)
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
		s.logger.Println(errInfo)
		ShowNotification("SSR连接测试失败", errInfo, "", -1)
		return
	}

	s.connStat.SetColorText("OK", "green")
}

// DataRefresh 更新config和nodes
func (s *SSRSwitchPanel) DataRefresh(conf *config.UserConfig, nodes []*parser.SSRNode) {
	// 停止旧的客户端运行
	if running := s.ssrClient.IsRunning(); running == nil {
		s.ssrClient.Stop()
		s.switchButton.SetChecked(false)
		ShowNotification("SSR客户端", "已关闭", "", -1)
	}
	s.conf = conf
	s.ssrClient = ssr.NewLauncher("python", s.conf)
	if s.ssrClient == nil {
		s.logger.Println("ssr switch DataRefresh: 初始化ssr客户端错误")
		// TODO 更详细的错误信息
		showErrorDialog("初始化ssr客户端错误", s)
		return
	}

	s.nodes = make([]*parser.SSRNode, len(nodes))
	copy(s.nodes, nodes)
	s.SortNode()

	//TODO 检测当前节点不在节点列表的情况
	s.currentNode = &parser.SSRNode{}
	nodeConfigPath, err := s.conf.SSRNodeConfigPath.AbsPath()
	if err != nil {
		// 节点配置获取失败，错误信息用信息框显示
		showErrorDialog(err.Error(), s)
		return
	}
	s.currentNode.Load(nodeConfigPath)
	s.nodeInfo.DataRefresh(s.currentNode)
	s.setConnStat()
}
