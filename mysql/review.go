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

	orm.Debug = true

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
		log.Printf("%s: %s", "Failed to insert a review", err)
		return 0, err
	}
	return id, nil
}

func (rds *ReviewDbService) Update(review *review.Review) (int64, error) {
	num, err := rds.orm.Update(review)
	if err != nil {
		log.Printf("%s: %s", "Failed to insert a review", err)
		return 0, err
	}
	return num, nil
}
