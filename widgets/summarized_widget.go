package widgets

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"

	"schannel-qt5/config"
	"schannel-qt5/geoip"
	"schannel-qt5/parser"
)

// SummarizedWidget 综合服务信息显示，包括用户信息，服务信息
type SummarizedWidget struct {
	widgets.QWidget

	// 收到数据变动
	_ func() `signal:"dataRefresh,auto"`
	// 发出数据变动，让上层控件更新service
	// 上层控件完成service的更新后发送DataRefresh信号，int值为当前的index
	_ func(int) `signal:"serviceNeedUpdate"`

	// 用户数据接口
	dataBridge UserDataBridge

	// 服务信息面板
	servicePanel *ServicePanel
	// ssr开关面板
	switchPanel *SSRSwitchPanel
	// 使用量信息
	usedPanel *UsedPanel
	// 是否需要付款
	invoicePanel *InvoicePanel
	// 下载GeoIP Database
	getGeoButton *widgets.QPushButton

	// 用户名-email
	user string
	// 用户配置
	conf *config.UserConfig
	// 服务信息
	service *parser.Service
	// 综合信息面板编号，因为可能不止一个服务，所以用来做身份区别
	// index与services数组中的索引相同
	index int
}

// NewSummarizedWidget2 创建综合信息控件
func NewSummarizedWidget2(index int,
	user string,
	service *parser.Service,
	conf *config.UserConfig,
	dataBridge UserDataBridge) *SummarizedWidget {
	if user == "" || dataBridge == nil {
		return nil
	}
	sw := NewSummarizedWidget(nil, 0)

	sw.user = user
	sw.dataBridge = dataBridge
	sw.service = service
	sw.conf = conf
	sw.index = index
	sw.InitUI()

	return sw
}

// InitUI 初始化UI
func (sw *SummarizedWidget) InitUI() {
	ssrInfo := sw.dataBridge.SSRInfos(sw.service)
	logger := sw.dataBridge.GetLogger()
	sw.servicePanel = NewServicePanel2(sw.user, ssrInfo)
	sw.invoicePanel = NewInvoicePanelWithData(sw.dataBridge)
	sw.switchPanel = NewSSRSwitchPanel2(sw.conf, ssrInfo.Nodes, logger)
	sw.usedPanel = NewUsedPanelWithInfo(sw.user, ssrInfo, logger)

	updateButton := widgets.NewQPushButton2("刷新", nil)
	// 通知上层控件更新sw的service
	updateButton.ConnectClicked(func(_ bool) {
		sw.ServiceNeedUpdate(sw.index)
	})
	leftLayout := widgets.NewQVBoxLayout()
	leftLayout.AddWidget(sw.servicePanel, 0, 0)
	leftLayout.AddWidget(sw.invoicePanel, 0, 0)
	leftLayout.AddStretch(0)
	buttonLayout := widgets.NewQHBoxLayout()
	buttonLayout.AddStretch(0)

	geoPath, err := geoip.GetGeoIPSavePath()
	dbPath := filepath.Join(geoPath, geoip.DatabaseName)
	if err != nil {
		logger.Println(err)
	} else if err := syscall.Access(dbPath, syscall.F_OK); err != nil {
		logger.Println("未下载GeoIP数据库")
		sw.getGeoButton = widgets.NewQPushButton2("下载GeoIP数据库", nil)
		sw.getGeoButton.ConnectClicked(sw.downloadGeoIPDatabase)
		buttonLayout.AddWidget(sw.getGeoButton, 0, core.Qt__AlignRight)
	} else if geoip.IsGeoIPOutdated(24 * time.Hour * 30) {
		// GeoIP数据有效期为30天，超过后需要更新
		logger.Println("GeoIP数据库需要更新")
		sw.getGeoButton = widgets.NewQPushButton2("更新GeoIP数据库", nil)
		sw.getGeoButton.ConnectClicked(sw.downloadGeoIPDatabase)
		buttonLayout.AddWidget(sw.getGeoButton, 0, core.Qt__AlignRight)
	}
	buttonLayout.AddWidget(updateButton, 0, core.Qt__AlignRight)
	leftLayout.AddLayout(buttonLayout, 0)

	rightLayout := widgets.NewQVBoxLayout()
	rightLayout.AddWidget(sw.switchPanel, 0, 0)
	rightLayout.AddWidget(sw.usedPanel, 0, 0)

	mainLayout := widgets.NewQHBoxLayout()
	mainLayout.AddLayout(leftLayout, 0)
	mainLayout.AddLayout(rightLayout, 0)
	sw.SetLayout(mainLayout)
}

