package widgets

import (
	"fmt"
	"log"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/therecipe/qt/charts"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"

	"schannel-qt5/models"
	"schannel-qt5/parser"
)

// ChartsDialog 显示上传下载对比pie chart和使用趋势line charts
type ChartsDialog struct {
	widgets.QDialog

	// 用于从数据库获取相关服务的数据
	user    string
	service string
	amounts []*models.UsedAmount
	// charts color
	downloadColor *gui.QColor
	uploadColor   *gui.QColor
	// 统计图的边界日期
	// TODO 未来可以自行选择日期，目前以GetCurrentDay为基准
	date   time.Time
	logger *log.Logger
}

// NewChartsDialog2 创建统计图表对话框，创建失败会将信息写入日志
// 对话框不能放大或拉伸，会导致charts无法正常显示
func NewChartsDialog2(user, service string, logger *log.Logger, parent widgets.QWidget_ITF) *ChartsDialog {
	dialog := NewChartsDialog(parent, 0)
	dialog.user = user
	dialog.service = service
	dialog.date = parser.GetCurrentDay()
	dialog.logger = logger
	dialog.downloadColor = gui.NewQColor6("red")
	dialog.uploadColor = gui.NewQColor6("orange")

	db := orm.NewOrm()
	var err error
	dialog.amounts, err = models.GetRecentUsedAmount(db, dialog.user, dialog.service, dialog.date)
	if err != nil {
		dialog.logger.Fatalf("ChartsDialog get recent data error: %v\n", err)
	}
	dialog.InitUI()

	return dialog
}

// InitUI 生成图表并显示
func (dialog *ChartsDialog) InitUI() {
	pieChart := dialog.CreatePieChart()
	uploadLineChart := dialog.CreateUploadLineChart()
	downloadLineChart := dialog.CreateDownloadLineChart()

	rightLayout := widgets.NewQVBoxLayout()
	rightLayout.AddWidget(downloadLineChart, 0, 0)
	rightLayout.AddWidget(uploadLineChart, 0, 0)
	mainLayout := widgets.NewQHBoxLayout()
	mainLayout.AddWidget(pieChart, 0, 0)
	mainLayout.AddLayout(rightLayout, 0)
	dialog.SetLayout(mainLayout)
	// 设置dialog大小，保证charts能正常绘制
	dialog.SetMinimumHeight(700)
	dialog.SetMinimumWidth(900)
	// 非模态dialog，设置关闭后销毁dialog对象
	dialog.SetAttribute(core.Qt__WA_DeleteOnClose, true)
	dialog.SetWindowTitle("数据用量统计")
}

// CreatePieChart 生成upload/download饼图
func (dialog *ChartsDialog) CreatePieChart() *charts.QChartView {
	// amounts倒序排列，0是今天
	todayAmount := dialog.amounts[0]
	pie := charts.NewQPieSeries(nil)
	pie.Append3("下载", float64(todayAmount.Download))
	pie.Append3("上传", float64(todayAmount.Upload))
	items := pie.Slices()
	downloadPer, uploadPer := computePercent(todayAmount.Download, todayAmount.Upload)
	items[0].SetColor(dialog.downloadColor)
	items[0].SetLabelVisible(true)
	items[0].SetLabelPosition(charts.QPieSlice__LabelOutside)
	items[0].SetLabel(fmt.Sprintf("下载:%.2f%%", downloadPer))
	items[1].SetColor(dialog.uploadColor)
	items[1].SetLabelVisible(true)
	items[1].SetLabelPosition(charts.QPieSlice__LabelOutside)
	items[1].SetLabel(fmt.Sprintf("上传:%.2f%%", uploadPer))

	chart := charts.NewQChart(nil, 0)
	chart.AddSeries(pie)
	chart.SetTitle("上传/下载对比图")
	chart.Legend().Hide()
	chartView := charts.NewQChartView2(chart, nil)
	chartView.SetRenderHints(gui.QPainter__Antialiasing)
	return chartView
}

