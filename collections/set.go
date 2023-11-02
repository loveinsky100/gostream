/**
 * @author leo
 * @date 2020/8/18 5:54 下午
 */
package collections

type Set[K any] interface {

	//
	// Add
	//  @Description: 添加
	//  @param key
	//
	Add(key K)

	//
	// Remove
	//  @Description: 移除
	//  @param key
	//
	Remove(key K)

	//
	// Contains
	//  @Description: 包含
	//  @param key
	//  @return bool
	//
	Contains(key K) bool

	//
	// Foreach
	//  @Description: 遍历
	//  @param consumer
	//
	Foreach(consumer func(key K) bool)
}
