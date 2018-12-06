package widgets

import (
	"fmt"
	"log"
	"net/http"

	"github.com/astaxie/beego/orm"
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

	username     *widgets.QComboBox
	password     *widgets.QLineEdit
	loginStatus  *ColorLabel
	showPassword *widgets.QCheckBox
	remember     *widgets.QCheckBox
	loginButton  *widgets.QPushButton
	indicator    *LoginIndicator

	// 用户数据
	conf   *config.UserConfig
	logger *log.Logger
	db     orm.Ormer
}

// NewLoginWidget2 根据config，logger，db生成登录控件
func NewLoginWidget2(conf *config.UserConfig, logger *log.Logger, db orm.Ormer) *LoginWidget {
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

	names := make([]string, 0, len(users))
	for _, v := range users {
		names = append(names, v.Name)
	}
	userList := widgets.NewQListWidget(nil)
	// 设置combobox代理
	l.username.SetModel(userList.Model())
	l.username.SetView(userList)

	maxViewWidth := 0
	for _, name := range names {
		accountItem := NewAccountItem2(name)
		// 设置下拉框宽度与宽度最大item一致
		if accountItem.SizeHint().Width() > maxViewWidth {
			maxViewWidth = accountItem.SizeHint().Width()
			l.username.View().SetFixedWidth(maxViewWidth)
		}
		// ComboBox处理item被选中和点击删除按钮
		accountItem.ConnectShowAccount(func(userName string) {
			l.username.HidePopup()
			for i, v := range names {
				if v == userName {
					l.username.SetCurrentIndex(i)
					l.username.SetCurrentText(userName)
					break
				}
			}
		})

		accountItem.ConnectRemoveAccount(func(userName string) {
			l.username.HidePopup()

			info := fmt.Sprintf("将删除用户：%s (同时删除使用数据)", userName)
			buttons := widgets.QMessageBox__Yes | widgets.QMessageBox__Cancel
			defaultButton := widgets.QMessageBox__Yes
			answer := widgets.QMessageBox_Question4(l, "是否删除记录", info, buttons, defaultButton)
			if answer != int(widgets.QMessageBox__Yes) {
				return
			}

			// listWidget中的顺序和names一致
			for i, v := range names {
				if userName == v {
					userList.TakeItem(i)
					break
				}
			}

			err := models.DelUser(l.db, userName)
			if err != nil {
				l.logger.Fatal("删除用户失败:", userName, err)
			}
		})

		listItem := widgets.NewQListWidgetItem(userList, 0)
		userList.SetItemWidget(listItem, accountItem)
	}

	l.password = widgets.NewQLineEdit(nil)
	l.password.SetPlaceholderText("密码")
	l.password.SetEchoMode(widgets.QLineEdit__Password)

	// 勾选是否明文显示密码
	l.showPassword = widgets.NewQCheckBox2("显示密码", nil)
	l.showPassword.SetChecked(false)
	l.showPassword.ConnectClicked(func(_ bool) {
		if l.showPassword.IsChecked() {
			l.password.SetEchoMode(widgets.QLineEdit__Normal)
			return
		}

		l.password.SetEchoMode(widgets.QLineEdit__Password)
	})

	l.remember = widgets.NewQCheckBox2("记住用户名和密码", nil)
	// 设置第一个记录用户的密码
	// 因为comboBox默认选择显示第一个name，不会触发信号
	if len(users) != 0 {
		info, err := models.GetUserPassword(l.db, names[0])
		if err != nil {
			l.logger.Println(err)
		} else if info.Passwd != "" {
			// 密码不为空，设置密码和选中记住密码
			l.password.SetText(string(info.Passwd))
			l.username.SetCurrentText(info.Name)
			l.remember.SetChecked(true)
		}
	}
	// 实现记住用户密码
	l.username.ConnectCurrentTextChanged(l.setPassword)

	// 空的ColorLabel，预备填充错误信息
	l.loginStatus = NewColorLabelWithColor("", "red")
	l.loginStatus.Hide()

	// login时显示busy进度条
	l.indicator = NewLoginIndicator2()
	l.indicator.Hide()

	l.loginButton = widgets.NewQPushButton2("登录", nil)
	l.loginButton.ConnectClicked(l.login)

	loginLayout := widgets.NewQHBoxLayout()
	loginLayout.AddWidget(l.remember, 0, 0)
	loginLayout.AddStretch(0)
	loginLayout.AddWidget(l.loginButton, 0, 0)

	mainLayout := widgets.NewQFormLayout(nil)
	mainLayout.AddRow5(l.loginStatus)
	mainLayout.AddRow3("用户名：", l.username)
	mainLayout.AddRow3("密码：", l.password)
	mainLayout.AddRow5(l.showPassword)
	mainLayout.AddRow6(loginLayout)
	mainLayout.AddRow5(l.indicator)
	l.SetLayout(mainLayout)
}

func (l *LoginWidget) login(_ bool) {
	l.indicator.Show()
	l.setEditAreaEnabled(false)

	go l.checkLogin()
}

// 控制输入区是否可编辑，禁止用户在登录过程中影响输入信息
func (l *LoginWidget) setEditAreaEnabled(enabled bool) {
	l.username.SetEnabled(enabled)
	l.password.SetEnabled(enabled)
	l.showPassword.SetEnabled(enabled)
	l.remember.SetEnabled(enabled)
	l.loginButton.SetEnabled(enabled)
}

// checkLogin 请求登录，用户名密码正确则登陆成功
// 勾选了remember时将会更新记录进数据库
// 登录失败显示失败信息
func (l *LoginWidget) checkLogin() {
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
		if err := models.SetUserPassword(l.db, user, passwd); err != nil {
			l.logger.Println(err)
		}
	} else {
		// 如果未勾选，表示用户不想记住密码，已经记住的将会被设置为null
		cond := orm.NewCondition()
		cond = cond.And("Passwd__isnull", false).And("Name", user)
		if l.db.QueryTable(&models.User{}).SetCond(cond).Exist() {
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
	l.indicator.Hide()
	l.setEditAreaEnabled(true)

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
	} else if info.Passwd != "" {
		l.password.SetText(string(info.Passwd))
		l.remember.SetChecked(true)
		return
	}

	// 记录了用户但是没记录密码
	l.password.SetText("")
	l.remember.SetChecked(false)
}
