package widgets

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// AccountItem 显示账户名和删除按钮，由QComboBox使用
type AccountItem struct {
	widgets.QWidget

	// 选中时显示账户名
	_ func(string) `signal:"showAccount"`
	// 删除所选用户
	_ func(string) `signal:"removeAccount"`

	accountName *widgets.QLabel
	delButton   *widgets.QPushButton
	userName    string
	// 用户是否进行左键单击
	mousePress bool
}

func NewAccountItem2(user string) *AccountItem {
	item := NewAccountItem(nil, 0)
	item.userName = user
	item.InitUI()
	return item
}

func (item *AccountItem) InitUI() {
	item.accountName = widgets.NewQLabel2(item.userName, nil, 0)
	delIcon := widgets.QApplication_Style().StandardIcon(widgets.QStyle__SP_DialogCloseButton, nil, nil)
	item.delButton = widgets.NewQPushButton3(delIcon, "", nil)
	item.delButton.SetStyleSheet("background:transparent;")
	// 设置icon大小
	iconHeight := item.accountName.FontMetrics().Height()
	iconWidth := iconHeight * 2
	item.delButton.SetIconSize(core.NewQSize2(iconWidth, iconHeight))
	item.delButton.ConnectClicked(func(_ bool) {
		item.RemoveAccount(item.userName)
	})

	mainLayout := widgets.NewQHBoxLayout()
	mainLayout.AddWidget(item.accountName, 0, core.Qt__AlignLeft)
	mainLayout.AddStretch(0)
	mainLayout.AddWidget(item.delButton, 0, 0)
	mainLayout.SetSpacing(5)
	mainLayout.SetContentsMargins(5, 5, 5, 5)
	item.SetLayout(mainLayout)

	// 处理click
	item.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		if event.Button() == core.Qt__LeftButton {
			item.mousePress = true
		}
	})
	item.ConnectMouseReleaseEvent(func(_ *gui.QMouseEvent) {
		if item.mousePress {
			item.ShowAccount(item.userName)
			item.mousePress = false
		}
	})
}
