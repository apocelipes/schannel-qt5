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

const (
	// 留给伸缩因子的空间
	spaceStreth = 40
	// 边界宽度
	leftBorder   = 5
	rightBorder  = 5
	topBorder    = 5
	bottomBorder = 5
)

func (item *AccountItem) InitUI() {
	item.accountName = widgets.NewQLabel2(item.userName, nil, 0)
	// account设置为button的3倍
	accountSizePolicy := item.accountName.SizePolicy()
	accountSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Expanding)
	accountSizePolicy.SetHorizontalStretch(3)
	item.accountName.SetSizePolicy(accountSizePolicy)

	delIcon := widgets.QApplication_Style().StandardIcon(widgets.QStyle__SP_DialogCloseButton, nil, nil)
	item.delButton = widgets.NewQPushButton3(delIcon, "", nil)
	item.delButton.SetStyleSheet("background:transparent;")
	// 设置icon大小
	iconHeight := item.accountName.FontMetrics().Height()
	iconWidth := iconHeight
	item.delButton.SetIconSize(core.NewQSize2(iconWidth, iconHeight))
	btnSizePolicy := item.delButton.SizePolicy()
	btnSizePolicy.SetHorizontalStretch(1)
	item.delButton.SetSizePolicy(btnSizePolicy)
	item.delButton.ConnectClicked(func(_ bool) {
		item.RemoveAccount(item.userName)
	})

	mainLayout := widgets.NewQHBoxLayout()
	mainLayout.AddWidget(item.accountName, 0, core.Qt__AlignLeft)
	mainLayout.AddStretch(0)
	mainLayout.AddWidget(item.delButton, 0, 0)
	mainLayout.SetSpacing(5)
	mainLayout.SetContentsMargins(leftBorder, topBorder, rightBorder, bottomBorder)
	item.SetLayout(mainLayout)

	// 处理click
	item.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		if event.Button() == core.Qt__LeftButton {
			// 记录左键按下
			item.mousePress = true
		}
	})
	item.ConnectMouseReleaseEvent(func(_ *gui.QMouseEvent) {
		if item.mousePress {
			item.ShowAccount(item.userName)
			item.mousePress = false
		}
	})

	// 设置大小
	item.ConnectSizeHint(func() *core.QSize {
		pointSize := item.accountName.Font().PointSize()
		textLength := len(item.userName) * pointSize
		// left + text + spacing + stretch + button + right
		width := leftBorder + textLength + 5*2 + spaceStreth + textLength/3 + rightBorder
		// top + text + bottom
		height := topBorder + item.accountName.FontMetrics().Height() + bottomBorder
		return core.NewQSize2(width, height)
	})

	// 设置自身大小策略
	sizePolicy := item.SizePolicy()
	sizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__MinimumExpanding)
}