// computePercent 计算download和upload的百分比
func computePercent(download, upload int) (float64, float64) {
	total := float64(download + upload)
	return float64(download) / total * 100, float64(upload) / total * 100
}

// CreateUploadLineChart 生成upload流量使用情况趋势图
func (dialog *ChartsDialog) CreateUploadLineChart() *charts.QChartView {
	line := charts.NewQLineSeries(nil)
	line.SetName("上传")
	line.SetColor(dialog.uploadColor)
	// dataSet用于计算计量单位和range
	dataSet := make([]int, 0, len(dialog.amounts))
	for _, v := range dialog.amounts {
		dataSet = append(dataSet, v.Upload)
	}
	ratio, unit := computeSizeUnit(dataSet)

	for i := len(dialog.amounts) - 1; i >= 0; i-- {
		date := dialog.amounts[i].Date
		qDate := core.NewQDate3(date.Year(), int(date.Month()), date.Day())
		datetime := core.NewQDateTime2(qDate)
		value := float64(dialog.amounts[i].Upload) / float64(ratio)
		line.Append(float64(datetime.ToMSecsSinceEpoch()), value)
	}
	chart := charts.NewQChart(nil, 0)
	chart.AddSeries(line)
	chart.SetTitle("上传使用趋势（月底清零）")

	axisX := charts.NewQDateTimeAxis(nil)
	axisX.SetTickCount(5)
	axisX.SetFormat("MM-dd")
	axisX.SetTitleText("日期")
	chart.AddAxis(axisX, core.Qt__AlignBottom)

	axisY := charts.NewQValueAxis(nil)
	axisY.SetTitleText(fmt.Sprintf("上传（%s）", unit))
	axisY.SetRange(computeRange(dataSet, ratio))
	chart.AddAxis(axisY, core.Qt__AlignRight)

	line.AttachAxis(axisX)
	line.AttachAxis(axisY)
	chartView := charts.NewQChartView2(chart, nil)
	chartView.SetRenderHints(gui.QPainter__Antialiasing)
	return chartView
}

// CreateDownloadLineChart 生成download流量使用情况趋势图
func (dialog *ChartsDialog) CreateDownloadLineChart() *charts.QChartView {
	line := charts.NewQLineSeries(nil)
	line.SetName("下载")
	line.SetColor(dialog.downloadColor)
	// dataSet用于计算计量单位和range
	dataSet := make([]int, 0, len(dialog.amounts))
	for _, v := range dialog.amounts {
		dataSet = append(dataSet, v.Download)
	}
	ratio, unit := computeSizeUnit(dataSet)

	for i := len(dialog.amounts) - 1; i >= 0; i-- {
		date := dialog.amounts[i].Date
		qDate := core.NewQDate3(date.Year(), int(date.Month()), date.Day())
		datetime := core.NewQDateTime2(qDate)
		value := float64(dialog.amounts[i].Download) / float64(ratio)
		line.Append(float64(datetime.ToMSecsSinceEpoch()), value)
	}
	chart := charts.NewQChart(nil, 0)
	chart.AddSeries(line)
	chart.SetTitle("下载使用趋势（月底清零）")

	axisX := charts.NewQDateTimeAxis(nil)
	axisX.SetTickCount(5)
	axisX.SetFormat("MM-dd")
	axisX.SetTitleText("日期")
	chart.AddAxis(axisX, core.Qt__AlignBottom)

	axisY := charts.NewQValueAxis(nil)
	axisY.SetTitleText(fmt.Sprintf("下载（%s）", unit))
	axisY.SetRange(computeRange(dataSet, ratio))
	chart.AddAxis(axisY, core.Qt__AlignLeft)

	line.AttachAxis(axisX)
	line.AttachAxis(axisY)
	chartView := charts.NewQChartView2(chart, nil)
	chartView.SetRenderHints(gui.QPainter__Antialiasing)
	return chartView
}
