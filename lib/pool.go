package lib

import "sync/atomic"
import "sync"

// JobFn is a type describing a generic worker pool task.
type JobFn[T any] func(params T)

type pool[T any] struct {
	queue chan T
	poolSize uint32
	jobFn JobFn[T]
	doneC chan struct{}
}

// NewPool returns a worker pool that has a defined size and defined generic taks function.
func NewPool[T any] (poolSize uint32, jobFn JobFn[T]) *pool[T] {
	return &pool[T]{
		queue: make(chan T),
		poolSize: poolSize,
		jobFn: jobFn,
		doneC: make(chan struct{}),
	}
}

// Start spawns the worker pool.
func (p *pool[T]) Start() {
	go p.worker()
}

// Do passes the task parameters to the task function and queues it for the execution in the pool.
func (p *pool[T]) Do(task T) {
	p.queue <- task
}

// Done informs the initiator when all the tasks provided to the pool are finished.
func (p *pool[T]) Done() <-chan struct{} {
	close(p.queue)
	return p.doneC
}

func (p *pool[T]) worker() {
	var jobCount uint32
	var wg sync.WaitGroup

	for {
		// Limiting the pool size with the atomic counter
		if atomic.LoadUint32(&jobCount) > p.poolSize -1 {
			continue
		}

		params, ok := <-p.queue
		if !ok {
			break
		}

		// We want to make sure that the tasks are all completed in our pool before reporting that it's done.
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer atomic.AddUint32(&jobCount, ^uint32(0))
			p.jobFn(params)
		}()

		atomic.AddUint32(&jobCount, 1)
	}

	wg.Wait()
	close(p.doneC)
}
