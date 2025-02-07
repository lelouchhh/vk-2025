package repository

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQRepository interface {
	PublishPingResult(result interface{}) error
	Close() error
}

type rabbitMQRepository struct {
	conn        *amqp.Connection
	ch          *amqp.Channel
	rabbitMQURL string
}

func NewRabbitMQRepository(rabbitMQURL string) (RabbitMQRepository, error) {
	repo := &rabbitMQRepository{
		rabbitMQURL: rabbitMQURL,
	}

	// Попытка подключения к RabbitMQ
	if err := repo.connect(); err != nil {
		return nil, err
	}

	// Объявление очереди
	if err := repo.declareQueue(); err != nil {
		repo.Close()
		return nil, err
	}

	return repo, nil
}

func (r *rabbitMQRepository) connect() error {
	var conn *amqp.Connection
	var err error

	conn, err = amqp.Dial(r.rabbitMQURL)

	if err != nil {
		return errors.New("failed to connect to RabbitMQ after multiple attempts")
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	r.conn = conn
	r.ch = ch
	log.Println("Successfully connected to RabbitMQ")
	return nil
}

func (r *rabbitMQRepository) declareQueue() error {
	if r.ch == nil {
		return errors.New("RabbitMQ channel is not initialized")
	}

	_, err := r.ch.QueueDeclare(
		"ping_results", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
		return err
	}

	log.Println("Queue 'ping_results' declared successfully")
	return nil
}

func (r *rabbitMQRepository) ensureConnection() error {
	if r.conn == nil || r.conn.IsClosed() {
		log.Println("Reconnecting to RabbitMQ...")
		if err := r.connect(); err != nil {
			log.Printf("Failed to reconnect to RabbitMQ: %v", err)
			return err
		}
	}
	return nil
}

func (r *rabbitMQRepository) PublishPingResult(result interface{}) error {
	// Убедимся, что соединение активно
	if err := r.ensureConnection(); err != nil {
		return err
	}

	body, err := json.Marshal(result)
	if err != nil {
		log.Printf("Failed to marshal result: %v", err)
		return err
	}

	err = r.ch.Publish(
		"",             // exchange
		"ping_results", // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	log.Println("Message published successfully")
	return nil
}

func (r *rabbitMQRepository) Close() error {
	if r.ch != nil {
		if err := r.ch.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ channel: %v", err)
			return err
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
			return err
		}
	}
	log.Println("RabbitMQ connection closed")
	return nil
}
