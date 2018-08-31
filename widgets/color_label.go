package widgets

import (
	"errors"
	"fmt"
	"github.com/therecipe/qt/widgets"
	"regexp"
)

const (
	red    = `<font color="#ff0000">%s</font>`
	orange = `<font color="#ff8000">%s</font>`
	yellow = `<font color="#ffff00">%s</font>`
	green  = `<font color="#00ff00">%s</font>`
	blue   = `<font color="#0080ff">%s</font>`
	indigo = `<font color="#0000a0">%s</font>`
	purple = `<font color="#8000ff">%s</font>`
	gray   = `<font color="#c0c0c0">%s</font>`
)

var (
	colors map[string]string = map[string]string{
		"red":    red,
		"orange": orange,
		"blue":   blue,
		"yellow": yellow,
		"green":  green,
		"indigo": indigo,
		"purple": purple,
		"gray":   gray,
	}
	textMatcher = regexp.MustCompile(`<font color=".+">(.+)</font>`)
)

type ColorLabel struct {
	widgets.QLabel

	_ func(string) `slot:"setText,auto"`

	defaultColor string
}

// NewColorLabelWithColor 生成colorlabel，设置default color为color
func NewColorLabelWithColor(color, text string) *ColorLabel {
	l := NewColorLabel(nil, 0)
	c, ok := colors[color]
	if !ok {
		if color == "" {
			c = "black"
		}

		return nil
	}

	l.defaultColor = c
	l.SetText(text)

	return l
}

// SetDefaultColor 设置defaultColor
// 不会改变现有text内容的颜色
// 颜色不存在时返回error
func (l *ColorLabel) SetDefaultColor(color string) error {
	c, ok := colors[color]
	if !ok {
		return errors.New("color does not support")
	}
	l.defaultColor = c

	return nil
}

// ChangeColor 改变现有text的颜色
// 并且设置defaultColor为新的颜色
// 颜色不存在时返回error
func (l *ColorLabel) ChangeColor(color string) error {
	if err := l.SetDefaultColor(color); err != nil {
		return err
	}

	text, err := l.PureText()
	if err != nil {
		return err
	}
	l.SetText(text)

	return nil
}

// PureText 获取纯text内容，不包括color部分
// 不许label的text为空
func (l *ColorLabel) PureText() (string, error) {
	text := l.Text()
	tmp := textMatcher.FindStringSubmatch(text)
	if len(tmp) < 2 {
		return "", errors.New("cannot found text")
	}

	return tmp[1], nil
}

// SetColorText 用color显示新的text
func (l *ColorLabel) SetColorText(color, text string) {
	if color == "" || l.defaultColor == "black" {
		l.SetTextDefault(text)
	}

	c, ok := colors[color]
	if !ok {
		return
	}

	newText := fmt.Sprintf(c, text)
	l.SetTextDefault(newText)
}

// setText 设置新的text值，并使其显示创建时指定的default color
func (l *ColorLabel) setText(text string) {
	l.SetColorText(l.defaultColor, text)
}
