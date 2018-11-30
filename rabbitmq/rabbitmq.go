package rabbitmq

import (
	"encoding/json"
	"github.com/bayansar/AdventureWorkReview/app"
	"github.com/streadway/amqp"
	"log"
)

type ReviewQueueService struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	q     amqp.Queue
}

func NewReviewQueueService(uri string, qname string) *ReviewQueueService {
	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	q, err := ch.QueueDeclare(
		qname, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	return &ReviewQueueService {
		conn: conn,
		ch:   ch,
		q:    q,
	}
}

func (rqs *ReviewQueueService) Publish(review *app.Review) error {
	b, err := json.Marshal(review)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to serialize review object!", err)
		return err
	}
	err = rqs.ch.Publish(
		"",         // exchange
		rqs.q.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})

	if err != nil {
		log.Fatalf("%s: %s", "Failed to publish a message", err)
	}
	return err
}

func (rqs *ReviewQueueService) Subscribe() (*app.Review, error) {
	return nil, nil
}
