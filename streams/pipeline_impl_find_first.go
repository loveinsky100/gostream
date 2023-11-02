/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type findFirstPipeline[T any] struct {
	ResultFinishOp[T, T]
	predicate Predicate[T]
	hasValue  bool
	value     T
}

func NewFindFirstPipeline[T any]() ResultFinishOp[T, T] {
	return &findFirstPipeline[T]{
		ResultFinishOp: NewResultFinishPipeline[T, T](),
	}
}

func (r *findFirstPipeline[T]) SetNext(pipeline Pipeline[T]) {
}

func (r *findFirstPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.ResultFinishOp.SetDelegate(delegate)
}

func (r *findFirstPipeline[T]) Accept(t T) {
	r.hasValue = true
	r.value = t
}

func (r *findFirstPipeline[T]) CancellationRequested() bool {
	return r.hasValue
}

func (r *findFirstPipeline[T]) Result() T {
	return r.value
}
