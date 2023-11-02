/**
 * @author leo
 * @date 2020/10/25 3:58 下午
 */
package collections

type Node[T any] struct {
	// next node
	Next *Node[T]

	// node value
	Value T
}
