// Package periodic contains utilities to run periodic tasks
package periodic

import "time"

// Serve asynchronously runs a task periodically until it is cancelled.
// To cancel Serve, send to the cancel channel.
// If cancel is nil, a new buffered channel of size 1 is created.
// Serve returns the cancel channel.
func Serve(period time.Duration, task func(), cancel chan bool) chan bool {
	if cancel == nil {
		cancel = make(chan bool, 1)
	}
	go func() {
		for {
			if !Sleep(period, cancel) {
				return
			}
			task()
		}
	}()
	return cancel
}

// Sleep will block until interval passes, or cancel is read from.
// To cancel the sleep, simply send a value to cancel. Cancel can be nil,
// in which case the caller will not be able to cancel the sleep.
// Sleep returns true if it completed sleeping,
// and returns false if it was cancelled.
func Sleep(interval time.Duration, cancel <-chan bool) bool {
	timeout := make(chan bool)
	go func() {
		time.Sleep(interval)
		timeout <- true
	}()
	select {
	case <-timeout:
		return true
	case <-cancel:
		return false
	}
}
