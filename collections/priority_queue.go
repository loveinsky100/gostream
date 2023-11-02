/**
 * @author leo
 * @date 2020/11/20 4:53 下午
 */
package collections

/**
 * @Description: 优先级比较器
 */
type Comparator interface {
	/**
	 * @Description: 比较两者的优先级
	 * @param comparator1
	 * @param comparator2
	 * @return int 0: 两者优先级相同，大于0: 当前的优先级高，小于0: comparator优先级高
	 */
	compare(comparator Comparator) int
}

type _PriorityQueue[T Comparator] struct {
	_LinkedQueue[T]
}

/**
 * @Description: 返回优先级队列，写入的数据必须实现Comparator接口
 * @return Queue
 */
func NewPriorityQueue[T Comparator]() Queue[T] {
	return &_PriorityQueue[T]{
		_LinkedQueue[T]{},
	}
}

func (queue *_PriorityQueue[T]) Offer(v T) {
	node := &Node[T]{
		Next:  nil,
		Value: v,
	}

	if nil == queue.root {
		queue.root = node
		return
	}

	pre := queue.root
	current := queue.root
	for nil != current {
		v2 := current.Value
		priority := v.compare(v2)
		if priority > 0 {
			break
		}

		pre = current
		current = current.Next
	}

	pre.Next = node
	pre.Next.Next = current
	queue.count++
}

func (queue *_PriorityQueue[T]) Poll() (T, bool) {
	return queue._LinkedQueue.Poll()
}

func (queue *_PriorityQueue[T]) Peek() (T, bool) {
	return queue._LinkedQueue.Peek()
}

func (queue *_PriorityQueue[T]) Foreach(consumer func(value T, index int) bool) {
	queue._LinkedQueue.Foreach(consumer)
}

func (queue *_PriorityQueue[T]) Size() int {
	return queue._LinkedQueue.Size()
}
