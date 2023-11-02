/**
 * @author leo
 * @date 2023/4/24 21:06
 */
package streams

type mapCollector[T any, K comparable, V any] struct {
	result map[K]V
	key    Function[T, K]
	value  Function[T, V]
}

func ToMap[T any, K comparable, V any](key FunctionHandler[T, K], value FunctionHandler[T, V]) Collector[T, map[K]V] {
	return &mapCollector[T, K, V]{
		result: nil,
		key:    NewHandlerFunction(key),
		value:  NewHandlerFunction(value),
	}
}

func (r *mapCollector[T, K, V]) Begin(size int) {
	r.result = map[K]V{}
}

func (r *mapCollector[T, K, V]) Accept(t T) {
	r.result[r.key.Apply(t)] = r.value.Apply(t)
}

func (r *mapCollector[T, K, V]) End() {

}

func (r *mapCollector[T, K, V]) Result() map[K]V {
	return r.result
}

func (r *mapCollector[T, K, V]) SetNext(pipeline Pipeline[T]) {
	//
}

func (r *mapCollector[T, K, V]) CancellationRequested() bool {
	return false
}
