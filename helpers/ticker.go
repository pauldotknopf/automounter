package helpers

import (
	"time"
)

// Every tick, run the given function
func Every(duration time.Duration, f func(time.Time)) chan bool {
	done := make(chan bool, 1)
	go func() {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()
		for {
			select {
			case t := <-ticker.C:
				f(t)
			case <-done:
				return
			}
		}
	}()
	return done
}
