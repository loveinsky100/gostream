/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type referencePipeline[T any] struct {
	next     Pipeline[T]
	delegate PipelineDelegate[T]
}

func newReferencePipeline[T any]() Pipeline[T] {
	return &referencePipeline[T]{}
}

func (r *referencePipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.delegate = delegate
}

func (r *referencePipeline[T]) PushError(source interface{}, err error) {
	if nil != r.delegate {
		r.delegate.OnError(source, err)
	}
}

func (r *referencePipeline[T]) SetNext(pipeline Pipeline[T]) {
	r.next = pipeline
}

func (r *referencePipeline[T]) Accept(t T) {
	r.next.Accept(t)
}

func (r *referencePipeline[T]) Begin(size int) {
	r.next.Begin(size)
}

func (r *referencePipeline[T]) End() {
	r.next.End()
}

func (r *referencePipeline[T]) CancellationRequested() bool {
	return r.next.CancellationRequested()
}

func (r *referencePipeline[T]) Stateless() bool {
	return false
}

//
// statelessPipeline
//  @Description: 无状态流水线
//
type statelessPipeline[T any] struct {
	Pipeline[T]
}

//
// NewStatelessPipeline
//  @Description: 无状态流水线
//  @return StatelessOp[T]
//
func NewStatelessPipeline[T any]() StatelessOp[T] {
	return &statelessPipeline[T]{
		Pipeline: newReferencePipeline[T](),
	}
}

func (r *statelessPipeline[T]) Accept(t T) {
	r.Pipeline.Accept(t)
}

func (r *statelessPipeline[T]) Stateless() bool {
	return true
}

//
// statelessPipeline
//  @Description: 有状态流水线
//
type statefulPipeline[T any] struct {
	Pipeline[T]
}

//
// NewStatefulPipeline
//  @Description: 有状态流水线
//  @return StateOp[T]
//
func NewStatefulPipeline[T any]() StatefulOp[T] {
	return &statefulPipeline[T]{
		Pipeline: newReferencePipeline[T](),
	}
}

func (r *statefulPipeline[T]) Accept(t T) {
	r.Pipeline.Accept(t)
}

func (r *statefulPipeline[T]) Stateless() bool {
	return false
}

//
// statelessPipeline
//  @Description: 有状态流水线
//
type finishPipeline[T any] struct {
	StatefulOp[T]
}

//
// NewFinishPipeline
//  @Description:
//  @return FinishOp[T]
//
func NewFinishPipeline[T any]() FinishOp[T] {
	return &finishPipeline[T]{
		StatefulOp: NewStatefulPipeline[T](),
	}
}

func (f *finishPipeline[T]) SetNext(pipeline Pipeline[T]) {
}

func (f *finishPipeline[T]) Begin(size int) {
}

func (f *finishPipeline[T]) Accept(t T) {
}

func (f *finishPipeline[T]) End() {
}

func (f *finishPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	f.StatefulOp.SetDelegate(delegate)
}

func (f *finishPipeline[T]) CancellationRequested() bool {
	return false
}

//
// statelessPipeline
//  @Description: 有状态流水线
//
type resultFinishPipeline[T any, R any] struct {
	StatefulOp[T]
}

//
// NewResultFinishPipeline
//  @Description:
//  @return FinishOp[T]
//
func NewResultFinishPipeline[T any, R any]() ResultFinishOp[T, R] {
	return &resultFinishPipeline[T, R]{
		StatefulOp: NewStatefulPipeline[T](),
	}
}

func (f *resultFinishPipeline[T, R]) Result() R {
	var r R
	return r
}

func (f *resultFinishPipeline[T, R]) SetNext(pipeline Pipeline[T]) {
}

func (f *resultFinishPipeline[T, R]) SetDelegate(delegate PipelineDelegate[T]) {
	f.StatefulOp.SetDelegate(delegate)
}

func (f *resultFinishPipeline[T, R]) Begin(size int) {
}

func (f *resultFinishPipeline[T, R]) Accept(t T) {
}

func (f *resultFinishPipeline[T, R]) End() {
}

func (f *resultFinishPipeline[T, R]) CancellationRequested() bool {
	return false
}
