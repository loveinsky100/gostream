/**
 * @author leo
 * @date 2020/8/26 3:57 下午
 */
package routines

import (
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

type FutureStatus int

const (
	// 新建状态
	NEW FutureStatus = iota

	// 被取消
	CANCEL FutureStatus = iota

	// 运行中
	RUNNING FutureStatus = iota

	// 执行完成
	FINISH FutureStatus = iota

	// 出现panic异常
	PANIC FutureStatus = iota
)

//
// Future
//  @Description: future，类似于java的future
//
type Future interface {
	//
	// Get
	//  @Description: wait until task done or panic
	//  @return interface{}
	//  @return error
	//
	Get() (interface{}, error)

	//
	// GetWithTimeout
	//  @Description: wait until task done or timeout, notice: when task timeout it will still running in background until return method
	//  @param timeout
	//  @return interface{}
	//  @return error
	//
	GetWithTimeout(timeout time.Duration) (interface{}, error)

	//
	// Cancel
	//  @Description: cancel task, when task not enter, you can cancel it, true: cancel success false:
	//  @return bool
	//
	Cancel() bool

	//
	// Status
	//  @Description: check is done or in other status
	//  @return FutureStatus
	//
	Status() FutureStatus

	//
	// Execute
	//  @Description: execute task, you need not call this when you put it in Routines pool
	//
	Execute()
}

type CallableResult struct {
	result       interface{}
	errorMessage error
	status       FutureStatus
}

//
// CallableFuture
//  @Description: future的默认实现
//
type CallableFuture struct {
	// 读写锁
	RW sync.RWMutex

	// 具体执行的逻辑
	callable Callable

	// 响应结果的channel
	respChannel chan *CallableResult

	// 当前执行状态
	status FutureStatus

	// 执行结果错误信息
	errorMessage error

	// 执行结果
	result interface{}

	// 开始执行时间
	startTime time.Time

	// 结束执行时间
	endTime time.Time
}

//
// NewCallableFuture
//  @Description: 构建默认的future
//  @param callable
//  @return Future
//
func NewCallableFuture(callable Callable) Future {
	statusChan := make(chan *CallableResult, 1)
	return &CallableFuture{
		callable:     callable,
		respChannel:  statusChan,
		status:       NEW,
		errorMessage: nil,
		result:       nil,
		startTime:    time.Now(),
		endTime:      time.Now(),
	}
}

func (future *CallableFuture) Get() (interface{}, error) {
	return future.GetWithTimeout(0)
}

func (future *CallableFuture) GetWithTimeout(timeout time.Duration) (interface{}, error) {
	future.RW.RLock()
	if future.status == CANCEL || future.status > RUNNING {
		future.RW.RUnlock()
		return future.result, future.errorMessage
	}

	future.RW.RUnlock()

	if 0 == timeout {
		select {
		case resp := <-future.respChannel:
			future.innerDone(resp)
			return future.innerReadResult()
		}
	} else {
		distance := time.Until(future.startTime)
		currentTimeout := timeout + distance
		timer := NewTimer(currentTimeout)
		defer Stop(timer)
		select {
		case resp := <-future.respChannel:
			future.innerDone(resp)
			return future.innerReadResult()
		case <-timer.C:
			return nil, errors.New(fmt.Sprintf("future task timeout(%s)", timeout.String()))
		}
	}
}

func (future *CallableFuture) Cancel() bool {
	future.RW.RLock()
	if future.status > CANCEL {
		future.RW.RUnlock()
		return false
	}

	future.RW.RUnlock()

	future.RW.Lock()
	defer future.RW.Unlock()

	// double check, when get write lock, status maybe changed
	if future.status > CANCEL {
		return false
	}

	future.status = CANCEL
	future.errorMessage = errors.New("future task canceled")
	future.endTime = time.Now()
	close(future.respChannel)
	return true
}

func (future *CallableFuture) Status() FutureStatus {
	future.RW.RLock()
	defer future.RW.RUnlock()

	return future.status
}

func (future *CallableFuture) Execute() {
	// check if it is canceled or in other status
	future.RW.RLock()
	status := future.status
	if status > NEW {
		future.RW.RUnlock()
		return
	}

	future.RW.RUnlock()
	future.RW.Lock()
	if future.status > NEW {
		future.RW.Unlock()
		return
	}

	future.status = RUNNING
	future.RW.Unlock()

	future.innerExecute()
}

func (future *CallableFuture) String() string {
	now := time.Now().Unix()
	return fmt.Sprintf("[Future] start: %+v, end: %+v, cost: %d status: %d, invocation: %+v", future.startTime, future.endTime, now-future.startTime.Unix(), future.status, future.callable)
}

func (future *CallableFuture) innerExecute() {
	callableResult := &CallableResult{}

	defer func() {
		err := recover()
		if err != nil {
			callableResult.status = PANIC
			callableResult.errorMessage = errors.New(fmt.Sprintf("an error occured when execute task: %+v, recover panic: %+v, debug stack: %+v", future.callable, err, string(debug.Stack())))
		}

		future.respChannel <- callableResult
		close(future.respChannel)
	}()

	if nil == future.callable {
		return
	}

	resp, err := future.callable.Call()
	callableResult.result = resp
	callableResult.errorMessage = err
	callableResult.status = FINISH
}

func (future *CallableFuture) innerDone(resp *CallableResult) {
	future.RW.Lock()
	defer future.RW.Unlock()

	future.status = resp.status
	future.errorMessage = resp.errorMessage
	future.result = resp.result
	future.endTime = time.Now()
}

func (future *CallableFuture) innerReadResult() (interface{}, error) {
	future.RW.RLock()
	defer future.RW.RUnlock()

	return future.result, future.errorMessage
}
