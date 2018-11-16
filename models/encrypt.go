package models

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
)

// genKey 根据用户名生成key
func genKey(user string) []byte {
	salt := user[:len(user)/2] + "models"
	data := salt[:len(salt)/2] + user + salt[len(salt)/2:]
	hash := sha256.New()
	return hash.Sum([]byte(data))[:aes.BlockSize]
}

// encryptPassword 加密用户名密码，返回加密后的数据
func encryptPassword(user string, password string) (string, error) {
	key := genKey(user)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	origData := PKCS5Padding([]byte(password), block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	res := base64.StdEncoding.EncodeToString(crypted)
	return res, nil
}

// decryptPassword 返回解密后的信息
func decryptPassword(user string, crypted string) (string, error) {
	key := genKey(user)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, key)
	unbase, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return "", err
	}

	n := len(unbase)
	origData := make([]byte, n)
	blockMode.CryptBlocks(origData, unbase)
	origData = PKCS5UnPadding(origData)
	return string(origData), nil
}

// PKCS5Padding 将数据填充至合适的大小，以便加密算法处理
func PKCS5Padding(origData []byte, blockSize int) []byte {
	padding := blockSize - len(origData)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(origData, padText...)
}

// PKCS5UnPadding 去除填充
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
