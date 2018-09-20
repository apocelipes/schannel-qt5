package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	Name   string `gorm:"type:varchar(255);not null;unique;primary_key"`
	Passwd []byte `gorm:"type:blob"`
}

// GetUserPassword 获取用户名以及密码
func GetUserPassword(db *gorm.DB, user string) (*User, error) {
	u := User{Name: user}
	if err := db.Where(&u).First(&u).Error; err != nil {
		return nil, err
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
func SetUserPassword(db *gorm.DB, user string, password []byte) error {
	u := User{
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

	tmp := User{}
	if err := db.Where("name = ?", u.Name).First(&tmp).Error; err == nil {
		db.Model(&u).Update(&u)
		return nil
	}

	if err := db.Create(&u).Error; err != nil {
		return err
	}

	return nil
}

// GetAllUsers 返回所有user，包括未
func GetAllUsers(db *gorm.DB) ([]*User, error) {
	users := make([]*User, 0)
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
