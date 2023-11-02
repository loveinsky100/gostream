/**
 * @author leo
 * @date 2020/8/18 5:32 下午
 */
package collections

type Map[K comparable, V any] interface {
	//
	// Put
	//  @Description: 添加元素
	//  @param key
	//  @param value
	//
	Put(key K, value V)

	//
	// Remove
	//  @Description: 移除元素
	//  @param key
	//
	Remove(key K)

	//
	// Get
	//  @Description: 获取数据
	//  @param key
	//  @return V
	//  @return bool
	//
	Get(key K) (V, bool)

	//
	// Foreach
	//  @Description: 遍历数据
	//  @param consumer
	//
	Foreach(consumer func(key K, value V) bool)
}
