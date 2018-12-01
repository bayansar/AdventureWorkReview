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
	*** Simulating send email line below ***
	res.emailService.sendEmail(review.Name, review.Email, message)
	*/
	log.Println("Email is sent successfully")
	return nil
}
