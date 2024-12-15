package utils

import (
	"fmt"
	_aws "main/utils/aws"
	"main/utils/db/mysql"
	_redis "main/utils/db/redis"
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
	if err := _aws.InitAws(); err != nil {
		fmt.Sprintf("aws 초기화 에러 : %s", err.Error())
		return err
	}

	if err := InitGoogleOauth(); err != nil {
		fmt.Sprintf("google oauth 초기화 에러 : %s", err.Error())
		return err
	}
	if err := _redis.InitRedis(); err != nil {
		fmt.Sprintf("redis 초기화 에러 : %s", err.Error())
		return err
	}

	if err := mysql.InitMySQL(); err != nil {
		fmt.Sprintf("db 초기화 에러 : %s", err.Error())
		return err
	}
	if err := InitNotice(); err != nil {
		fmt.Sprintf("notice 초기화 에러 : %s", err.Error())
		return err
	}
	if err := InitRabbitMQ(); err != nil {
		fmt.Sprintf("rabbitmq 초기화 에러 : %s", err.Error())
		return err
	}
	if !Env.IsLocal {
		if err := InitLogging(); err != nil {
			return err
		}
	}
	return nil
}
