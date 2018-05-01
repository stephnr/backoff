package backoff

import (
	"errors"
	"time"
)

// Service is the backoff service for retrying an operation.
type Service struct {
	Policy
	ServiceAPI
}

// A Policy represents the requirements for a backoff retry operation.
type Policy struct {
	Algorithm           Algorithm
	RandomizationFactor float64
	IntervalMultiplier  float64
	StartInterval       time.Duration
	MaxInterval         time.Duration
	MaxElapsedTime      time.Duration
}

// New constructs an instance of the Backoff Service for retrying an operation.
func New(policy *Policy) (ServiceAPI, error) {
	service := &Service{Policy: *policy}

	switch policy.Algorithm {
	case AlgorithmExponential:
		service.ServiceAPI = &ExponentialBackoff{}
	default:
		return nil, errors.New("The selected backoff algorithm is not valid")
	}

	service.Setup(policy)

	return service, nil
}
