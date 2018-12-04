package widgets

import (
	"github.com/therecipe/qt/widgets"

	"schannel-qt5/config"
)

// ConfigWidget 显示和设置本客户端的配置
type ConfigWidget struct {
	widgets.QWidget

	// 通知conf已经更新
	_ func(*config.UserConfig) `signal:"configChanged"`

	// client设置
	clientConfigWidget *ClientConfigWidget

	// ssr client设置
	ssrClientConfigWidget *SSRConfigWidget

	// 配置数据和接口
	conf *config.UserConfig

	// 变动的配置是否已经保存
	saved bool
}

// NewConfigWidget2 根据conf生成ConfigWidget
func NewConfigWidget2(conf *config.UserConfig) *ConfigWidget {
	if conf == nil || conf.SSRClientConfig == nil {
		return nil
	}
	widget := NewConfigWidget(nil, 0)
	widget.conf = conf
	widget.saved = true
	widget.InitUI()

	return widget
}

// InitUI 初始化并显示
func (w *ConfigWidget) InitUI() {
	w.clientConfigWidget = NewClientConfigWidget2(w.conf)
	w.clientConfigWidget.ConnectValueChanged(func() {
		w.setSaved(false)
	})
	w.ssrClientConfigWidget = NewSSRConfigWidget2(w.conf.SSRClientConfig)
	w.ssrClientConfigWidget.ConnectValueChanged(func() {
		w.setSaved(false)
	})

	saveButton := widgets.NewQPushButton2("保存", nil)
	saveButton.ConnectClicked(func(_ bool) {
		w.SaveConfig()
	})

	// 大小策略，client和ssrClient大小2:1
	clientConfigSizePolicy := w.clientConfigWidget.SizePolicy()
	clientConfigSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Expanding)
	clientConfigSizePolicy.SetHorizontalStretch(2)
	w.clientConfigWidget.SetSizePolicy(clientConfigSizePolicy)
	ssrSizePolicy := w.ssrClientConfigWidget.SizePolicy()
	ssrSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Expanding)
	ssrSizePolicy.SetHorizontalStretch(1)
	w.ssrClientConfigWidget.SetSizePolicy(ssrSizePolicy)
	topLayout := widgets.NewQHBoxLayout()
	topLayout.AddWidget(w.clientConfigWidget, 0, 0)
	topLayout.AddWidget(w.ssrClientConfigWidget, 0, 0)
	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.AddLayout(topLayout, 0)
	mainLayout.AddWidget(saveButton, 0, 0)
	w.SetLayout(mainLayout)
}

// saveConfig 验证并保存配置
func (w *ConfigWidget) SaveConfig() {
	// 保存时不可修改设置信息
	w.SetEnabled(false)
	defer w.SetEnabled(true)

	var err error
	// 更新ssr client config
	err = w.ssrClientConfigWidget.UpdateSSRClientConfig()
	if err != nil {
		// 错误信息已经显示，无需使用showErrorDialog
		return
	}
	// 更新client config
	err = w.clientConfigWidget.UpdateClientConfig()
	if err != nil {
		return
	}

	if err := w.conf.StoreConfig(); err != nil {
		showErrorDialog("保存出错: " + err.Error(), w)
		return
	}

	// 设置状态为已保存
	w.setSaved(true)
	// 通知其他组件配置发生变化
	w.ConfigChanged(w.conf)
}

// Saved 获取配置是否已经保存
func (w *ConfigWidget) Saved() bool {
	return w.saved
}

// setSaved 设置配置是否已经保存
func (w *ConfigWidget) setSaved(saved bool) {
	w.saved = saved
}
