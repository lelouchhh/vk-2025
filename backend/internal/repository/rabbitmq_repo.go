package repository

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"backend/domain"
	"github.com/streadway/amqp"
)

type RabbitMQRepository interface {
	ConsumePingResults(ctx context.Context) (<-chan domain.PingResult, error)
	Close() error
}

type rabbitMQRepository struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQRepository(rabbitMQURL string) (RabbitMQRepository, error) {
	var conn *amqp.Connection
	var err error

	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(rabbitMQURL)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, errors.New("failed to connect to RabbitMQ after multiple attempts")
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		"ping_results", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &rabbitMQRepository{
		conn: conn,
		ch:   ch,
	}, nil
}

func (r *rabbitMQRepository) ConsumePingResults(ctx context.Context) (<-chan domain.PingResult, error) {
	msgs, err := r.ch.Consume(
		"ping_results", // queue
		"",             // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		return nil, err
	}

	results := make(chan domain.PingResult)
	go func() {
		defer close(results)
		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping RabbitMQ consumption due to context cancellation")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("RabbitMQ channel closed")
					return
				}
				var result domain.PingResult
				err := json.Unmarshal(msg.Body, &result)
				if err != nil {
					log.Printf("Failed to unmarshal message: %v", err)
					continue
				}
				results <- result
				msg.Ack(false)
			}
		}
	}()

	return results, nil
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
