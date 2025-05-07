package utils

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_aws "main/utils/aws"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	instance     *RabbitMQManager
	once         sync.Once
	reconnectMux sync.Mutex
)

type RabbitMQManager struct {
	Connection *amqp.Connection
	Channels   map[string]*amqp.Channel // 채널 맵
	Queues     map[string]*amqp.Queue   // 큐 맵
}

func GetRabbitMQManager() *RabbitMQManager {
	once.Do(func() {
		instance = &RabbitMQManager{
			Channels: make(map[string]*amqp.Channel),
			Queues:   make(map[string]*amqp.Queue),
		}
		if err := instance.connect(); err != nil {
			log.Fatalf("Failed to initialize RabbitMQ: %v", err)
		}
		go instance.monitorConnection()
	})
	return instance
}

func (r *RabbitMQManager) connect() error {
	var connURL string
	if Env.IsLocal {
		connURL = fmt.Sprintf("amqp://%s:%s@localhost:5672/",
			os.Getenv("RABBITMQ_USER"),
			os.Getenv("RABBITMQ_PASSWORD"),
		)
	} else {
		dbInfos, err := _aws.AwsSsmGetParams([]string{"dev_frog_rabbitmq_user", "dev_frog_rabbitmq_password", "dev_common_mysql_host", "dev_frog_rabbitmq_port"})
		if err != nil {
			return fmt.Errorf("Failed to fetch RabbitMQ credentials: %v", err)
		}
		connURL = fmt.Sprintf("amqp://%s:%s@%s:%s/",
			dbInfos[3], // user
			dbInfos[1], // password
			dbInfos[0], // host
			dbInfos[2], // port
		)
	}

	conn, err := amqp.Dial(connURL)
	if err != nil {
		return fmt.Errorf("Failed to connect to RabbitMQ: %v", err)
	}

	r.Connection = conn

	// `wingspan` 채널 및 큐 초기화
	if err := r.initChannelAndQueue("wingspan"); err != nil {
		return err
	}

	// `frog` 채널 및 큐 초기화
	if err := r.initChannelAndQueue("frog"); err != nil {
		return err
	}
	// `frog` 채널 및 큐 초기화
	if err := r.initChannelAndQueue("find-it"); err != nil {
		return err
	}
	if err := r.initChannelAndQueue("slime-war"); err != nil {
		return err
	}
	if err := r.initChannelAndQueue("sequence"); err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQManager) initChannelAndQueue(queueName string) error {
	channel, err := r.Connection.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open RabbitMQ channel for %s: %v", queueName, err)
	}

	queue, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		true,      // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("Failed to declare RabbitMQ queue %s: %v", queueName, err)
	}

	r.Channels[queueName] = channel
	r.Queues[queueName] = &queue

	log.Printf("Initialized RabbitMQ channel and queue: %s", queueName)
	return nil
}

func (r *RabbitMQManager) monitorConnection() {
	for {
		if r.Connection.IsClosed() {
			log.Println("RabbitMQ connection lost. Reconnecting...")
			reconnectMux.Lock()
			err := r.connect()
			if err != nil {
				log.Printf("Failed to reconnect to RabbitMQ: %v", err)
				reconnectMux.Unlock()
				time.Sleep(5 * time.Second)
				continue
			}
			// 연결 복구 후 채널 및 큐 재초기화
			for queueName, queue := range r.Queues {
				count, err := r.Channels[queueName].QueuePurge(queue.Name, false)
				if err != nil {
					log.Printf("Failed to purge queue %s: %v", queueName, err)
				} else {
					log.Printf("Successfully purged queue: %s, count: %d", queueName, count)
				}
				// 채널 및 큐 재초기화
				if err := r.initChannelAndQueue(queueName); err != nil {
					log.Printf("Failed to reinitialize channel and queue for %s: %v", queueName, err)
				} else {
					log.Printf("Successfully reinitialized channel and queue for %s", queueName)
				}
			}
			reconnectMux.Unlock()
			log.Println("RabbitMQ connection reestablished.")
		}
		time.Sleep(1 * time.Second)
	}
}

func (r *RabbitMQManager) GetChannel(queueName string) (*amqp.Channel, error) {
	channel, ok := r.Channels[queueName]
	if !ok || channel.IsClosed() {
		log.Printf("Channel for queue %s is not available or closed. Reinitializing...", queueName)
		if err := r.initChannelAndQueue(queueName); err != nil {
			return nil, fmt.Errorf("failed to reinitialize channel for queue %s: %v", queueName, err)
		}
		channel = r.Channels[queueName]
	}
	return channel, nil
}
func (r *RabbitMQManager) PublishMessage(queueName string, message []byte) error {
	channel, ok := r.Channels[queueName]
	if !ok {
		return fmt.Errorf("Channel for queue %s not initialized", queueName)
	}

	queue, ok := r.Queues[queueName]
	if !ok {
		return fmt.Errorf("Queue %s not initialized", queueName)
	}

	err := channel.Publish(
		"",         // exchange
		queue.Name, // routing key (queue name)
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("Failed to publish message to RabbitMQ queue %s: %v", queueName, err)
	}
	return nil
}
