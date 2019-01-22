package widgets

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// ShadeWidget 半透明遮罩层
type ShadeWidget struct {
	widgets.QWidget
}

// NewShadeWidget2 返回parent的遮罩，将会遮住parent，Close后资源自动释放
// 创建后自动调用Show()
func NewShadeWidget2(parent *widgets.QWidget) *ShadeWidget {
	shade := NewShadeWidget(parent, 0)
	// 子控件设置Qt::WA_StyledBackground后才可设置背景
	shade.SetAttribute(core.Qt__WA_StyledBackground, true)
	shade.SetAttribute(core.Qt__WA_DeleteOnClose, true)
	// alpha max is 255, so 40% is 102
	shade.SetStyleSheet("background-color:rgba(0,0,0,102);")
	shade.SetWindowFlags(core.Qt__FramelessWindowHint)
	shade.SetGeometry2(0, 0, shade.ParentWidget().Width(), shade.ParentWidget().Height())
	shade.Show()

	return shade
}
