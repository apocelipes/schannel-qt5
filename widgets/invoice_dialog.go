package widgets

import (
	"fmt"
	"strconv"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"

	"schannel-qt5/crawler"
	"schannel-qt5/parser"
)

// InvoiceDialog 显示全部的账单信息
type InvoiceDialog struct {
	widgets.QDialog

	// goroutine中触发显示ErrorDialog及显示下载成功
	_ func(errInfo string) `signal:"errorHappened,auto"`
	_ func(file string)    `signal:"downloadFinish,auto"`

	table   *widgets.QTableWidget
	infoBar *widgets.QStatusBar
	// 是否在选中时复制到剪贴板
	copy2Clipboard *widgets.QCheckBox
	// 选中的行数
	selected *widgets.QLabel
	// 选中的链接
	link *widgets.QLabel

	invoices   []*parser.Invoice
	dataBridge UserDataBridge
}

var (
	// 表头
	cols = []string{
		"账单编号",
		"链接",
		"开始日期",
		"结束日期",
		"金额（元）",
		"支付状态",
	}
)

// NewInvoiceDialogWithData 生成dialog
// bridge用户获取登录信息
func NewInvoiceDialogWithData(bridge UserDataBridge, data []*parser.Invoice) *InvoiceDialog {
	dialog := NewInvoiceDialog(nil, 0)
	dialog.invoices = data
	dialog.dataBridge = bridge

	// 设置infobar，选中内容时显示账单链接
	dialog.infoBar = widgets.NewQStatusBar(nil)
	dialog.selected = widgets.NewQLabel2("未选中", nil, 0)
	dialog.link = widgets.NewQLabel(nil, 0)
	dialog.infoBar.AddPermanentWidget(dialog.selected, 0)
	dialog.infoBar.AddPermanentWidget(dialog.link, 0)

	dialog.copy2Clipboard = widgets.NewQCheckBox2("将链接复制到剪贴板", nil)
	dialog.copy2Clipboard.SetChecked(false)
	// 选中时如果已经选择了Link则进行复制
	dialog.copy2Clipboard.ConnectClicked(func(_ bool) {
		link := dialog.link.Text()
		if dialog.copy2Clipboard.IsChecked() && link != "" {
			dialog.copyLink(link)
		}
	})

	// 初始化table，数据已经被排序
	dialog.table = widgets.NewQTableWidget(nil)
	// 设置行数，不设置将不显示任何数据
	dialog.table.SetRowCount(len(dialog.invoices))
	// 设置表头
	dialog.table.SetColumnCount(len(cols))
	dialog.table.SetHorizontalHeaderLabels(cols)
	// 去除边框
	dialog.table.SetShowGrid(false)
	dialog.table.SetFrameShape(widgets.QFrame__NoFrame)
	// 去除行号
	dialog.table.VerticalHeader().SetVisible(false)
	// 设置选择整行内容
	dialog.table.SetSelectionBehavior(widgets.QAbstractItemView__SelectRows)
	dialog.table.SetSelectionMode(widgets.QAbstractItemView__SingleSelection)
	// 设置table的数据项目
	dialog.setTable()

	dialog.table.ConnectCellClicked(func(row, col int) {
		dialog.setLink(dialog.invoices[row])
	})
	dialog.table.ConnectCellDoubleClicked(func(row int, column int) {
		dialog.showInvoiceView(dialog.invoices[row])
	})

	dialog.table.ConnectContextMenuEvent(dialog.invoiceContextMenu)

	// 设置不可编辑table
	dialog.table.SetEditTriggers(widgets.QAbstractItemView__NoEditTriggers)
	// Qt5.12.0无法正常折叠link text
	dialog.table.ResizeColumnsToContents()

	vbox := widgets.NewQVBoxLayout()
	vbox.AddWidget(dialog.table, 0, 0)
	vbox.AddWidget(dialog.copy2Clipboard, 0, core.Qt__AlignLeft)
	vbox.AddStretch(0)
	vbox.AddWidget(dialog.infoBar, 0, 0)
	dialog.SetLayout(vbox)
	dialog.setDialogSize()
	dialog.SetWindowTitle("账单详情")
	dialog.SetAttribute(core.Qt__WA_DeleteOnClose, true)

	return dialog
}

// setDialogSize 设置dialog宽度与table一致
func (dialog *InvoiceDialog) setDialogSize() {
	width := 0
	for i := 0; i < len(cols); i++ {
		width += dialog.table.ColumnWidth(i)
	}
	dialog.SetMinimumWidth(width)
}

// setTable 设置table
func (dialog *InvoiceDialog) setTable() {
	for row := 0; row < len(dialog.invoices); row++ {
		invoice := dialog.invoices[row]

		number := widgets.NewQTableWidgetItem2(invoice.Number, 0)
		dialog.table.SetItem(row, 0, number)
		link := widgets.NewQTableWidgetItem2(invoice.Link, 0)
		dialog.table.SetItem(row, 1, link)

		startTime := time2string(invoice.StartDate)
		start := widgets.NewQTableWidgetItem2(startTime, 0)
		dialog.table.SetItem(row, 2, start)
		expireTime := time2string(invoice.ExpireDate)
		expire := widgets.NewQTableWidgetItem2(expireTime, 0)
		dialog.table.SetItem(row, 3, expire)

		payment := strconv.FormatInt(invoice.Payment, 10)
		pay := widgets.NewQTableWidgetItem2(payment, 0)
		dialog.table.SetItem(row, 4, pay)

		text := ""
		color := ""
		if invoice.State == parser.NeedPay {
			text = "未付款"
			color = "red"
		} else if invoice.State == parser.FinishedPay {
			text = "已付款"
			color = "green"
		}
		label := NewColorLabelWithColor(text, color)
		dialog.table.SetCellWidget(row, 5, label)
	}
}

