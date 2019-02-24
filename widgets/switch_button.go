package widgets

import (
	"math"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// 开关按钮，带有动画过度效果
type SwitchButton struct {
	widgets.QWidget

	_ func(bool) `signal:"clicked"`

	// 动画的起始和结束值
	startValue float64
	endValue   float64
	animation  *core.QVariantAnimation

	checked bool

	uncheckedColor *gui.QColor
	checkedColor   *gui.QColor
	indicatorColor *gui.QColor
}

// 生成一个带有默认状态的SwitchButton
func NewSwitchButton2(checked bool) *SwitchButton {
	button := NewSwitchButton(nil, 0)
	button.checked = checked
	button.InitUI()

	return button
}

func (s *SwitchButton) InitUI() {
	s.startValue = 0.0
	s.endValue = 1.0
	s.animation = core.NewQVariantAnimation(s)
	s.animation.ConnectValueChanged(func(_ *core.QVariant) {
		s.Update()
	})

	s.checkedColor = gui.NewQColor3(0x00, 0xee, 0x00, 255)
	s.uncheckedColor = gui.NewQColor3(0xee, 0xe9, 0xe9, 255)
	s.indicatorColor = gui.NewQColor2(core.Qt__white)

	sizePolicy := s.SizePolicy()
	sizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Fixed)
	sizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Fixed)
	s.SetSizePolicy(sizePolicy)
	s.ConnectSizeHint(func() *core.QSize {
		return core.NewQSize2(60, 30)
	})

	s.ConnectPaintEvent(s.paintEvent)
	// 模拟按钮点击信号
	s.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		if event.Button() == core.Qt__LeftButton {
			s.SetChecked(!s.IsChecked())
			// 触发clicked信号，checked已经更新
			s.Clicked(s.IsChecked())
			event.Accept()
			return
		}

		s.QWidget.MousePressEventDefault(event)
	})
}

// 设置SwitchButton的点击状态并启动动画效果
func (s *SwitchButton) SetChecked(checked bool) {
	if s.IsChecked() == checked {
		return
	}

	s.checked = checked
	// checked改变为true，则按钮动画从左向右移动；改变为false则从右向左移动
	if checked {
		s.animation.SetStartValue(core.NewQVariant12(s.startValue))
		s.animation.SetEndValue(core.NewQVariant12(s.endValue))
	} else {
		s.animation.SetStartValue(core.NewQVariant12(s.endValue))
		s.animation.SetEndValue(core.NewQVariant12(s.startValue))
	}

	s.animation.Start(core.QAbstractAnimation__KeepWhenStopped)
}

func (s *SwitchButton) paintEvent(_ *gui.QPaintEvent) {
	painter := gui.NewQPainter2(s)
	painter.SetRenderHints(gui.QPainter__Antialiasing|gui.QPainter__SmoothPixmapTransform, true)
	var background *gui.QColor
	if s.IsChecked() {
		background = s.checkedColor
	} else {
		background = s.uncheckedColor
	}

	border := 1.0
	heightF := float64(s.Height())
	widthF := float64(s.Width())

	backgroundPath := gui.NewQPainterPath()
	// 背景色条区域，起点：x为0+边框宽度，y为高度减去上下边框后的上部1/4处
	// 宽为宽度减去左右边框，高为高度减去上下边框后的1/2
	backgroundRect := core.NewQRectF4(border,
		(heightF-2*border)/4+border,
		widthF-2*border,
		(heightF-2*border)/2,
	)
	backgroundPath.AddRoundedRect(backgroundRect,
		(heightF-2*border)/4,
		(heightF-2*border)/4,
		core.Qt__AbsoluteSize,
	)
	backgroundPath.CloseSubpath()

	// 获得当前动画产生的移动距离的比例
	var moveRatio float64
	if s.animation.State() == core.QAbstractAnimation__Stopped {
		if !s.IsChecked() {
			moveRatio = s.startValue
		} else {
			moveRatio = s.endValue
		}
	} else {
		valid := false // 无实际用途，只是为了正常调用ToDouble
		moveRatio = s.animation.CurrentValue().ToDouble(&valid)
	}

	// 按钮可移动距离，宽减去高（按钮大小和高一致，不包括边框）
	moveWidth := (widthF - 2*border) - (heightF - 2*border)
	indicatorSize := heightF - 2*border
	indicatorPath := gui.NewQPainterPath()
	// 按钮绘制的起点为边框+可移动距离*当前需要移动的比例
	indicatorRect := core.NewQRectF4(moveWidth*moveRatio+border,
		border,
		indicatorSize,
		indicatorSize,
	)
	indicatorPath.AddEllipse(indicatorRect)
	indicatorPath.CloseSubpath()

	indicatorCenterPath := gui.NewQPainterPath()
	// 位于按钮正中的小圆点宽高均为按钮1/2也就是整体控件高度的1/4
	// 为了突出白色按钮的显示效果而添加，颜色与背景色相同
	indicatorCenterRect := core.NewQRectF4(
		moveWidth*moveRatio+heightF*3.0/8,
		heightF*3/8,
		heightF/4,
		heightF/4,
	)
	indicatorCenterPath.AddEllipse(indicatorCenterRect)
	indicatorCenterPath.CloseSubpath()

	// 位于背景色条的宽度25%/75%处的白色分隔符
	// checked为false时位于宽度75%处
	// 上下距离背景色条的距离均为3/8的height，高度为height的1/4暨背景色条的1/2
	sepPath := gui.NewQPainterPath()
	sepRect := core.NewQRectF4(widthF*math.Abs(moveRatio-3.0/4),
		heightF*3/8,
		widthF/20,
		heightF/4,
	)
	sepPath.AddRect(sepRect)
	sepPath.CloseSubpath()

	// 背景色条的边框
	backgroundBorderRect := core.NewQRectF4(0,
		heightF/4,
		widthF,
		heightF/2,
	)
	painter.SetPen2(gui.NewQColor2(core.Qt__lightGray))
	painter.DrawRoundedRect(backgroundBorderRect,
		heightF/4,
		heightF/4,
		core.Qt__AbsoluteSize,
	)

	painter.FillPath(backgroundPath, gui.NewQBrush3(background, core.Qt__SolidPattern))
	painter.FillPath(sepPath, gui.NewQBrush3(s.indicatorColor, core.Qt__SolidPattern))
	painter.FillPath(indicatorPath, gui.NewQBrush3(s.indicatorColor, core.Qt__SolidPattern))
	painter.FillPath(indicatorCenterPath, gui.NewQBrush3(background, core.Qt__SolidPattern))

	// 按钮的边框
	indicatorBorderRect := core.NewQRectF4((widthF-heightF)*moveRatio,
		0,
		heightF,
		heightF,
	)
	painter.DrawEllipse(indicatorBorderRect)

	// 因为golang没有析构函数，需要手动调用End
	painter.End()
}

// 返回按钮是否处于打开状态
func (s *SwitchButton) IsChecked() bool {
	return s.checked
}
