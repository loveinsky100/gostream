/**
 * @author leo
 * @date 2020/10/25 3:52 下午
 */
package collections

type Queue[T any] interface {
	//
	// Offer
	//  @Description: 往队列中塞入元素
	//  @param v
	//
	Offer(v T)

	//
	// Poll
	//  @Description: 取出队列中的第一个元素，并删除
	//  @return T
	//  @return bool
	//
	Poll() (T, bool)

	//
	// Peek
	//  @Description: 取出队列中的第一个元素
	//  @return T
	//  @return bool
	//
	Peek() (T, bool)

	//
	// Foreach
	//  @Description: loop items
	//  @param consumer
	//
	Foreach(consumer func(value T, index int) bool)

	//
	// Size
	//  @Description: queue size
	//  @return int
	//
	Size() int
}
