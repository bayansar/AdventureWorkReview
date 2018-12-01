package review

import (
	"log"
	"net/http"
	"fmt"
	"encoding/json"
	"strconv"
)

type Api struct {
	Queue ReviewQueueService
	DB    ReviewDbService
}

func (a *Api) CreateReview() http.HandlerFunc {
	return errorHandler(func(w http.ResponseWriter, r *http.Request) error {
		defer r.Body.Close()

		newReview := Review{}
		if err := json.NewDecoder(r.Body).Decode(&newReview); err != nil {
			return fmt.Errorf("cannot decoding rule : %v", err)
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

		json.NewEncoder(w).Encode(map[string]string{
			"success":  "true",
			"reviewId": strconv.FormatInt(newReview.ID, 10),
		})
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
