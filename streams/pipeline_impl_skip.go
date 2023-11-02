/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type skipPipeline[T any] struct {
	StatefulOp[T]
	skip  int
	count int
}

func NewSkipPipeline[T any](skip int) StatefulOp[T] {
	return &skipPipeline[T]{
		StatefulOp: NewStatefulPipeline[T](),
		skip:       skip,
	}
}

func (r *skipPipeline[T]) SetNext(pipeline Pipeline[T]) {
	r.StatefulOp.SetNext(pipeline)
}

func (r *skipPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.StatefulOp.SetDelegate(delegate)
}

func (r *skipPipeline[T]) Accept(t T) {
	r.count++
	if r.count > r.skip {
		r.StatefulOp.Accept(t)
	}
}
