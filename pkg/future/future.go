package future

import (
	"sync"
)

type Future chan error

func New() Future {
	return make(Future)
}

func Run(fn func() error) Future {
	future := New()
	go func() {
		defer close(future)
		future <- fn()
	}()
	return future
}

func Error(err error) Future {
	return Run(func() error { return err })
}

func All(futures ...Future) Future {
	result := New()
	resultEmpty := true
	wg := sync.WaitGroup{}
	wg.Add(len(futures))

	for _, f := range futures {
		go func(f Future) {
			err := <-f
			wg.Done()
			if resultEmpty && err != nil {
				resultEmpty = false
				result <- err
				close(result)
			}
		}(f)
	}

	go func() {
		wg.Wait()
		if resultEmpty {
			result <- nil
			close(result)
		}
	}()

	return result
}
