package utils

import (
	"fmt"
	"log"
	"os"

	_aws "main/utils/aws"

	amqp "github.com/rabbitmq/amqp091-go"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var V2MQ *amqp.Queue
var V1MQ *amqp.Queue
var V1MQCH *amqp.Channel
var V2MQCH *amqp.Channel

func InitRabbitMQ() error {
	var connURL string
	if Env.IsLocal {
		// MySQL 연결 문자열
		connURL = fmt.Sprintf("amqp://%s:%s@localhost:5672/",
			os.Getenv("RABBITMQ_USER"),
			os.Getenv("RABBITMQ_PASSWORD"),
		)
	} else {
		dbInfos, err := _aws.AwsSsmGetParams([]string{"dev_frog_rabbitmq_user", "dev_frog_rabbitmq_password", "dev_common_mysql_host", "dev_frog_rabbitmq_port"})
		if err != nil {
			return err
		}
		fmt.Println(dbInfos)
		connURL = fmt.Sprintf("amqp://%s:%s@%s:%s/",
			dbInfos[3], //user
			dbInfos[1], //password
			dbInfos[0], //host
			dbInfos[2], //port
		)
	}
	fmt.Println(connURL)
	conn, err := amqp.Dial(connURL)
	if err != nil {
		fmt.Errorf("Failed to connect to RabbitMQ: %v", err)
		return err
	}

	V1MQCH, err = conn.Channel()
	if err != nil {
		fmt.Errorf("Failed to open a channel: %v", err)
		return err
	}
	V2MQCH, err = conn.Channel()
	if err != nil {
		fmt.Errorf("Failed to open a channel: %v", err)
		return err
	}
	v2queue, err := V2MQCH.QueueDeclare(
		"wingspan", // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		fmt.Errorf("Failed to declare a queue: %v", err)
		return err
	}

	V2MQ = &v2queue

	v1queue, err := V1MQCH.QueueDeclare(
		"frog", // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		fmt.Errorf("Failed to declare a queue: %v", err)
		return err
	}

	V1MQ = &v1queue
	return nil
}
