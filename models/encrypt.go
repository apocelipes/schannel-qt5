package models

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
)

// genKey 根据用户名生成key
func genKey(user string) []byte {
	salt := user[:len(user)/2] + "models"
	data := salt[:len(salt)/2] + user + salt[len(salt)/2:]
	hash := md5.New()
	return hash.Sum([]byte(data))[:des.BlockSize]
}

// encryptPassword 加密用户名密码，返回加密后的数据
func encryptPassword(user string, password []byte) ([]byte, error) {
	key := genKey(user)
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	origData := PKCS5Padding(password, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// decryptPassword 返回解密后的信息
func decryptPassword(user string, crypted []byte) ([]byte, error) {
	key := genKey(user)
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

// PKCS5Padding 将数据填充至合适的大小，以便加密算法处理
func PKCS5Padding(origData []byte, blockSize int) []byte {
	padding := blockSize - len(origData)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(origData, padtext...)
}

// PKCS5UnPadding 去除填充
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
