/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

import (
	"github.com/loveinsky100/gostreams/routines"
	"time"
)

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type parallelPipeline[T any] struct {
	StatelessOp[T]
	timeout time.Duration
	group   routines.GoRoutineGroup
}

func NewParallelPipeline[T any](timeout time.Duration) (StatelessOp[T], routines.GoRoutineGroup) {
	group := routines.NewGoRoutineGroup()
	return &parallelPipeline[T]{
		StatelessOp: NewStatelessPipeline[T](),
		timeout:     timeout,
		group:       group,
	}, group
}

func (r *parallelPipeline[T]) SetNext(pipeline Pipeline[T]) {
	r.StatelessOp.SetNext(pipeline)
}

func (r *parallelPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.StatelessOp.SetDelegate(delegate)
}

func (r *parallelPipeline[T]) Accept(t T) {
	r.group.AddWithTimeout(r.timeout, func() error {
		r.StatelessOp.Accept(t)
		return nil
	})
}
