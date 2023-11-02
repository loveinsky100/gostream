/**
 * @author leo
 * @date 2023/4/17 12:07
 */
package streams

type Number interface {
	int8 | int16 | int | int32 | int64 | float32 | float64 |
		uint8 | uint16 | uint | uint32 | uint64
}

type ComparatorHandler[T any] func(o1, o2 T) bool

type handlerComparator[T any] struct {
	handler ComparatorHandler[T]
}

func NewHandlerComparator[T any](handler ComparatorHandler[T]) Comparator[T] {
	return &handlerComparator[T]{
		handler: handler,
	}
}

func (c *handlerComparator[T]) Compare(o1 T, o2 T) bool {
	return c.handler(o1, o2)
}

//
// NumericDesc
//  @Description: 自然顺序倒序
//  @return ComparatorHandler[T]
//
func NumericDesc[T Number]() ComparatorHandler[T] {
	return func(o1, o2 T) bool {
		return o1 > o2
	}
}

//
// Numeric
//  @Description: 自然顺序顺序
//  @return ComparatorHandler[T]
//
func Numeric[T Number]() ComparatorHandler[T] {
	return func(o1, o2 T) bool {
		return o1 < o2
	}
}
