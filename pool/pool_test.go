package pool

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MarcoVitangeli/workerpool/worker"
)

func TestNumCalls(t *testing.T) {
	var (
		c    = uint64(0)
		buff = make([]interface{}, 1000)
		p    = NewPool(10, buff...)
		w    = func(a interface{}) error {
			atomic.AddUint64(&c, 1)
			time.Sleep(time.Millisecond)
			return nil
		}
	)
	p.Run(worker.FromFunc(w))

	if c != 1000 {
		t.Fatalf("expected 1000 calls, found %d", c)
	}
}

func TestNumErrs(t *testing.T) {
	var (
		buff = make([]interface{}, 1000)
		p    = NewPool(10, buff...)
		w    = func(a interface{}) error {
			time.Sleep(time.Millisecond)
			return errors.New("error")
		}
	)
	errs := p.Run(worker.FromFunc(w))

	if len(errs) != 1000 {
		t.Errorf("expected 1000 errors, found %d", len(errs))
	}

	for _, err := range errs {
		if err == nil || err.Error() != "error" {
			t.Errorf("expected 'error', found %s", err)
		}
	}
}
