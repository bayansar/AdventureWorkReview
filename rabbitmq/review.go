package rabbitmq

import (
	"encoding/json"
	"log"
	"github.com/streadway/amqp"
	"github.com/bayansar/AdventureWorkReview/review"
)

type ReviewQueueService struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func NewReviewQueueService(uri string, qname string) *ReviewQueueService {
	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Printf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("%s: %s", "Failed to open a channel", err)
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
		log.Printf("%s: %s", "Failed to declare a queue", err)
	}

	return &ReviewQueueService{
		conn: conn,
		ch:   ch,
		q:    q,
	}
}

func (rqs *ReviewQueueService) Publish(review *review.Review) error {
	b, err := json.Marshal(review)
	if err != nil {
		log.Printf("%s: %s", "Failed to serialize review object!", err)
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
		log.Printf("%s: %s", "Failed to publish a message", err)
	}
	log.Printf("the review is pushed to queue: %s with id: %d,", rqs.q.Name, review.ID)
	return err
}

func (rqs *ReviewQueueService) Subscribe() (<-chan review.Review, error) {
	messages, err := rqs.ch.Consume(
		rqs.q.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	if err != nil {
		log.Printf("%s: %s", "Failed to register a consumer", err)
		return nil, err
	}

	reviewMessages := make(chan review.Review)

	go func() {
		for m := range messages {
			r := review.Review{}
			err := json.Unmarshal(m.Body, &r)
			if err != nil {
				log.Printf("%s: %s", "Failed to deserialize review", err)
				continue
			}
			log.Printf("review is retrieved from queue: %s  with id: %d,", rqs.q.Name, r.ID)
			reviewMessages <- r
		}
	}()

	return reviewMessages, nil
}
