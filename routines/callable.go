/**
 * @author leo
 * @date 2020/8/27 8:52 下午
 */
package routines

import "errors"

//
// Callable
//  @Description: 类似java的callable
//
type Callable interface {
	//
	// Call
	//  @Description: 调用具体的函数
	//  @return interface{}
	//  @return error
	//
	Call() (interface{}, error)
}

//
// HandlerCallable
//  @Description: 默认的callable实现
//
type HandlerCallable struct {
	//
	//  Handler
	//  @Description: 具体的handler
	//
	Handler func() (interface{}, error)
}

func (callable *HandlerCallable) Call() (interface{}, error) {
	if nil != callable.Handler {
		return callable.Handler()
	}

	return nil, errors.New("HandlerCallable handler is nil")
}
