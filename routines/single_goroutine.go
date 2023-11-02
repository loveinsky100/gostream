/**
 * @author leo
 * @date 2021/4/26 3:15 下午
 */
package routines

import (
	"fmt"
	"github.com/loveinsky100/gostreams/collections"
	"runtime/debug"
	"sync"
	"time"
)

type GoroutineStatus int

const (
	// 未创建
	NOT_CREATE GoroutineStatus = iota

	// 准备中
	READY GoroutineStatus = iota

	// 空闲中
	IDLE GoroutineStatus = iota

	// 运行中
	WORK GoroutineStatus = iota
)

type _SingleGoroutine struct {
	// 任务队列
	queue  collections.Queue[func()]
	rw     sync.RWMutex
	idle   time.Duration
	alive  chan bool
	status GoroutineStatus
	timer  time.Timer
}

func NewSingleGoroutine(idle time.Duration) Goroutine {
	aliveChain := make(chan bool, 1)
	return &_SingleGoroutine{
		queue:  collections.NewLinkedQueue[func()](),
		idle:   idle,
		alive:  aliveChain,
		status: NOT_CREATE,
	}
}

func (goroutine *_SingleGoroutine) Go(runnable func()) {
	goroutine.rw.Lock()
	defer goroutine.rw.Unlock()
	goroutine.queue.Offer(runnable)
	if goroutine.status == NOT_CREATE {
		goroutine.status = READY
		go goroutine.run()
	} else if goroutine.status == IDLE {
		goroutine.alive <- true
		goroutine.status = READY
	}
}

func (goroutine *_SingleGoroutine) pool() (func(), bool) {
	goroutine.rw.RLock()
	defer goroutine.rw.RUnlock()

	runnable, ok := goroutine.queue.Poll()
	if !ok {
		return nil, false
	}

	return runnable, true
}

func (goroutine *_SingleGoroutine) run() {
LOOP:
	{
		// read runnable
		runnable, ok := goroutine.pool()
		if ok {
			goroutine.execute(runnable)
			goto LOOP
		}

		// current status is: READY
		goroutine.rw.Lock()
		// double check, READY - IDLE - NOT_CREATE
		_, ok = goroutine.queue.Peek()
		if ok {
			goroutine.rw.Unlock()
			goto LOOP
		}

		goroutine.status = IDLE
		goroutine.rw.Unlock()

		timer := NewTimer(goroutine.idle)
		// wait or timeout, current status is: IDLE
		select {
		case <-goroutine.alive:
			{
				Stop(timer)
				goto LOOP
			}
		case <-timer.C:
			{
				Stop(timer)
				// timeout
				goroutine.rw.Lock()
				// double check, current status is: IDLE
				if goroutine.status == READY {
					// try remove alive
					select {
					case <-goroutine.alive:
						{
						}
					default:
						{
						}
					}

					goroutine.rw.Unlock()
					goto LOOP
				}

				goroutine.status = NOT_CREATE
				goroutine.rw.Unlock()
				return
			}
		}
	}
}

func (goroutine *_SingleGoroutine) execute(runnable func()) {
	goroutine.rw.Lock()
	goroutine.status = WORK
	goroutine.rw.Unlock()

	goroutine.safeExecuteRunnable(runnable)

	goroutine.rw.Lock()
	goroutine.status = READY
	goroutine.rw.Unlock()
}

func (goroutine *_SingleGoroutine) safeExecuteRunnable(runnable func()) {
	defer func() {
		err := recover()
		if nil != err {
			fmt.Println(fmt.Sprintf("an error occured when execute task: %+v, recover panic: %+v, debug stack: %+v", goroutine, err, string(debug.Stack())))
		}
	}()

	runnable()
}
