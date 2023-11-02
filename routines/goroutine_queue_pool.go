/**
 * @author leo
 * @date 2020/10/25 3:39 下午
 */
package routines

import (
	"errors"
	"github.com/loveinsky100/gostreams/collections"
	"sync"
)

type QueuedGoRoutinePool struct {
	mutex            sync.Mutex
	corePoolSize     int
	runningTaskCount int
	queueCapacity    int
	queue            collections.Queue[Future]
	reject           RejectedHandler
	delegate         GoRoutinePoolDelegate
}

//
// NewQueuedGoRoutinePool
// @Description: create new blocking go Routine pool use QueuedGoRoutinePool
// @param corePoolSize int the max pool size with the pool
// @param queueCapacity int the waiting queue size
// @param reject RejectedHandler reject when pool is full, call add will reject
// @return GoRoutinePool
//
func NewQueuedGoRoutinePool(corePoolSize int, queueCapacity int, reject RejectedHandler) GoRoutinePool {
	if corePoolSize <= 0 {
		corePoolSize = 1
	}

	if queueCapacity < 0 {
		corePoolSize = 0
	}

	return &QueuedGoRoutinePool{
		corePoolSize:     corePoolSize,
		queueCapacity:    queueCapacity,
		queue:            collections.NewLinkedQueue[Future](),
		reject:           reject,
		runningTaskCount: 0,
	}
}

func (pool *QueuedGoRoutinePool) Add(callable Callable) (Future, error) {
	future := NewCallableFuture(callable)
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	// running task is full
	if pool.runningTaskCount >= pool.corePoolSize {
		// put task into waiting queue
		success := pool.offerFuture(future)
		if !success {
			// if queue is full, reject it
			if nil != pool.reject {
				return nil, pool.reject.Reject(callable)
			}

			return nil, errors.New("add pool failed due to queue is full")
		}

		// add future into queue
		return future, nil
	}

	// submit future to pool
	return pool.submitFuture(future)
}

func (pool *QueuedGoRoutinePool) AddHandler(handler func() (interface{}, error)) (Future, error) {
	callable := &HandlerCallable{
		Handler: handler,
	}

	return pool.Add(callable)
}

func (pool *QueuedGoRoutinePool) SetDelegate(delegate GoRoutinePoolDelegate) {
	pool.delegate = delegate
}

func (pool *QueuedGoRoutinePool) NumGoroutine() int {
	return pool.runningTaskCount
}

func (pool *QueuedGoRoutinePool) execute(future Future) {
	defer pool.finishFuture()

	future.Execute()
}

func (pool *QueuedGoRoutinePool) offerFuture(future Future) bool {
	if pool.queue.Size() >= pool.queueCapacity {
		return false
	}

	pool.queue.Offer(future)
	return true
}

func (pool *QueuedGoRoutinePool) pollFuture() Future {
	data, ok := pool.queue.Poll()
	if !ok {
		return nil
	}

	future, ok := data.(Future)
	if !ok {
		return nil
	}

	return future
}

func (pool *QueuedGoRoutinePool) submitFuture(future Future) (Future, error) {
	pool.runningTaskCount++
	if nil == pool.delegate {
		go pool.execute(future)
	} else {
		pool.delegate.Execute(func() {
			pool.execute(future)
		})
	}
	return future, nil
}

func (pool *QueuedGoRoutinePool) finishFuture() {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	pool.runningTaskCount--
	// peek task from queue
	peek := pool.pollFuture()
	if nil != peek {
		pool.submitFuture(peek)
	}
}
