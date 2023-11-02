/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type reducePipeline[T any, R any] struct {
	ResultFinishOp[T, R]
	reduce Reduce[T, R]
}

func NewReducePipeline[T any, R any](reduce Reduce[T, R]) ResultFinishOp[T, R] {
	return &reducePipeline[T, R]{
		ResultFinishOp: NewResultFinishPipeline[T, R](),
		reduce:         reduce,
	}
}

func (r *reducePipeline[T, R]) SetNext(pipeline Pipeline[T]) {
	//
}

func (r *reducePipeline[T, R]) SetDelegate(delegate PipelineDelegate[T]) {
	r.ResultFinishOp.SetDelegate(delegate)
}

func (r *reducePipeline[T, R]) Begin(size int) {
	r.reduce.Begin(size)
}

func (r *reducePipeline[T, R]) Accept(t T) {
	r.reduce.Accept(t)
}

func (r *reducePipeline[T, R]) End() {
	r.reduce.End()
}

func (r *reducePipeline[T, R]) CancellationRequested() bool {
	return r.reduce.CancellationRequested()
}

func (r *reducePipeline[T, R]) Result() R {
	return r.reduce.Result()
}
