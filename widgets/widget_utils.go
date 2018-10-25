package widgets

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
