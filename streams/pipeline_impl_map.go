/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type mapperPipeline[T any] struct {
	StatelessOp[T]
	mapper Function[T, T]
}

func NewMapperPipeline[T any](mapper FunctionHandler[T, T]) StatelessOp[T] {
	return &mapperPipeline[T]{
		StatelessOp: NewStatelessPipeline[T](),
		mapper:      NewHandlerFunction(mapper),
	}
}

func (r *mapperPipeline[T]) SetNext(pipeline Pipeline[T]) {
	r.StatelessOp.SetNext(pipeline)
}

func (r *mapperPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.StatelessOp.SetDelegate(delegate)
}

func (r *mapperPipeline[T]) Accept(t T) {
	r.StatelessOp.Accept(r.mapper.Apply(t))
}
