/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type distinctPipeline[T any] struct {
	StatefulOp[T]
	exists map[interface{}]bool
	mapper Function[T, interface{}]
}

func NewDistinctPipeline[T any](mapper FunctionHandler[T, interface{}]) StatefulOp[T] {
	return &distinctPipeline[T]{
		StatefulOp: NewStatefulPipeline[T](),
		mapper:     NewHandlerFunction(mapper),
	}
}

func (r *distinctPipeline[T]) SetNext(pipeline Pipeline[T]) {
	r.StatefulOp.SetNext(pipeline)
}

func (r *distinctPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.StatefulOp.SetDelegate(delegate)
}

func (r *distinctPipeline[T]) Begin(size int) {
	r.exists = map[interface{}]bool{}
}

func (r *distinctPipeline[T]) Accept(t T) {
	value := r.mapper.Apply(t)
	if !r.exists[value] {
		r.StatefulOp.Accept(t)
		r.exists[value] = true
	}
}
