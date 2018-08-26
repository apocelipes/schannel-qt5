package widgets

import (
	"strconv"

	"github.com/therecipe/qt/widgets"
)

type UsedProgressBar struct {
	widgets.QWidget

	pbar *widgets.QProgressBar
	used *widgets.QLineEdit

	min     int
	max     int
	current int
}

func (p *UsedProgressBar) InitUI(min, max, current int) {
	p.min = min
	p.max = max
	p.current = current

	p.SetLayout(widgets.NewQVBoxLayout())

	p.pbar = widgets.NewQProgressBar(p)
	p.pbar.SetMaximum(p.max)
	p.pbar.SetMinimum(p.min)
	p.pbar.SetValue(p.current)

	p.used = widgets.NewQLineEdit2(strconv.FormatInt(int64(p.current), 10), p)

	p.Layout().AddWidget(p.used)
	p.Layout().AddWidget(p.pbar)

	p.used.ConnectTextChanged(p.usedChanged)
	p.pbar.ConnectValueChanged(p.checkValue)

	p.Show()
}

func (p *UsedProgressBar) usedChanged(text string) {
	v, err := strconv.Atoi(text)
	if err != nil {
		return
	}

	if v >= p.pbar.Minimum() && v <= p.pbar.Maximum() {
		p.pbar.SetValue(v)
	}
}

func (p *UsedProgressBar) checkValue(v int) {
	if float64(v)/float64(p.pbar.Maximum()) > 0.8 {
		p.pbar.SetStyleSheet("QProgressBar{text-align:center}QProgressBar::chunk{background:red}")
	} else {
		p.pbar.SetStyleSheet("QProgressBar{text-align:center}QProgressBar::chunk{background:green}")
	}
}
