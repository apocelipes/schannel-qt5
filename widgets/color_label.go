package widgets

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/therecipe/qt/widgets"
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
	// 找不到对应颜色
	ErrColorNotFound = errors.New("cannot found text")
)

// ColorLabel 使用HTML显示彩色文字
type ColorLabel struct {
	widgets.QLabel

	defaultColor string
}

// NewColorLabelWithColor 生成colorlabel，设置default color为color
// color为空则设置为黑色
func NewColorLabelWithColor(text, color string) *ColorLabel {
	l := NewColorLabel(nil, 0)

	if err := l.SetDefaultColor(color); err != nil {
		return nil
	}
	l.SetDefaultColorText(text)

	return l
}

// getColor 从colors中获取对应的color并返回
func getColor(color string) string {
	c, ok := colors[color]
	if !ok {
		c = "none"
	}

	return c
}

// SetDefaultColor 设置defaultColor
// color为""时设置为黑色
// 不会改变现有text内容的颜色
// 颜色不存在时返回error
func (l *ColorLabel) SetDefaultColor(color string) error {
	if color == "" {
		l.defaultColor = "black"
		return nil
	}

	c := getColor(color)
	if c == "none" {
		return ErrColorNotFound
	}
	l.defaultColor = color

	return nil
}

// ChangeColor 改变现有text的颜色
// 并且设置defaultColor为新的颜色
// color为""时设置为黑色
// 颜色不存在时返回error
func (l *ColorLabel) ChangeColor(color string) error {
	if err := l.SetDefaultColor(color); err != nil {
		return err
	}

	text, err := l.PureText()
	if err != nil {
		return err
	}
	l.SetDefaultColorText(text)

	return nil
}

// PureText 获取纯text内容，不包括color部分
// 不许label的text为空
func (l *ColorLabel) PureText() (string, error) {
	text := l.Text()
	tmp := textMatcher.FindStringSubmatch(text)
	if len(tmp) < 2 {
		return "", ErrColorNotFound
	}

	return tmp[1], nil
}

// SetColorText 用color显示新的text
// color为""时显示黑色
func (l *ColorLabel) SetColorText(text, color string) error {
	if color == "black" {
		l.SetText(text)
		return nil
	} else if color == "" {
		if l.defaultColor != "" {
			// 递归调用，下次调用时color不可能是空，所以一定可以设置颜色
			// 不在这里直接使用SetDefaultColorText，避免调用链不清晰
			return l.SetColorText(text, l.defaultColor)
		}

		return errors.New("no default color set")
	}

	c := getColor(color)
	if c == "none" {
		return ErrColorNotFound
	}

	newText := fmt.Sprintf(c, text)
	l.SetText(newText)

	return nil
}

// SetDefaultColorText 设置新的text值，并使其显示创建时指定的default color
func (l *ColorLabel) SetDefaultColorText(text string) error {
	return l.SetColorText(text, l.defaultColor)
}
