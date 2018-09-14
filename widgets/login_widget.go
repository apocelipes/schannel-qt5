package widgets

import (
	"net/http"

	"github.com/therecipe/qt/widgets"

	"schannel-qt5/config"
	"schannel-qt5/crawler"
)

type LoginWidget struct {
	widgets.QWidget

	_ func() `constructor:"init"`

	_ func()               `signal:"loginFailed,auto"`
	_ func([]*http.Cookie) `signal:"logined"`

	username    *widgets.QLineEdit
	password    *widgets.QLineEdit
	loginStatus *ColorLabel
	remember    *widgets.QCheckBox
	conf        *config.UserConfig
}

func (l *LoginWidget) init() {
	l.conf = new(config.UserConfig)
	err := l.conf.LoadConfig()
	if err != nil {
		panic(err)
	}

	userLabel := widgets.NewQLabel2("&username:", nil, 0)
	l.username = widgets.NewQLineEdit(nil)
	l.username.SetPlaceholderText("用户名/邮箱")
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
	loginButton := widgets.NewQPushButton2("login", nil)
	loginButton.ConnectClicked(l.checkLogin)
	loginLayout := widgets.NewQHBoxLayout()
	loginLayout.AddWidget(l.remember, 0, 0)
	loginLayout.AddStretch(0)
	loginLayout.AddWidget(loginButton, 0, 0)

	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.AddWidget(l.loginStatus, 0, 0)
	mainLayout.AddLayout(userInputLayout, 0)
	mainLayout.AddLayout(passwdInputLayout, 0)
	mainLayout.AddLayout(loginLayout, 0)
	l.SetLayout(mainLayout)
}

func (l *LoginWidget) checkLogin(_ bool) {
	// 防止多次点击登录按钮或在登录时改变lineedit内容
	l.Layout().SetEnabled(false)
	defer l.Layout().SetEnabled(true)

	passwd := l.password.Text()
	user := l.username.Text()
	if user == "" || passwd == "" {
		l.LoginFailed()
		return
	}
	// 记住密码
	if l.remember.IsChecked() &&
		(passwd != l.conf.Passwd || user != l.conf.UserName) {
		l.conf.Passwd = passwd
		l.conf.UserName = user
		l.conf.StoreConfig()
	}
	// 登录
	cookies, err := crawler.GetAuth(user, passwd, l.conf.Proxy)
	if err != nil {
		l.LoginFailed()
		return
	}
	l.Logined(cookies)
}

func (l *LoginWidget) loginFailed() {
	if l.loginStatus.IsHidden() {
		l.loginStatus.Show()
	}
}
