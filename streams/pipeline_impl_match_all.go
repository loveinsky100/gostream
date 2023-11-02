/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type matchAllPipeline[T any] struct {
	*matchPipeline[T]
	allMatch bool
}

func NewMatchAllPipeline[T any](predicate PredicateHandler[T]) ResultFinishOp[T, bool] {
	return &matchAllPipeline[T]{
		matchPipeline: NewMatchPipeline[T](predicate),
	}
}

func (r *matchAllPipeline[T]) Accept(t T) {
	r.allMatch = r.predicate.Test(t)
	if !r.allMatch {
		r.stop = true
	}
}

func (r *matchAllPipeline[T]) End() {
	r.result = r.allMatch
}
