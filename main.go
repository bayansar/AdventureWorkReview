package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"

	"github.com/bayansar/AdventureWorkReview/review"
	"github.com/bayansar/AdventureWorkReview/mysql"
	"github.com/bayansar/AdventureWorkReview/email"
	"github.com/bayansar/AdventureWorkReview/rabbitmq"
)

type Config struct {
	RabbitUri         string   `env:"RABBIT_URI"`
	DbHost            string   `env:"DB_HOST"`
	DbUser            string   `env:"MYSQL_USER"`
	DbPassword        string   `env:"MYSQL_PASSWORD"`
	DbName            string   `env:"DB_NAME"`
	ValidateQueueName string   `env:"VALIDATE_QUEUE_NAME"`
	NotifyQueueName   string   `env:"NOTIFY_QUEUE_NAME"`
	BadWords          []string `env:"BAD_WORDS" envSeparator:","`
}

type Application struct {
	ReviewQueueService review.ReviewQueueService
	ReviewDbService    review.ReviewDbService
}

func main() {

	cfg := Config{}
	env.Parse(&cfg)

	reviewValidateQueueService := rabbitmq.NewReviewQueueService(cfg.RabbitUri, cfg.ValidateQueueName)
	reviewNotifyQueueService := rabbitmq.NewReviewQueueService(cfg.RabbitUri, cfg.NotifyQueueName)
	reviewDbService := mysql.NewReviewDbService(cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbHost)
	notifyService := email.NewReviewEmailService()

	reviewApi := &review.Api{
		Queue: reviewValidateQueueService,
		DB:    reviewDbService,
	}

	reviewValidator := &review.Validator{
		ConsumerQueue:  reviewValidateQueueService,
		PublisherQueue: reviewNotifyQueueService,
		DB:             reviewDbService,
		BadWords:       cfg.BadWords,
	}
	reviewValidator.Run()

	reviewNotifier := &review.Notifier{
		ConsumerQueue: reviewNotifyQueueService,
		NotifyService: notifyService,
	}
	reviewNotifier.Run()

	r := mux.NewRouter()
	r.HandleFunc("/api/reviews", reviewApi.CreateReview()).Methods("POST")

	log.Println("Started listening on 8888...")
	log.Fatal(http.ListenAndServe(":8888", r))
}
