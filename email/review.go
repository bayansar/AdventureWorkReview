package email

import (
	"log"
	"github.com/bayansar/AdventureWorkReview/review"
)

type ReviewEmailService struct {
	emailService interface{}
}

func NewReviewEmailService() *ReviewEmailService{
	return &ReviewEmailService{}
}

func (res *ReviewEmailService) Notify(review *review.Review, message string) error {

	/*
	*** Send email is simulated below line***
	res.emailService.sendEmail(review.Name, review.Email, message)
	*/
	log.Println("email is sent successfully")
	return nil
}
