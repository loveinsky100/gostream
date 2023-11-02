/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type filterPipeline[T any] struct {
	StatelessOp[T]
	predicate Predicate[T]
}

func NewFilterPipeline[T any](predicate PredicateHandler[T]) StatelessOp[T] {
	return &filterPipeline[T]{
		StatelessOp: NewStatelessPipeline[T](),
		predicate:   NewHandlerPredicate(predicate),
	}
}

func (r *filterPipeline[T]) SetNext(pipeline Pipeline[T]) {
	r.StatelessOp.SetNext(pipeline)
}

func (r *filterPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.StatelessOp.SetDelegate(delegate)
}

func (r *filterPipeline[T]) Accept(t T) {
	if r.predicate.Test(t) {
		r.StatelessOp.Accept(t)
	}
}