// setLink 当选中row中的单元格时将链接更新到infoBar
func (dialog *InvoiceDialog) setLink(invoice *parser.Invoice) {
	index := 0
	for i, v := range dialog.invoices {
		if v == invoice {
			index = i
			break
		}
	}
	dialog.selected.SetText(fmt.Sprintf("选中第%d行", index+1))
	dialog.link.SetText(invoice.Link)
	dialog.copyLink(invoice.Link)
}

// copyLink 如果勾选了copy2Clipboard则将link复制到系统剪贴板
func (dialog *InvoiceDialog) copyLink(link string) {
	if dialog.copy2Clipboard.IsChecked() {
		dialog.copy(link)
	}
}

// copy 将值复制进剪贴板
func (dialog *InvoiceDialog) copy(text string) {
	clip := gui.QGuiApplication_Clipboard()
	clip.SetText(text, gui.QClipboard__Clipboard)
}

// showInvoiceView 显示invoice对应的InvoiceViewWidget
func (dialog *InvoiceDialog) showInvoiceView(invoice *parser.Invoice) {
	dialog.setLink(invoice)

	data, err := crawler.GetInvoiceInfoHTML(invoice,
		dialog.dataBridge.GetCookies(),
		dialog.dataBridge.GetProxy())
	if err != nil {
		dialog.dataBridge.GetLogger().Println("GetInvoiceInfoHTML error:", err)
		return
	}
	viewHTML := parser.GetInvoiceViewHTML(string(data))

	invoiceView := NewInvoiceViewWidget(dialog, 0)
	invoiceView.SetHTML(viewHTML)
	invoiceView.Show()
}

// invoiceContextMenu 显示table中invoice的右键菜单选项
func (dialog *InvoiceDialog) invoiceContextMenu(_ *gui.QContextMenuEvent) {
	invoice := dialog.invoices[dialog.table.CurrentItem().Row()]
	dialog.setLink(invoice)

	menu := widgets.NewQMenu(dialog)
	menu.AddAction2(gui.NewQIcon5(":/image/ic_copy_link.svg"), "复制")
	menu.AddAction2(gui.NewQIcon5(":/image/download.svg"), "下载")
	menu.ConnectTriggered(func(action *widgets.QAction) {
		switch action.Text() {
		case "下载":
			dialog.download(invoice)
		case "复制":
			dialog.copy(invoice.Link)
		}
	})

	menu.Exec2(gui.QCursor_Pos(), nil)
	menu.DestroyQMenu()
}

// errorHappened 收到goroutine返回的err并显示
func (dialog *InvoiceDialog) errorHappened(errInfo string) {
	showErrorDialog(errInfo, dialog)
}

// 下载完成，显示成功信息
func (dialog *InvoiceDialog) downloadFinish(file string) {
	info := fmt.Sprintf("%s下载成功", file)
	ShowNotification("账单", info, "", -1)
}

// download 下载选定的账单
// 更新statusbar，启动另一个goroutine进行下载并反馈进度
func (dialog *InvoiceDialog) download(invoice *parser.Invoice) {
	defaultName := fmt.Sprintf("账单-%s.pdf", invoice.Number)
	filter := "PDF Files(*.pdf)"
	// 获取上次保存文件的目录
	savePath, err := getFileSavePath("invoice", defaultName, filter, dialog)
	if err == ErrCanceled {
		return
	} else if err != nil {
		showErrorDialog("保存路径获取失败："+err.Error(), dialog)
		return
	}

	cookies := dialog.dataBridge.GetCookies()
	proxy := dialog.dataBridge.GetProxy()
	html, err := crawler.GetInvoiceInfoHTML(invoice, cookies, proxy)
	if err != nil {
		logger := dialog.dataBridge.GetLogger()
		logger.Printf("GetInvoiceInfoHTML error: %v\n", err)
		dialog.ErrorHappened("获取下载地址失败：" + err.Error())
		return
	}
	downloadURL := parser.GetInvoiceDownloadURL(html)
	downloader, err := NewHTTPDownloader2(downloadURL, invoice.Link, proxy, cookies)
	if err != nil {
		logger := dialog.dataBridge.GetLogger()
		logger.Printf("NewHTTPDownloader2 error: %v\n", err)
		showErrorDialog("获取下载器失败："+err.Error(), dialog)
		return
	}
	downloader.SetParent(dialog)

	progressDialog := getProgressDialog("保存账单", "账单下载进度：", dialog)
	progressDialog.ConnectCanceled(func() {
		downloader.Stop()
		progressDialog.Cancel()
		dialog.ErrorHappened("下载已取消")
	})
	downloader.ConnectUpdateProgress(func(size int) {
		// 已经cancel的dialog不能调用setValue，避免dialog反复出现
		if progressDialog.WasCanceled() {
			return
		}

		progressDialog.SetValue(size)
	})
	downloader.ConnectUpdateMax(progressDialog.SetMaximum)
	downloader.ConnectFailed(func(err error) {
		dialog.ErrorHappened("下载发生错误：" + err.Error())
	})
	downloader.ConnectDone(func() {
		progressDialog.Cancel()
		dialog.DownloadFinish(savePath)
	})

	go downloader.Download(savePath)
	progressDialog.Exec()
}
