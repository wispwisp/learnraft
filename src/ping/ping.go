package ping

import (
	"time"

	log "github.com/wispwisp/learnraft/logger"
)

func StartPing() (statuses chan int) {
	statuses = make(chan int)

	go func() {
		defer close(statuses)

		counter := 0
		ticker := time.NewTicker(1 * time.Second)
		for {
			<-ticker.C
			statuses <- counter
			counter++

			log.Info("Status send:", counter)
		}
	}()

	return
}

func RecievePing(statuses chan int) {
	go func() {
		for status := range statuses {
			log.Info("Status recieved:", status)

		}
	}()
}
