/**
 * @author leo
 * @date 2023/4/26 11:39
 */
package streams

import (
	"sync"
)

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type sequentialPipeline[T any] struct {
	StatefulOp[T]
	lock    sync.Mutex
	values  []T
	wait    func() error
	timeout bool
}

func NewSequentialPipeline[T any](wait func() error) StatefulOp[T] {
	return &sequentialPipeline[T]{
		StatefulOp: NewStatefulPipeline[T](),
		wait:       wait,
	}
}

func (r *sequentialPipeline[T]) SetNext(pipeline Pipeline[T]) {
	r.StatefulOp.SetNext(pipeline)
}

func (r *sequentialPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.StatefulOp.SetDelegate(delegate)
}

func (r *sequentialPipeline[T]) Begin(size int) {
	r.values = make([]T, 0, size)
}

func (r *sequentialPipeline[T]) Accept(t T) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.timeout {
		return
	}
	r.values = append(r.values, t)
}

func (r *sequentialPipeline[T]) End() {
	err := r.wait()
	if nil != err {
		r.StatefulOp.PushError(r, err)
	}

	r.lock.Lock()
	defer r.lock.Unlock()
	r.timeout = true
	r.StatefulOp.Begin(len(r.values))
	for _, v := range r.values {
		if r.StatefulOp.CancellationRequested() {
			break
		}

		r.StatefulOp.Accept(v)
	}

	r.StatefulOp.End()
}
