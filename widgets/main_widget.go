package widgets

import (
	"fmt"
	"log"
	"net/http"

	"github.com/astaxie/beego/orm"
	"github.com/therecipe/qt/widgets"

	"schannel-qt5/config"
)

// MainWidget 客户端主界面，并处理子widget信号
type MainWidget struct {
	widgets.QMainWindow

	login *LoginWidget
	// 可能具有多个服务
	summary []*SummarizedWidget
	// 客户端配置widget
	setting *ConfigWidget
	// 包含三种widget的tabWidget
	tab *widgets.QTabWidget

	// 用户数据接口
	dataBridge UserDataBridge
	// 用户配置
	conf *config.UserConfig
	// 存储用户登录信息和使用信息
	db orm.Ormer
	// 用户名
	user string
	// 日志记录
	logger *log.Logger
}

// NewMainWidget2 创建主界面
func NewMainWidget2(conf *config.UserConfig, logger *log.Logger, db orm.Ormer) *MainWidget {
	widget := NewMainWidget(nil, 0)
	widget.conf = conf
	widget.db = db
	widget.logger = logger
	widget.InitUI()

	return widget
}

// InitUI 初始化UI
// 先初始化login，登录成功后login隐藏，再初始化summary和setting
func (m *MainWidget) InitUI() {
	m.tab = widgets.NewQTabWidget(nil)
	m.login = NewLoginWidget2(m.conf, m.logger, m.db)
	m.tab.AddTab(m.login, "登录")
	m.login.ConnectLoginSuccess(m.finishLogin)
	m.SetWindowTitle("schannel-qt5")
	m.SetCentralWidget(m.tab)
}

// finishLogin 登录成功后隐藏LoginWidget，显示summary和setting
func (m *MainWidget) finishLogin(user string, cookies []*http.Cookie) {
	m.user = user
	m.dataBridge = NewDataBridge(cookies, m.conf.Proxy.String(), m.logger)
	// 删除login，因为目前只有login一个widget所以index是0
	m.tab.RemoveTab(0)
	// 移动到左上角，避免窗口因较长显示不完整
	m.Move2(0, 0)

	m.setting = NewConfigWidget2(m.conf)
	// 可能存在多个服务
	services := m.dataBridge.ServiceInfos()
	for i, v := range services {
		widget := NewSummarizedWidget2(i, m.user, v, m.conf, m.dataBridge)
		// 处理更新请求
		widget.ConnectServiceNeedUpdate(func(index int) {
			services := m.dataBridge.ServiceInfos()
			widget.SetService(services[index])
			widget.DataRefresh()
			m.logger.Printf("服务%d 数据已更新\n", index+1)
		})
		// 处理配置更新
		m.setting.ConnectConfigChanged(widget.UpdateConfig)

		m.tab.AddTab(widget, fmt.Sprintf("服务%d", i+1))
		m.summary = append(m.summary, widget)
		m.logger.Printf("已添加综合信息面板：服务%d\n", i+1)
	}
	m.tab.AddTab(m.setting, "设置")
}
