package widgets

import (
	"fmt"
	"strconv"

	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"

	"schannel-qt5/parser"
)

// InvoiceDialog 显示全部的账单信息
type InvoiceDialog struct {
	widgets.QDialog

	table   *widgets.QTableWidget
	infoBar *widgets.QStatusBar
	// 是否在选中时复制到剪贴板
	copy2Clipboard *widgets.QCheckBox
	// 选中的行数
	selected *widgets.QLabel
	// 选中的链接
	link *widgets.QLabel

	invoices []*parser.Invoice
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
func NewInvoiceDialogWithData(data []*parser.Invoice) *InvoiceDialog {
	dialog := NewInvoiceDialog(nil, 0)
	dialog.invoices = data

	// 设置infobar，选中内容时显示账单链接
	dialog.infoBar = widgets.NewQStatusBar(nil)
	dialog.selected = widgets.NewQLabel2("未选中", nil, 0)
	dialog.link = widgets.NewQLabel(nil, 0)
	dialog.infoBar.AddPermanentWidget(dialog.selected, 0)
	dialog.infoBar.AddPermanentWidget(dialog.link, 0)

	dialog.copy2Clipboard = widgets.NewQCheckBox2("将链接复制到剪贴板", nil)
	dialog.copy2Clipboard.SetChecked(false)

	// 初始化table，数据已经被排序
	dialog.table = widgets.NewQTableWidget(nil)
	// 设置行数，不设置将不显示任何数据
	dialog.table.SetRowCount(len(dialog.invoices))
	// 设置表头
	dialog.table.SetColumnCount(len(cols))
	dialog.table.SetHorizontalHeaderLabels(cols)
	// 设置链接列的列宽，以显示更完整的内容
	linkColWidth := dialog.table.ColumnWidth(1) * 2
	dialog.table.SetColumnWidth(1, linkColWidth)
	// 去除边框
	dialog.table.SetShowGrid(false)
	dialog.table.SetFrameShape(widgets.QFrame__NoFrame)
	// 去除行号
	dialog.table.VerticalHeader().SetVisible(false)
	// 设置table的数据项目
	dialog.setTable()
	dialog.table.ConnectItemClicked(dialog.setLink)

	dialog.table.ConnectCellClicked(func(row, col int) {
		invoice := dialog.invoices[row]
		dialog.selected.SetText(fmt.Sprintf("选中第%d行", row+1))
		dialog.link.SetText(invoice.Link)
		dialog.copyLink(row)
	})

	// 设置不可编辑table
	dialog.table.SetEditTriggers(widgets.QAbstractItemView__NoEditTriggers)

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
func (dialog *InvoiceDialog) setLink(item *widgets.QTableWidgetItem) {
	index := item.Row()
	invoice := dialog.invoices[index]
	dialog.selected.SetText(fmt.Sprintf("选中第%d行", index+1))
	dialog.link.SetText(invoice.Link)
	dialog.copyLink(index)
}

// copyLink 如果勾选了copy2Clipboard则将link复制到系统剪贴板
func (dialog *InvoiceDialog) copyLink(index int) {
	if dialog.copy2Clipboard.IsChecked() {
		clip := gui.QGuiApplication_Clipboard()
		clip.SetText(dialog.invoices[index].Link, gui.QClipboard__Clipboard)
	}
}
