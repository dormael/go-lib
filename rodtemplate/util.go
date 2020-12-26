package rodtemplate

import (
	"fmt"
	"log"
	"time"
)

type TimeoutError struct {
	Timout  time.Duration
	Started time.Time
	Message string
}

func WaitFor(timeout, retryDuration time.Duration, checkFunc func() bool, retryFunc func()) error {
	started := time.Now()
	lastRetry := time.Now()

	for timeout > time.Now().Sub(started) {
		if true == checkFunc() {
			return nil
		}

		sleepDuration := retryDuration - time.Now().Sub(lastRetry)
		if sleepDuration < 0 {
			sleepDuration = retryDuration
		}
		time.Sleep(sleepDuration)
		retryFunc()
		lastRetry = time.Now()

		log.Println("retry after sleep", sleepDuration, "for retryDuration", retryDuration)
	}

	message := fmt.Sprintf("timeout %s exceeded after %s", timeout, started)
	log.Println(message)

	return &TimeoutError{Timout: timeout, Started: started, Message: message}
}

func (e TimeoutError) Error() string {
	return e.Message
}

func (e TimeoutError) Timeout() bool {
	return true
}