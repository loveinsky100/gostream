/**
 * @author leo
 * @date 2023/4/17 11:47
 */
package streams

import (
	"encoding/json"
	"fmt"
	"github.com/loveinsky100/gostreams/routines"
	"reflect"
	"time"
)

type streamImpl[T any] struct {
	runner          func()
	current         Pipeline[T]
	parallel        bool
	requireParallel bool
	parallelGroup   routines.GoRoutineGroup
	debug           bool
	tracker         []string
	debugger        Consumer[string]
	errorHandler    StreamErrorHandler
	pool            routines.GoRoutinePool
}

func Of[T any](args ...T) Stream[T] {
	header := NewSliceHeader[T]()
	return &streamImpl[T]{
		runner: func() {
			header.Begin(len(args))
			for _, data := range args {
				if header.CancellationRequested() {
					break
				}

				header.Accept(data)
			}

			header.End()
		},
		current: header,
	}
}

func SliceOf[T any](args []T) Stream[T] {
	return Of(args...)
}

func MapOf[K comparable, V any](data map[K]V) Stream[Entry[K, V]] {
	header := NewMapHeader[K, V]()
	return &streamImpl[Entry[K, V]]{
		runner: func() {
			header.Begin(len(data))
			for k, v := range data {
				if header.CancellationRequested() {
					break
				}

				entry := NewEntry(k, v)

				header.Accept(entry)
			}

			header.End()
		},
		current: header,
	}
}

func (s *streamImpl[T]) Debug(action ConsumeHandler[string]) Stream[T] {
	if s.debug {
		return s
	}

	s.debug = true
	s.tracker = make([]string, 0)
	s.debugger = NewHandlerConsumer(action)
	return s
}

func (s *streamImpl[T]) Error(errorHandler StreamErrorHandler) Stream[T] {
	s.errorHandler = errorHandler
	return s
}

func (s *streamImpl[T]) Pool(pool routines.GoRoutinePool) Stream[T] {
	s.pool = pool
	return s
}

func (s *streamImpl[T]) Recover() Stream[T] {
	s.Join(NewRecoverPipeline[T]())
	return s
}

func (s *streamImpl[T]) Parallel() Stream[T] {
	return s.ParallelWith(0)
}

func (s *streamImpl[T]) ParallelWith(timeout time.Duration) Stream[T] {
	if s.parallel {
		return s
	}

	s.parallel = true
	s.requireParallel = false
	pipeline, group := NewParallelPipeline[T](timeout)
	s.simpleJoin(pipeline)
	s.parallelGroup = group
	return s
}

func (s *streamImpl[T]) Sequential() Stream[T] {
	return s.sequentialWithRequire(false)
}

func (s *streamImpl[T]) sequentialWithRequire(require bool) Stream[T] {
	if !s.parallel {
		return s
	}

	s.requireParallel = require
	s.parallel = false
	_group := s.parallelGroup
	s.simpleJoin(NewSequentialPipeline[T](func() error {
		if nil == s.pool {
			return _group.Execute()
		}

		return _group.ExecuteInPool(s.pool)
	}))

	return s
}

func (s *streamImpl[T]) Filter(predicate PredicateHandler[T]) Stream[T] {
	s.Join(NewFilterPipeline(predicate))
	return s
}

func (s *streamImpl[T]) Skip(index int) Stream[T] {
	// 跳过即过滤出非此条件的
	s.Join(NewSkipPipeline[T](index))
	return s
}

func (s *streamImpl[T]) Sort(comparator ComparatorHandler[T]) Stream[T] {
	s.Join(NewSortPipeline(comparator))
	return s
}

func (s *streamImpl[T]) Limit(maxSize int) Stream[T] {
	s.Join(NewLimitPipeline[T](maxSize))
	return s
}

func (s *streamImpl[T]) Peek(action ConsumeHandler[T]) Stream[T] {
	s.Join(NewPeekPipeline[T](action))
	return s
}

func (s *streamImpl[T]) Distinct() Stream[T] {
	s.Join(NewDistinctPipeline[T](func(t T) interface{} {
		return t
	}))
	return s
}

func (s *streamImpl[T]) DistinctWith(unique FunctionHandler[T, interface{}]) Stream[T] {
	s.Join(NewDistinctPipeline[T](unique))
	return s
}

func (s *streamImpl[T]) Map(mapper FunctionHandler[T, T]) Stream[T] {
	s.Join(NewMapperPipeline[T](mapper))
	return s
}

func (s *streamImpl[T]) MapTo(mapper func(s Stream[T]) interface{}) interface{} {
	return mapper(s)
}

