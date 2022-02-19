package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"strconv"
)

//Advanced Encryption Standard, AES
const KeyFile = "common/AESKEY"

// PKCS7 padding mode
func PKCS7Padding(cipherText []byte, blockSize int) []byte {
	padCnt := blockSize - len(cipherText)%blockSize
	padding := bytes.Repeat([]byte{byte(padCnt)}, padCnt)
	return append(cipherText, padding...)
}

func PKCS7UnPadding(origData []byte) ([]byte, error) {
	// get data length
	length := len(origData)
	if length == 0 {
		return nil, errors.New("wrong length of string")
	} else {
		unpadding := int(origData[length-1])
		if length < unpadding {
			return nil, errors.New(
				"Abnormal padding count : " + strconv.Itoa(unpadding) +
					" while pwd length: " + strconv.Itoa(length))
		}
		// get original slice, delete the padded bytes
		return origData[:(length - unpadding)], nil
	}
}

func AesEcrypt(origData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// get the block size, the block size would be 128bit given our key of length 16
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	// use CBC ecrypt method
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDeCrypt(crypted []byte, key []byte) ([]byte, error) {
	// create instance of aes
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// get the block size
	blockSize := block.BlockSize()
	// create an instance of decipher
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	// remove padding
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, err
}

// EnPwdCode ecrypts the password with base64
func EnPwdCode(pwd []byte) (string, error) {
	aeskeys, err := ReadPrivateFile(KeyFile)
	if err != nil {
		return "", err
	}
	key := []byte(aeskeys[0])
	result, err := AesEcrypt(pwd, key)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(result), err
}

func DePwdCode(pwd string) ([]byte, error) {
	// base64 decoding
	pwdByte, err := base64.URLEncoding.DecodeString(pwd)
	if err != nil {
		return nil, err
	}
	aeskeys, err := ReadPrivateFile(KeyFile)
	if err != nil {
		return nil, err
	}
	key := []byte(aeskeys[0])
	return AesDeCrypt(pwdByte, key)
}
