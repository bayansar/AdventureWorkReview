package review

import (
	"log"
	"regexp"
	"strings"
)

type Validator struct {
	PublisherQueue QueueService
	ConsumerQueue  QueueService
	DB             DbService
	BadWords       []string
}

func (v *Validator) Run() error {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Printf("%s: %s", "Failed compile regex", err)
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
			_, err := v.DB.Update(&m)
			if err != nil {
				continue
			}
			err = v.PublisherQueue.Publish(&m)
			if err != nil {
				continue
			}
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
