/**
 * @author leo
 * @date 2023/4/17 13:24
 */
package streams

import (
	"fmt"
)

type ConsumeHandler[T any] func(data T)

type handlerConsumer[T any] struct {
	handler ConsumeHandler[T]
}

func NewHandlerConsumer[T any](handler ConsumeHandler[T]) Consumer[T] {
	return &handlerConsumer[T]{
		handler: handler,
	}
}

func Println[T any]() ConsumeHandler[T] {
	return func(data T) {
		fmt.Println(data)
	}
}

func (h *handlerConsumer[T]) Accept(data T) {
	if nil != h.handler {
		h.handler(data)
	}
}
