/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type countPipeline[T any] struct {
	ResultFinishOp[T, int]
	count int
}

func NewCountPipeline[T any]() ResultFinishOp[T, int] {
	return &countPipeline[T]{
		ResultFinishOp: NewResultFinishPipeline[T, int](),
	}
}

func (r *countPipeline[T]) SetNext(pipeline Pipeline[T]) {
}

func (r *countPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.ResultFinishOp.SetDelegate(delegate)
}

func (r *countPipeline[T]) Accept(t T) {
	r.count++
}

func (r *countPipeline[T]) Result() int {
	return r.count
}
