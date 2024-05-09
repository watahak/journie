package utility

import (
	"time"
)

func Debounce(fn func(args ...interface{}), duration time.Duration) func(args ...interface{}) {
	var timer *time.Timer

	return func(args ...interface{}) {
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
