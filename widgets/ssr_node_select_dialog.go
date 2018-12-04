package widgets

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/therecipe/qt/widgets"

	"schannel-qt5/parser"
)

// NodeSelectDialog 显示所有节点信息，并选择设置节点
// 需要使用模态运行
type NodeSelectDialog struct {
	widgets.QDialog

	// dialog功能按钮
	okButton, cancelButton *widgets.QPushButton
	// 节点列表
	list *widgets.QListWidget
	// node详细信息
	detail *NodeDetailWidget

	// 选择的节点，将作为结果被使用
	CurrentNode *parser.SSRNode
	// 全部可用节点，已排序
	nodes []*parser.SSRNode
}

// NewNodeSelectDialog2 生成node选择对话框
func NewNodeSelectDialog2(current *parser.SSRNode, nodes []*parser.SSRNode) *NodeSelectDialog {
	if nodes == nil {
		return nil
	}

	dialog := NewNodeSelectDialog(nil, 0)
	dialog.CurrentNode = current
	// dialog运行于模态，nodes不会被修改
	dialog.nodes = nodes
	dialog.InitUI()

	return dialog
}

// InitUI 初始化界面
func (dialog *NodeSelectDialog) InitUI() {
	dialog.detail = NewNodeDetailWidgetWithNode(dialog.CurrentNode)
	dialog.list = widgets.NewQListWidget(nil)
	dialog.list.AddItems(dialog.getNodeNames())
	// 绑定选择事件
	dialog.list.ConnectCurrentRowChanged(func(index int) {
		// 因为list和nodes顺序一致，所以可以直接找到当前选中节点
		if index == -1 {
			return
		}
		dialog.CurrentNode = dialog.nodes[index]
		dialog.detail.SetNodeDetail(dialog.CurrentNode)
	})
	// 设置当前选择的节点
	// 如果节点不在当前列表中，则默认选择index 0
	if dialog.CurrentNode != nil {
		for i := range dialog.nodes {
			if dialog.nodes[i].NodeName == dialog.CurrentNode.NodeName {
				dialog.list.SetCurrentRow(i)
				break
			}
		}
	}

	listLayout := widgets.NewQVBoxLayout()
	listLabel := widgets.NewQLabel2("可用节点列表：", nil, 0)
	listLayout.AddWidget(listLabel, 0, 0)
	listLayout.AddWidget(dialog.list, 0, 0)

	dialog.okButton = widgets.NewQPushButton2("选择", nil)
	dialog.okButton.ConnectClicked(func(_ bool) {
		dialog.Accept()
	})
	dialog.cancelButton = widgets.NewQPushButton2("取消", nil)
	dialog.cancelButton.ConnectClicked(func(_ bool) {
		dialog.Reject()
	})
	saveNodeButton := widgets.NewQPushButton2("保存至文件", nil)
	saveNodeButton.ConnectClicked(dialog.saveNode)

	mainLayout := widgets.NewQGridLayout2()
	mainLayout.AddLayout(listLayout, 0, 0, 0)
	mainLayout.AddWidget3(dialog.detail, 0, 1, 1, 3, 0)
	// 水平分割线
	hFrame := widgets.NewQFrame(nil, 0)
	hFrame.SetFrameStyle(int(widgets.QFrame__HLine) | int(widgets.QFrame__Sunken))
	mainLayout.AddWidget3(hFrame, 1, 0, 1, 4, 0)
	mainLayout.AddWidget(saveNodeButton, 2, 1, 0)
	mainLayout.AddWidget(dialog.cancelButton, 2, 2, 0)
	mainLayout.AddWidget(dialog.okButton, 2, 3, 0)
	dialog.SetLayout(mainLayout)
	dialog.SetWindowTitle("选择节点")
}

// getNodeNames 返回所有节点的名称，顺序与nodes一致
func (dialog *NodeSelectDialog) getNodeNames() []string {
	names := make([]string, 0, 10)
	for _, v := range dialog.nodes {
		names = append(names, v.NodeName)
	}

	return names
}

// saveNode 保存节点信息至文件
func (dialog *NodeSelectDialog) saveNode(_ bool) {
	home := os.Getenv("HOME")
	if home == "" {
		showErrorDialog("未找到$HOME，无法保存", dialog)
		return
	}
	jsonFileFilter := "JSON Files(*.json)"
	nodeFileName := fmt.Sprintf("%s.json", dialog.CurrentNode.NodeName)
	defaultSavePath := filepath.Join(home, nodeFileName)

	savePath := widgets.QFileDialog_GetSaveFileName(dialog,
		"保存",
		defaultSavePath,
		jsonFileFilter,
		"",
		0)
	if savePath == "" {
		return
	}

	data, err := json.MarshalIndent(dialog.CurrentNode, "", "\t")
	if err != nil {
		showErrorDialog("配置解析失败：" + err.Error(), dialog)
		return
	}
	f, err := os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		showErrorDialog("文件创建失败：" + err.Error(), dialog)
		return
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		showErrorDialog("写入配置失败：" + err.Error(), dialog)
		return
	}

	infoButton := widgets.QMessageBox__Yes
	widgets.QMessageBox_Information4(dialog,
		"保存成功",
		savePath+"保存成功",
		infoButton,
		infoButton)
}
