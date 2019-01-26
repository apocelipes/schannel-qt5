package widgets

import (
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// TransparentButton 透明按钮，当鼠标移动到按钮上时才会显示内容，否则处于透明状态
type TransparentButton struct {
	widgets.QPushButton

	_ func() `constructor:"init"`

	hoverStyle   string
	// defaultStyle 按钮不处于hover状态时的样式
	defaultStyle string
}

func (tb *TransparentButton) init() {
	tb.SetAttribute(core.Qt__WA_StyledBackground, true)
	tb.defaultStyle = "QPushButton{color:rgba(0,0,0,0);background:rgba(0,0,0,0);border:0px solid rgba(0,0,0,0);}"
	sizePolicy := tb.SizePolicy()
	sizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Fixed)
	sizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Fixed)
	tb.SetSizePolicy(sizePolicy)
	tb.SetStyleSheet(tb.defaultStyle)
}

// NewTransparentButtonWithStyle 创建带有hover style的透明按钮
// style参数的格式详见SetHoverStyle
func NewTransparentButtonWithStyle(hoverStyle string) *TransparentButton {
	button := NewTransparentButton(nil)
	button.SetHoverStyle(hoverStyle)

	return button
}

// SetHoverStyle 设置按钮处于hover状态时的qss
// style的格式须为{attr1:value1;attr2:value2; ...}
func (tb *TransparentButton) SetHoverStyle(style string) {
	if !strings.HasPrefix(style, "{") || !strings.HasSuffix(style, "}") {
		return
	}

	tb.hoverStyle = "QPushButton:hover" + style
	// setStyleSheet会覆盖所有的qss，所以拼接后添加
	tb.SetStyleSheet(tb.defaultStyle + tb.hoverStyle)
}
