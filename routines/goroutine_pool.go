/**
 * @author leo
 * @date 2020/8/27 8:52 下午
 */
package routines

import (
	"errors"
	"sync"
)

//
// GoRoutinePoolDelegate
//  @Description: 执行的代理
//
type GoRoutinePoolDelegate interface {
	//
	// Execute
	// @Description: 执行future
	// @param future runnable
	//
	Execute(runnable func())
}

type DefaultPoolDelegate struct {
}

func NewPoolDelegate() GoRoutinePoolDelegate {
	return &DefaultPoolDelegate{}
}

func (*DefaultPoolDelegate) Execute(runnable func()) {
	go runnable()
}

type GoRoutinePool interface {
	//
	// Add
	// @Description: add callable
	// @param callable Callable
	// @return Future
	// @return error
	//
	Add(callable Callable) (Future, error)

	//
	// AddHandler
	// @Description: add func
	// @param handler func() (interface{}, error)
	// @return Future
	// @return error
	//
	AddHandler(handler func() (interface{}, error)) (Future, error)

	//
	// SetDelegate
	// @Description: set pool delegate
	// @param delegate GoRoutinePoolDelegate
	//
	SetDelegate(delegate GoRoutinePoolDelegate)

	//
	// NumGoroutine
	// @Description: current goroutine count in pool
	// @return int
	//
	NumGoroutine() int
}

type DefaultGoRoutinePool struct {
	mutex    sync.Mutex
	poolSize int
	current  int
	reject   RejectedHandler
	delegate GoRoutinePoolDelegate
}

//
// NewGoRoutinePool
// @Description: create new goRoutine pool use DefaultGoRoutinePool
// @param poolSize int the max pool size with the pool
// @param reject RejectedHandler reject when pool is full, call add will reject
// @return GoRoutinePool
//
func NewGoRoutinePool(poolSize int, reject RejectedHandler) GoRoutinePool {
	pool := &DefaultGoRoutinePool{
		poolSize: poolSize,
		reject:   reject,
	}

	return pool
}

func (pool *DefaultGoRoutinePool) Add(callable Callable) (Future, error) {
	pool.mutex.Lock()
	defer func() {
		pool.mutex.Unlock()
	}()

	if pool.current >= pool.poolSize {
		if nil != pool.reject {
			return nil, pool.reject.Reject(callable)
		}

		return nil, errors.New("add pool failed due to pool full")
	}

	pool.current++
	future := NewCallableFuture(callable)
	if nil == pool.delegate {
		go pool.execute(future)
	} else {
		pool.delegate.Execute(func() {
			pool.execute(future)
		})
	}

	return future, nil
}

func (pool *DefaultGoRoutinePool) AddHandler(handler func() (interface{}, error)) (Future, error) {
	callable := &HandlerCallable{
		Handler: handler,
	}

	return pool.Add(callable)
}

func (pool *DefaultGoRoutinePool) SetDelegate(delegate GoRoutinePoolDelegate) {
	pool.delegate = delegate
}

func (pool *DefaultGoRoutinePool) NumGoroutine() int {
	return pool.current
}

func (pool *DefaultGoRoutinePool) execute(future Future) {
	defer func() {
		pool.mutex.Lock()
		pool.current--
		pool.mutex.Unlock()
	}()

	future.Execute()
}