func (s *streamImpl[T]) Foreach(action ConsumeHandler[T]) {
	s.run(NewForeachPipeline(action))
}

func (s *streamImpl[T]) AnyMatch(predicate PredicateHandler[T]) bool {
	p := NewMatchAnyPipeline(predicate)
	s.run(p)
	return p.Result()
}

func (s *streamImpl[T]) AllMatch(predicate PredicateHandler[T]) bool {
	p := NewMatchAllPipeline(predicate)
	s.run(p)
	return p.Result()
}
func (s *streamImpl[T]) NoneMatch(predicate PredicateHandler[T]) bool {
	p := NewMatchNonePipeline(predicate)
	s.run(p)
	return p.Result()
}

func (s *streamImpl[T]) FindFirst() T {
	p := NewFindFirstPipeline[T]()
	s.run(p)
	return p.Result()
}

func (s *streamImpl[T]) Collect(collector Collector[T, interface{}]) interface{} {
	p := NewCollectPipeline[T, interface{}](collector)
	s.run(p)
	return p.Result()
}

func (s *streamImpl[T]) CollectToList() []T {
	p := NewCollectPipeline[T, []T](ToList[T]())
	s.run(p)
	return p.Result()
}

func (s *streamImpl[T]) GroupBy(keyMapper FunctionHandler[T, string]) map[string]T {
	p := NewCollectPipeline[T, map[string]T](ToMap[T, string, T](keyMapper, Identity[T]()))
	s.run(p)
	return p.Result()
}

func (s *streamImpl[T]) Reduce(reduce Reduce[T, interface{}]) interface{} {
	p := NewReducePipeline[T, interface{}](reduce)
	s.run(p)
	return p.Result()
}

func (s *streamImpl[T]) Max(comparator ComparatorHandler[T]) T {
	p := NewMaxPipeline[T](comparator)
	s.run(p)
	return p.Result()
}

func (s *streamImpl[T]) Min(comparator ComparatorHandler[T]) T {
	p := NewMinPipeline[T](comparator)
	s.run(p)
	return p.Result()
}

func (s *streamImpl[T]) Count() int {
	p := NewCountPipeline[T]()
	s.run(p)
	return p.Result()
}

func (s *streamImpl[T]) Join(pipeline Pipeline[T]) {
	if !s.parallel && !s.requireParallel {
		s.simpleJoin(pipeline)
		return
	}

	// 并行流，当加入有状态节点时需要暂停并行
	if !pipeline.Stateless() {
		// 加入一个串行流，并标记之后需要转回并行
		s.sequentialWithRequire(true)
		s.simpleJoin(pipeline)
	} else {
		// 判断是否需要转回并行
		if s.requireParallel {
			s.requireParallel = false
			s.Parallel()
		}
		s.simpleJoin(pipeline)
	}
}

func (s *streamImpl[T]) run(pipeline FinishOp[T]) {
	if !s.parallel {
		s.simpleRun(pipeline)
		return
	}

	s.Sequential()
	s.simpleRun(pipeline)
}

func (s *streamImpl[T]) simpleJoin(pipeline Pipeline[T]) {
	pipeline.SetDelegate(s)
	s.current.SetNext(pipeline)
	s.current = pipeline
	if s.debug {
		s.tracker = append(s.tracker, s.pipelineName(pipeline))
	}
}

func (s *streamImpl[T]) simpleRun(pipeline FinishOp[T]) {
	pipeline.SetDelegate(s)
	s.current.SetNext(pipeline)
	s.current = pipeline
	s.runner()
	if s.debug {
		s.tracker = append(s.tracker, s.pipelineName(pipeline))
	}

	if nil != s.debugger {
		s.debugger.Accept(s.String())
	}
}

func (s *streamImpl[T]) pipelineName(pipeline Pipeline[T]) string {
	p := reflect.TypeOf(pipeline)
	if p.Kind() == reflect.Ptr {
		p = p.Elem()
	}

	state := "stateful"
	if pipeline.Stateless() {
		state = "stateless"
	}
	return fmt.Sprintf("%s{%s}", p.Name(), state)
}

func (s *streamImpl[T]) sourceName(source interface{}) string {
	p := reflect.TypeOf(source)
	if p.Kind() == reflect.Ptr {
		p = p.Elem()
	}

	return p.Name()
}

func (s *streamImpl[T]) String() string {
	val, _ := json.Marshal(s.tracker)
	return string(val)
}

func (s *streamImpl[T]) OnError(source interface{}, err error) {
	if nil != s.errorHandler {
		s.errorHandler(s.sourceName(source), err)
	}
}
