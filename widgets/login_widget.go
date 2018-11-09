package widgets

import (
	"log"
	"net/http"

	"github.com/go-xorm/xorm"
	"github.com/therecipe/qt/widgets"

	"schannel-qt5/config"
	"schannel-qt5/crawler"
	"schannel-qt5/models"
)

// LoginWidget 登录界面
type LoginWidget struct {
	widgets.QWidget

	// loginFailed 登录失败显示错误信息
	// loginSuccess 将登录成功的用户名和cookies传递给父控件
	_ func(string)                 `signal:"loginFailed,auto"`
	_ func(string, []*http.Cookie) `signal:"loginSuccess"`

	username    *widgets.QComboBox
	password    *widgets.QLineEdit
	loginStatus *ColorLabel
	remember    *widgets.QCheckBox

	// 用户数据
	conf   *config.UserConfig
	logger *log.Logger
	db     *xorm.Engine
}

// NewLoginWidget2 根据config，logger，db生成登录控件
func NewLoginWidget2(conf *config.UserConfig, logger *log.Logger, db *xorm.Engine) *LoginWidget {
	if conf == nil || logger == nil {
		return nil
	}

	widget := NewLoginWidget(nil, 0)
	widget.conf = conf
	widget.logger = logger
	widget.db = db
	widget.InitUI()

	return widget
}

func (l *LoginWidget) InitUI() {
	l.username = widgets.NewQComboBox(nil)
	l.username.SetEditable(true)
	users, err := models.GetAllUsers(l.db)
	if err != nil {
		l.logger.Fatalln(err)
	}
	// 第一项为空
	names := make([]string, 1, len(users)+1)
	for _, v := range users {
		names = append(names, v.Name)
	}
	l.username.AddItems(names)
	// 实现记住用户密码
	l.username.ConnectCurrentTextChanged(l.setPassword)

	l.password = widgets.NewQLineEdit(nil)
	l.password.SetPlaceholderText("密码")
	l.password.SetEchoMode(widgets.QLineEdit__Password)

	// 空的ColorLabel，预备填充错误信息
	l.loginStatus = NewColorLabelWithColor("", "red")
	l.loginStatus.Hide()

	l.remember = widgets.NewQCheckBox2("记住用户名和密码", nil)
	loginButton := widgets.NewQPushButton2("登录", nil)
	loginButton.ConnectClicked(l.checkLogin)

	loginLayout := widgets.NewQHBoxLayout()
	loginLayout.AddStretch(0)
	loginLayout.AddWidget(loginButton, 0, 0)

	mainLayout := widgets.NewQFormLayout(nil)
	mainLayout.AddRow5(l.loginStatus)
	mainLayout.AddRow3("用户名：", l.username)
	mainLayout.AddRow3("密码：", l.password)
	mainLayout.AddRow6(loginLayout)
	l.SetLayout(mainLayout)
}

func (l *LoginWidget) checkLogin(_ bool) {
	// 防止多次点击登录按钮或在登录时改变lineedit内容
	l.SetEnabled(false)
	defer l.SetEnabled(true)

	passwd := l.password.Text()
	user := l.username.CurrentText()
	if user == "" || passwd == "" {
		l.LoginFailed("用户名/密码不能为空")
		return
	}

	// 登录
	cookies, err := crawler.GetAuth(user, passwd, l.conf.Proxy.String())
	if err != nil {
		l.logger.Printf("crawler failed: %v\n", err)
		l.LoginFailed("用户名或密码错误")
		return
	}

	// 登陆成功，记住密码
	if l.remember.IsChecked() {
		if err := models.SetUserPassword(l.db, user, []byte(passwd)); err != nil {
			l.logger.Println(err)
		}
	} else {
		if has, err := l.db.Where("password is not null").Exist(&models.User{Name: user}); err != nil {
			l.logger.Println(err)
		} else if has {
			if err := models.DelPassword(l.db, user); err != nil {
				l.logger.Printf("delete %v password failed: %v\n", user, err)
			}
		}
	}

	// 传递登录信息
	l.logger.Printf("logined as [%s] success\n", user)
	l.LoginSuccess(user, cookies)
}

// 更新并显示错误信息
func (l *LoginWidget) loginFailed(errInfo string) {
	l.loginStatus.SetDefaultColorText(errInfo)
	if l.loginStatus.IsHidden() {
		l.loginStatus.Show()
	}
}

// setPassword 将密码不为null的用户显示
func (l *LoginWidget) setPassword(user string) {
	info, err := models.GetUserPassword(l.db, user)
	if err != nil {
		// 输入的用户名未记录
		l.logger.Println(err)
		l.username.SetCurrentText(user)
	} else if info.Passwd != nil {
		l.password.SetText(string(info.Passwd))
		l.remember.SetChecked(true)
		return
	}

	// 记录了用户但是没记录密码
	l.password.SetText("")
	l.remember.SetChecked(false)
}
