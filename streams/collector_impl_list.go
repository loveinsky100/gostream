/**
 * @author leo
 * @date 2023/4/24 21:06
 */
package streams

type listCollector[T any] struct {
	result []T
}

func ToList[T any]() Collector[T, []T] {
	return &listCollector[T]{}
}

func (r *listCollector[T]) Begin(size int) {
	r.result = make([]T, 0, size)
}

func (r *listCollector[T]) Accept(t T) {
	r.result = append(r.result, t)
}

func (r *listCollector[T]) End() {

}

func (r *listCollector[T]) Result() []T {
	return r.result
}

func (r *listCollector[T]) CancellationRequested() bool {
	return false
}
