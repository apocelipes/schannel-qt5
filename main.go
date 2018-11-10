package main

import (
	"log"
	"os"

	"github.com/go-xorm/xorm"
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
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	logger := log.New(logFile, prefix, log.LstdFlags|log.Lshortfile)

	// 初始化用户数据库
	dbPath, err := models.GetDBPath()
	if err != nil {
		logger.Fatalf("获取数据库存放路径失败: %v\n", err)
	}
	db, err := xorm.NewEngine("sqlite3", dbPath)
	if err != nil {
		logger.Fatalf("数据库初始化失败: %v\n", err)
	}
	defer db.Close()
	if err := db.Sync2(&models.User{}); err != nil {
		logger.Fatalf("数据库同步失败：%v\n", err)
	}

	// 初始化GUI
	mainWindow := widgets.NewMainWidget2(conf, logger, db)
	mainWindow.Show()

	app.Exec()
}
