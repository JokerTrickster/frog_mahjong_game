package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	_aws "main/utils/aws"
)

var AESKey []byte

func InitCrypto() error {
	// aws에서 aes 키를 가져온다,
	key, err := _aws.AwsSsmGetParam("frog_aes_key")
	if err != nil {
		panic(err)
	}
	AESKey = []byte(key)
	return nil
}

// 암호화 함수
func EncryptAES(plainText string) (string, error) {
	block, err := aes.NewCipher(AESKey)
	if err != nil {
		return "", err
	}

	// 초기화 벡터(IV) 생성
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// 암호화
	cipherText := make([]byte, len(plainText))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText, []byte(plainText))

	// IV + 암호화된 텍스트를 Base64로 인코딩하여 반환
	finalCipher := append(iv, cipherText...)
	return base64.StdEncoding.EncodeToString(finalCipher), nil
}

// 복호화 함수
func DecryptAES(cipherText string) (string, error) {
	block, err := aes.NewCipher(AESKey)
	if err != nil {
		return "", err
	}

	// Base64 디코딩
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	// IV와 암호화된 데이터 분리
	iv := data[:aes.BlockSize]
	cipherTextBytes := data[aes.BlockSize:]

	// 복호화
	plainText := make([]byte, len(cipherTextBytes))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plainText, cipherTextBytes)

	return string(plainText), nil
}
