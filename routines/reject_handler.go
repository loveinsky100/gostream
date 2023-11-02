/**
 * @author leo
 * @date 2020/8/27 9:43 下午
 */
package routines

import "errors"

type RejectedHandler interface {
	Reject(callable Callable) error
}

type DefaultRejectedHandler struct {
	Handler func(callable Callable) error
}

func (rejectedHandler *DefaultRejectedHandler) Reject(callable Callable) error {
	if nil != rejectedHandler.Handler {
		return rejectedHandler.Handler(callable)
	}

	return errors.New("add pool failed due to pool full")
}

type CallableRejectedHandler struct {
	Handler func(callable Callable)
}

func (rejectedHandler *CallableRejectedHandler) Reject(callable Callable) error {
	if nil != rejectedHandler.Handler {
		rejectedHandler.Handler(callable)
		return errors.New("add pool failed due to pool full")
	}

	return nil
}
