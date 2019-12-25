package retry

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/jitter"
	"github.com/Rican7/retry/strategy"
)

func TestRetry(t *testing.T) {
	action := func(attempt uint) error {
		return nil
	}

	err := Retry(action)

	if nil != err {
		t.Error("expected a nil error")
	}
}

func TestRetryRetriesUntilNoErrorReturned(t *testing.T) {
	const errorUntilAttemptNumber = 5

	var attemptsMade uint

	action := func(attempt uint) error {
		attemptsMade = attempt

		if errorUntilAttemptNumber == attempt {
			return nil
		}

		return errors.New("erroring")
	}

	err := Retry(action)

	if nil != err {
		t.Error("expected a nil error")
	}

	if errorUntilAttemptNumber != attemptsMade {
		t.Errorf(
			"expected %d attempts to be made, but %d were made instead",
			errorUntilAttemptNumber,
			attemptsMade,
		)
	}
}

func TestShouldAttempt(t *testing.T) {
	shouldAttempt := shouldAttempt(1)

	if !shouldAttempt {
		t.Error("expected to return true")
	}
}

func TestShouldAttemptWithStrategy(t *testing.T) {
	const attemptNumberShouldReturnFalse = 7

	strategy := func(attempt uint) bool {
		return (attemptNumberShouldReturnFalse != attempt)
	}

	should := shouldAttempt(1, strategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1+attemptNumberShouldReturnFalse, strategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(attemptNumberShouldReturnFalse, strategy)

	if should {
		t.Error("expected to return false")
	}
}

func TestShouldAttemptWithMultipleStrategies(t *testing.T) {
	trueStrategy := func(attempt uint) bool {
		return true
	}

	falseStrategy := func(attempt uint) bool {
		return false
	}

	should := shouldAttempt(1, trueStrategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1, falseStrategy)

	if should {
		t.Error("expected to return false")
	}

	should = shouldAttempt(1, trueStrategy, trueStrategy, trueStrategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1, falseStrategy, falseStrategy, falseStrategy)

	if should {
		t.Error("expected to return false")
	}

	should = shouldAttempt(1, trueStrategy, trueStrategy, falseStrategy)

	if should {
		t.Error("expected to return false")
	}
}

func Example() {
	Retry(func(attempt uint) error {
		return nil // Do something that may or may not cause an error
	})
}

func Example_fileOpen() {
	const logFilePath = "/var/log/myapp.log"

	var logFile *os.File

	err := Retry(func(attempt uint) error {
		var err error

		logFile, err = os.Open(logFilePath)

		return err
	})

	if nil != err {
		log.Fatalf("Unable to open file %q with error %q", logFilePath, err)
	}
}

func Example_httpGetWithStrategies() {
	var response *http.Response

	action := func(attempt uint) error {
		var err error

		response, err = http.Get("https://api.github.com/repos/Rican7/retry")

		if nil == err && nil != response && response.StatusCode > 200 {
			err = fmt.Errorf("failed to fetch (attempt #%d) with status code: %d", attempt, response.StatusCode)
		}

		return err
	}

	err := Retry(
		action,
		strategy.Limit(5),
		strategy.Backoff(backoff.Fibonacci(10*time.Millisecond)),
	)

	if nil != err {
		log.Fatalf("Failed to fetch repository with error %q", err)
	}
}

func Example_withBackoffJitter() {
	action := func(attempt uint) error {
		return errors.New("something happened")
	}

	seed := time.Now().UnixNano()
	random := rand.New(rand.NewSource(seed))

	Retry(
		action,
		strategy.Limit(5),
		strategy.BackoffWithJitter(
			backoff.BinaryExponential(10*time.Millisecond),
			jitter.Deviation(random, 0.5),
		),
	)
}
