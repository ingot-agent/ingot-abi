package lifecycle_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/ingot-agent/ingot-abi/lifecycle"
)

type recorder struct {
	mu       sync.Mutex
	requests []error
	cancel   chan struct{}
}

func newRecorder() *recorder { return &recorder{cancel: make(chan struct{})} }

func (r *recorder) RequestShutdown(err error) {
	r.mu.Lock()
	r.requests = append(r.requests, err)
	first := len(r.requests) == 1
	r.mu.Unlock()
	if first {
		close(r.cancel)
	}
}

func (r *recorder) shutdowns() []error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]error(nil), r.requests...)
}

var _ lifecycle.Controller = (*recorder)(nil)

func TestRequestShutdownIsNonBlockingAndConcurrent(t *testing.T) {
	t.Parallel()
	value := newRecorder()
	var wait sync.WaitGroup
	for index := 0; index < 64; index++ {
		wait.Add(1)
		go func(err error) {
			defer wait.Done()
			value.RequestShutdown(err)
		}(nil)
	}
	wait.Wait()
	if len(value.shutdowns()) != 64 {
		t.Fatalf("recorded %d shutdowns", len(value.shutdowns()))
	}
	select {
	case <-value.cancel:
	default:
		t.Fatal("first request did not trigger the runtime cancellation signal")
	}
}

func TestNilAndNonNilReasonsAreBothRecorded(t *testing.T) {
	t.Parallel()
	value := newRecorder()
	sentinel := errors.New("fatal")
	value.RequestShutdown(nil)
	value.RequestShutdown(sentinel)
	value.RequestShutdown(sentinel)
	requests := value.shutdowns()
	if len(requests) != 3 || requests[0] != nil || !errors.Is(requests[1], sentinel) || !errors.Is(requests[2], sentinel) {
		t.Fatalf("requests = %#v", requests)
	}
}
