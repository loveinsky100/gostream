/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type mapperFlatPipeline[T any, R any] struct {
	StatelessOp[T]
	relay  Pipeline[R]
	mapper Function[T, []R]
}

func NewMapperFlatPipeline[T any, R any](relay StatelessOp[R], mapper FunctionHandler[T, []R]) StatelessOp[T] {
	return &mapperFlatPipeline[T, R]{
		StatelessOp: NewStatelessPipeline[T](),
		relay:       relay,
		mapper:      NewHandlerFunction(mapper),
	}
}

func (r *mapperFlatPipeline[T, R]) SetNext(pipeline Pipeline[T]) {
	r.StatelessOp.SetNext(pipeline)
}

func (r *mapperFlatPipeline[T, R]) Accept(t T) {
	values := r.mapper.Apply(t)
	for _, v := range values {
		r.relay.Accept(v)
	}
}

func (r *mapperFlatPipeline[T, R]) Begin(size int) {
	r.relay.Begin(size)
}

func (r *mapperFlatPipeline[T, R]) End() {
	r.relay.End()
}

func (r *mapperFlatPipeline[T, R]) CancellationRequested() bool {
	return r.relay.CancellationRequested()
}

func (r *mapperFlatPipeline[T, R]) SetDelegate(delegate PipelineDelegate[T]) {
	r.StatelessOp.SetDelegate(delegate)
}