// dataRefresh 处理数据更新
// 一般在SetService之后调用，直接调用将更新servicePanel以外的数据
func (sw *SummarizedWidget) dataRefresh() {
	// sw.service已经被外部更新
	ssrInfo := sw.dataBridge.SSRInfos(sw.service)
	sw.servicePanel.UpadteInfo(sw.user, ssrInfo)
	sw.invoicePanel.UpdateInvoices(sw.dataBridge.Invoices())
	sw.switchPanel.DataRefresh(sw.conf, ssrInfo.Nodes)
	sw.usedPanel.DataRefresh(ssrInfo)
	ShowNotification("数据更新", "数据更新成功", "", -1)
}

// SetService 重新设置service，可用于更新数据
// 调用后一般需要出发DataRefresh信号
func (sw *SummarizedWidget) SetService(service *parser.Service) {
	sw.service = service
}

// UpdateConfig 当config更新时刷新switchPanel
// 一般用作ConfigWidget的ConfigChanged信号处理函数
func (sw *SummarizedWidget) UpdateConfig(conf *config.UserConfig) {
	sw.conf = conf
	nodes := sw.dataBridge.SSRInfos(sw.service).Nodes
	sw.switchPanel.DataRefresh(sw.conf, nodes)
	ShowNotification("配置更新", "配置更新成功", "", -1)
}

// 下载GeoIP数据库的回调函数
func (sw *SummarizedWidget) downloadGeoIPDatabase(_ bool) {
	geoPath, err := geoip.GetGeoIPSavePath()
	if err != nil {
		info := fmt.Sprintf("GetGeoIPSavePath error: %v", err)
		showErrorDialog(info, sw)
		return
	}

	if err := syscall.Access(geoPath, syscall.F_OK); err != nil {
		err = os.MkdirAll(geoPath, 0775)
		if err != nil {
			info := fmt.Sprintf("make download dir: %s error: %v", geoPath, err)
			showErrorDialog(info, sw)
			return
		}
	}
	savePath := filepath.Join(geoPath, "GeoLite2-City.mmdb.gz")

	downloader, err := NewHTTPDownloader2(geoip.DownloadPath, "", "", nil)
	if err != nil {
		info := fmt.Sprintf("downloader error: %v", err)
		showErrorDialog(info, sw)
		return
	}
	downloader.SetParent(sw)

	progressDialog := getProgressDialog("下载地理信息", "GeoIP Database下载进度：", sw)
	progressDialog.ConnectCanceled(func() {
		downloader.Stop()
		progressDialog.Cancel()
		showErrorDialog("下载已取消", sw)
	})
	downloader.ConnectUpdateProgress(func(size int) {
		if progressDialog.WasCanceled() {
			return
		}

		progressDialog.SetValue(size)
	})
	downloader.ConnectUpdateMax(progressDialog.SetMaximum)
	downloader.ConnectFailed(func(err error) {
		progressDialog.Cancel()
		info := fmt.Sprintf("下载发生错误: %v", err)
		showErrorDialog(info, sw)
	})
	downloader.ConnectDone(func() {
		progressDialog.Cancel()

		// 解压缩数据库
		f, err := os.Open(savePath)
		if err != nil {
			showErrorDialog(err.Error(), sw)
			return
		}
		defer f.Close()
		greader, err := gzip.NewReader(f)
		if err != nil {
			showErrorDialog(err.Error(), sw)
			return
		}
		defer greader.Close()

		dbPath := filepath.Join(geoPath, geoip.DatabaseName)
		dbFile, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			showErrorDialog(err.Error(), sw)
			return
		}
		defer dbFile.Close()

		buf, err := ioutil.ReadAll(greader)
		if err != nil {
			showErrorDialog(err.Error(), sw)
			return
		}
		_, err = dbFile.Write(buf)
		if err != nil {
			showErrorDialog(err.Error(), sw)
			return
		}

		ShowNotification("下载", "地理信息数据下载完成", "", -1)
		os.Remove(savePath)
		sw.getGeoButton.Hide()
		// 更新switch面板的GeoIP信息
		ssrInfo := sw.dataBridge.SSRInfos(sw.service)
		sw.switchPanel.DataRefresh(sw.conf, ssrInfo.Nodes)
		sw.dataBridge.GetLogger().Println("GeoIP Database下载成功")
	})

	go downloader.Download(savePath)
	progressDialog.Exec()
}
