package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"

	"github.com/bayansar/AdventureWorkReview/email"
	"github.com/bayansar/AdventureWorkReview/mysql"
	"github.com/bayansar/AdventureWorkReview/rabbitmq"
	"github.com/bayansar/AdventureWorkReview/review"
	"os"
)

type Config struct {
	RabbitUri         string   `env:"RABBIT_URI,required"`
	DbHost            string   `env:"DB_HOST,required"`
	DbUser            string   `env:"MYSQL_USER,required"`
	DbPassword        string   `env:"MYSQL_PASSWORD,required"`
	DbName            string   `env:"DB_NAME,required"`
	ValidateQueueName string   `env:"VALIDATE_QUEUE_NAME,required"`
	NotifyQueueName   string   `env:"NOTIFY_QUEUE_NAME,required"`
	BadWords          []string `env:"BAD_WORDS,required" envSeparator:","`
}

func main() {

	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("%s : %v","Couldn't parse environment variables!" , err)
	}

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
	err = reviewValidator.Run()
	if err != nil {
		os.Exit(1)
	}

	reviewNotifier := &review.Notifier{
		ConsumerQueue: reviewNotifyQueueService,
		NotifyService: notifyService,
	}
	err = reviewNotifier.Run()
	if err != nil {
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/reviews", reviewApi.CreateReview()).Methods("POST")
	r.HandleFunc("/api/reviews/approved", reviewApi.GetApprovedReviews()).Methods("GET")
	r.HandleFunc("/api/reviews/{id}", reviewApi.GetReviewById()).Methods("GET")

	log.Println("Server listening on port 8888...")
	log.Fatal(http.ListenAndServe(":8888", r))
}
