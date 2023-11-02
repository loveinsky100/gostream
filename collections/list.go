/**
 * @author leo
 * @date 2020/8/18 5:24 下午
 */
package collections

type List[T any] interface {
	//
	// Add
	//  @Description: 添加元素
	//  @param value
	//
	Add(value T)

	//
	// Get
	//  @Description: 获取数据
	//  @param index
	//  @return T
	//  @return bool
	//
	Get(index int) (T, bool)

	//
	// Remove
	//  @Description: 删除元素
	//  @param index
	//  @return bool
	//
	Remove(index int) bool

	//
	// Foreach
	//  @Description: 遍历
	//  @param consumer
	//
	Foreach(consumer func(value T, index int) bool)

	//
	// Size
	//  @Description: 获取数组大小
	//  @return int
	//
	Size() int

	//
	// Data
	//  @Description: 获取元数据
	//  @return []T
	//
	Data() []T
}
