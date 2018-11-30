package main

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"

	"github.com/bayansar/AdventureWorkReview/app"
	"github.com/bayansar/AdventureWorkReview/rabbitmq"
)

type Application struct {
	ReviewQueueService *rabbitmq.ReviewQueueService
	ReviewDbService    *app.ReviewDbService
}

func main() {

	appContext := &Application{
		ReviewQueueService: rabbitmq.NewReviewQueueService("amqp://guest:guest@127.0.0.1:5672", "test"),
		ReviewDbService:    nil,
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/reviews", appContext.review()).Methods("POST")
	http.ListenAndServe(":8888", r)
}

func (appContext *Application) review() http.HandlerFunc {
	return errorHandler(func(w http.ResponseWriter, r *http.Request) error {
		defer r.Body.Close()

		newReview := app.Review{}
		if err := json.NewDecoder(r.Body).Decode(&newReview); err != nil {
			return fmt.Errorf("cannot decoding rule : %v", err)
		}

		newReview.ID = ksuid.New().String()
		err := appContext.ReviewQueueService.Publish(&newReview)
		if err != nil{
			json.NewEncoder(w).Encode(map[string]string{
				"success": "false",
				"message": err.Error(),
			})
			return err
		}

		json.NewEncoder(w).Encode(map[string]string{
			"success": "true",
			"reviewId": newReview.ID,
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
