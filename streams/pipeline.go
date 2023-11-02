/**
 * @author leo
 * @date 2023/4/17 14:00
 */
package streams

//
// Sink
//  @Description: 处理器
//
type Sink[T any] interface {
	Consumer[T]

	//
	// Begin
	//  @Description: 开始处理
	//  @param size
	//
	Begin(size int)

	//
	// End
	//  @Description: 结束处理
	//
	End()

	//
	// CancellationRequested
	//  @Description: 是否停止
	//  @return bool
	//
	CancellationRequested() bool
}

type PipelineDelegate[T any] interface {
	//
	// OnError
	//  @Description: 发生错误
	//  @param pipeline
	//  @param err
	//
	OnError(source interface{}, err error)
}

//
// Pipeline
//  @Description: 流水线
//
type Pipeline[T any] interface {
	Sink[T]
	SetDelegate(delegate PipelineDelegate[T])
	SetNext(pipeline Pipeline[T])
	Stateless() bool
	PushError(source interface{}, err error)
}

//
// StatelessOp
//  @Description: 无状态
//
type StatelessOp[T any] interface {
	Pipeline[T]
}

//
// StatefulOp
//  @Description: 有状态
//
type StatefulOp[T any] interface {
	Pipeline[T]
}

//
// FinishOp
//  @Description: 结束操作
//
type FinishOp[T any] interface {
	Pipeline[T]
}

//
// ResultFinishOp
//  @Description: 返回结果的结束操作
//
type ResultFinishOp[T any, R any] interface {
	FinishOp[T]
	Result() R
}
