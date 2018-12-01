package review

import (
	"time"
)

type Review struct {
	ID           int64     `json:"id" orm:"auto;column(ProductReviewID)"`
	ProductId    string    `json:"productid" orm:"column(ProductID)"`
	Name         string    `json:"name" orm:"column(ReviewerName)"`
	CreatedAt    time.Time `json:"createdAt" orm:"null;column(ReviewDate)"`
	Email        string    `json:"email" orm:"column(EmailAddress)"`
	Review       string    `json:"review" orm:"column(Comments)"`
	Status       string    `json:"status" orm:"null;column(Status)"`
	LastUpdateAt time.Time `json:"lastUpdateAt" orm:"null;column(ModifiedDate)"`
}

func (u *Review) TableName() string {
	return "productreview"
}

type ReviewQueueService interface {
	Publish(review *Review) error
	Subscribe() (<-chan Review, error)
}

type ReviewDbService interface {
	//Get(id string) (*Review, error)
	Insert(review *Review) (int64, error)
	Update(review *Review) (int64, error)
}

type ReviewNotifyService interface {
	Notify(review *Review, message string) error
}
