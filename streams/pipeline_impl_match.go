/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type matchPipeline[T any] struct {
	ResultFinishOp[T, bool]
	predicate Predicate[T]
	result    bool
	stop      bool
}

func NewMatchPipeline[T any](predicate PredicateHandler[T]) *matchPipeline[T] {
	return &matchPipeline[T]{
		ResultFinishOp: NewResultFinishPipeline[T, bool](),
		predicate:      NewHandlerPredicate(predicate),
	}
}

func (r *matchPipeline[T]) SetNext(pipeline Pipeline[T]) {
}

func (r *matchPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.ResultFinishOp.SetDelegate(delegate)
}

func (r *matchPipeline[T]) Accept(t T) {

}

func (r *matchPipeline[T]) CancellationRequested() bool {
	return r.stop
}

func (r *matchPipeline[T]) Result() bool {
	return r.result
}
