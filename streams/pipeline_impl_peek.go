/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type peekPipeline[T any] struct {
	StatelessOp[T]
	consumer Consumer[T]
}

func NewPeekPipeline[T any](consumer ConsumeHandler[T]) StatelessOp[T] {
	return &peekPipeline[T]{
		StatelessOp: NewStatelessPipeline[T](),
		consumer:    NewHandlerConsumer(consumer),
	}
}

func (r *peekPipeline[T]) SetNext(pipeline Pipeline[T]) {
	r.StatelessOp.SetNext(pipeline)
}

func (r *peekPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.StatelessOp.SetDelegate(delegate)
}

func (r *peekPipeline[T]) Accept(t T) {
	r.consumer.Accept(t)
	r.StatelessOp.Accept(t)
}
