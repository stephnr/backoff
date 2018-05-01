package backoff

import (
	"errors"
	"fmt"
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
	IntervalMultiplier  float64
	MaxElapsedTime      time.Duration
	MaxInterval         time.Duration
	MaxRetryCount       int64
	RandomizationFactor float64
	StartInterval       time.Duration
}

func validPolicy(policy *Policy) error {
	if policy.RandomizationFactor < 0 || policy.RandomizationFactor > 1 {
		return fmt.Errorf("the provided randomization factor of [ %f ] is not allowed. The allowed range of values is from 0 to 1", policy.RandomizationFactor)
	}

	return nil
}

// New constructs an instance of the Backoff Service for retrying an operation.
func New(policy *Policy) (ServiceAPI, error) {
	service := &Service{Policy: *policy}

	if err := validPolicy(policy); err != nil {
		return nil, err
	}

	switch policy.Algorithm {
	case AlgorithmExponential:
		service.ServiceAPI = &ExponentialBackoff{}
	default:
		return nil, errors.New("The selected backoff algorithm is not valid")
	}

	service.Setup(policy)

	return service, nil
}
