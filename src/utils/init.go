package utils

import (
	"fmt"
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

	// if err := mysql.InitMySQL(); err != nil {
	// 	fmt.Sprintf("db 초기화 에러 : %s", err.Error())
	// 	return err
	// }

	return nil
}
