package widgets

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// showErrorMsg 控制error label的显示
// err为nil则代表没有错误发生，如果label可见则设为隐藏
// err不为nil时设置label可见
// 设置label可见时返回true，否则返回false（不受label原有状态影响）
func showErrorMsg(label *ColorLabel, err error) bool {
	if err != nil {
		label.Show()
		return true
	}

	label.Hide()
	return false
}

// showErrorDialog 显示错误信息
func showErrorDialog(info string, parent widgets.QWidget_ITF) {
	errMsg := widgets.NewQErrorMessage(parent)
	errMsg.ShowMessage(info)
	errMsg.SetAttribute(core.Qt__WA_DeleteOnClose, true)
	errMsg.Exec()
}
