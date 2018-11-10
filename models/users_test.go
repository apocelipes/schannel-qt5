package models

import (
	"testing"

	"bytes"
	"crypto/rand"
	"encoding/hex"
	"os"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// 测试数据库存储路径
	dbPath = "/tmp/db_users_test.db"
)

// initDB 初始化测试数据
func initDB(t *testing.T) (*xorm.Engine, []*User) {
	users := []*User{
		{
			Name:   "test@test.com",
			Passwd: nil,
		},
		{
			Name:   "test@example.com",
			Passwd: genPassword(),
		},
		{
			Name:   "example",
			Passwd: genPassword(),
		},
	}

	os.Remove(dbPath)
	var err error
	db, err := xorm.NewEngine("sqlite3", dbPath)
	if err != nil {
		t.Errorf("initDB: %v\n", err)
	}

	db.Sync2(&User{})
	for _, v := range users {
		err := SetUserPassword(db, v.Name, v.Passwd)
		if err != nil {
			t.Errorf("initDB: %v\n", err)
		}
	}

	return db, users
}

// genPassword 生成随机密码
func genPassword() []byte {
	origData := make([]byte, 16)
	n, err := rand.Read(origData)
	if err != nil {
		panic(err)
	}

	pw := make([]byte, hex.EncodedLen(n))
	hex.Encode(pw, origData[:n])
	return pw
}

func TestGetUserPassword(t *testing.T) {
	db, users := initDB(t)
	defer db.Close()

	for _, v := range users {
		user, err := GetUserPassword(db, v.Name)
		if err != nil {
			t.Errorf("get user: %s error: %v\n", v.Name, err)
		}
		if !bytes.Equal(user.Passwd, v.Passwd) {
			format := "get user: %s password different\nhave: %v\n\twant: %v\n"
			t.Errorf(format, v.Name, user.Passwd, v.Passwd)
		}
	}
}

func TestSetUserPassword(t *testing.T) {
	db, _ := initDB(t)
	defer db.Close()

	testData := []*struct {
		// 用户对象
		u *User
		// 是否insert成功
		inserted bool
	}{
		{
			u: &User{
				Name:   "example",
				Passwd: nil,
			},
			inserted: false,
		},
		{
			u: &User{
				Name:   "a",
				Passwd: genPassword(),
			},
			inserted: true,
		},
		{
			u: &User{
				Name:   "b",
				Passwd: nil,
			},
			inserted: true,
		},
	}

	for _, v := range testData {
		err := SetUserPassword(db, v.u.Name, v.u.Passwd)
		if (err == nil) != v.inserted {
			t.Errorf("set user: %v error: %v\n", v.u, err)
		}
	}
}

func TestGetAllUsers(t *testing.T) {
	db, users := initDB(t)
	defer db.Close()

	u, err := GetAllUsers(db)
	if err != nil {
		t.Errorf("get all users error: %v\n", err)
	}

	// 取得的数据量是否相同
	if len(u) != len(users) {
		format := "length error: have: %v\n\twant: %v\n"
		t.Errorf(format, len(u), len(users))
	}
}

func TestDelPassword(t *testing.T) {
	db, users := initDB(t)
	defer db.Close()

	for _, v := range users {
		if v.Passwd != nil {
			err := DelPassword(db, v.Name)
			if err != nil {
				format := "del %s password error: %v\n"
				t.Errorf(format, v.Name, err)
			}

			// 查看是否已将密码设置为null
			u, err := GetUserPassword(db, v.Name)
			if err != nil {
				format := "del %s password error: %v\n"
				t.Errorf(format, v.Name, err)
			}
			if u.Passwd != nil {
				t.Errorf("del password failed: %v\n", u.Passwd)
			}
		}
	}
}

func TestGetDBPath(t *testing.T) {
	testData := []*struct {
		// 设置环境变量HOME的值
		home string
		res  string
	}{
		{
			home: "/home/test",
			res:  "/home/test/" + databasePath,
		},
		{
			home: "/home/test/",
			res:  "/home/test/" + databasePath,
		},
	}

	for _, v := range testData {
		err := os.Setenv("HOME", v.home)
		if err != nil {
			t.Fatalf("无法设置$HOME: %v\n", err)
		}
		res, err := GetDBPath()
		if err != nil {
			t.Fatalf("获取DB Path错误：%v\n", err)
		}
		if v.res != res {
			format := "不正确的DB Path:\n\twant: %s\n\thave: %v\n"
			t.Errorf(format, v.res, res)
		}
	}
}
