package backoff

import (
	"errors"
	"math/rand"
	"time"
)

const (
	// DefaultExponentialInitialInterval is a constant defining the starting wait time per wait cycle.
	DefaultExponentialInitialInterval = 500 * time.Millisecond
	// DefaultExponentialRandFactor is a constant defining the factor by which time fluctuates per wait cycle.
	DefaultExponentialRandFactor = 0.5
	// DefaultExponentialMultiplier is a constant defining the multiplication factor for the interval.
	DefaultExponentialMultiplier = 1.5
	// DefaultExponentialMaxInterval is the maximum wait duration that the backoff wait interval can have per cycle.
	DefaultExponentialMaxInterval = 60 * time.Second
	// DefaultExponentialMaxElapsedTime is the maximum time that the backoff can run for before exiting.
	DefaultExponentialMaxElapsedTime = 15 * time.Minute
)

// ExponentialBackoff is a BackoffServiceAPI that implements the exponential backoff algorithm.
type ExponentialBackoff struct {
	InitialInterval    time.Duration
	RandFactor         float64
	IntervalMultiplier float64
	MaxInterval        time.Duration
	MaxElapsedTime     time.Duration
	Clock              Clock

	currentInterval time.Duration
	startTime       time.Time
	random          *rand.Rand
}

// Setup accepts a Policy struct and sets the backoff execution properties with defaults for properties undeclared.
func (bckOff *ExponentialBackoff) Setup(policy *Policy) {
	bckOff.ResetDefaults()

	if int64(policy.StartInterval) != 0 {
		bckOff.InitialInterval = policy.StartInterval
	}

	if policy.RandomizationFactor != 0 {
		bckOff.RandFactor = policy.RandomizationFactor
	}

	if policy.IntervalMultiplier != 0 {
		bckOff.IntervalMultiplier = policy.IntervalMultiplier
	}

	if int64(policy.MaxInterval) != 0 {
		bckOff.MaxInterval = policy.MaxInterval
	}

	if int64(policy.MaxElapsedTime) != 0 {
		bckOff.MaxElapsedTime = policy.MaxElapsedTime
	}

}

// ExecuteFunction accepts an input of type Function to be executed via an exponential backoff
// algorithm. After the function is called successfully, the output of the function is returned.
func (bckOff *ExponentialBackoff) ExecuteFunction(op Function) (interface{}, error) {
	bckOff.Reset()

	var (
		dur    time.Duration
		output interface{}
		valid  bool
		err    error
	)

	valid = true

	for valid {
		output, err = op()
		if dur = bckOff.NextBackOff(); err != nil && dur != -1 {
			time.Sleep(dur)
		} else {
			valid = false
		}
	}

	if dur == -1 {
		return nil, errors.New("Backoff function has reached a stop condition and failed")
	}

	return output, nil
}

// ExecuteAction accepts an input of type Action to be executed via an exponential backoff
// algorithm. The execution of this function is equivalent to a fire and forget strategy.
func (bckOff *ExponentialBackoff) ExecuteAction(op Action) error {
	bckOff.Reset()

	var dur time.Duration
	valid := true

	for valid {
		err := op()
		if dur = bckOff.NextBackOff(); err != nil && dur != -1 {
			time.Sleep(dur)
		} else {
			valid = false
		}
	}

	if dur == -1 {
		return errors.New("Backoff function has reached a stop condition and failed")
	}

	return nil
}

// NextBackOff determines the amount of time to wait before executing again.
// The function also updates certain properties for the next backoff cycle.
func (bckOff *ExponentialBackoff) NextBackOff() time.Duration {
	if bckOff.MaxElapsedTime != 0 && bckOff.ElapsedTime() > bckOff.MaxElapsedTime {
		return -1
	}

	delta := bckOff.RandFactor * float64(bckOff.currentInterval)
	minInterval := float64(bckOff.currentInterval) - delta
	maxInterval := float64(bckOff.currentInterval) + delta

	bckOff.currentInterval = time.Duration(minInterval + (bckOff.random.Float64() * (maxInterval - minInterval + 1)))

	if float64(bckOff.currentInterval) >= float64(bckOff.MaxInterval)/bckOff.IntervalMultiplier {
		bckOff.currentInterval = bckOff.MaxInterval
		return -1
	}

	bckOff.currentInterval = time.Duration(float64(bckOff.currentInterval) * bckOff.IntervalMultiplier)

	return bckOff.currentInterval
}

// ElapsedTime determines the total amount of time that the backoff execution has been running.
func (bckOff *ExponentialBackoff) ElapsedTime() time.Duration {
	return bckOff.Clock.Now().Sub(bckOff.startTime)
}

// ResetDefaults sets the properties of the backoff service to their API defaults. See constants prefixed with DefaultExponential.
func (bckOff *ExponentialBackoff) ResetDefaults() {
	bckOff.InitialInterval = DefaultExponentialInitialInterval
	bckOff.RandFactor = DefaultExponentialRandFactor
	bckOff.IntervalMultiplier = DefaultExponentialMultiplier
	bckOff.MaxInterval = DefaultExponentialMaxInterval
	bckOff.MaxElapsedTime = DefaultExponentialMaxElapsedTime
	bckOff.Clock = &SystemClock{}
}

// Reset is ran prior to executing the start of a backoff cycle in order to properly calculate for elapsed time and other variables.
func (bckOff *ExponentialBackoff) Reset() {
	bckOff.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	bckOff.currentInterval = bckOff.InitialInterval
	bckOff.startTime = bckOff.Clock.Now()
}
