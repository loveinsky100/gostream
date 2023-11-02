/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type matchNonePipeline[T any] struct {
	*matchPipeline[T]
}

func NewMatchNonePipeline[T any](predicate PredicateHandler[T]) ResultFinishOp[T, bool] {
	pipeline := &matchNonePipeline[T]{
		matchPipeline: NewMatchPipeline[T](predicate),
	}

	pipeline.result = true
	return pipeline
}

func (r *matchNonePipeline[T]) Accept(t T) {
	if r.predicate.Test(t) {
		r.stop = true
		r.result = false
	}
}
