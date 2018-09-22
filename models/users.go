package models

import (
	"github.com/go-xorm/xorm"
)

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
