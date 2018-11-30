package app

type Review struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	ProductId string `json:"productid"`
	Review    string `json:"review"`
}

type ReviewQueueService interface {
	Publish(review *Review) error
	Subscribe() (*Review, error)
}

type ReviewDbService interface {
	Review(id string) (*Review, error)
	CreateReview(review *Review) (*Review, error)
	UpdateReview(review *Review) (*Review, error)
}