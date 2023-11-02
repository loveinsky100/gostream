/**
 * @author leo
 * @date 2020/10/25 9:46 下午
 */
package collections

import "sync"

type _BlockingQueue[T any] struct {
	data   Queue[T]
	signal chan int
	l      sync.Mutex
}

func NewBlockingQueue[T any](data Queue[T]) Queue[T] {
	return &_BlockingQueue[T]{
		data:   data,
		signal: make(chan int, 1), // buffered chan
	}
}

func (queue *_BlockingQueue[T]) Offer(v T) {
	queue.l.Lock()
	defer queue.l.Unlock()

	queue.data.Offer(v)
	if queue.Size() == 1 {
		queue.signal <- 1
	}
}

func (queue *_BlockingQueue[T]) Poll() (T, bool) {
	// 防止删除一个元素的时候写数据进来，保证Offer和Poll是串行的
	queue.l.Lock()
	data, ok := queue.data.Poll()
	if ok {
		queue.l.Unlock()
		return data, ok
	}

	queue.l.Unlock()
	select {
	case <-queue.signal:
		{
			return queue.data.Poll()
		}
	}
}

func (queue *_BlockingQueue[T]) Peek() (T, bool) {
	return queue.data.Peek()
}

func (queue *_BlockingQueue[T]) Foreach(consumer func(value T, index int) bool) {
	queue.data.Foreach(consumer)
}

func (queue *_BlockingQueue[T]) Size() int {
	return queue.data.Size()
}
