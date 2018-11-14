package widgets

import (
	"fmt"

	"github.com/therecipe/qt/widgets"
)

var (
	// 控制颜色的qss模板
	colorStyle = "QLabel{color:%s;}"
	// 默认颜色-黑色
	defaultStyle = "QLabel{color:black;}"
)

// ColorLabel 使用QSS显示彩色文字
type ColorLabel struct {
	widgets.QLabel

	// color style sheet
	defaultColor string
}

// NewColorLabelWithColor 生成colorlabel，设置default color为color
// color为空则设置为黑色
// color可以是颜色对应的名字，例如"black", "green"
// 也可以是16进制的RGB值，例如 #ffffff, #ff08ff, #000000
func NewColorLabelWithColor(text, color string) *ColorLabel {
	l := NewColorLabel(nil, 0)

	l.SetDefaultColor(color)
	l.SetDefaultColorText(text)

	return l
}

// SetDefaultColor 设置defaultColor
// color为""时设置为黑色
// 不会改变现有text内容的颜色
func (l *ColorLabel) SetDefaultColor(color string) {
	if color == "" {
		l.defaultColor = defaultStyle
		return
	}

	l.defaultColor = fmt.Sprintf(colorStyle, color)
}

// ChangeColor 改变现有text的颜色
// 并且设置defaultColor为新的颜色
// color为""时设置为defaultStyle
func (l *ColorLabel) ChangeColor(color string) {
	l.SetDefaultColor(color)
	text := l.Text()
	l.SetDefaultColorText(text)
}

// SetColorText 用color显示新的text
// color为""时显示defaultStyle
func (l *ColorLabel) SetColorText(text, color string) {
	var style string
	if color == "" {
		style = defaultStyle
	} else {
		style = fmt.Sprintf(colorStyle, color)
	}

	l.SetText(text)
	l.SetStyleSheet(style)
}

// SetDefaultColorText 设置新的text值，并使其显示设置的default color
func (l *ColorLabel) SetDefaultColorText(text string) {
	l.SetText(text)
	l.SetStyleSheet(l.defaultColor)
}

// DropColor 去除自定义颜色，显示系统主题默认的颜色
func (l *ColorLabel) DropColor() {
	// 空字符串去除stylesheet
	l.SetStyleSheet("")
	l.defaultColor = ""
}
