package backoff_test

import (
	"errors"
	"testing"
	"time"

	"github.com/defaltd/backoff"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ExponentialTestSuite struct {
	suite.Suite

	service  backoff.ServiceAPI
	mockFunc *MockFunction
}

func (st *ExponentialTestSuite) SetupTest() {
	st.service, _ = backoff.New(&backoff.Policy{
		Algorithm:           backoff.AlgorithmExponential,
		StartInterval:       time.Duration(time.Millisecond * 100),
		RandomizationFactor: 0.5,
		IntervalMultiplier:  0.5,
		MaxInterval:         time.Duration(time.Millisecond * 500),
		MaxElapsedTime:      time.Duration(time.Second * 2),
	})

	st.mockFunc = &MockFunction{}
}

type MockFunction struct {
	mock.Mock
}

func (fun *MockFunction) BadAction() error {
	return errors.New("FAILED")
}

func (fun *MockFunction) GoodAction() error {
	return nil
}

func (fun *MockFunction) BadFunction() (interface{}, error) {
	return nil, errors.New("FAILED")
}

func (fun *MockFunction) GoodFunction() (interface{}, error) {
	return "Hello World", nil
}

func (st *ExponentialTestSuite) TestBadAction() {
	startTime := time.Now()
	err := st.service.ExecuteAction(st.mockFunc.BadAction)

	assert.Error(st.T(), err)
	assert.True(st.T(), time.Since(startTime).Nanoseconds() >= 500000000)
}

func (st *ExponentialTestSuite) TestGoodAction() {
	startTime := time.Now()
	err := st.service.ExecuteAction(st.mockFunc.GoodAction)

	assert.Nil(st.T(), err)
	assert.True(st.T(), time.Since(startTime).Nanoseconds() < 500000000)
}

func (st *ExponentialTestSuite) TestBadFunction() {
	startTime := time.Now()
	_, err := st.service.ExecuteFunction(st.mockFunc.BadFunction)

	assert.Error(st.T(), err)
	assert.True(st.T(), time.Since(startTime).Nanoseconds() >= 500000000)
}

func (st *ExponentialTestSuite) TestGoodFunction() {
	startTime := time.Now()
	output, err := st.service.ExecuteFunction(st.mockFunc.GoodFunction)

	assert.Nil(st.T(), err)
	assert.NotNil(st.T(), output)
	assert.True(st.T(), time.Since(startTime).Nanoseconds() < 500000000)
}

func (st *ExponentialTestSuite) TestMaxIntervial() {
	st.service, _ = backoff.New(&backoff.Policy{
		Algorithm:      backoff.AlgorithmExponential,
		MaxInterval:    time.Duration(time.Second * 1),
		MaxElapsedTime: time.Duration(time.Second * 10),
	})

	err := st.service.ExecuteAction(st.mockFunc.BadAction)

	assert.Error(st.T(), err)
}

func (st *ExponentialTestSuite) TestInvalidAlgorithm() {
	_, err := backoff.New(&backoff.Policy{
		Algorithm: -1,
	})

	assert.Error(st.T(), err)
}

func TestExponentialTestSuite(t *testing.T) {
	suite.Run(t, new(ExponentialTestSuite))
}
