package backoff

import "time"

// A Function is a function with no inputs but returns some data
// and is retried using a backoff policy if it returns an error.
type Function func() (interface{}, error)

// An Action is a function with no inputs and requires no return data.
// The action will be retried using a backoff policy if it returns an error.
type Action func() error

// Algorithm is a custom type for defining constants that
// represent available backoff algorithms to choose from.
type Algorithm int

const (
	// AlgorithmExponential is a constant defining the [Exponential Backoff Algorithm](https://en.wikipedia.org/wiki/Exponential_backoff)
	AlgorithmExponential Algorithm = iota
)

// ServiceAPI defines the format of a backoff service. This interface may be used to help mock certain functionality.
type ServiceAPI interface {
	Setup(policy *Policy)
	ExecuteFunction(op Function) (interface{}, error)
	ExecuteAction(op Action) error
	NextBackOff() time.Duration
	ElapsedTime() time.Duration
	ResetDefaults()
	Reset()
}
