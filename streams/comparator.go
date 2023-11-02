/**
 * @author leo
 * @date 2023/4/17 12:07
 */
package streams

type Comparator[T any] interface {
	//
	// Compare
	//  @Description: 比较两个数据, o1 > o2: true, o1 <= o2: false
	//  @param o1
	//  @param o2
	//  @return bool
	//
	Compare(o1 T, o2 T) bool
}
