/**
 * @author leo
 * @date 2020/9/28 9:33 下午
 */
package routines

import (
	"sort"
	"time"
)

//
// GoRoutineGroup
//  @Description: 批量处理请求group
//
type GoRoutineGroup interface {
	//
	// Add
	//  @Description: 添加请求
	//  @param handler处理逻辑
	//
	Add(handler func() error)

	//
	// AddWithTimeout
	//  @Description: 添加请求并配置超时时间
	//  @param timeout 超时时间
	//  @param handler 处理逻辑
	//
	AddWithTimeout(timeout time.Duration, handler func() error)

	//
	// Execute
	//  @Description: 执行
	//  @return error 错误信息，可能是超时/handler中的错误/panic信息
	//
	Execute() error

	//
	// ExecuteInPool
	//  @Description: 通过池子进行支持
	//  @param pool 协程池
	//  @return error 错误信息，可能是超时/handler中的错误/panic信息
	//
	ExecuteInPool(pool GoRoutinePool) error
}

//
// GoRoutineGroupTask
//  @Description: 执行任务信息
//
type GoRoutineGroupTask struct {
	// 执行的处理逻辑
	callable Callable

	// 超时时间，0表示无限
	timeout time.Duration

	// 执行过程中产生的future
	future Future
}

//
// DefaultGoRoutineGroup
//  @Description: 默认的Group实现
//
type DefaultGoRoutineGroup struct {
	// 优先级队列
	PriorityList []*GoRoutineGroupTask
}

//
// NewGoRoutineGroup
//  @Description: 构建
//  @return *DefaultGoRoutineGroup
//
func NewGoRoutineGroup() *DefaultGoRoutineGroup {
	return &DefaultGoRoutineGroup{
		PriorityList: make([]*GoRoutineGroupTask, 0),
	}
}

func (group *DefaultGoRoutineGroup) Add(handler func() error) {
	group.AddWithTimeout(0, handler)
}

func (group *DefaultGoRoutineGroup) AddWithTimeout(timeout time.Duration, handler func() error) {
	callable := &HandlerCallable{
		Handler: func() (interface{}, error) {
			if nil != handler {
				return nil, handler()
			}

			return nil, nil
		},
	}

	if nil == group.PriorityList {
		group.PriorityList = make([]*GoRoutineGroupTask, 0)
	}

	group.PriorityList = append(group.PriorityList, &GoRoutineGroupTask{
		callable: callable,
		timeout:  timeout,
	})
}

func (group *DefaultGoRoutineGroup) Execute() error {
	if len(group.PriorityList) == 0 {
		return nil
	}

	pool := NewGoRoutinePool(len(group.PriorityList), nil)
	return group.ExecuteInPool(pool)
}

func (group *DefaultGoRoutineGroup) ExecuteInPool(pool GoRoutinePool) error {
	if len(group.PriorityList) == 0 {
		return nil
	}

	// sort
	sort.Slice(group.PriorityList, func(i, j int) bool {
		// 0 is max timeout
		if group.PriorityList[i].timeout == 0 {
			return false
		}

		if group.PriorityList[j].timeout == 0 {
			return true
		}

		return group.PriorityList[i].timeout < group.PriorityList[j].timeout
	})

	for _, task := range group.PriorityList {
		future, err := pool.Add(task.callable)
		if nil != err {
			return err
		}

		task.future = future
	}

	for _, task := range group.PriorityList {
		if nil == task.future {
			continue
		}

		_, err := task.future.GetWithTimeout(task.timeout)
		if nil != err {
			return err
		}
	}

	return nil
}
