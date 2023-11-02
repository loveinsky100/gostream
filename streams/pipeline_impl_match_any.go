/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type matchAnyPipeline[T any] struct {
	*matchPipeline[T]
}

func NewMatchAnyPipeline[T any](predicate PredicateHandler[T]) ResultFinishOp[T, bool] {
	return &matchAnyPipeline[T]{
		matchPipeline: NewMatchPipeline[T](predicate),
	}
}

func (r *matchAnyPipeline[T]) Accept(t T) {
	if r.predicate.Test(t) {
		r.result = true
		r.stop = true
	}
}
