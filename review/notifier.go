package review

import (
	"log"
)

type Notifier struct {
	ConsumerQueue QueueService
	NotifyService NotifyService
	approveMes    string
	declineMes    string
}

func (n *Notifier) Run() error {
	n.approveMes = "Your review is approved and released!"
	n.declineMes = "Your review is declined because of bad language!"

	messages, err := n.ConsumerQueue.Subscribe()
	if err != nil {
		return err
	}

	go func() {
		for m := range messages {
			var notifyMes string
			if m.Status == "APPROVED" {
				notifyMes = n.approveMes
			} else if m.Status == "DECLINED" {
				notifyMes = n.declineMes
			} else {
				log.Printf("%s: %s", "Illegal status to notify which is ", m.Status)
			}
			err := n.NotifyService.Notify(&m, notifyMes)
			if err != nil {
				log.Printf("%s: %s", "Failed to notify user", err)
			}
		}
	}()

	return nil
}
