package queue

import (
	"auth-service/internal/config"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type queueService struct {
	ch    *amqp.Channel
	queue map[string]amqp.Queue
	conn  *amqp.Connection
}

type QueueService interface {
	SendMsg(msg string) error
	Close()
}

func NewClientQueue(cnf *config.ConfigQueue) QueueService {
	amqpServerURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", cnf.User, cnf.Password, cnf.Host, cnf.Port)
	conn, err := amqp.Dial(amqpServerURL) // Создаем подключение к RabbitMQ
	if err != nil {
		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		fmt.Println(err)
	}

	q, err := ch.QueueDeclare(
		"event", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	mapQueue := map[string]amqp.Queue{"event": q}
	return &queueService{ch, mapQueue, conn}
}

func (s *queueService) SendMsg(msg string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.ch.PublishWithContext(ctx,
		"",                    // exchange
		s.queue["event"].Name, // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	if err != nil {
		return err
	}
	return nil
}

func (s *queueService) Close() {
	s.conn.Close()
}
