package debounce

import (
	"time"
)

func Debounce(interval time.Duration, cb func()) (debouncedfn func(), cancel func()) {
	ch := make(chan struct{})
	timer := time.NewTimer(interval)
	go func() {
		for {
			select {
			case _, ok := <-ch:
				if !ok {
					return
				}
				timer.Reset(interval)
			case <-timer.C:
				cb()
			}
		}
	}()

	return func() {
			ch <- struct{}{}
		}, func() {
			close(ch)
		}
}
