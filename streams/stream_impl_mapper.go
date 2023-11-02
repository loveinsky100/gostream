/**
 * @author leo
 * @date 2023/4/24 20:10
 */
package streams

type mapperStreamImpl[T any, R any] struct {
	source *streamImpl[T]
}

func Mapper[T any, R any](stream Stream[T]) interface{} {
	source, _ := stream.(*streamImpl[T])
	return &mapperStreamImpl[T, R]{
		source: source,
	}
}

func (r *mapperStreamImpl[T, R]) Map(mapper FunctionHandler[T, R]) Stream[R] {
	// 中继节点，连接T和R
	// pipeline[T] -> map[T, R] -> relay[R] -> ....
	relay := NewStatelessPipeline[R]()
	pipeline := NewMapperToPipeline[T, R](relay, mapper)
	// 自己作为上一个节点的下一个节点
	r.source.Join(pipeline)
	return r.relaySource(relay)
}

func (r *mapperStreamImpl[T, R]) FlatMap(mapper FunctionHandler[T, []R]) Stream[R] {
	relay := NewStatelessPipeline[R]()
	pipeline := NewMapperFlatPipeline[T, R](relay, mapper)
	r.source.Join(pipeline)
	return r.relaySource(relay)
}

func (r *mapperStreamImpl[T, R]) relaySource(relay Pipeline[R]) Stream[R] {
	// 新增一个中继节点
	return &streamImpl[R]{
		runner:          r.source.runner,
		current:         relay,
		parallel:        r.source.parallel,
		requireParallel: r.source.requireParallel,
		parallelGroup:   r.source.parallelGroup,
		debug:           r.source.debug,
		tracker:         r.source.tracker,
		debugger:        r.source.debugger,
		errorHandler:    r.source.errorHandler,
		pool:            r.source.pool,
	}
}
