package main

import (
	"log"
	"os"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	std_widgets "github.com/therecipe/qt/widgets"

	"schannel-qt5/config"
	"schannel-qt5/models"
	_ "schannel-qt5/pyclient"
	"schannel-qt5/ssr"
	"schannel-qt5/widgets"
)

const (
	// 日志主前缀
	prefix = "schannel-qt5: "
)

func init() {
	dbPath, err := models.GetDBPath()
	if err != nil {
		panic(err)
	}
	orm.RegisterDataBase("default", "sqlite3", dbPath)
	err = orm.RunSyncdb("default", false, false)
	if err != nil {
		panic(err)
	}
}

func main() {
	app := std_widgets.NewQApplication(len(os.Args), os.Args)

	// 初始化用户配置
	conf := &config.UserConfig{}
	conf.SSRClientConfig = ssr.NewClientConfig("python")
	err := conf.LoadConfig()
	if err != nil {
		panic(err)
	}

	logPath, err := conf.LogFile.AbsPath()
	if err != nil {
		panic(err)
	}
	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	logger := log.New(logFile, prefix, log.LstdFlags|log.Lshortfile)

	// 获取用户数据库连接
	db := orm.NewOrm()

	// 初始化GUI
	mainWindow := widgets.NewMainWidget2(conf, logger, db)
	mainWindow.Show()

	app.Exec()
}
