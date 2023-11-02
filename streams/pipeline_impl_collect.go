/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type collectPipeline[T any, R any] struct {
	ResultFinishOp[T, R]
	collector Collector[T, R]
}

func NewCollectPipeline[T any, R any](collector Collector[T, R]) ResultFinishOp[T, R] {
	return &collectPipeline[T, R]{
		ResultFinishOp: NewResultFinishPipeline[T, R](),
		collector:      collector,
	}
}

func (r *collectPipeline[T, R]) SetNext(pipeline Pipeline[T]) {
	//
}

func (r *collectPipeline[T, R]) SetDelegate(delegate PipelineDelegate[T]) {
	r.ResultFinishOp.SetDelegate(delegate)
}

func (r *collectPipeline[T, R]) Begin(size int) {
	r.collector.Begin(size)
}

func (r *collectPipeline[T, R]) Accept(t T) {
	r.collector.Accept(t)
}

func (r *collectPipeline[T, R]) End() {
	r.collector.End()
}

func (r *collectPipeline[T, R]) CancellationRequested() bool {
	return r.collector.CancellationRequested()
}

func (r *collectPipeline[T, R]) Result() R {
	return r.collector.Result()
}
