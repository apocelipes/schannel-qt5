package widgets

import (
	"github.com/therecipe/qt/widgets"
)

const (
	// Red progressbar进度条背景色红色
	Red = "QProgressBar{text-align:center}QProgressBar::chunk{background:red}"
	// Green progressbar进度条背景色绿色
	Green = "QProgressBar{text-align:center}QProgressBar::chunk{background:green}"
)

type ProgressBar struct {
	widgets.QProgressBar

	highMark int
}

// NewProgressBarWithMaxium 返回progressbar
// 并且设置最小值为0，最大值为max
// 设置当前值为current
// 设置颜色变换标志位数值为mark
func NewProgressBarWithMark(max, current, mark int) *ProgressBar {
	p := NewProgressBar(nil)
	p.SetMark(mark)
	p.SetRange(0, max)
	p.ConnectValueChanged(p.setColor)
	p.SetValue(current)

	return p
}

// setColor 根据新的value判断该选用哪种颜色
// 小于等于highMark为绿色，大于则为红色
func (p *ProgressBar) setColor(newValue int) {
	if newValue > p.highMark {
		p.SetStyleSheet(Red)
		return
	}
	p.SetStyleSheet(Green)
}

// SetMark 更改highmark
func (p *ProgressBar) SetMark(mark int) {
	p.highMark = mark
}
