/**
 * @author leo
 * @date 2023/4/17 12:00
 */
package streams

type FunctionHandler[T any, R any] func(data T) R

type handlerFunction[T any, R any] struct {
	handler FunctionHandler[T, R]
}

func NewHandlerFunction[T any, R any](handler FunctionHandler[T, R]) Function[T, R] {
	return &handlerFunction[T, R]{
		handler: handler,
	}
}

func (h *handlerFunction[T, R]) Apply(data T) R {
	return h.handler(data)
}
