package pool

import (
	"sync"

	"github.com/MarcoVitangeli/workerpool/worker"
)

type Pool struct {
	size uint
	jobs []interface{}
}

func NewPool(s uint, j ...interface{}) *Pool {
	return &Pool{
		size: s,
		jobs: j,
	}
}

func (p *Pool) Run(w worker.Worker) []error {
	var (
		workers = min(p.size, uint(len(p.jobs)))
		wQueue  = make(chan struct{}, workers)
		errCh   = make(chan error, len(p.jobs))
		wg      = &sync.WaitGroup{}
	)

	for _, j := range p.jobs {
		wQueue <- struct{}{}
		wg.Add(1)
		go func(j any, wg *sync.WaitGroup) {
			defer wg.Done()
			errCh <- w.Do(j)
			<-wQueue
		}(j, wg)
	}

	wg.Wait()
	close(errCh)
	return p.drain(errCh)
}

// drain reads from the channel until it is closed
// and returns an array with its values
func (p *Pool) drain(errCh chan error) []error {
	errs := make([]error, 0, len(p.jobs))

	for err := range errCh {
		errs = append(errs, err)
	}

	return errs
}

func min(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}
