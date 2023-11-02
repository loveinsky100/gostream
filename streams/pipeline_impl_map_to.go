/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type mapperToPipeline[T any, R any] struct {
	StatelessOp[T]
	relay  Pipeline[R]
	mapper Function[T, R]
}

func NewMapperToPipeline[T any, R any](relay StatelessOp[R], mapper FunctionHandler[T, R]) StatelessOp[T] {
	return &mapperToPipeline[T, R]{
		StatelessOp: NewStatelessPipeline[T](),
		relay:       relay,
		mapper:      NewHandlerFunction(mapper),
	}
}

func (r *mapperToPipeline[T, R]) SetNext(pipeline Pipeline[T]) {
	r.StatelessOp.SetNext(pipeline)
}

func (r *mapperToPipeline[T, R]) Accept(t T) {
	r.relay.Accept(r.mapper.Apply(t))
}

func (r *mapperToPipeline[T, R]) Begin(size int) {
	r.relay.Begin(size)
}

func (r *mapperToPipeline[T, R]) End() {
	r.relay.End()
}

func (r *mapperToPipeline[T, R]) CancellationRequested() bool {
	return r.relay.CancellationRequested()
}

func (r *mapperToPipeline[T, R]) SetDelegate(delegate PipelineDelegate[T]) {
	r.StatelessOp.SetDelegate(delegate)
}
