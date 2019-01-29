package widgets

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/webengine"
	"github.com/therecipe/qt/widgets"
)

// InvoiceViewWidget 使用QWebEngineView展示账单详情
type InvoiceViewWidget struct {
	widgets.QMainWindow

	_ func() `constructor:"init"`

	// 移动前的位置，用于计算需要位移的距离
	lastPosition *core.QPoint
	mousePressed bool

	webEngine *webengine.QWebEngineView
	// 接受QWebEngineView事件的子对象
	webEngineChild *core.QObject

	// 窗口移动是否有效
	moveAble bool
}

func (i *InvoiceViewWidget) init() {
	i.SetAttribute(core.Qt__WA_DeleteOnClose, true)
	i.SetWindowFlag(core.Qt__FramelessWindowHint, true)

	i.webEngine = webengine.NewQWebEngineView(nil)
	i.SetCentralWidget(i.webEngine)

	closeButton := NewTransparentButtonWithStyle("{background:url(:/image/close.png);}")
	closeButton.Resize2(36, 36)
	closeButton.SetParent(i)
	closeButton.Move2(0, 0)
	closeButton.ConnectClicked(func(_ bool) {
		i.Close()
	})
	closeButton.SetToolTip("关闭窗口")
	// fix window position
	fixButton := NewTransparentButtonWithStyle("{background-color:blue;background-image:url(:/image/Removefixed.svg);}")
	fixButton.Resize2(36, 36)
	fixButton.SetParent(i)
	fixButton.Move2(0, 36)
	fixButton.ConnectClicked(func(_ bool) {
		i.moveAble = !i.moveAble
	})
	fixButton.SetToolTip("禁止/允许拖动窗口")

	// window move event
	i.webEngine.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		if event.Button() == core.Qt__LeftButton {
			i.lastPosition = event.GlobalPos()
			i.mousePressed = true
			return
		}

		i.webEngine.MousePressEventDefault(event)
	})

	i.webEngine.ConnectMouseReleaseEvent(func(event *gui.QMouseEvent) {
		i.mousePressed = false
		i.webEngine.MouseReleaseEventDefault(event)
	})

	i.webEngine.ConnectMouseMoveEvent(func(event *gui.QMouseEvent) {
		if i.mousePressed {
			movementX := event.GlobalX() - i.lastPosition.X()
			movementY := event.GlobalY() - i.lastPosition.Y()
			i.lastPosition = event.GlobalPos()
			i.Move2(i.X()+movementX, i.Y()+movementY)
			return
		}

		i.webEngine.MouseMoveEventDefault(event)
	})

	// 获取WebEngineView子对象
	i.webEngine.ConnectEvent(func(event *core.QEvent) bool {
		if event.Type() == core.QEvent__ChildPolished {
			childEvent := core.NewQChildEventFromPointer(event.Pointer())
			if childEvent.Child() != nil {
				i.webEngineChild = childEvent.Child()
				i.webEngineChild.InstallEventFilter(i.webEngine)
			}
		}

		return i.webEngine.EventDefault(event)
	})

	// 处理鼠标事件
	i.webEngine.ConnectEventFilter(func(watched *core.QObject, event *core.QEvent) bool {
		if watched.Pointer() == i.webEngineChild.Pointer() {
			// 窗口移动被禁止
			if !i.moveAble {
				return false
			}

			switch event.Type() {
			case core.QEvent__MouseButtonPress:
				mouseEvent := gui.NewQMouseEventFromPointer(event.Pointer())
				if mouseEvent.Button() == core.Qt__LeftButton {
					i.webEngine.MousePressEvent(mouseEvent)
					return true
				}
				return false
			case core.QEvent__MouseButtonRelease:
				mouseEvent := gui.NewQMouseEventFromPointer(event.Pointer())
				if mouseEvent.Button() == core.Qt__LeftButton {
					i.webEngine.MouseReleaseEvent(mouseEvent)
					return true
				}
				return false
			case core.QEvent__MouseMove:
				mouseEvent := gui.NewQMouseEventFromPointer(event.Pointer())
				i.webEngine.MouseMoveEvent(mouseEvent)
				return true
			}
		}

		return i.webEngine.EventFilterDefault(watched, event)
	})

	i.webEngine.ConnectContextMenuEvent(func(_ *gui.QContextMenuEvent) {
		// 去除自带的右键菜单
	})

	// 设置合适的宽度完整显示invoice表格内容
	i.Resize2(800, 650)
}

// 设置view展示的内容
func (i *InvoiceViewWidget) SetHTML(html string) {
	i.webEngine.SetHtml(html, core.NewQUrl())
}
