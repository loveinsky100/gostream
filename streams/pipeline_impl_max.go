/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type maxPipeline[T any] struct {
	ResultFinishOp[T, T]
	comparator Comparator[T]
	max        T
	accepted   bool
}

func NewMaxPipeline[T any](comparator ComparatorHandler[T]) ResultFinishOp[T, T] {
	return &maxPipeline[T]{
		ResultFinishOp: NewResultFinishPipeline[T, T](),
		comparator:     NewHandlerComparator(comparator),
	}
}

func (r *maxPipeline[T]) SetNext(pipeline Pipeline[T]) {
}

func (r *maxPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.ResultFinishOp.SetDelegate(delegate)
}

func (r *maxPipeline[T]) Accept(t T) {
	if !r.accepted {
		r.accepted = true
		r.max = t
		return
	}

	// t 比 r.max大
	if r.comparator.Compare(r.max, t) {
		r.max = t
	}
}

func (r *maxPipeline[T]) Result() T {
	return r.max
}
