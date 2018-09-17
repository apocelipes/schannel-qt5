package widgets

import (
	"log"
	"net/http"

	"github.com/therecipe/qt/widgets"

	"schannel-qt5/config"
	"schannel-qt5/crawler"
)

type LoginWidget struct {
	widgets.QWidget

	// loginUser 将登录所用的用户名传递给父控件
	_ func(string)         `signal:"loginUser"`
	_ func()               `signal:"loginFailed,auto"`
	_ func([]*http.Cookie) `signal:"logined"`

	username    *widgets.QLineEdit
	password    *widgets.QLineEdit
	loginStatus *ColorLabel
	remember    *widgets.QCheckBox
	conf        *config.UserConfig
	logger      *log.Logger
}

func NewLoginWidget2(conf *config.UserConfig, logger *log.Logger) *LoginWidget {
	if conf == nil || logger == nil {
		return nil
	}

	widget := NewLoginWidget(nil, 0)
	widget.conf = conf
	widget.logger = logger
	widget.InitUI()

	return widget
}

func (l *LoginWidget) InitUI() {
	userLabel := widgets.NewQLabel2("&username:", nil, 0)
	l.username = widgets.NewQLineEdit(nil)
	l.username.SetPlaceholderText("邮箱")
	if l.conf.UserName != "" {
		l.username.SetText(l.conf.UserName)
	}
	userLabel.SetBuddy(l.username)
	userInputLayout := widgets.NewQHBoxLayout()
	userInputLayout.AddWidget(userLabel, 0, 0)
	userInputLayout.AddWidget(l.username, 0, 0)

	passwdLabel := widgets.NewQLabel2("&password:", nil, 0)
	l.password = widgets.NewQLineEdit(nil)
	l.password.SetPlaceholderText("密码")
	l.password.SetEchoMode(widgets.QLineEdit__Password)
	if l.conf.Passwd != "" {
		l.password.SetText(l.conf.Passwd)
	}

	passwdLabel.SetBuddy(l.password)
	passwdInputLayout := widgets.NewQHBoxLayout()
	passwdInputLayout.AddWidget(passwdLabel, 0, 0)
	passwdInputLayout.AddWidget(l.password, 0, 0)

	l.loginStatus = NewColorLabelWithColor("用户名或密码错误，请重试", "red")
	l.loginStatus.Hide()

	l.remember = widgets.NewQCheckBox2("记住用户名和密码", nil)
	echoCheck := widgets.NewQCheckBox2("密码可见", nil)
	echoCheck.ConnectClicked(func(_ bool) {
		if echoCheck.IsChecked() {
			l.password.SetEchoMode(widgets.QLineEdit__Normal)
			return
		}

		l.password.SetEchoMode(widgets.QLineEdit__Password)
	})
	loginButton := widgets.NewQPushButton2("登录", nil)
	loginButton.ConnectClicked(l.checkLogin)

	checkLayout := widgets.NewQHBoxLayout()
	checkLayout.AddWidget(echoCheck, 0, 0)
	checkLayout.AddWidget(l.remember, 0, 0)
	loginLayout := widgets.NewQHBoxLayout()
	loginLayout.AddStretch(0)
	loginLayout.AddWidget(loginButton, 0, 0)

	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.AddWidget(l.loginStatus, 0, 0)
	mainLayout.AddLayout(userInputLayout, 0)
	mainLayout.AddLayout(passwdInputLayout, 0)
	mainLayout.AddLayout(checkLayout, 0)
	mainLayout.AddLayout(loginLayout, 0)
	l.SetLayout(mainLayout)
}

func (l *LoginWidget) checkLogin(_ bool) {
	// 防止多次点击登录按钮或在登录时改变lineedit内容
	l.SetEnabled(false)
	defer l.SetEnabled(true)

	passwd := l.password.Text()
	user := l.username.Text()
	if user == "" || passwd == "" {
		l.LoginFailed()
		return
	}

	// 登录
	cookies, err := crawler.GetAuth(user, passwd, l.conf.Proxy.String())
	if err != nil {
		l.logger.Printf("crawler failed: %v\n", err)
		l.LoginFailed()
		return
	}

	// 登陆成功，记住密码
	if l.remember.IsChecked() &&
		(passwd != l.conf.Passwd || user != l.conf.UserName) {
		l.conf.Passwd = passwd
		l.conf.UserName = user
		if err := l.conf.StoreConfig(); err != nil {
			log.Fatalf("set user and password failed: %v\n", err)
		}
	}

	// 传递登录信息
	l.logger.Printf("logined as [%s] success\n", user)
	l.LoginUser(user)
	l.Logined(cookies)
}

func (l *LoginWidget) loginFailed() {
	if l.loginStatus.IsHidden() {
		l.loginStatus.Show()
	}
}
