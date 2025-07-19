package _redis

import (
	"context"
	"fmt"
	"log"
	"os"

	_aws "main/utils/aws"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

const RankingKey = "food:rankings"

func InitRedis() error {
	ctx := context.Background()
	isLocal := os.Getenv("IS_LOCAL")
	var connectionString string
	if isLocal == "true" {
		connectionString = fmt.Sprintf("redis://%s:%s@localhost:6379/1", os.Getenv("REDIS_USER"), os.Getenv("REDIS_PASSWORD"))
	} else {
		dbInfos, err := _aws.AwsSsmGetParams([]string{"dev_frog_redis_user", "dev_frog_redis_password", "dev_common_redis_host", "dev_common_redis_port", "dev_frog_redis_db"})
		if err != nil {
			return err
		}
		fmt.Println(dbInfos)
		connectionString = fmt.Sprintf("redis://:%s@%s:%s/%s",
			dbInfos[3], //password
			dbInfos[0], //host
			dbInfos[1], //port
			dbInfos[2], //db
		)
		fmt.Println(connectionString)
	}

	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		log.Println(err)
		return err
	}

	Client = redis.NewClient(opt)

	_, err = Client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Println("Connected to Redis!")

	return nil
}
