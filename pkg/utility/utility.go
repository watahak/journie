package utility

import (
	"sync"
	"time"
)

type Throttle struct {
	mutex sync.Mutex
	last  time.Time
	delay time.Duration
}

func Debounce(fn func(args ...interface{}), duration time.Duration) func(args ...interface{}) {
	var timer *time.Timer
	var mu sync.Mutex // Mutex for thread safety

	return func(args ...interface{}) {
		mu.Lock()
		defer mu.Unlock() // Ensure critical section is thread-safe

		if timer != nil {
			timer.Stop() // Stop any pending timer
		}

		// Capture arguments for later execution
		capturedArgs := args

		// Reset the timer to trigger the function after the duration
		timer = time.NewTimer(duration)
		go func() {
			<-timer.C
			fn(capturedArgs...) // Call the original function with captured arguments
		}()
	}
}

func NewThrottle(delay time.Duration) *Throttle {
	return &Throttle{delay: delay}
}

func (t *Throttle) Process() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	elapsed := time.Since(t.last)
	if elapsed < t.delay {
		time.Sleep(t.delay - elapsed)
	}
	t.last = time.Now()
}
