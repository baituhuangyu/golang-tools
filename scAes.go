package common

import (
	"bytes"
	"crypto/cipher"
	"crypto/aes"
	"encoding/base64"
	"strings"
)

// AES 加密,采用CBC 加密方式
// origData 需要加密的原始数据, []byte 类型
// key 加密秘钥, []byte 类型
// iv 加密偏移量, []byte 类型
func AesEncrypt(origData, key []byte, iv []byte) ([]byte, error) {
	// 每次必须new一个新的block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AES 解密,采用CBC 加密方式
// crypted 需要解密的原始数据, []byte 类型
// key 解密秘钥, []byte 类型
// iv 解密偏移量, []byte 类型
func AesDecrypt(crypted, key []byte, iv []byte) ([]byte, error) {
	// 每次必须new一个新的block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

// Padding, 采用PKCS5Padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// UnPadding, 采用PKCS5Padding
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length - 1])
	return origData[:(length - unpadding)]
}

// 加密部分, Cc = ChinaClearing, 中金支付, 采用AES-128 CBC pad5加密，结果采用base64加密
func ScAesEncrypt(rawByte []byte, encodeKey string, aesCbcIv string) (string, error) {
	/*
	 AES-128。key长度：16
	 encodeKey 加密/解密秘钥
	 aesCbcIv CBC加密时的偏移量
	*/
	key := []byte(encodeKey)
	rst, err := AesEncrypt(rawByte, key, []byte(aesCbcIv))
	if err != nil {
		return "", err
	}
	base64rst := base64.StdEncoding.EncodeToString(rst)
	return base64rst, nil
}

// 解密部分, Cc = ChinaClearing, 中金支付, 采用AES-128 CBC pad5解密，输入为加密之后的base64的字符串
func ScAesDecrypt(rawStr string, encodeKey string, aesCbcIv string) (string, error) {
	/*
	 AES-128。key长度：16
	 encodeKey 加密/解密秘钥
	 aesCbcIv CBC加密时的偏移量
	*/
	key := []byte(encodeKey)
	debase64result, err := base64.StdEncoding.DecodeString(rawStr)
	if err != nil {
		return "", err
	}
	origData, err := AesDecrypt(debase64result, key, []byte(aesCbcIv))
	if err != nil {
		return "", err
	}
	return string(origData), nil
}

// 加密 字符串加密 AES 加密后进行 Base64 转码
func AesStringBase64(rawString string, encodeKey string, iv string) (string, error){
	return ScAesEncrypt([]byte(rawString), encodeKey, iv)
}

// 输入输出都是Sting
// 如果需要在url中传输，需要注意+/符号，这2个符号有时会引起一些异常。
//简单做法可以在标准base64后将+/换成-*，然后需要decode的时候，先将-*换成+/，再进行decode
func Base64EncodeSafeUrl(rawString string) string {
	unsafeStr := Base64EncodeStr2Str(rawString)
	unsafeStr = strings.Replace(unsafeStr, "+", "-", -1)
	safeStr := strings.Replace(unsafeStr, "/", "*", -1)
	return safeStr
}

// 输入输出都是Sting
// 如果需要在url中传输，需要注意+/符号，这2个符号有时会引起一些异常。
//简单做法可以在标准base64后将+/换成-*，然后需要decode的时候，先将-*换成+/，再进行decode
func Base64DecodeSafeUrl(rawString string) (string, error) {
	rawString = strings.Replace(rawString, "-", "+", -1)
	rawString = strings.Replace(rawString, "*", "/", -1)
	return Base64DecodeStr2Str(rawString)
}

// base64 encode string to sting
func Base64EncodeStr2Str(rawString string) string {
	return base64.StdEncoding.EncodeToString([]byte(rawString))
}

// base64 string to sting
func Base64DecodeStr2Str(rawString string) (string, error) {
	debase64result, err := base64.StdEncoding.DecodeString(rawString)
	if err != nil {
		return "", err
	}
	return string(debase64result), nil
}
