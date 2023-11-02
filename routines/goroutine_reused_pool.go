/**
 * @author leo
 * @date 2020/8/27 8:52 下午
 */
package routines

import (
	"time"
)

type ReusedGoRoutineDelegate struct {
	goroutine Goroutine
}

func NewReusedGoRoutineDelegate(workCount int, idle time.Duration, wait bool) GoRoutinePoolDelegate {
	return &ReusedGoRoutineDelegate{
		goroutine: NewReusedWaitGoRoutine(workCount, idle, wait),
	}
}

func (delegate *ReusedGoRoutineDelegate) Execute(runnable func()) {
	delegate.goroutine.Go(runnable)
}

type ReusedGoRoutinePool struct {
	executorPool   GoRoutinePool
	coreReusedSize int
	coreReusedIdle time.Duration
	wait           bool
}

//
// NewReusedGoRoutinePool
// @Description: 创建一个可复用协程的协程池，如果可复用协程资源不足，则等待复用协程执行完毕继续执行
// @param coreReusedSize int
// @param coreReusedIdle time.Duration
// @param RoutinePool GoRoutinePool
// @return GoRoutinePool
//
func NewReusedGoRoutinePool(coreReusedSize int, coreReusedIdle time.Duration, RoutinePool GoRoutinePool) GoRoutinePool {
	RoutinePool.SetDelegate(NewReusedGoRoutineDelegate(coreReusedSize, coreReusedIdle, true))
	return &ReusedGoRoutinePool{
		executorPool:   RoutinePool,
		coreReusedSize: coreReusedSize,
		coreReusedIdle: coreReusedIdle,
		wait:           true,
	}
}

//
// NewSimpleReusedGoRoutinePool
// @Description: 创建一个可复用协程的协程池，如果可复用协程资源不足，则直接创建新的协程执行
// @param coreReusedSize int
// @param coreReusedIdle time.Duration
// @param RoutinePool GoRoutinePool
// @return GoRoutinePool
//
func NewSimpleReusedGoRoutinePool(coreReusedSize int, coreReusedIdle time.Duration, RoutinePool GoRoutinePool) GoRoutinePool {
	RoutinePool.SetDelegate(NewReusedGoRoutineDelegate(coreReusedSize, coreReusedIdle, false))
	return &ReusedGoRoutinePool{
		executorPool:   RoutinePool,
		coreReusedSize: coreReusedSize,
		coreReusedIdle: coreReusedIdle,
		wait:           false,
	}
}

func (pool *ReusedGoRoutinePool) Add(callable Callable) (Future, error) {
	return pool.executorPool.Add(callable)
}

func (pool *ReusedGoRoutinePool) AddHandler(handler func() (interface{}, error)) (Future, error) {
	return pool.executorPool.AddHandler(handler)
}

func (pool *ReusedGoRoutinePool) SetDelegate(delegate GoRoutinePoolDelegate) {
	pool.executorPool.SetDelegate(delegate)
}

func (pool *ReusedGoRoutinePool) NumGoroutine() int {
	if !pool.wait {
		return pool.executorPool.NumGoroutine()
	}

	if pool.executorPool.NumGoroutine() > pool.coreReusedSize {
		return pool.coreReusedSize
	}

	return pool.executorPool.NumGoroutine()
}