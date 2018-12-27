package review

import (
	"log"
	"net/http"
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/gorilla/mux"
)

type Api struct {
	Queue QueueService
	DB    DbService
}

func (a *Api) CreateReview() http.HandlerFunc {
	return errorHandler(func(w http.ResponseWriter, r *http.Request) error {
		defer r.Body.Close()

		newReview := Review{}
		if err := json.NewDecoder(r.Body).Decode(&newReview); err != nil {
			return fmt.Errorf("cannot decoding review : %v", err)
		}

		id, err := a.DB.Insert(&newReview)
		if err != nil {
			return err
		}

		newReview.ID = id
		err = a.Queue.Publish(&newReview)
		if err != nil {
			return err
		}

		err = json.NewEncoder(w).Encode(map[string]string{
			"success":  "true",
			"reviewId": strconv.FormatInt(newReview.ID, 10),
		})
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *Api) GetReviewById() http.HandlerFunc {
	return errorHandler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse variable: %v", err)
		}
		review, err := a.DB.GetById(id)
		if err != nil {
			return err
		}
		err = json.NewEncoder(w).Encode(review)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *Api) GetApprovedReviews() http.HandlerFunc {
	return errorHandler(func(w http.ResponseWriter, r *http.Request) error {
		reviews, err := a.DB.GetApprovedReviews()
		if err != nil {
			return nil
		}
		err = json.NewEncoder(w).Encode(reviews)
		if err != nil {
			return err
		}
		return nil
	})
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			log.Printf("handling %q: %v", r.RequestURI, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
