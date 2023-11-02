/**
 * @author leo
 * @date 2023/4/24 17:40
 */
package streams

type sliceHead[T any] struct {
	StatelessOp[T]
}

func NewSliceHeader[T any]() StatefulOp[T] {
	return &sliceHead[T]{
		StatelessOp: NewStatelessPipeline[T](),
	}
}

func (r *sliceHead[T]) SetNext(pipeline Pipeline[T]) {
	r.StatelessOp.SetNext(pipeline)
}

func (r *sliceHead[T]) SetDelegate(delegate PipelineDelegate[T]) {
	r.StatelessOp.SetDelegate(delegate)
}

func (r *sliceHead[T]) Accept(t T) {
	r.StatelessOp.Accept(t)
}
