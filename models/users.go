package models

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"

	"schannel-qt5/config"
)

const (
	// 数据库存放路径
	databasePath = ".local/share/schannel-users.db"
)

// 注册模型
func init() {
	orm.RegisterModel(&User{})
}

// GetDBPath 获取数据库存放路径
func GetDBPath() (string, error) {
	home, exist := os.LookupEnv("HOME")
	if !exist {
		return "", config.ErrHOME
	}

	return filepath.Join(home, databasePath), nil
}

// User 用户表，将和使用量表关联
type User struct {
	Name   string `orm:"size(100);pk"`
	Passwd string `orm:"size(100);null"`
}

// GetUserPassword 获取用户名以及密码
func GetUserPassword(db orm.Ormer, user string) (*User, error) {
	u := &User{Name: user}

	if err := db.Read(u); err != nil {
		return nil, err
	}

	// 密码解密
	if u.Passwd != "" {
		data, err := decryptPassword(u.Name, u.Passwd)
		if err != nil {
			return nil, err
		}
		u.Passwd = data
	}

	return u, nil
}

// SetUserPassword 将用户名密码保存
// password为nil表示不记住密码
func SetUserPassword(db orm.Ormer, user string, password string) error {
	u := &User{
		Name:   user,
		Passwd: password,
	}

	if u.Passwd != "" {
		data, err := encryptPassword(u.Name, u.Passwd)
		if err != nil {
			return err
		}
		u.Passwd = data
	}

	if db.QueryTable(u).Filter("Name", u.Name).Exist() {
		old := &User{Name: u.Name}
		db.QueryTable(old).Filter("Name", old.Name).One(old)
		// 和旧值一样，不更新，返回error
		if u.Passwd == old.Passwd {
			return errors.New("insert same values")
		}
		_, err := db.QueryTable(u).Filter("Name", u.Name).Update(orm.Params{
			"Passwd": u.Passwd,
		})
		if err != nil {
			return err
		}
		return nil
	}

	if _, err := db.Insert(u); err != nil {
		return err
	}

	return nil
}

// GetAllUsers 返回所有user，包括未
func GetAllUsers(db orm.Ormer) ([]*User, error) {
	users := make([]*User, 0)

	if _, err := db.QueryTable(&User{}).All(&users); err != nil {
		return nil, err
	}

	return users, nil
}

// DelPassword 将指定user的password设置为null
func DelPassword(db orm.Ormer, user string) error {
	u := &User{Name: user}

	_, err := db.QueryTable(u).Filter("Name", u.Name).Update(orm.Params{
		"Passwd": "",
	})
	if err != nil {
		return err
	}

	return nil
}

// DelUser 删除名字与name相同的User，同时会删除UserAmount记录
func DelUser(db orm.Ormer, name string) error {
	user := &User{Name: name}
	_, err := db.QueryTable(user).Filter("Name", user.Name).Delete()
	return err
}
