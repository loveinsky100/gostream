/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

import (
	"errors"
	"fmt"
	"runtime/debug"
)

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type recoverPipeline[T any] struct {
	StatelessOp[T]
}

func NewRecoverPipeline[T any]() StatelessOp[T] {
	return &recoverPipeline[T]{
		StatelessOp: NewStatelessPipeline[T](),
	}
}

func (r *recoverPipeline[T]) SetNext(pipeline Pipeline[T]) {
	r.StatelessOp.SetNext(pipeline)
}

func (r *recoverPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.StatelessOp.SetDelegate(delegate)
}

func (r *recoverPipeline[T]) Accept(t T) {
	defer r.tryRecover()
	r.StatelessOp.Accept(t)
}

func (r *recoverPipeline[T]) tryRecover() {
	err := recover()
	if err != nil {
		err := errors.New(fmt.Sprintf("an error occured when execute stream task, recover panic: %+v, debug stack: %+v", err, string(debug.Stack())))
		r.StatelessOp.PushError(r, err)
	}
}
