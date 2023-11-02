/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

import "sort"

//
// Pipeline
//  @Description: 流水线，实际上是一个链表，每个pipeline存储一个sink用于处理数据
//
type sortPipeline[T any] struct {
	StatefulOp[T]
	comparator Comparator[T]
	values     []T
}

func NewSortPipeline[T any](comparator ComparatorHandler[T]) StatefulOp[T] {
	return &sortPipeline[T]{
		StatefulOp: NewStatefulPipeline[T](),
		comparator: NewHandlerComparator(comparator),
	}
}

func (r *sortPipeline[T]) SetNext(pipeline Pipeline[T]) {
	r.StatefulOp.SetNext(pipeline)
}

func (r *sortPipeline[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.StatefulOp.SetDelegate(delegate)
}

func (r *sortPipeline[T]) Begin(size int) {
	r.values = make([]T, 0, size)
}

func (r *sortPipeline[T]) Accept(t T) {
	r.values = append(r.values, t)
}

func (r *sortPipeline[T]) End() {
	sort.Slice(r.values, func(i, j int) bool {
		dataI := r.values[i]
		dataJ := r.values[j]
		return r.comparator.Compare(dataI, dataJ)
	})

	r.StatefulOp.Begin(len(r.values))
	for _, v := range r.values {
		if r.StatefulOp.CancellationRequested() {
			break
		}
		r.StatefulOp.Accept(v)
	}

	r.StatefulOp.End()
}
