package widgets

import "github.com/therecipe/qt/widgets"

// login时显示正在登录的busy indicator
type LoginIndicator struct {
	widgets.QWidget

	infoLabel *widgets.QLabel
	progress  *widgets.QProgressBar
}

func NewLoginIndicator2() *LoginIndicator {
	indicator := NewLoginIndicator(nil, 0)
	indicator.InitUI()

	return indicator
}

func (i *LoginIndicator) InitUI() {
	i.infoLabel = widgets.NewQLabel2("登录中：", nil, 0)
	i.progress = widgets.NewQProgressBar(nil)
	i.progress.SetRange(0, 0)

	mainLayout := widgets.NewQHBoxLayout()
	mainLayout.AddWidget(i.infoLabel, 0, 0)
	mainLayout.AddWidget(i.progress, 0, 0)
	i.SetLayout(mainLayout)
}
