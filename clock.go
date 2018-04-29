package backoff

import "time"

// Clock defines the format of a clock struct for computing the current time.
type Clock interface {
	Now() time.Time
}

// SystemClock is a collections of methods focused on using the default system clock
type SystemClock struct{}

// Now returns the exact current system time
func (clock *SystemClock) Now() time.Time {
	return time.Now()
}
