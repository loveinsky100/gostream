/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type limitPipeline[T any] struct {
	StatefulOp[T]
	limit int
	count int
}

func NewLimitPipeline[T any](limit int) StatefulOp[T] {
	return &limitPipeline[T]{
		StatefulOp: NewStatefulPipeline[T](),
		limit:      limit,
	}
}

func (r *limitPipeline[T]) SetNext(pipeline Pipeline[T]) {
	r.StatefulOp.SetNext(pipeline)
}

func (r *limitPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.StatefulOp.SetDelegate(delegate)
}

func (r *limitPipeline[T]) Accept(t T) {
	r.count++
	r.StatefulOp.Accept(t)
}

func (r *limitPipeline[T]) CancellationRequested() bool {
	return r.count >= r.limit
}
