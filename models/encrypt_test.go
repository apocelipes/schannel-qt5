package models

import "testing"

func TestEncryptDecrypt(t *testing.T) {
	testData := []*struct {
		user     string
		password string
	}{
		{
			user:     "test",
			password: "abcdefg123",
		},
		{
			user:     "example user",
			password: "foratest1990812",
		},
		{
			user:     "用户A1@",
			password: "worldHello17.,",
		},
		{
			user:     "sample1",
			password: genPassword(),
		},
	}

	for _, v := range testData {
		crypted, err := encryptPassword(v.user, v.password)
		if err != nil {
			t.Errorf("encrypto: %v\n", err)
		}
		password, err := decryptPassword(v.user, crypted)
		if err != nil {
			t.Errorf("decrypto: %v\n", err)
		}
		if password != v.password {
			format := "password not equal, old: %s, new: %s\n"
			t.Errorf(format, v.password, password)
		}
	}
}
