package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
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

	// IV 생성 (12바이트)
	iv := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// AES-GCM 생성
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 암호화
	cipherText := aesGCM.Seal(nil, iv, []byte(plainText), nil)

	// IV + 암호화된 데이터를 Base64로 인코딩하여 반환
	finalCipher := append(iv, cipherText...)
	return base64.StdEncoding.EncodeToString(finalCipher), nil
}

// 복호화 함수
func DecryptAES(cipherText string) (string, error) {
	block, err := aes.NewCipher(AESKey)
	if err != nil {
		return "", err
	}

	// AES-GCM 생성
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Base64 디코딩
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	// IV와 암호화된 데이터 분리
	iv := data[:12] // AES-GCM에서 권장되는 12바이트 IV
	cipherTextBytes := data[12:]

	// 복호화
	plainText, err := aesGCM.Open(nil, iv, cipherTextBytes, nil)
	if err != nil {
		return "", errors.New("failed to decrypt")
	}

	return string(plainText), nil
}
