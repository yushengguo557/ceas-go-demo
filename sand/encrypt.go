package sand

import (
	"bytes"
	"ceas-go-demo/utils"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
)

// AESEncrypt AES加密(AES/ECB/PKCS5Padding)  data: 原始数据
func AESEncrypt(data []byte) (string, string, error) {
	// 1.加密数据
	// 1.1 生成16位AES密钥
	aesKey, err := utils.RandomBytes(16)
	if err != nil {
		return "", "", err
	}

	// 1.2 创建 AES 加密块
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", "", err
	}

	// Generate a random IV with the same size as the block
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", "", err
	}

	// 1.3 对原始数据进行填充
	padding := aes.BlockSize - (len(data) % aes.BlockSize)
	paddedData := append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)

	// 1.4 创建 ECB 分组模式
	mode := cipher.NewCBCEncrypter(block, iv)

	// 1.5 加密数据
	cipherText := make([]byte, len(paddedData))
	mode.CryptBlocks(cipherText, paddedData)

	// 1.6 Base64编码
	encryptValue := base64.StdEncoding.EncodeToString(cipherText)

	// 2.对AES Key 进行RSA加密
	// 2.1 加载公钥
	rsaPublicKey, err := LoadPublicKey("./cert/sand_public.cer")
	if err != nil {
		return "", "", fmt.Errorf("load public key, err: %w", err)
	}
	// 2.2 公钥加密
	encryptKeyBytes, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, aesKey)
	if err != nil {
		return "", "", fmt.Errorf("aes encrypt, err: %w", err)
	}
	// 2.3 Base64编码
	sandEncryptKey := base64.StdEncoding.EncodeToString(encryptKeyBytes)

	// 3. 返回加密后的结果 包括value和key
	return encryptValue, sandEncryptKey, nil
}

// AESEncrypt AES加密
//func AESEncrypt(data []byte) (string, string, error) {
//	// 生成16字节的AES密钥
//	aesKey, err := utils.RandomBytes(16)
//	if err != nil {
//		return "", "", err
//	}
//
//	// 对数据进行AES加密
//	//plainValue, err := json.Marshal(data)
//	//if err != nil {
//	//	return "", "", err
//	//}
//	block, err := aes.NewCipher(aesKey)
//	if err != nil {
//		return "", "", err
//	}
//	paddedValue := Pad(data, block.BlockSize())
//	iv := make([]byte, aes.BlockSize)
//	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
//		return "", "", err
//	}
//	mode := cipher.NewCBCEncrypter(block, iv)
//	encryptValueBytes := make([]byte, len(paddedValue))
//	mode.CryptBlocks(encryptValueBytes, paddedValue)
//	encryptValue := base64.StdEncoding.EncodeToString(encryptValueBytes)
//
//	// 对AES Key 进行RSA加密
//	rsaPublicKey, err := LoadPublicKey("./cert/sand_public.cer")
//	if err != nil {
//		return "", "", fmt.Errorf("load public key, err: %w", err)
//	}
//	encryptKeyBytes, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, aesKey)
//	if err != nil {
//		return "", "", fmt.Errorf("aes encrypt, err: %w", err)
//	}
//	sandEncryptKey := base64.StdEncoding.EncodeToString(encryptKeyBytes)
//
//	// 返回加密后的结果 包括value和key
//	return encryptValue, sandEncryptKey, nil
//}
//
//// Pad 对数据进行填充
//func Pad(data []byte, blockSize int) []byte {
//	padding := blockSize - len(data)%blockSize
//	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
//	return append(data, padtext...)
//}
