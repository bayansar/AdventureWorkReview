package review

import (
	"strings"
	"regexp"
	"log"
)

type Validator struct {
	PublisherQueue ReviewQueueService
	ConsumerQueue  ReviewQueueService
	DB             ReviewDbService
	BadWords       []string
}

func (v *Validator) Run() error {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Printf("%s: %s", "Failed parse review", err)
		return err
	}

	messages, err := v.ConsumerQueue.Subscribe()
	if err != nil {
		return err
	}

	go func() {
		for m := range messages {
			wordList := v.split(m.Review, reg)
			if v.checkForBadWords(wordList, v.BadWords) == true {
				m.Status = "APPROVED"
			} else {
				m.Status = "DECLINED"
			}
			v.DB.Update(&m)
			v.PublisherQueue.Publish(&m)
		}
	}()

	return nil
}

func (v *Validator) split(review string, reg *regexp.Regexp) []string {
	processedReview := reg.ReplaceAllString(review, " ")
	return strings.Split(processedReview, " ")
}

func (v *Validator) checkForBadWords(wordList []string, badWords []string) bool {
	// A map might be used in case of high number of bad words to avoid from quadratic complexity
	for _, word := range wordList {
		for _, badWord := range v.BadWords {
			if strings.EqualFold(word, badWord) == true {
				return false
			}
		}
	}
	return true
}
