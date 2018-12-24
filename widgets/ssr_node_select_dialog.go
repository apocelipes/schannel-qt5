package widgets

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/therecipe/qt/core"
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
	//list *widgets.QListWidget
	tree      *widgets.QTreeView
	nodeModel *NodeTreeModel
	// node详细信息
	detail *NodeDetailWidget

	// 选择的节点，将作为结果被使用
	CurrentNode *parser.SSRNode
}

// NewNodeSelectDialog2 生成node选择对话框
func NewNodeSelectDialog2(current *parser.SSRNode, nodes []*parser.SSRNode) *NodeSelectDialog {
	if nodes == nil {
		return nil
	}

	dialog := NewNodeSelectDialog(nil, 0)
	dialog.CurrentNode = current
	// dialog运行于模态，nodes不会被修改
	dialog.nodeModel = NewNodeTreeModel2(nodes)
	dialog.InitUI()

	return dialog
}

// InitUI 初始化界面
func (dialog *NodeSelectDialog) InitUI() {
	dialog.detail = NewNodeDetailWidgetWithNode(dialog.CurrentNode)

	dialog.tree = widgets.NewQTreeView(nil)
	dialog.tree.SetModel(dialog.nodeModel)
	dialog.tree.SetAnimated(true)
	// 设置选择当前节点
	currentIndex := dialog.nodeModel.FindNodeIndex(dialog.CurrentNode)
	dialog.tree.SetCurrentIndex(currentIndex)
	dialog.tree.Expand(currentIndex)
	dialog.tree.ConnectClicked(func(index *core.QModelIndex) {
		item := NewNodeTreeItemFromPointer(index.InternalPointer())
		if item.ChildCount() == 0 {
			dialog.detail.SetNodeDetail(item.Node())
			dialog.CurrentNode = item.Node()
		} else {
			nodeItem := item.LatestChild(0)
			dialog.detail.SetNodeDetail(nodeItem.Node())
			dialog.CurrentNode = nodeItem.Node()
		}
	})
	dialog.tree.Clicked(currentIndex)

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
	contentLayout := widgets.NewQHBoxLayout()
	contentLayout.AddWidget(dialog.tree, 1, 0)
	contentLayout.AddWidget(dialog.detail, 2, 0)
	mainLayout.AddLayout2(contentLayout, 0, 0, 1, 4, 0)
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

// saveNode 保存节点信息至文件
func (dialog *NodeSelectDialog) saveNode(_ bool) {
	jsonFileFilter := "JSON Files(*.json)"
	nodeFileName := fmt.Sprintf("%s.json", dialog.CurrentNode.NodeName)
	savePath, err := getFileSavePath("node", nodeFileName, jsonFileFilter, dialog)
	if err == ErrCanceled {
		return
	} else if err != nil {
		showErrorDialog("保存路径获取失败："+err.Error(), dialog)
		return
	}

	data, err := json.MarshalIndent(dialog.CurrentNode, "", "\t")
	if err != nil {
		showErrorDialog("配置解析失败："+err.Error(), dialog)
		return
	}
	f, err := os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		showErrorDialog("文件创建失败："+err.Error(), dialog)
		return
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		showErrorDialog("写入配置失败："+err.Error(), dialog)
		return
	}

	ShowNotification("节点", savePath+"保存成功", "", -1)
}
