package utils

import (
	"fmt"
	"main/utils/db"
)

func InitServer() error {
	if err := InitEnv(); err != nil {
		fmt.Sprintf("서버 에러 발생 : %s", err.Error())
		return err
	}

	if err := InitJwt(); err != nil {
		fmt.Sprintf("jwt 초기화 에러 : %s", err.Error())
		return err
	}

	if err := db.InitMySQL(); err != nil {
		fmt.Sprintf("db 초기화 에러 : %s", err.Error())
		return err
	}

	return nil
}
