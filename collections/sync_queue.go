/**
 * @author leo
 * @date 2020/10/25 3:57 下午
 */
package collections

import "sync"

type _SyncQueue[T any] struct {
	queue Queue[T]
	l     sync.RWMutex
}

func NewSyncQueue[T any](queue Queue[T]) Queue[T] {
	return &_SyncQueue[T]{queue: queue}
}

func (queue *_SyncQueue[T]) Offer(v T) {
	queue.l.Lock()
	defer queue.l.Unlock()

	queue.queue.Offer(v)
}

func (queue *_SyncQueue[T]) Poll() (T, bool) {
	queue.l.Lock()
	defer queue.l.Unlock()

	return queue.queue.Poll()
}

func (queue *_SyncQueue[T]) Peek() (T, bool) {
	queue.l.RLocker()
	defer queue.l.RUnlock()

	return queue.queue.Peek()
}

func (queue *_SyncQueue[T]) Foreach(consumer func(value T, index int) bool) {
	queue.l.RLocker()
	defer queue.l.RUnlock()

	queue.queue.Foreach(consumer)
}

func (queue *_SyncQueue[T]) Size() int {
	queue.l.RLocker()
	defer queue.l.RUnlock()

	return queue.queue.Size()
}
