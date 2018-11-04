package widgets

import (
	"fmt"

	"github.com/therecipe/qt/widgets"

	"schannel-qt5/parser"
)

type UsedPanel struct {
	widgets.QWidget

	// 数据更新时触发的信号，更新展示的数据
	_ func(*parser.SSRInfo) `signal:"dataRefresh,auto"`

	usedBar     *ProgressBar
	uploadBar   *ProgressBar
	downloadBar *ProgressBar

	totalLabel    *widgets.QLabel
	usedLabel     *widgets.QLabel
	uploadLabel   *widgets.QLabel
	downloadLabel *widgets.QLabel

	// 套餐数据量信息
	total    int
	used     int
	upload   int
	download int
}

// NewUsedPanelWithInfo 创建使用进度widget
func NewUsedPanelWithInfo(info *parser.SSRInfo) *UsedPanel {
	u := NewUsedPanel(nil, 0)
	u.InitUI(info)

	return u
}

// InitUI 根据info初始化ui，并连接信号
func (u *UsedPanel) InitUI(info *parser.SSRInfo) {
	// 将string信息转换成int
	u.total = convertToKb(info.TotalData)
	u.used = convertToKb(info.UsedData)
	u.upload = convertToKb(info.Upload)
	u.download = convertToKb(info.Download)

	// 初始化groupbox
	group := widgets.NewQGroupBox2("流量使用情况", u)

	// 初始化labels
	u.totalLabel = widgets.NewQLabel(nil, 0)
	u.usedLabel = widgets.NewQLabel(nil, 0)
	u.uploadLabel = widgets.NewQLabel(nil, 0)
	u.downloadLabel = widgets.NewQLabel(nil, 0)
	u.setLabels(info)

	// 初始化progressbar
	u.usedBar = NewProgressBarWithMark(u.total, u.used, computeRatio(u.total))
	u.uploadBar = NewProgressBarWithMark(u.total, u.upload, computeRatio(u.total))
	u.downloadBar = NewProgressBarWithMark(u.total, u.download, computeRatio(u.total))

	// 布局管理
	vbox := widgets.NewQVBoxLayout()
	vbox.AddWidget(u.totalLabel, 0, 0)
	vbox.AddSpacing(0)
	vbox.AddWidget(u.usedLabel, 0, 0)
	vbox.AddWidget(u.usedBar, 0, 0)
	vbox.AddSpacing(0)
	vbox.AddWidget(u.downloadLabel, 0, 0)
	vbox.AddWidget(u.downloadBar, 0, 0)
	vbox.AddSpacing(0)
	vbox.AddWidget(u.uploadLabel, 0, 0)
	vbox.AddWidget(u.uploadBar, 0, 0)

	group.SetLayout(vbox)
	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.AddWidget(group, 0, 0)
	u.SetLayout(mainLayout)
	u.Show()
}

func (u *UsedPanel) setLabels(info *parser.SSRInfo) {
	u.totalLabel.SetText(fmt.Sprintf("套餐总量：%v", info.TotalData))
	u.usedLabel.SetText(fmt.Sprintf("已使用：%v", info.UsedData))
	u.uploadLabel.SetText(fmt.Sprintf("已上传：%v", info.Upload))
	u.downloadLabel.SetText(fmt.Sprintf("已下载：%v", info.Download))
}

// dataRefresh 刷新数据显示
func (u *UsedPanel) dataRefresh(info *parser.SSRInfo) {
	u.total = convertToKb(info.TotalData)
	u.used = convertToKb(info.UsedData)
	u.upload = convertToKb(info.Upload)
	u.download = convertToKb(info.Download)

	u.setLabels(info)
	// 更新progressbar
	u.usedBar.SetMaximum(u.total)
	u.usedBar.SetMark(computeRatio(u.total))
	u.usedBar.SetValue(u.used)
	u.uploadBar.SetMaximum(u.total)
	u.uploadBar.SetMark(computeRatio(u.total))
	u.uploadBar.SetValue(u.upload)
	u.downloadBar.SetMaximum(u.total)
	u.downloadBar.SetMark(computeRatio(u.total))
	u.downloadBar.SetValue(u.download)
}
