/**
 * @author leo
 * @date 2023/4/24 17:40
 */
package streams

type Entry[K comparable, V any] interface {
	Key() K
	Value() V
}

type entryImpl[K comparable, V any] struct {
	k K
	v V
}

func NewEntry[K comparable, V any](k K, v V) Entry[K, V] {
	return &entryImpl[K, V]{
		k: k,
		v: v,
	}
}

func (e *entryImpl[K, V]) Key() K {
	return e.k
}

func (e *entryImpl[K, V]) Value() V {
	return e.v
}

type mapHead[K comparable, V any] struct {
	StatelessOp[Entry[K, V]]
}

func NewMapHeader[K comparable, V any]() StatefulOp[Entry[K, V]] {
	return &mapHead[K, V]{
		StatelessOp: NewStatelessPipeline[Entry[K, V]](),
	}
}

func (r *mapHead[K, V]) SetDelegate(delegate PipelineDelegate[Entry[K, V]]) {
	r.StatelessOp.SetDelegate(delegate)
}

func (r *mapHead[K, V]) SetNext(pipeline Pipeline[Entry[K, V]]) {
	r.StatelessOp.SetNext(pipeline)
}

func (r *mapHead[K, V]) Accept(t Entry[K, V]) {
	r.StatelessOp.Accept(t)
}
