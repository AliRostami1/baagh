package debounce

import (
	"time"
)

func Debounce(interval time.Duration, cb func()) (debouncedfn func()) {
	ch := make(chan struct{})
	timer := time.NewTimer(interval)
	go func() {
		for {
			select {
			case <-ch:
				timer.Reset(interval)
			case <-timer.C:
				cb()
			}
		}
	}()

	return func() {
		ch <- struct{}{}
	}
}
