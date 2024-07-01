package usecase

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSecret(userID uint) string {
	// 32바이트의 랜덤 데이터 생성
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}

	// 바이트 배열을 16진수 문자열로 변환하여 출력
	keyString := hex.EncodeToString(key)
	return keyString
}
