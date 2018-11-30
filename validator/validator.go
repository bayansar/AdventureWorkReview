package validator

import "github.com/bayansar/AdventureWorkReview/app"

type Validator struct {
	BadWords []string
	SubQname string
	PubQname string

	PublisherQueue app.ReviewQueueService
	ConsumerQueue app.ReviewQueueService
}

func (*Validator) run() error{

}