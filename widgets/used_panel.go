package widgets

import (
	"fmt"
	"log"

	"github.com/astaxie/beego/orm"
	"github.com/therecipe/qt/widgets"

	"schannel-qt5/models"
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

	// 登录用户名
	user string
	info *parser.SSRInfo

	logger *log.Logger
}

// NewUsedPanelWithInfo 创建使用进度widget
func NewUsedPanelWithInfo(user string, info *parser.SSRInfo, logger *log.Logger) *UsedPanel {
	u := NewUsedPanel(nil, 0)
	u.user = user
	u.info = info
	u.logger = logger
	u.InitUI()

	return u
}

// InitUI 根据info初始化ui，并连接信号
func (u *UsedPanel) InitUI() {
	// 将string信息转换成int
	u.setData()

	// 初始化groupbox
	group := widgets.NewQGroupBox2("流量使用情况", u)

	// 初始化labels
	u.totalLabel = widgets.NewQLabel(nil, 0)
	u.usedLabel = widgets.NewQLabel(nil, 0)
	u.uploadLabel = widgets.NewQLabel(nil, 0)
	u.downloadLabel = widgets.NewQLabel(nil, 0)
	u.setLabels()

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

// saveUsedAmount 将使用量信息保存至数据库
func (u *UsedPanel) saveUsedAmount() {
	db := orm.NewOrm()
	now := parser.GetCurrentDay()
	err := models.SetUsedAmount(db, u.user, u.info.Service, u.total, u.upload, u.download, now)
	if err != nil {
		u.logger.Fatalf("save used amound error: %v\n", err)
	}
	u.logger.Printf("saved used_amount success")
}

// setData 将used转换成int并设置给widget，随后存入数据库
func (u *UsedPanel) setData() {
	format := "[%s] convert error: %s"
	u.total = convertToKb(u.info.TotalData)
	if u.total == -1 {
		u.logger.Fatalf(format, "total", u.info.TotalData)
	}

	u.used = convertToKb(u.info.UsedData)
	if u.used == -1 {
		u.logger.Fatalf(format, "used", u.info.UsedData)
	}

	u.upload = convertToKb(u.info.Upload)
	if u.upload == -1 {
		u.logger.Fatalf(format, "upload", u.info.Upload)
	}

	u.download = convertToKb(u.info.Download)
	if u.download == -1 {
		u.logger.Fatalf(format, "download", u.info.Download)
	}

	u.saveUsedAmount()
}

func (u *UsedPanel) setLabels() {
	u.totalLabel.SetText(fmt.Sprintf("套餐总量：%v", u.info.TotalData))
	u.usedLabel.SetText(fmt.Sprintf("已使用：%v", u.info.UsedData))
	u.uploadLabel.SetText(fmt.Sprintf("已上传：%v", u.info.Upload))
	u.downloadLabel.SetText(fmt.Sprintf("已下载：%v", u.info.Download))
}

// dataRefresh 刷新数据显示
func (u *UsedPanel) dataRefresh(info *parser.SSRInfo) {
	u.info = info
	u.setData()

	u.setLabels()
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
