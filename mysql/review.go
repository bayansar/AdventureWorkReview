package mysql

import (
	"log"
	"time"
	"fmt"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/bayansar/AdventureWorkReview/review"
)

type ReviewDbService struct {
	orm orm.Ormer
}

func NewReviewDbService(username, password, dbName, host string) *ReviewDbService {
	orm.RegisterModel(new(review.Review))

	orm.RegisterDataBase(
		"default",
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", username, password, host, dbName), 30)
	orm.RunSyncdb("default", false, true)

	return &ReviewDbService{
		orm: orm.NewOrm(),
	}
}

func (rds *ReviewDbService) Insert(review *review.Review) (int64, error) {
	review.CreatedAt = time.Now()
	review.LastUpdateAt = time.Now()
	review.Status = "SUBMITTED"
	id, err := rds.orm.Insert(review)
	if err != nil {
		return 0, err
	}
	log.Printf("new review is created in db with id: %d,", review.ID)
	return id, nil
}

func (rds *ReviewDbService) Update(review *review.Review) (int64, error) {
	num, err := rds.orm.Update(review)
	if err != nil {
		log.Printf("%s: %s", "Failed to update a review", err)
		return 0, err
	}
	log.Printf("the review is updated in db with id: %d,", review.ID)
	return num, nil
}

func (rds *ReviewDbService) GetById(id int64) (*review.Review, error) {
	r := &review.Review{
		ID: id,
	}
	err := rds.orm.Read(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (rds *ReviewDbService) GetApprovedReviews() ([]*review.Review, error) {
	var reviews []*review.Review
	_, err := rds.orm.QueryTable("productreview").Filter("status", "APPROVED").All(&reviews)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
