/**
 * @author leo
 * @date 2023/4/17 11:47
 */
package streams

import (
	"github.com/loveinsky100/gostreams/routines"
	"time"
)

type StreamErrorHandler func(source string, err error)

type Stream[T any] interface {
	//
	// Debug
	//  @Description: debug，打印出流的结构信息
	//  @param action
	//  @return Stream[T]
	//
	Debug(action ConsumeHandler[string]) Stream[T]

	//
	// Error
	//  @Description: 设置错误监听
	//  @param listener
	//  @return Stream[T]
	//
	Error(listener StreamErrorHandler) Stream[T]

	//
	// Recover
	//  @Description: 尝试Recover，注意Parallel会默认Recover
	//  @return Stream[T]
	//
	Recover() Stream[T]

	//
	// Pool
	//  @Description: 设置并行池
	//  @param pool
	//  @return Stream[T]
	//
	Pool(pool routines.GoRoutinePool) Stream[T]

	//
	// Filter
	//  @Description: 过滤，流入的数据必须通过测试才能流出
	//  @param predicate
	//  @return Stream[T]
	//
	Filter(predicate PredicateHandler[T]) Stream[T]

	//
	// Skip
	//  @Description: 流入的数据跳过前面N个
	//  @param index
	//  @return Stream[T]
	//
	Skip(index int) Stream[T]

	//
	// Sort
	//  @Description: 对流入的数据进行排序
	//  @param comparator
	//  @return Stream[T]
	//
	Sort(comparator ComparatorHandler[T]) Stream[T]

	//
	// Limit
	//  @Description: 限制最多N个数据流出
	//  @param maxSize
	//  @return Stream[T]
	//
	Limit(maxSize int) Stream[T]

	//
	// Distinct
	//  @Description: 去除重复的数据
	//  @return Stream[T]
	//
	Distinct() Stream[T]

	//
	// DistinctWith
	//  @Description: 去除重复数组，自定义唯一函数
	//  @param unique
	//  @param interface{}]
	//  @return Stream[T]
	//
	DistinctWith(unique FunctionHandler[T, interface{}]) Stream[T]

	//
	// Parallel
	//  @Description: 接下来并行处理
	//  @return Stream[T]
	//
	Parallel() Stream[T]

	//
	// ParallelWith
	//  @Description: 接下来并行处理
	//  @param timeout
	//  @return Stream[T]
	//
	ParallelWith(timeout time.Duration) Stream[T]

	//
	// Sequential
	//  @Description: 串行处理
	//  @return Stream[T]
	//
	Sequential() Stream[T]

	//
	// Peek
	//  @Description: 对每一个流入的数据执行一次action
	//  @param action
	//  @return Stream[T]
	//
	Peek(action ConsumeHandler[T]) Stream[T]

	//
	// Foreach
	//  @Description: 遍历数据
	//  @param action
	//
	Foreach(action ConsumeHandler[T])

	//
	// AnyMatch
	//  @Description: 判断是否存在任意一个流入的数据满足条件
	//  @param predicate
	//  @return bool
	//
	AnyMatch(predicate PredicateHandler[T]) bool

	//
	// AllMatch
	//  @Description: 判断流入的数据是都满足条件
	//  @param predicate
	//  @return bool
	//
	AllMatch(predicate PredicateHandler[T]) bool

	//
	// NoneMatch
	//  @Description: 判断流入的数据是都不满足条件
	//  @param predicate
	//  @return bool
	//
	NoneMatch(predicate PredicateHandler[T]) bool

	//
	// FindFirst
	//  @Description: 返回第一个流入的数据
	//  @return T
	//
	FindFirst() T

	//
	// Max
	//  @Description: 返回最大的数据
	//  @param comparator
	//  @return T
	//
	Max(comparator ComparatorHandler[T]) T

	//
	// Min
	//  @Description: 返回最小的数据
	//  @param comparator
	//  @return T
	//
	Min(comparator ComparatorHandler[T]) T

	//
	// Count
	//  @Description: 统计数量
	//  @return int
	//
	Count() int

	//
	// Collect
	//  @Description: 收集数据
	//  @param collector
	//  @return interface{}
	//
	Collect(collector Collector[T, interface{}]) interface{}

	//
	// CollectToList
	//  @Description: 转化为数组
	//  @return []T
	//
	CollectToList() []T

	//
	// GroupBy
	//  @Description: 转化map
	//  @param keyMapper
	//  @param string]
	//  @return map[string]T
	//
	GroupBy(keyMapper FunctionHandler[T, string]) map[string]T

	//
	// Reduce
	//  @Description:
	//  @param reduce
	//  @return interface{}
	//
	Reduce(reduce Reduce[T, interface{}]) interface{}

	//
	// Map
	//  @Description: 转化数据
	//  @param mapper
	//  @return Stream[T]
	//
	Map(mapper FunctionHandler[T, T]) Stream[T]

	//
	// MapTo
	//  @Description: 转化数据，使用xxx.Map(Mapper[T, R]).(MapperStream[T, R]).xxx
	//  @param mapper
	//  @return interface{}
	//
	MapTo(mapper func(s Stream[T]) interface{}) interface{}

	//
	// Join
	//  @Description: 加入数据处理节点，非必要请勿使用
	//  @param pipeline
	//
	Join(pipeline Pipeline[T])
}

type MapperStream[T any, R any] interface {
	//
	// Map
	//  @Description: 转化数据，将流入的数据转化为其他数据
	//  @param R]
	//  @return Stream[R]
	//
	Map(mapper FunctionHandler[T, R]) Stream[R]

	//
	// FlatMap
	//  @Description: 转化数据，将流入的数据转化为其他数据
	//  @param mapper
	//  @return Stream[R]
	//
	FlatMap(mapper FunctionHandler[T, []R]) Stream[R]
}
