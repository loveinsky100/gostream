/**
 * @author leo
 * @date 2020/10/25 3:57 下午
 */
package collections

type _LinkedQueue[T any] struct {
	// query items
	root *Node[T]

	// last node
	last *Node[T]

	// count items
	count int
}

func NewLinkedQueue[T any]() Queue[T] {
	return &_LinkedQueue[T]{}
}

func (queue *_LinkedQueue[T]) Offer(v T) {
	node := &Node[T]{
		Next:  nil,
		Value: v,
	}

	if nil == queue.root {
		queue.root = node
	} else {
		queue.last.Next = node
	}

	queue.count++
	queue.last = node
}

func (queue *_LinkedQueue[T]) Poll() (T, bool) {
	node := queue.pollNode()
	if nil == node {
		var zero T
		return zero, false
	}

	return node.Value, true
}

func (queue *_LinkedQueue[T]) Peek() (T, bool) {
	node := queue.peekNode()
	if nil == node {
		var zero T
		return zero, false
	}

	return node.Value, true
}

func (queue *_LinkedQueue[T]) Foreach(consumer func(value T, index int) bool) {
	if nil == queue.root || nil == consumer {
		return
	}

	current := queue.root
	index := 0
	for nil != current {
		ok := consumer(current.Value, index)
		if !ok {
			break
		}

		index++
		current = current.Next
	}
}

func (queue *_LinkedQueue[T]) Size() int {
	return queue.count
}

func (queue *_LinkedQueue[T]) pollNode() *Node[T] {
	root := queue.peekNode()
	if nil == root {
		return root
	}

	next := root.Next
	queue.root = next
	if nil == queue.root {
		queue.last = nil
	}

	queue.count--
	return root
}

func (queue *_LinkedQueue[T]) peekNode() *Node[T] {
	return queue.root
}
