package models

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"

	"github.com/go-xorm/xorm"

	"schannel-qt5/config"
)

const (
	// 数据库存放路径
	databasePath = ".local/share/schannel-users.db"
)

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
	Name   string `xorm:"pk"`
	Passwd []byte `xorm:"blob"`
}

// GetUserPassword 获取用户名以及密码
func GetUserPassword(db *xorm.Engine, user string) (*User, error) {
	u := User{Name: user}
	if has, err := db.Get(&u); err != nil {
		return nil, err
	} else if !has {
		return nil, xorm.ErrNotExist
	}

	// 密码解密
	if u.Passwd != nil {
		data, err := decryptPassword(u.Name, u.Passwd)
		if err != nil {
			return nil, err
		}
		u.Passwd = data
	}

	return &u, nil
}

// SetUserPassword 将用户名密码保存
// password为nil表示不记住密码
func SetUserPassword(db *xorm.Engine, user string, password []byte) error {
	u := &User{
		Name:   user,
		Passwd: password,
	}

	if u.Passwd != nil {
		data, err := encryptPassword(u.Name, u.Passwd)
		if err != nil {
			return err
		}
		u.Passwd = data
	}

	if has, err := db.Exist(&User{Name: u.Name}); err != nil {
		return err
	} else if has {
		old := &User{Name: u.Name}
		db.Get(&old)
		// 和旧值一样，不更新，返回error
		if bytes.Equal(u.Passwd, old.Passwd) {
			return errors.New("insert same values")
		}
		db.Where("name = ?", u.Name).Cols("passwd").Update(u)
		return nil
	}

	if _, err := db.InsertOne(u); err != nil {
		return err
	}

	return nil
}

// GetAllUsers 返回所有user，包括未
func GetAllUsers(db *xorm.Engine) ([]*User, error) {
	users := make([]*User, 0)
	if err := db.Find(&users); err != nil {
		return nil, err
	}

	return users, nil
}

// DelPassword 将指定user的password设置为null
func DelPassword(db *xorm.Engine, user string) error {
	u := &User{Name: user}
	_, err := db.Where("name = ?", user).Cols("passwd").Update(u)
	if err != nil {
		return err
	}

	return nil
}
