/**
 * @author leo
 * @date 2023/4/17 13:20
 */
package streams

import "unsafe"

type eface struct {
	rtype unsafe.Pointer
	data  unsafe.Pointer
}

type PredicateHandler[T any] func(data T) bool

type handlerPredicate[T any] struct {
	handler PredicateHandler[T]
}

func NewHandlerPredicate[T any](handler PredicateHandler[T]) Predicate[T] {
	return &handlerPredicate[T]{
		handler: handler,
	}
}

func (h *handlerPredicate[T]) Test(data T) bool {
	if nil != h.handler {
		return h.handler(data)
	}

	return false
}

func NilValue[T any]() PredicateHandler[T] {
	return func(data T) bool {
		return IsNil(data)
	}
}

func NotNil[T any]() PredicateHandler[T] {
	return func(data T) bool {
		return !IsNil(data)
	}
}

func IsNil(obj interface{}) bool {
	if obj == nil {
		return true
	}
	return unpackEFace(obj).data == nil
}

func unpackEFace(obj interface{}) *eface {
	return (*eface)(unsafe.Pointer(&obj))
}
