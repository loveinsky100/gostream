/**
 * @author leo
 * @date 2021/4/26 3:13 下午
 */
package routines

import (
	"fmt"
	"github.com/loveinsky100/gostreams/collections"
	"runtime/debug"
	"sync"
	"time"
)

type _ReusedGoRoutine struct {
	// 空闲时间
	idle time.Duration

	// 并发数
	workCount int

	// 空闲的任务协程
	idleGoroutine map[int]Goroutine

	// 运行中的任务协程
	workGoroutine map[int]Goroutine

	// 等待的任务
	waitQueue collections.Queue[interface{}]

	// 任务协程的锁
	goroutineMutex sync.Mutex

	// 等待队列的读写锁
	waitQueueRw sync.RWMutex

	// 是否开启等待队列
	wait bool
}

func NewReusedGoRoutine(workCount int, idle time.Duration) Goroutine {
	return &_ReusedGoRoutine{
		idle:          idle,
		workCount:     workCount,
		waitQueue:     collections.NewLinkedQueue[interface{}](),
		idleGoroutine: make(map[int]Goroutine),
		workGoroutine: make(map[int]Goroutine),
	}
}

func NewReusedWaitGoRoutine(workCount int, idle time.Duration, wait bool) Goroutine {
	return &_ReusedGoRoutine{
		idle:          idle,
		workCount:     workCount,
		waitQueue:     collections.NewLinkedQueue[interface{}](),
		idleGoroutine: make(map[int]Goroutine),
		workGoroutine: make(map[int]Goroutine),
		wait:          wait,
	}
}

func (goroutine *_ReusedGoRoutine) Go(runnable func()) {
	// 添加到等待队列中，然后尝试执行队列中的任务
	if goroutine.wait {
		goroutine.waitQueueRw.Lock()
		goroutine.waitQueue.Offer(runnable)
		goroutine.waitQueueRw.Unlock()

		goroutine.checkAndRunWaitTask()
		return
	}

	// 直接取出空闲的协程进行运行，如果没有空闲的协程则再创建一个
	workGoroutine, workIndex := goroutine.popWorkGoroutine()
	if nil == workGoroutine {
		go goroutine.safeExecuteRunnable(runnable)
		return
	}

	workGoroutine.Go(func() {
		goroutine.run(workIndex, runnable)
	})
}

func (goroutine *_ReusedGoRoutine) run(goroutineIndex int, runnable func()) {
	defer func() {
		goroutine.recycle(goroutineIndex)
		goroutine.checkAndRunWaitTask()
	}()

	goroutine.safeExecuteRunnable(runnable)
}

func (goroutine *_ReusedGoRoutine) popWorkGoroutine() (Goroutine, int) {
	goroutine.goroutineMutex.Lock()
	defer goroutine.goroutineMutex.Unlock()

	var workGoroutine Goroutine
	var workIndex int
	for index := 0; index < goroutine.workCount; index++ {
		idleGoroutine, ok := goroutine.idleGoroutine[index]
		if ok {
			workGoroutine = idleGoroutine
			delete(goroutine.idleGoroutine, index)
			goroutine.workGoroutine[index] = workGoroutine
			workIndex = index
			break
		}

		_, ok = goroutine.workGoroutine[index]
		if !ok {
			// create new work goroutine
			workGoroutine = NewSingleGoroutine(goroutine.idle)
			goroutine.workGoroutine[index] = workGoroutine
			workIndex = index
			break
		}
	}

	return workGoroutine, workIndex
}

func (goroutine *_ReusedGoRoutine) recycle(workIndex int) {
	goroutine.goroutineMutex.Lock()
	defer goroutine.goroutineMutex.Unlock()

	current, ok := goroutine.workGoroutine[workIndex]
	if !ok {
		return
	}

	delete(goroutine.workGoroutine, workIndex)
	goroutine.idleGoroutine[workIndex] = current

	// 不再删除idleGoroutine了，普通的对象占用内存不高
}

func (goroutine *_ReusedGoRoutine) checkAndRunWaitTask() {
	// 获取一个可执行的协程
	workGoroutine, workIndex := goroutine.popWorkGoroutine()
	if nil == workGoroutine {
		// 如果不存在则说明当前协程都在执行中，只要其中一个执行完毕，又会重新判断
		return
	}

	// 如果存在，则取出队列中的第一个任务
	goroutine.waitQueueRw.RLock()
	waitRunnable, ok := goroutine.waitQueue.Poll()
	goroutine.waitQueueRw.RUnlock()
	if !ok {
		// 不存在可执行的任务，协程重新回收
		goroutine.recycle(workIndex)
		return
	}

	workGoroutine.Go(func() {
		goroutine.run(workIndex, waitRunnable.(func()))
	})

	goroutine.checkAndRunWaitTask()
}

func (goroutine *_ReusedGoRoutine) safeExecuteRunnable(runnable func()) {
	defer func() {
		err := recover()
		if nil != err {
			fmt.Println(fmt.Sprintf("an error occured when execute task: %+v, recover panic: %+v, debug stack: %+v", goroutine, err, string(debug.Stack())))
		}
	}()

	runnable()
}
