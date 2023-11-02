/**
 * @author leo
 * @date 2023/4/27 13:01
 */
package streams

type sumReduce[T Number] struct {
	from  T
	adder T
}

func Sum[T Number]() Reduce[T, interface{}] {
	return SumFrom[T](0)
}

func SumFrom[T Number](from T) Reduce[T, interface{}] {
	return &sumReduce[T]{
		from:  from,
		adder: 0,
	}
}

func (r *sumReduce[T]) Begin(size int) {
}

func (r *sumReduce[T]) Accept(t T) {
	r.adder += t
}

func (r *sumReduce[T]) End() {

}

func (r *sumReduce[T]) Result() interface{} {
	return r.adder + r.from
}

func (r *sumReduce[T]) CancellationRequested() bool {
	return false
}
