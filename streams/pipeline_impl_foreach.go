/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type foreachPipeline[T any] struct {
	FinishOp[T]
	consumer Consumer[T]
}

func NewForeachPipeline[T any](consumer ConsumeHandler[T]) FinishOp[T] {
	return &foreachPipeline[T]{
		FinishOp: NewFinishPipeline[T](),
		consumer: NewHandlerConsumer(consumer),
	}
}

func (r *foreachPipeline[T]) SetNext(pipeline Pipeline[T]) {
	r.FinishOp.SetNext(pipeline)
}

func (r *foreachPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.FinishOp.SetDelegate(delegate)
}

func (r *foreachPipeline[T]) Accept(t T) {
	r.consumer.Accept(t)
}
