/**
 * @author leo
 * @date 2020/8/26 3:57 下午
 */
package routines

import (
	"errors"
	"time"
)

type FutureWrapper[T any] struct {
	Future
}

func Wrapper[T any](future Future) *FutureWrapper[T] {
	return &FutureWrapper[T]{
		future,
	}
}

func (f *FutureWrapper[T]) Get() (T, error) {
	return f.GetWithTimeout(0)
}

func (f *FutureWrapper[T]) GetWithTimeout(timeout time.Duration) (T, error) {
	data, err := f.Future.GetWithTimeout(timeout)
	if nil != err {
		var t T
		return t, err
	}

	result, ok := data.(T)
	if !ok {
		var t T
		return t, errors.New("type T error")
	}

	return result, nil
}
