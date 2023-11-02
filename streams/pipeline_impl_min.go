/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type minPipeline[T any] struct {
	ResultFinishOp[T, T]
	comparator Comparator[T]
	min        T
	accepted   bool
}

func NewMinPipeline[T any](comparator ComparatorHandler[T]) ResultFinishOp[T, T] {
	return &minPipeline[T]{
		ResultFinishOp: NewResultFinishPipeline[T, T](),
		comparator:     NewHandlerComparator(comparator),
	}
}

func (r *minPipeline[T]) SetNext(pipeline Pipeline[T]) {
}

func (r *minPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.ResultFinishOp.SetDelegate(delegate)
}

func (r *minPipeline[T]) Accept(t T) {
	if !r.accepted {
		r.accepted = true
		r.min = t
		return
	}

	// !(t 比 r.min 大)
	if !r.comparator.Compare(r.min, t) {
		r.min = t
	}
}

func (r *minPipeline[T]) Result() T {
	return r.min
}
